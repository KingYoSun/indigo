// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.server.revokeAppPassword

import (
	"context"

	"github.com/KingYoSun/indigo/xrpc"
)

// ServerRevokeAppPassword_Input is the input argument to a com.atproto.server.revokeAppPassword call.
type ServerRevokeAppPassword_Input struct {
	Name string `json:"name" cborgen:"name"`
}

// ServerRevokeAppPassword calls the XRPC method "com.atproto.server.revokeAppPassword".
func ServerRevokeAppPassword(ctx context.Context, c *xrpc.Client, input *ServerRevokeAppPassword_Input) error {
	if err := c.Do(ctx, xrpc.Procedure, "application/json", "com.atproto.server.revokeAppPassword", nil, input, nil); err != nil {
		return err
	}

	return nil
}
