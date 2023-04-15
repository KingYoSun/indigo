package bgs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/models"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"

	"github.com/ipfs/go-cid"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm"
)

type MeiliFeedPost struct {
	Cid 	string						`json:"cid"`
	Tid		string						`json:"tid"`
	Post	*appbsky.FeedPost	`json:"post"`
	User	*User							`json:"user"`
}

type MeiliSlurper struct {
	meili *meilisearch.Client

	db *gorm.DB

	repoman *repomgr.RepoManager

	lk     sync.Mutex
	active map[string]*models.PDS
}

func NewMeiliSlurper(db *gorm.DB, meili *meilisearch.Client, repoman *repomgr.RepoManager) *MeiliSlurper {
	return &MeiliSlurper{
		meili: 		meili,
		db:     	db,
		repoman: 	repoman,
		active: make(map[string]*models.PDS),
	}
}

func (s *MeiliSlurper) PdsToMeili(ctx context.Context, host string, reg bool) error {
	// TODO: for performance, lock on the hostname instead of global
	s.lk.Lock()
	defer s.lk.Unlock()

	_, ok := s.active[host]
	if ok {
		return nil
	}

	var peering models.PDS
	if err := s.db.Find(&peering, "host = ?", host).Error; err != nil {
		return err
	}

	if peering.ID == 0 {
		return errors.New("PDS is not found")
	}

	if !peering.Registered && reg {
		peering.Registered = true
		if err := s.db.Model(models.PDS{}).Where("id = ?", peering.ID).Update("registered", true).Error; err != nil {
			return err
		}
	}

	s.active[host] = &peering

	go s.copyRecordsToMeili(ctx, &peering)

	return nil
}


func (s *MeiliSlurper) copyRecordsToMeili(ctx context.Context, host *models.PDS) {
	defer func() {
		s.lk.Lock()
		defer s.lk.Unlock()

		delete(s.active, host.Host)
	}()

	rows, err := s.db.
		Debug().
		Model(&models.FeedPost{}).
		Select("feed_posts.*, users.id, users.pds").
		Joins("JOIN users ON feed_posts.author = users.id AND users.pds = ?", int64(host.ID)).
		Order("feed_posts.id ASC").
		Where("feed_posts.missing = ?", false).
		Rows()
	if err != nil {
		log.Errorf("[PdsToMeili] %v", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var feedPost *models.FeedPost

		s.db.ScanRows(rows, &feedPost)

		var user *User

		if feedPost.Cid == "" || feedPost.Missing {
			log.Errorf("[PdsToMeili] feed_post is missing")
			continue
		}

		buf := new(bytes.Buffer)
		targetCid, err := cid.Decode(feedPost.Cid)
		if err != nil {
			log.Errorf("[PdsToMeili] %v", err.Error())
			continue
		}

		if err := s.db.Find(&user, "id = ?", feedPost.Author).Error; err != nil {
			log.Errorf("[PdsToMeili] %v", err.Error())
			continue
		}

		if err := s.repoman.ReadRepoAtCid(ctx, feedPost.Author, targetCid, buf); err != nil {
			log.Errorf("[PdsToMeili] %v", err.Error())
			continue
		}

		sliceRepo, err := repo.ReadRepoFromCar(ctx, buf)
		if err != nil {
			log.Errorf("[PdsToMeili] %v", err.Error())
			continue
		}
		rpath := "app.bsky.feed.post/" + feedPost.Rkey

		_, rec, err := sliceRepo.GetRecord(ctx, rpath)
		if err != nil {
			log.Errorf("[PdsToMeili] %v", err.Error())
			continue
		}

		post, suc := rec.(*appbsky.FeedPost)
		if !suc {
			log.Errorf("[PdsToMeili] failed to deserialize post")
			continue
		}

		document := &MeiliFeedPost{
			Cid: 	feedPost.Cid,
			Tid: 	"app.bsky.feed.post/" + feedPost.Rkey,
			Post:	post,
			User:	user,
		}

		encoded, err := json.Marshal(document)
		if err != nil {
			log.Errorf("[PdsToMeili] %v", err.Error())
			continue
		}

		if _, err = s.meili.Index("feed_posts").AddDocuments(document, "cid"); err != nil {
			log.Errorf("[PdsToMeili] %v, %s", err.Error(), encoded)
			continue
		}
	}
}