package personalspace

// Policy defines a space policy for confidential user data.
// This is the simplest policy type - only the space authority's own account is authorized to read, and the authority tracks no other writers.
//
// NSID: com.atproto.personalspace.policy
//
// # Use Cases
//
// Personal spaces are designed for private user data such as:
//
//	• Bookmarks
//	• Mutes and blocks
//	• Drafts
//	• Private notes
//	• User preferences
//
// # Characteristics
//
//	• Owner-only access: Only the space owner can read and write
//	• No membership: No other users can be granted access
//	• No configuration: The policy has no configurable fields
//	• Maximum privacy: Data never leaves the owner's control
//
// # Example (Go)
//
//	import (
//		"github.com/reiver/go-atproto/com/atproto/personalspace"
//		"github.com/reiver/ggo-atproto/com/atproto/space"
//	)
//	
//	// Create a personal space for private bookmarks
//	policy := personalspace.Policy{}
//	
//	var resp space.CreateSpaceResponse
//	err := space.CreateSpace(
//		&resp,
//		"https://pds.example.com",
//		bearerToken,
//		"app.bsky.bookmarks",     // space type
//		"self",                   // skey
//		policy,
//	)
//
// # HTTP/JSON Representation
//
//	{
//	  "$type": "com.atproto.personalspace.policy"
//	}
//
// # Comparison with Other Policy Types
//
//	• personalspace.Policy: Owner-only, no configuration (this type)
//	• memberspace.Policy: Group access with membership management
//	• managedspace.Policy: App-controlled access decisions
//
// # Notes
//
// Since this policy grants access only to the owner, attempting to use
// membership-related operations (putMember, removeMember, listMembers)
// on a personal space will fail with PolicyMismatch error.
type Policy struct {
	// Type is the lexicon type identifier.
	// This should always be "com.atproto.personalspace.policy" when serialized.
	Type string `json:"$type,omitempty"`
}

// NewPolicy creates a new personal space policy.
//
// # Example
//
//	policy := personalspace.NewPolicy()
//	
//	var resp space.CreateSpaceResponse
//	err := space.CreateSpace(&resp, host, token, "app.bsky.bookmarks", "self", policy)
func NewPolicy() Policy {
	return Policy{
		Type: "com.atproto.personalspace.policy",
	}
}
