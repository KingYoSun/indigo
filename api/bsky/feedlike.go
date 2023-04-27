// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package bsky

// schema: app.bsky.feed.like

import (
	comatprototypes "github.com/KingYoSun/indigo/api/atproto"
	"github.com/KingYoSun/indigo/lex/util"
)

func init() {
	util.RegisterType("app.bsky.feed.like", &FeedLike{})
} //
// RECORDTYPE: FeedLike
type FeedLike struct {
	LexiconTypeID string                         `json:"$type,const=app.bsky.feed.like" cborgen:"$type,const=app.bsky.feed.like"`
	CreatedAt     string                         `json:"createdAt" cborgen:"createdAt"`
	Subject       *comatprototypes.RepoStrongRef `json:"subject" cborgen:"subject"`
}
