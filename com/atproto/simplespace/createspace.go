package simplespace

import (
	"github.com/reiver/go-xrpc"
)

// CreateSpace creates a new ATProto space.
// The caller becomes the authority of the space.
//
// I.e., NSID = com.atproto.simplespace.createSpace
//
// The 'host' parameter should be the user's PDS server.
// The 'readPolicy' parameter can be one of: [ReadPolicyPublic], [ReadPolicyMemberList], or [ReadPolicyManagingApp].
// The 'writePolicy' parameter can be one of: [WritePolicyPublic], [WritePolicyMemberList], or [WritePolicyManagingApp].
// The 'appAccess' parameter can be either "open" (default) or "allowList".
// The 'managingApp' parameter is a service identifier used when policies delegate decisions.
func CreateSpace(dst any, host string, bearerToken string, readPolicy string, writePolicy string, appAccess string, managingApp string) error {

	const nsid string = "com.atproto.simplespace.createSpace"

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]any{
		"readPolicy":readPolicy,
		"writePolicy":writePolicy,
	}

	if "" != appAccess {
		src["appAccess"] = appAccess
	}
	if "" != managingApp {
		src["managingApp"] = managingApp
	}

	return xrpc.AuthorizedExecute(dst, bearerToken, url, src)
}

type CreateSpaceResponse struct {
	URI         string `json:"uri"`
	CID         string `json:"cid"`
	ReadPolicy  string `json:"readPolicy"`
	WritePolicy string `json:"writePolicy"`
	AppAccess   string `json:"appAccess"`
	ManagingApp string `json:"managingApp"`
}
