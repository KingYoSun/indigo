// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package bsky

// schema: app.bsky.graph.muteActor

import (
	"context"

	"github.com/KingYoSun/indigo/xrpc"
)

// GraphMuteActor_Input is the input argument to a app.bsky.graph.muteActor call.
type GraphMuteActor_Input struct {
	Actor string `json:"actor" cborgen:"actor"`
}

// GraphMuteActor calls the XRPC method "app.bsky.graph.muteActor".
func GraphMuteActor(ctx context.Context, c *xrpc.Client, input *GraphMuteActor_Input) error {
	if err := c.Do(ctx, xrpc.Procedure, "application/json", "app.bsky.graph.muteActor", nil, input, nil); err != nil {
		return err
	}

	return nil
}
