package simplespace

import (
	"github.com/reiver/go-xrpc"
)

// DeleteSpace deletes an existing ATProto space.
//
// I.e., NSID = com.atproto.simplespace.deleteSpace
//
// The 'host' parameter should be the user's PDS server.
func DeleteSpace(host string, bearerToken string, space string) error {

	const nsid string = "com.atproto.simplespace.deleteSpace"

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]any{
		"space":space,
	}

	return xrpc.AuthorizedExecute(nil, bearerToken, url, src)
}
