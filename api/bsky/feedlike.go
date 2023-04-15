package bsky

import (
	comatprototypes "github.com/KingYoSun/indigo/api/atproto"
	"github.com/KingYoSun/indigo/lex/util"
)

// schema: app.bsky.feed.like

func init() {
	util.RegisterType("app.bsky.feed.like", &FeedLike{})
}

// RECORDTYPE: FeedLike
type FeedLike struct {
	LexiconTypeID string                         `json:"$type,const=app.bsky.feed.like" cborgen:"$type,const=app.bsky.feed.like"`
	CreatedAt     string                         `json:"createdAt" cborgen:"createdAt"`
	Subject       *comatprototypes.RepoStrongRef `json:"subject" cborgen:"subject"`
}
