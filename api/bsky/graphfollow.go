// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package bsky

// schema: app.bsky.graph.follow

import (
	"github.com/KingYoSun/indigo/lex/util"
)

func init() {
	util.RegisterType("app.bsky.graph.follow", &GraphFollow{})
} //
// RECORDTYPE: GraphFollow
type GraphFollow struct {
	LexiconTypeID string `json:"$type,const=app.bsky.graph.follow" cborgen:"$type,const=app.bsky.graph.follow"`
	CreatedAt     string `json:"createdAt" cborgen:"createdAt"`
	Subject       string `json:"subject" cborgen:"subject"`
}
