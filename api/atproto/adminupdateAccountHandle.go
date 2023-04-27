// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.admin.updateAccountHandle

import (
	"context"

	"github.com/KingYoSun/indigo/xrpc"
)

// AdminUpdateAccountHandle_Input is the input argument to a com.atproto.admin.updateAccountHandle call.
type AdminUpdateAccountHandle_Input struct {
	Did    string `json:"did" cborgen:"did"`
	Handle string `json:"handle" cborgen:"handle"`
}

// AdminUpdateAccountHandle calls the XRPC method "com.atproto.admin.updateAccountHandle".
func AdminUpdateAccountHandle(ctx context.Context, c *xrpc.Client, input *AdminUpdateAccountHandle_Input) error {
	if err := c.Do(ctx, xrpc.Procedure, "application/json", "com.atproto.admin.updateAccountHandle", nil, input, nil); err != nil {
		return err
	}

	return nil
}
