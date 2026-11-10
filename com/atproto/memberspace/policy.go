package memberspace

import (
	"encoding/json"
	"fmt"
)

// policyNSID is the Namespaced Identifier for the member space policy.
const policyNSID string = "com.atproto.memberspace.policy"

// memberNSID is the Namespaced Identifier for the member record.
const memberNSID string = "com.atproto.memberspace.member"

// Policy defines a space policy with membership-based access control.
// Access is governed by a member list maintained by the space host, managed via PutMember and RemoveMember operations.
//
// NSID: com.atproto.memberspace.policy
//
// # Use Cases
//
// Member spaces are designed for collaborative scenarios such as:
//
//	• Group chats and direct messages
//	• Collaborative projects
//	• Communities and forums
//	• Shared workspaces
//	• Team documents
//
// # Configuration Options
//
//	• PublicRead: When true, anyone can read (credentials issued to all)
//	• PublicWrite: When true, anyone's writes are tracked and forwarded
//	• AllowedApps: Restrict credentials to specific OAuth client IDs
//	• PublishPolicy: Make the policy visible as a com.atproto.space.policy record
//	• PublishMembers: Make the member list visible as com.atproto.memberspace.member records
//
// # Example (Go)
//
//	import (
//		"github.com/reiver/go-atproto/com/atproto/memberspace"
//		"github.com/reiver/go-atproto/com/atproto/space"
//	)
//	
//	// Create a group chat space - members only
//	policy := memberspace.Policy{
//		PublicRead:     false,
//		PublicWrite:    false,
//		PublishPolicy:  true,
//		PublishMembers: false,  // keep member list private
//	}
//	
//	var resp space.CreateSpaceResponse
//	err := space.CreateSpace(
//		&resp,
//		"https://pds.example.com",
//		bearerToken,
//		"app.bsky.chat.group",    // space type
//		"my-group-chat",          // skey
//		policy,
//	)
//
// # HTTP/JSON Representation
//
//	{
//	  "$type": "com.atproto.memberspace.policy",
//	  "publicRead": false,
//	  "publicWrite": false,
//	  "allowedApps": ["did:web:chat.example.com"],
//	  "publishPolicy": true,
//	  "publishMembers": false
//	}
//
// # Access Control Flow
//
// When PublicRead is false:
//
//	1. User requests space credential via GetSpaceCredential
//	2. Space authority checks if user is in member list with read=true
//	3. If member, credential is issued; otherwise UserNotAuthorized error
//
// When PublicWrite is false:
//
//	1. User writes to their repo in the space
//	2. Repo host notifies space host via NotifyWrite
//	3. Space host checks if user is in member list with write=true
//	4. If member, write is tracked and forwarded to syncers
//
// # Comparison with Other Policy Types
//
//	• personalspace.Policy: Owner-only, no configuration
//	• memberspace.Policy: Group access with membership management (this type)
//	• managedspace.Policy: App-controlled access decisions
//
// The $type field is automatically handled during JSON marshaling/unmarshaling.
type Policy struct {
	// PublicRead controls whether non-members can read.
	// When true, the authority issues a space credential to anyone.
	// When false, only members with read=true can obtain credentials.
	// Default: false
	PublicRead bool

	// PublicWrite controls whether non-members' writes are tracked.
	// When true, the authority tracks and forwards notifications for any user's writes.
	// When false, only members with write=true have their writes tracked.
	// Default: false
	PublicWrite bool

	// AllowedApps restricts which OAuth clients can obtain credentials.
	// Each entry is an OAuth client ID (typically a DID).
	// This restriction applies to read access only, not writes.
	// If empty or nil, any app can obtain credentials (subject to other restrictions).
	AllowedApps []string

	// PublishPolicy controls whether the policy is published as a record.
	// When true, the policy is written as a com.atproto.space.policy record
	// in the authority's repo, allowing syncers to view governance rules.
	// Default: false
	PublishPolicy bool

	// PublishMembers controls whether the member list is published.
	// When true, members are written as com.atproto.memberspace.member records
	// in the authority's repo, making membership publicly visible to readers.
	// Default: false
	PublishMembers bool
}

// internalPolicy is the JSON wire format for Policy.
type internalPolicy struct {
	Type           string   `json:"$type"`
	PublicRead     bool     `json:"publicRead,omitempty"`
	PublicWrite    bool     `json:"publicWrite,omitempty"`
	AllowedApps    []string `json:"allowedApps,omitempty"`
	PublishPolicy  bool     `json:"publishPolicy,omitempty"`
	PublishMembers bool     `json:"publishMembers,omitempty"`
}

// MarshalJSON implements json.Marshaler.
// It automatically includes the $type field set to "com.atproto.memberspace.policy".
func (p Policy) MarshalJSON() ([]byte, error) {
	return json.Marshal(internalPolicy{
		Type:           policyNSID,
		PublicRead:     p.PublicRead,
		PublicWrite:    p.PublicWrite,
		AllowedApps:    p.AllowedApps,
		PublishPolicy:  p.PublishPolicy,
		PublishMembers: p.PublishMembers,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
// It validates that the $type field is "com.atproto.memberspace.policy".
func (p *Policy) UnmarshalJSON(data []byte) error {
	var wire internalPolicy
	if err := json.Unmarshal(data, &wire); err != nil {
		return fmt.Errorf("memberspace.Policy: failed to unmarshal: %w", err)
	}

	if wire.Type != policyNSID {
		return fmt.Errorf("memberspace.Policy: unexpected $type: got %q, want %q", wire.Type, policyNSID)
	}

	p.PublicRead = wire.PublicRead
	p.PublicWrite = wire.PublicWrite
	p.AllowedApps = wire.AllowedApps
	p.PublishPolicy = wire.PublishPolicy
	p.PublishMembers = wire.PublishMembers

	return nil
}

// Member represents a membership record in a memberspace.
// The space authority publishes these when PublishMembers is enabled.
//
// NSID: com.atproto.memberspace.member
//
// # Record Key
//
// Members are keyed by their DID in the authority's repo within the space.
// For example: at://did:plc:authority/com.atproto.memberspace.member/did:plc:member123
//
// # Example (Go)
//
//	member := memberspace.Member{
//		DID:   "did:plc:user123",
//		Read:  true,
//		Write: true,
//	}
//
// # HTTP/JSON Representation
//
//	{
//	  "$type": "com.atproto.memberspace.member",
//	  "did": "did:plc:user123",
//	  "read": true,
//	  "write": true
//	}
//
// The $type field is automatically handled during JSON marshaling/unmarshaling.
type Member struct {
	// DID is the member's decentralized identifier.
	DID string

	// Read indicates whether the member may read when publicRead is false.
	Read bool

	// Write indicates whether the member's writes are tracked when publicWrite is false.
	Write bool
}

// internalMember is the JSON wire format for Member.
type internalMember struct {
	Type  string `json:"$type"`
	DID   string `json:"did"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}

// MarshalJSON implements json.Marshaler.
// It automatically includes the $type field set to "com.atproto.memberspace.member".
func (m Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(internalMember{
		Type:  memberNSID,
		DID:   m.DID,
		Read:  m.Read,
		Write: m.Write,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
// It validates that the $type field is "com.atproto.memberspace.member".
func (m *Member) UnmarshalJSON(data []byte) error {
	var wire internalMember
	if err := json.Unmarshal(data, &wire); err != nil {
		return fmt.Errorf("memberspace.Member: failed to unmarshal: %w", err)
	}

	if wire.Type != memberNSID {
		return fmt.Errorf("memberspace.Member: unexpected $type: got %q, want %q", wire.Type, memberNSID)
	}

	m.DID = wire.DID
	m.Read = wire.Read
	m.Write = wire.Write

	return nil
}
