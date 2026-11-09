package simplespace

import (
	"net/url"
	"strconv"

	"github.com/reiver/go-xrpc"
)

// ListMembers lists the current members of an ATProto space and their read and write access.
//
// I.e., NSID = com.atproto.simplespace.listMembers
//
// The 'host' parameter should be the user's PDS server.
// This is a query that accepts an OAuth read_self grant or a space credential.
// The cursor parameter is used for pagination.
func ListMembers(dst any, host string, bearerToken string, space string, cursor string, limit int) error {

	const nsid string = "com.atproto.simplespace.listMembers"

	var query string = "space=" + url.QueryEscape(space)

	if "" != cursor {
		query = query + "&cursor=" + url.QueryEscape(cursor)
	}
	if 0 < limit {
		query = query + "&limit=" + url.QueryEscape(strconv.Itoa(limit))
	}

	var requestURL string = xrpc.ConstructURL(host, nsid, query).String()

	return xrpc.AuthorizedQuery(dst, bearerToken, requestURL)
}

type ListMembersResponse struct {
	Members []Member `json:"members"`
	Cursor  string   `json:"cursor"`
}

type Member struct {
	DID   string `json:"did"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}
