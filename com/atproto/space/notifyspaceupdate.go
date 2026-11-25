package space

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDNotifySpaceUpdate is the NSID of the procedure "com.atproto.space.notifySpaceUpdate".
const NSIDNotifySpaceUpdate = "com.atproto.space.notifySpaceUpdate"

// The "$type" values of the members of [NotifySpaceUpdateUpdate] defined in "com.atproto.space.notifySpaceUpdate".
const (
	TypeNotifySpaceUpdateSpaceDeleted        = NSIDNotifySpaceUpdate + "#spaceDeleted"
	TypeNotifySpaceUpdateRegistrationExpired = NSIDNotifySpaceUpdate + "#registrationExpired"
)

// NotifySpaceUpdateInput is the input of "com.atproto.space.notifySpaceUpdate".
type NotifySpaceUpdateInput struct {
	Space string `json:"space"`

	// Token is the token the endpoint was registered with.
	Token string `json:"token"`

	Update NotifySpaceUpdateUpdate `json:"update"`
}

// NotifySpaceUpdateSpaceDeleted is "com.atproto.space.notifySpaceUpdate#spaceDeleted" — the space has been deleted.
// Confirm with authorizeSpaceSync before dropping its data.
type NotifySpaceUpdateSpaceDeleted struct{}

// NotifySpaceUpdateRegistrationExpired is "com.atproto.space.notifySpaceUpdate#registrationExpired" — the endpoint's registration has expired.
// It is sent in place of the first update after it expired. Call authorizeSpaceSync to keep receiving notifications.
type NotifySpaceUpdateRegistrationExpired struct {
	SpaceRev string `json:"spaceRev"`
}

// NotifySpaceUpdateUpdate is [NotifySpaceUpdateInput].Update — a union of
// [RepoWrite], [RepoRemoved], [NotifySpaceUpdateSpaceDeleted], and [NotifySpaceUpdateRegistrationExpired].
//
// Exactly one field is set. The union is open: a member of another type is kept, as is, in Unknown.
type NotifySpaceUpdateUpdate struct {
	RepoWrite           *RepoWrite
	RepoRemoved         *RepoRemoved
	SpaceDeleted        *NotifySpaceUpdateSpaceDeleted
	RegistrationExpired *NotifySpaceUpdateRegistrationExpired

	// Unknown is a member of a type not listed above (its "$type" is in it). See [NotifySpaceUpdateUpdate.Type].
	Unknown json.RawMessage
}

var _ json.Marshaler = NotifySpaceUpdateUpdate{}
var _ json.Unmarshaler = new(NotifySpaceUpdateUpdate)

// Type returns the "$type" of the member. It returns "" if no field is set (or if Unknown has no "$type").
func (receiver NotifySpaceUpdateUpdate) Type() string {
	switch {
	case nil != receiver.RepoWrite:
		return TypeRepoWrite
	case nil != receiver.RepoRemoved:
		return TypeRepoRemoved
	case nil != receiver.SpaceDeleted:
		return TypeNotifySpaceUpdateSpaceDeleted
	case nil != receiver.RegistrationExpired:
		return TypeNotifySpaceUpdateRegistrationExpired
	case nil != receiver.Unknown:
		typ, _ := unionType(receiver.Unknown)
		return typ
	default:
		return ""
	}
}

func (receiver NotifySpaceUpdateUpdate) MarshalJSON() ([]byte, error) {
	if count := countSet(nil != receiver.RepoWrite, nil != receiver.RepoRemoved, nil != receiver.SpaceDeleted, nil != receiver.RegistrationExpired, nil != receiver.Unknown); 1 != count {
		return nil, errNotOneMember("NotifySpaceUpdateUpdate", count)
	}

	switch {
	case nil != receiver.RepoWrite:
		return marshalWithType(TypeRepoWrite, receiver.RepoWrite)
	case nil != receiver.RepoRemoved:
		return marshalWithType(TypeRepoRemoved, receiver.RepoRemoved)
	case nil != receiver.SpaceDeleted:
		return marshalWithType(TypeNotifySpaceUpdateSpaceDeleted, receiver.SpaceDeleted)
	case nil != receiver.RegistrationExpired:
		return marshalWithType(TypeNotifySpaceUpdateRegistrationExpired, receiver.RegistrationExpired)
	default:
		if _, err := unionType(receiver.Unknown); nil != err {
			return nil, err
		}
		return receiver.Unknown, nil
	}
}

func (receiver *NotifySpaceUpdateUpdate) UnmarshalJSON(data []byte) error {
	typ, err := unionType(data)
	if nil != err {
		return err
	}

	var member NotifySpaceUpdateUpdate
	switch typ {
	case TypeRepoWrite:
		member.RepoWrite = new(RepoWrite)
		err = json.Unmarshal(data, member.RepoWrite)
	case TypeRepoRemoved:
		member.RepoRemoved = new(RepoRemoved)
		err = json.Unmarshal(data, member.RepoRemoved)
	case TypeNotifySpaceUpdateSpaceDeleted:
		member.SpaceDeleted = new(NotifySpaceUpdateSpaceDeleted)
		err = json.Unmarshal(data, member.SpaceDeleted)
	case TypeNotifySpaceUpdateRegistrationExpired:
		member.RegistrationExpired = new(NotifySpaceUpdateRegistrationExpired)
		err = json.Unmarshal(data, member.RegistrationExpired)
	default:
		member.Unknown = copyRaw(data)
	}
	if nil != err {
		return err
	}

	*receiver = member
	return nil
}

// NotifySpaceUpdate calls "com.atproto.space.notifySpaceUpdate" — notify a syncing service of an update to a space.
//
// The space host sends it to the notify endpoint registered with authorizeSpaceSync.
//
// Implemented by the syncing service. ('endpoint' is the notify endpoint.)
func NotifySpaceUpdate(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input NotifySpaceUpdateInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDNotifySpaceUpdate, ""), header, input)
}
