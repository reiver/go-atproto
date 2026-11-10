package personalspace

import (
	"encoding/json"
	"fmt"
)

// nsid is the Namespaced Identifier for the personal space policy.
const nsid string = "com.atproto.personalspace.policy"

// Policy defines a space policy for confidential user data.
// This is the simplest policy type - only the space authority's own account is authorized to read, and the authority tracks no other writers.
//
// nsid: com.atproto.personalspace.policy
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
//		"github.com/reiver/go-atproto/com/atproto/space"
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
//
// The $type field is automatically handled during JSON marshaling/unmarshaling.
type Policy struct{}

// MarshalJSON implements json.Marshaler.
// It automatically includes the $type field set to "com.atproto.personalspace.policy".
func (Policy) MarshalJSON() ([]byte, error) {
	type wireFormat struct {
		Type string `json:"$type"`
	}
	return json.Marshal(wireFormat{
		Type: nsid,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
// It validates that the $type field is "com.atproto.personalspace.policy".
func (p *Policy) UnmarshalJSON(data []byte) error {
	type wireFormat struct {
		Type string `json:"$type"`
	}

	var wire wireFormat
	if err := json.Unmarshal(data, &wire); err != nil {
		return fmt.Errorf("personalspace.Policy: failed to unmarshal: %w", err)
	}

	if wire.Type != nsid {
		return fmt.Errorf("personalspace.Policy: unexpected $type: got %q, want %q", wire.Type, nsid)
	}

	return nil
}
