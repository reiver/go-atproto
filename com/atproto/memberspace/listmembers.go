package memberspace

import (
	"net/url"
	"strconv"

	"github.com/reiver/go-xrpc"
)

// ListMembers lists the members of a space governed by the memberspace policy,
// including their read and write access permissions.
// This is a host query that requires OAuth authentication with a read_self grant.
// Note: A space credential is not sufficient for this operation.
//
// NSID: com.atproto.memberspace.listMembers
//
// # Parameters
//
// The 'host' parameter should be the space host server (e.g., "https://space.example.com").
//
// The 'bearerToken' parameter is the OAuth access token with read_self scope.
//
// The 'req' parameter contains:
//
//	• Space: The space reference URI (e.g., "at://did:plc:abc123/space/app.bsky.chat.group/my-group")
//	• Cursor: Pagination cursor (empty string for first request)
//	• Limit: Maximum number of members to return (1-1000, default 100; 0 for server default)
//
// # Example (Go)
//
//	var resp memberspace.ListMembersResponse
//	err := memberspace.ListMembers(
//		&resp,
//		"https://space.example.com",
//		"eyJhbGciOiJFUzI1NiIs...",  // OAuth bearer token
//		memberspace.ListMembersRequest{
//			Space: "at://did:plc:abc123/space/app.bsky.chat.group/my-group",
//			Limit: 50,
//		},
//	)
//	if nil != err {
//		log.Fatal(err)
//	}
//	for _, member := range resp.Members {
//		fmt.Printf("Member: %s (read=%v, write=%v)\n",
//			member.DID, member.Read, member.Write)
//	}
//	if "" != resp.Cursor {
//		fmt.Println("More members available, use cursor:", resp.Cursor)
//	}
//
// # HTTP Request
//
//	GET /xrpc/com.atproto.memberspace.listMembers?space=at%3A%2F%2Fdid%3Aplc%3Aabc123%2Fspace%2Fapp.bsky.chat.group%2Fmy-group&limit=50 HTTP/1.1
//	Host: space.example.com
//	Authorization: Bearer eyJhbGciOiJFUzI1NiIs...
//
// # HTTP Response
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//
//	{
//	  "members": [
//	    {
//	      "did": "did:plc:member123",
//	      "read": true,
//	      "write": true
//	    },
//	    {
//	      "did": "did:plc:member456",
//	      "read": true,
//	      "write": false
//	    }
//	  ],
//	  "cursor": "did:plc:member456"
//	}
//
// # Errors
//
//	• SpaceNotFound: The referenced space doesn't exist.
//	• PolicyMismatch: The space is not governed by the memberspace policy.
//
// # Notes
//
// This endpoint requires OAuth authentication, not just a space credential.
// This ensures only authorized parties (typically the space owner or members
// with appropriate privileges) can enumerate the full member list.
//
// If PublishMembers is enabled in the space policy, the member list is also
// available as com.atproto.memberspace.member records in the authority's repo,
// which can be read by anyone with a space credential.
func ListMembers(dst any, host string, bearerToken string, req ListMembersRequest) error {

	const nsid string = "com.atproto.memberspace.listMembers"

	var query string = "space=" + url.QueryEscape(req.Space)

	if "" != req.Cursor {
		query = query + "&cursor=" + url.QueryEscape(req.Cursor)
	}
	if 0 < req.Limit {
		query = query + "&limit=" + url.QueryEscape(strconv.Itoa(req.Limit))
	}

	var requestURL string = xrpc.ConstructURL(host, nsid, query).String()

	return xrpc.AuthorizedQuery(dst, bearerToken, requestURL)
}

// ListMembersRequest contains the parameters for ListMembers.
type ListMembersRequest struct {
	// Space is the space reference URI (e.g., "at://did:plc:abc123/space/app.bsky.chat.group/my-group").
	Space string

	// Cursor is the pagination cursor from a previous response.
	// Pass empty string for the first request.
	Cursor string

	// Limit is the maximum number of members to return (1-1000, default 100).
	// Pass 0 to use the server default.
	Limit int
}

// ListMembersResponse is the response from ListMembers.
type ListMembersResponse struct {
	// Members is the list of members with their access permissions.
	Members []MemberInfo `json:"members"`

	// Cursor is the pagination cursor for the next page.
	// Empty if there are no more results.
	Cursor string `json:"cursor,omitempty"`
}

// MemberInfo represents a member's information as returned by ListMembers.
type MemberInfo struct {
	// DID is the member's decentralized identifier.
	DID string `json:"did"`

	// Read indicates whether the member may read when publicRead is false.
	Read bool `json:"read"`

	// Write indicates whether the member's writes are tracked when publicWrite is false.
	Write bool `json:"write"`
}
