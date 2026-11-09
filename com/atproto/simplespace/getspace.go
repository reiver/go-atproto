package simplespace

import (
	"net/url"

	"github.com/reiver/go-xrpc"
)

// GetSpace describes a ATProto space and its configuration.
//
// I.e., NSID = com.atproto.simplespace.getSpace
//
// The 'host' parameter should be the user's PDS server.
// This is a query that accepts an OAuth read_self grant or a space credential.
func GetSpace(dst any, host string, bearerToken string, space string) error {

	const nsid string = "com.atproto.simplespace.getSpace"

	var query string = "space=" + url.QueryEscape(space)

	var requestURL string = xrpc.ConstructURL(host, nsid, query).String()

	return xrpc.AuthorizedQuery(dst, bearerToken, requestURL)
}

type GetSpaceResponse struct {
	URI         string `json:"uri"`
	CID         string `json:"cid"`
	ReadPolicy  string `json:"readPolicy"`
	WritePolicy string `json:"writePolicy"`
	AppAccess   string `json:"appAccess"`
	ManagingApp string `json:"managingApp"`
}
