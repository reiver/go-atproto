package simplespace

import (
	"github.com/reiver/go-xrpc"
)

// RemoveMember removes a member (by DID) from a space.
//
// I.e., NSID = com.atproto.simplespace.removeMember
//
// The 'host' parameter should be the user's PDS server.
func RemoveMember(host string, bearerToken string, space string, did string) error {

	const nsid string = "com.atproto.simplespace.removeMember"

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]any{
		"space":space,
		"did":did,
	}

	return xrpc.AuthorizedExecute(nil, bearerToken, url, src)
}
