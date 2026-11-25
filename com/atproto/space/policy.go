package space

import (
	"encoding/json"

	"github.com/reiver/go-atproto/com/atproto/managedspace"
	"github.com/reiver/go-atproto/com/atproto/memberspace"
	"github.com/reiver/go-atproto/com/atproto/personalspace"
)

// NSIDPolicy is the NSID of the record type "com.atproto.space.policy" — see [PolicyRecord].
const NSIDPolicy = "com.atproto.space.policy"

// The "$type" values of the members of [Policy].
const (
	TypePersonalSpacePolicy = "com.atproto.personalspace.policy"
	TypeMemberSpacePolicy   = "com.atproto.memberspace.policy"
	TypeManagedSpacePolicy  = "com.atproto.managedspace.policy"
)

// Policy is a space's policy — a union of the personalspace, memberspace, and managedspace policies.
//
// Exactly one field is set. The union is open: a policy of another type is kept, as is, in Unknown.
type Policy struct {
	Personal *personalspace.Policy
	Member   *memberspace.Policy
	Managed  *managedspace.Policy

	// Unknown is a policy of a type not listed above (its "$type" is in it). See [Policy.Type].
	Unknown json.RawMessage
}

var _ json.Marshaler = Policy{}
var _ json.Unmarshaler = new(Policy)

// Type returns the "$type" of the policy. It returns "" if no field is set (or if Unknown has no "$type").
func (receiver Policy) Type() string {
	switch {
	case nil != receiver.Personal:
		return TypePersonalSpacePolicy
	case nil != receiver.Member:
		return TypeMemberSpacePolicy
	case nil != receiver.Managed:
		return TypeManagedSpacePolicy
	case nil != receiver.Unknown:
		typ, _ := unionType(receiver.Unknown)
		return typ
	default:
		return ""
	}
}

func (receiver Policy) MarshalJSON() ([]byte, error) {
	if count := countSet(nil != receiver.Personal, nil != receiver.Member, nil != receiver.Managed, nil != receiver.Unknown); 1 != count {
		return nil, errNotOneMember("Policy", count)
	}

	switch {
	case nil != receiver.Personal:
		return json.Marshal(receiver.Personal)
	case nil != receiver.Member:
		return json.Marshal(receiver.Member)
	case nil != receiver.Managed:
		return json.Marshal(receiver.Managed)
	default:
		if _, err := unionType(receiver.Unknown); nil != err {
			return nil, err
		}
		return receiver.Unknown, nil
	}
}

func (receiver *Policy) UnmarshalJSON(data []byte) error {
	typ, err := unionType(data)
	if nil != err {
		return err
	}

	var policy Policy
	switch typ {
	case TypePersonalSpacePolicy:
		policy.Personal = new(personalspace.Policy)
		err = json.Unmarshal(data, policy.Personal)
	case TypeMemberSpacePolicy:
		policy.Member = new(memberspace.Policy)
		err = json.Unmarshal(data, policy.Member)
	case TypeManagedSpacePolicy:
		policy.Managed = new(managedspace.Policy)
		err = json.Unmarshal(data, policy.Managed)
	default:
		policy.Unknown = copyRaw(data)
	}
	if nil != err {
		return err
	}

	*receiver = policy
	return nil
}

// PolicyRecord is the record "com.atproto.space.policy" — a published copy of a space's policy.
// Its record key is "self". It is only authoritative in the space authority's repo.
type PolicyRecord struct {
	Policy Policy `json:"policy"`
}

func (receiver PolicyRecord) MarshalJSON() ([]byte, error) {
	type plain PolicyRecord
	return marshalWithType(NSIDPolicy, plain(receiver))
}
