package managedspace

// Policy defines a space policy that defers each user access decision
// to an external managing app via CheckUserAccess requests.
//
// NSID: com.atproto.managedspace.policy
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
type Policy struct {
	// Type is the lexicon type identifier.
	// This should always be "com.atproto.managedspace.policy" when serialized.
	Type string `json:"$type,omitempty"`

	// ManagingApp is the service identifier of the managing app.
	// This is a DID with an optional service fragment, e.g., "did:web:example.com#forum".
	// The space authority will call CheckUserAccess on this service
	// to determine whether to grant access to each user.
	// Required.
	ManagingApp string `json:"managingApp"`

	// PublishPolicy controls whether the policy is published as a record.
	// When true, the policy is written as a com.atproto.space.policy record
	// in the authority's repo, allowing syncers to view governance rules.
	// Default: false
	PublishPolicy bool `json:"publishPolicy,omitempty"`
}

// NewPolicy creates a new managed space policy with the specified managing app.
//
// # Parameters
//
// The 'managingApp' parameter is the service identifier that will make access decisions.
// This should be a DID, optionally with a service fragment (e.g., "did:web:example.com#forum").
//
// # Example
//
//	policy := managedspace.NewPolicy("did:web:game.example.com#state")
//	
//	var resp space.CreateSpaceResponse
//	err := space.CreateSpace(&resp, host, token, "com.example.game", "data", policy)
func NewPolicy(managingApp string) Policy {
	return Policy{
		Type:        "com.atproto.managedspace.policy",
		ManagingApp: managingApp,
	}
}
