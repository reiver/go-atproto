package memberspace

import (
	"github.com/reiver/go-xrpc"
)

// RemoveMember removes a member from the member list of a space governed by the memberspace policy.
// After removal, the member can no longer access the space (unless PublicRead/PublicWrite is enabled).
// This is a host procedure that requires OAuth authentication as the space owner.
//
// NSID: com.atproto.memberspace.removeMember
//
// # Parameters
//
// The 'host' parameter should be the space host server (e.g., "https://space.example.com").
//
// The 'bearerToken' parameter is the OAuth access token (must be space owner).
//
// The 'req' parameter contains:
//
//	• Space: The space reference URI (e.g., "at://did:plc:abc123/space/app.bsky.chat.group/my-group")
//	• DID: The DID of the member to remove
//
// # Example (Go)
//
//	err := memberspace.RemoveMember(
//		"https://space.example.com",
//		"eyJhbGciOiJFUzI1NiIs...",  // OAuth bearer token (owner)
//		memberspace.RemoveMemberRequest{
//			Space: "at://did:plc:abc123/space/app.bsky.chat.group/my-group",
//			DID:   "did:plc:formerMember789",
//		},
//	)
//	if nil != err {
//		log.Fatal(err)
//	}
//	fmt.Println("Member removed successfully")
//
// # HTTP Request
//
//	POST /xrpc/com.atproto.memberspace.removeMember HTTP/1.1
//	Host: space.example.com
//	Authorization: Bearer eyJhbGciOiJFUzI1NiIs...
//	Content-Type: application/json
//
//	{
//	  "space": "at://did:plc:abc123/space/app.bsky.chat.group/my-group",
//	  "did": "did:plc:formerMember789"
//	}
//
// # HTTP Response
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//
//	{}
//
// # Errors
//
//	• SpaceNotFound: The referenced space doesn't exist.
//	• NotSpaceOwner: The caller is not the space owner.
//	• PolicyMismatch: The space is not governed by the memberspace policy.
//
// # Notes
//
// Removing a member does not delete any records they previously created in the space.
// To remove their content, use space.ClearRecords or space.DeleteRecord.
//
// If the removed member has active space credentials, consider calling
// space.NotifyCredentialRevoked to invalidate them immediately rather than
// waiting for expiration.
//
// If PublishMembers is enabled in the space policy, removing a member
// will also delete their com.atproto.memberspace.member record from the authority's repo.
func RemoveMember(host string, bearerToken string, req RemoveMemberRequest) error {

	const nsid string = "com.atproto.memberspace.removeMember"

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]any{
		"space": req.Space,
		"did":   req.DID,
	}

	return xrpc.AuthorizedExecute(nil, bearerToken, url, src)
}

// RemoveMemberRequest contains the parameters for RemoveMember.
type RemoveMemberRequest struct {
	// Space is the space reference URI (e.g., "at://did:plc:abc123/space/app.bsky.chat.group/my-group").
	Space string

	// DID is the decentralized identifier of the member to remove.
	DID string
}
