package atproto

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// schema: com.atproto.server.createInviteCodes

func init() {
}

type ServerCreateInviteCodes_AccountCodes struct {
	Account string   `json:"account" cborgen:"account"`
	Codes   []string `json:"codes" cborgen:"codes"`
}

type ServerCreateInviteCodes_Input struct {
	CodeCount   int64    `json:"codeCount" cborgen:"codeCount"`
	ForAccounts []string `json:"forAccounts,omitempty" cborgen:"forAccounts,omitempty"`
	UseCount    int64    `json:"useCount" cborgen:"useCount"`
}

type ServerCreateInviteCodes_Output struct {
	Codes []*ServerCreateInviteCodes_AccountCodes `json:"codes" cborgen:"codes"`
}

func ServerCreateInviteCodes(ctx context.Context, c *xrpc.Client, input *ServerCreateInviteCodes_Input) (*ServerCreateInviteCodes_Output, error) {
	var out ServerCreateInviteCodes_Output
	if err := c.Do(ctx, xrpc.Procedure, "application/json", "com.atproto.server.createInviteCodes", nil, input, &out); err != nil {
		return nil, err
	}

	return &out, nil
}