package atproto

import (
	"context"

	"github.com/KingYoSun/indigo/xrpc"
)

// schema: com.atproto.server.deleteSession

func init() {
}
func ServerDeleteSession(ctx context.Context, c *xrpc.Client) error {
	if err := c.Do(ctx, xrpc.Procedure, "", "com.atproto.server.deleteSession", nil, nil, nil); err != nil {
		return err
	}

	return nil
}
