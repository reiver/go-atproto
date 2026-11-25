package space

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDListWriters is the NSID of the query "com.atproto.space.listWriters".
const NSIDListWriters = "com.atproto.space.listWriters"

// ListWritersParams are the params of "com.atproto.space.listWriters".
type ListWritersParams struct {
	Space string `query:"space"`

	// Limit is from 1 to 1000. If 0, it is left out (the default is 100).
	Limit int64 `query:"limit"`

	// Cursor, if not "", lists only the repos updated after this space revision.
	//
	// Removed repos may be garbage collected after 72 hours,
	// so a syncer whose last seen space revision is older should list the whole writer set instead.
	Cursor string `query:"cursor"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver ListWritersParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDListWriters, encodeQuery(receiver))
}

// ListWritersOutput is the output of "com.atproto.space.listWriters".
type ListWritersOutput struct {
	// Repos is the latest update to each repo, in ascending order of spaceRev. (Their PrevSpaceRev is not set.)
	Repos  []ListWritersRepo `json:"repos"`
	Cursor string            `json:"cursor,omitempty"`
}

// ListWritersRepo is one of [ListWritersOutput].Repos — a union of [RepoWrite] and [RepoRemoved].
//
// Exactly one field is set. The union is open: a member of another type is kept, as is, in Unknown.
type ListWritersRepo struct {
	RepoWrite   *RepoWrite
	RepoRemoved *RepoRemoved

	// Unknown is a member of a type not listed above (its "$type" is in it). See [ListWritersRepo.Type].
	Unknown json.RawMessage
}

var _ json.Marshaler = ListWritersRepo{}
var _ json.Unmarshaler = new(ListWritersRepo)

// Type returns the "$type" of the member. It returns "" if no field is set (or if Unknown has no "$type").
func (receiver ListWritersRepo) Type() string {
	switch {
	case nil != receiver.RepoWrite:
		return TypeRepoWrite
	case nil != receiver.RepoRemoved:
		return TypeRepoRemoved
	case nil != receiver.Unknown:
		typ, _ := unionType(receiver.Unknown)
		return typ
	default:
		return ""
	}
}

func (receiver ListWritersRepo) MarshalJSON() ([]byte, error) {
	if count := countSet(nil != receiver.RepoWrite, nil != receiver.RepoRemoved, nil != receiver.Unknown); 1 != count {
		return nil, errNotOneMember("ListWritersRepo", count)
	}

	switch {
	case nil != receiver.RepoWrite:
		return marshalWithType(TypeRepoWrite, receiver.RepoWrite)
	case nil != receiver.RepoRemoved:
		return marshalWithType(TypeRepoRemoved, receiver.RepoRemoved)
	default:
		if _, err := unionType(receiver.Unknown); nil != err {
			return nil, err
		}
		return receiver.Unknown, nil
	}
}

func (receiver *ListWritersRepo) UnmarshalJSON(data []byte) error {
	typ, err := unionType(data)
	if nil != err {
		return err
	}

	var member ListWritersRepo
	switch typ {
	case TypeRepoWrite:
		member.RepoWrite = new(RepoWrite)
		err = json.Unmarshal(data, member.RepoWrite)
	case TypeRepoRemoved:
		member.RepoRemoved = new(RepoRemoved)
		err = json.Unmarshal(data, member.RepoRemoved)
	default:
		member.Unknown = copyRaw(data)
	}
	if nil != err {
		return err
	}

	*receiver = member
	return nil
}

// ListWriters calls "com.atproto.space.listWriters" — list the repos in a space's writer set,
// with the latest known state of each: its current revision, or its removal from the writer set.
//
// Implemented by the space host. Requires a space credential.
//
// Errors: SpaceNotFound, RepoTakendown (the space authority), RepoDeactivated (the space authority).
func ListWriters(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params ListWritersParams) (ListWritersOutput, error) {
	return query[ListWritersOutput](ctx, client, params.URL(endpoint), header)
}
