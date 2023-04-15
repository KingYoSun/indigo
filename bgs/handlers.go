package bgs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	comatprototypes "github.com/KingYoSun/indigo/api/atproto"
	appbsky "github.com/KingYoSun/indigo/api/bsky"
	"github.com/KingYoSun/indigo/repo"
	"github.com/ipfs/go-cid"
)

func (s *BGS) handleComAtprotoSyncGetCheckout(ctx context.Context, commit string, did string) (io.Reader, error) {
	/*
		u, err := s.Index.LookupUserByDid(ctx, did)
		if err != nil {
			return nil, err
		}

		c, err := cid.Decode(commit)
		if err != nil {
			return nil, err
		}

		// TODO: need to enable a 'write to' interface for codegenned things...
		buf := new(bytes.Buffer)
		if err := s.repoman.GetCheckout(ctx, u.Uid, c, buf); err != nil {
			return nil, err
		}

		return buf, nil
	*/
	panic("nyi")
}

func (s *BGS) handleComAtprotoSyncGetCommitPath(ctx context.Context, did string, earliest string, latest string) (*comatprototypes.SyncGetCommitPath_Output, error) {
	panic("nyi")
}

func (s *BGS) handleComAtprotoSyncGetHead(ctx context.Context, did string) (*comatprototypes.SyncGetHead_Output, error) {
	u, err := s.Index.LookupUserByDid(ctx, did)
	if err != nil {
		return nil, err
	}

	root, err := s.repoman.GetRepoRoot(ctx, u.Uid)
	if err != nil {
		return nil, err
	}

	return &comatprototypes.SyncGetHead_Output{
		Root: root.String(),
	}, nil
}

func (s *BGS) handleComAtprotoSyncGetRecord(ctx context.Context, collection string, commit string, did string, rkey string) (io.Reader, error) {
	panic("nyi")
}

func (s *BGS) handleComAtprotoSyncGetRepo(ctx context.Context, did string, earliest string, latest string) (io.Reader, error) {
	u, err := s.Index.LookupUserByDid(ctx, did)
	if err != nil {
		return nil, err
	}

	var earlyCid, lateCid cid.Cid
	if earliest != "" {
		c, err := cid.Decode(earliest)
		if err != nil {
			return nil, err
		}

		earlyCid = c
	}

	if latest != "" {
		c, err := cid.Decode(latest)
		if err != nil {
			return nil, err
		}

		lateCid = c
	}

	// TODO: stream the response
	buf := new(bytes.Buffer)
	if err := s.repoman.ReadRepo(ctx, u.Uid, earlyCid, lateCid, buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *BGS) handleComAtprotoSyncGetBlocks(ctx context.Context, cids []string, did string) (io.Reader, error) {
	panic("NYI")
}

func (s *BGS) handleComAtprotoSyncRequestCrawl(ctx context.Context, host string) error {
	if host == "" {
		return fmt.Errorf("must pass valid hostname")
	}

	log.Warnf("TODO: host validation for crawl requests")
	return s.slurper.SubscribeToPds(ctx, host, true)
}

func (s *BGS) handleComAtprotoSyncNotifyOfUpdate(ctx context.Context, hostname string) error {
	panic("NYI")
	//return s.slurper.SubscribeToPds(ctx, host, false)
}

func (s *BGS) handleComAtprotoSyncGetBlob(ctx context.Context, cid string, did string) (io.Reader, error) {
	if s.blobs == nil {
		return nil, fmt.Errorf("blob store disabled")
	}

	b, err := s.blobs.GetBlob(ctx, cid, did)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func (s *BGS) handleComAtprotoSyncListBlobs(ctx context.Context, did string, earliest string, latest string) (*comatprototypes.SyncListBlobs_Output, error) {
	panic("NYI")
}

func (s *BGS) handleDebugGetRepoJson(ctx context.Context, did string, bcid string, rkey string) ([]byte, error) {
	u, err := s.Index.LookupUserByDid(ctx, did)
	if err != nil {
		return nil, err
	}

	var decodedCid cid.Cid
	if bcid != "" {
		c, err := cid.Decode(bcid)
		if err != nil {
			return nil, err
		}

		decodedCid = c
	}

	// TODO: stream the response
	buf := new(bytes.Buffer)
	if err := s.repoman.ReadRepoAtCid(ctx, u.Uid, decodedCid, buf); err != nil {
		return nil, err
	}

	sliceRepo, err := repo.ReadRepoFromCar(ctx, buf)
	if err != nil {
		return nil, err
	}
	rpath := "app.bsky.feed.post/" + rkey

	_, rec, err := sliceRepo.GetRecord(ctx, rpath)
	if err != nil {
		return nil, err
	}

	post, suc := rec.(*appbsky.FeedPost)
	if !suc {
		return nil, errors.New("failed to deserialize post")
	}

	postJson, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	return postJson, nil
}

func (s *BGS) handleComAtprotoSyncListRepos(ctx context.Context, cursor string, limit int) (*comatprototypes.SyncListRepos_Output, error) {
	panic("NYI")
}

func (s *BGS) handleMeiliRequestCopyRecord(ctx context.Context, host string) error {
	if host == "" {
		return fmt.Errorf("must pass valid hostname")
	}

	return s.meilislur.PdsToMeili(ctx, host, true)
}
