package managedspace

import (
	"encoding/json"
	"fmt"
)

// nsid is the Namespaced Identifier for the managed space policy.
const nsid string = "com.atproto.managedspace.policy"

// Policy defines a space policy that defers each user access decision
// to an external managing app via CheckUserAccess requests.
//
// nsid: com.atproto.managedspace.policy
//
// # Use Cases
//
// Managed spaces are designed for app-controlled scenarios such as:
//
//	• Game state and progress
//	• App-specific user data
//	• Subscription-gated content
//	• Dynamic access control based on external criteria
//	•Forum or community moderation
//
// # How It Works
//
//	1. User requests a space credential
//	2. Space authority calls CheckUserAccess on the managing app
//	3. Managing app returns authorized=true/false based on its own logic
// 	4. Space authority issues or denies the credential accordingly
//
// This allows complex access control logic without embedding it in the space policy.
//
// # Configuration Options
//
//	• ManagingApp: The service that makes access decisions (required)
//	• PublishPolicy: Make the policy visible as a com.atproto.space.policy record
//
// # Example (Go)
//
//	import (
//		"github.com/reiver/go-atproto/com/atproto/managedspace"
//		"github.com/reiver/go-atproto/com/atproto/space"
//	)
//	
//	// Create an app-managed space
//	policy := managedspace.Policy{
//		ManagingApp:   "did:web:game.example.com#state",
//		PublishPolicy: true,
//	}
//	
//	var resp space.CreateSpaceResponse
//	err := space.CreateSpace(
//		&resp,
//		"https://pds.example.com",
//		bearerToken,
//		"com.example.game.state",   // space type
//		"player-data",              // skey
//		policy,
//	)
//
// # HTTP/JSON Representation
//
//	{
//	  "$type": "com.atproto.managedspace.policy",
//	  "managingApp": "did:web:game.example.com#state",
//	  "publishPolicy": true
//	}
//
// # Comparison with Other Policy Types
//
//	• personalspace.Policy: Owner-only, no configuration
//	• memberspace.Policy: Group access with membership management
//	• managedspace.Policy: App-controlled access decisions (this type)
//
// # Important Notes
//
// Managed spaces have no members. Operations like putMember, addMember,
// createInvite, and acceptInvite will fail with InvalidRequest error
// when called on a managed space.
//
// The managing app must implement the CheckUserAccess endpoint to
// respond to access queries from the space authority.
//
// The $type field is automatically handled during JSON marshaling/unmarshaling.
type Policy struct {
	// ManagingApp is the service identifier of the managing app.
	// This is a DID with an optional service fragment, e.g., "did:web:example.com#forum".
	// The space authority will call CheckUserAccess on this service
	// to determine whether to grant access to each user.
	// Required.
	ManagingApp string

	// PublishPolicy controls whether the policy is published as a record.
	// When true, the policy is written as a com.atproto.space.policy record
	// in the authority's repo, allowing syncers to view governance rules.
	// Default: false
	PublishPolicy bool
}

// internalPolicy is the JSON wire format for Policy.
type internalPolicy struct {
	Type          string `json:"$type"`
	ManagingApp   string `json:"managingApp"`
	PublishPolicy bool   `json:"publishPolicy,omitempty"`
}

// MarshalJSON implements json.Marshaler.
// It automatically includes the $type field set to "com.atproto.managedspace.policy".
func (p Policy) MarshalJSON() ([]byte, error) {
	return json.Marshal(internalPolicy{
		Type:          nsid,
		ManagingApp:   p.ManagingApp,
		PublishPolicy: p.PublishPolicy,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
// It validates that the $type field is "com.atproto.managedspace.policy".
func (p *Policy) UnmarshalJSON(data []byte) error {
	var wire internalPolicy
	if err := json.Unmarshal(data, &wire); err != nil {
		return fmt.Errorf("managedspace.Policy: failed to unmarshal: %w", err)
	}

	if wire.Type != nsid {
		return fmt.Errorf("managedspace.Policy: unexpected $type: got %q, want %q", wire.Type, nsid)
	}

	p.ManagingApp = wire.ManagingApp
	p.PublishPolicy = wire.PublishPolicy

	return nil
}
