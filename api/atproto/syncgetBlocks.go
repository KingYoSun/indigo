// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.sync.getBlocks

import (
	"bytes"
	"context"

	"github.com/KingYoSun/indigo/xrpc"
)

// SyncGetBlocks calls the XRPC method "com.atproto.sync.getBlocks".
//
// did: The DID of the repo.
func SyncGetBlocks(ctx context.Context, c *xrpc.Client, cids []string, did string) ([]byte, error) {
	buf := new(bytes.Buffer)

	params := map[string]interface{}{
		"cids": cids,
		"did":  did,
	}
	if err := c.Do(ctx, xrpc.Query, "", "com.atproto.sync.getBlocks", params, nil, buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
