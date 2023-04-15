package bsky

import (
	comatprototypes "github.com/KingYoSun/indigo/api/atproto"
	"github.com/KingYoSun/indigo/lex/util"
)

// schema: app.bsky.feed.repost

func init() {
	util.RegisterType("app.bsky.feed.repost", &FeedRepost{})
}

// RECORDTYPE: FeedRepost
type FeedRepost struct {
	LexiconTypeID string                         `json:"$type,const=app.bsky.feed.repost" cborgen:"$type,const=app.bsky.feed.repost"`
	CreatedAt     string                         `json:"createdAt" cborgen:"createdAt"`
	Subject       *comatprototypes.RepoStrongRef `json:"subject" cborgen:"subject"`
}
