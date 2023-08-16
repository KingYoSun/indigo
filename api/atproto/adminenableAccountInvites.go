// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.admin.enableAccountInvites

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// AdminEnableAccountInvites_Input is the input argument to a com.atproto.admin.enableAccountInvites call.
type AdminEnableAccountInvites_Input struct {
	Account string `json:"account" cborgen:"account"`
	// note: Additionally add a note describing why the invites were enabled
	Note *string `json:"note,omitempty" cborgen:"note,omitempty"`
}

// AdminEnableAccountInvites calls the XRPC method "com.atproto.admin.enableAccountInvites".
func AdminEnableAccountInvites(ctx context.Context, c *xrpc.Client, input *AdminEnableAccountInvites_Input) error {
	if err := c.Do(ctx, xrpc.Procedure, "application/json", "com.atproto.admin.enableAccountInvites", nil, input, nil); err != nil {
		return err
	}

	return nil
}
