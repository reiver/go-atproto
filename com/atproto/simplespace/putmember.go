package simplespace

import (
	"github.com/reiver/go-xrpc"
)

// PutMember adds a member (by DID) to an ATProto space or updates their read and write access.
//
// I.e., NSID = com.atproto.simplespace.putMember
//
// The 'host' parameter should be the user's PDS server.
// The 'read' and 'write' parameters control the member's access permissions.
// This operation replaces both values together.
func PutMember(dst any, host string, bearerToken string, space string, did string, read bool, write bool) error {

	const nsid string = "com.atproto.simplespace.putMember"

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]any{
		"space":space,
		"did":did,
		"read":read,
		"write":write,
	}

	return xrpc.AuthorizedExecute(dst, bearerToken, url, src)
}

type PutMemberResponse struct {
	DID   string `json:"did"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}
