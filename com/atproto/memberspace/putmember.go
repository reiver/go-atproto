package memberspace

import (
	"github.com/reiver/go-xrpc"
)

// PutMember adds a member to the member list of a space governed by the memberspace policy,
// or replaces their read and write access if they already exist.
// This is a host procedure that requires OAuth authentication as the space owner.
//
// NSID: com.atproto.memberspace.putMember
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
//	• DID: The DID of the member to add or update
//	• Read: Whether the member may read when publicRead is false
//	• Write: Whether the member's writes are tracked when publicWrite is false
//
// # Example (Go)
//
//	// Add a new member with read and write access
//	err := memberspace.PutMember(
//		"https://space.example.com",
//		"eyJhbGciOiJFUzI1NiIs...",  // OAuth bearer token (owner)
//		memberspace.PutMemberRequest{
//			Space: "at://did:plc:abc123/space/app.bsky.chat.group/my-group",
//			DID:   "did:plc:newmember456",
//			Read:  true,
//			Write: true,
//		},
//	)
//	if nil != err {
//		log.Fatal(err)
//	}
//	fmt.Println("Member added successfully")
//
//	// Update existing member to read-only
//	err = memberspace.PutMember(
//		"https://space.example.com",
//		bearerToken,
//		memberspace.PutMemberRequest{
//			Space: "at://did:plc:abc123/space/app.bsky.chat.group/my-group",
//			DID:   "did:plc:existingmember789",
//			Read:  true,
//			Write: false,
//		},
//	)
//
// # HTTP Request
//
//	POST /xrpc/com.atproto.memberspace.putMember HTTP/1.1
//	Host: space.example.com
//	Authorization: Bearer eyJhbGciOiJFUzI1NiIs...
//	Content-Type: application/json
//
//	{
//	  "space": "at://did:plc:abc123/space/app.bsky.chat.group/my-group",
//	  "did": "did:plc:newmember456",
//	  "read": true,
//	  "write": true
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
// If PublishMembers is enabled in the space policy, adding or updating a member
// will also create/update a com.atproto.memberspace.member record in the authority's repo.
func PutMember(host string, bearerToken string, req PutMemberRequest) error {

	const nsid string = "com.atproto.memberspace.putMember"

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]any{
		"space": req.Space,
		"did":   req.DID,
		"read":  req.Read,
		"write": req.Write,
	}

	return xrpc.AuthorizedExecute(nil, bearerToken, url, src)
}

// PutMemberRequest contains the parameters for PutMember.
type PutMemberRequest struct {
	// Space is the space reference URI (e.g., "at://did:plc:abc123/space/app.bsky.chat.group/my-group").
	Space string

	// DID is the decentralized identifier of the member to add or update.
	DID string

	// Read controls whether the member may read when publicRead is false.
	Read bool

	// Write controls whether the member's writes are tracked when publicWrite is false.
	Write bool
}
