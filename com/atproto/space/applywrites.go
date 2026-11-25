package space

import (
	"context"
	"encoding/json"
	"net/http"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"github.com/reiver/go-xrpc"
)

// NSIDApplyWrites is the NSID of the procedure "com.atproto.space.applyWrites".
const NSIDApplyWrites = "com.atproto.space.applyWrites"

// The "$type" values of the members of [ApplyWritesWrite] and [ApplyWritesResult].
const (
	TypeApplyWritesCreate       = NSIDApplyWrites + "#create"
	TypeApplyWritesUpdate       = NSIDApplyWrites + "#update"
	TypeApplyWritesDelete       = NSIDApplyWrites + "#delete"
	TypeApplyWritesCreateResult = NSIDApplyWrites + "#createResult"
	TypeApplyWritesUpdateResult = NSIDApplyWrites + "#updateResult"
	TypeApplyWritesDeleteResult = NSIDApplyWrites + "#deleteResult"
)

// ApplyWritesInput is the input of "com.atproto.space.applyWrites".
type ApplyWritesInput struct {
	Space string `json:"space"`
	Repo  string `json:"repo"`

	// Validate: false skips Lexicon schema validation of the records; true requires it; nil validates only records of known Lexicons.
	Validate *bool `json:"validate,omitempty"`

	Writes []ApplyWritesWrite `json:"writes"`
}

// ApplyWritesOutput is the output of "com.atproto.space.applyWrites".
type ApplyWritesOutput struct {
	Results []ApplyWritesResult `json:"results,omitempty"`
}

// ApplyWritesCreate is "com.atproto.space.applyWrites#create" — an operation that creates a new record.
type ApplyWritesCreate struct {
	Collection string          `json:"collection"`
	Rkey       string          `json:"rkey,omitempty"`
	Value      json.RawMessage `json:"value"`
}

// ApplyWritesUpdate is "com.atproto.space.applyWrites#update" — an operation that updates an existing record.
type ApplyWritesUpdate struct {
	Collection string          `json:"collection"`
	Rkey       string          `json:"rkey"`
	Value      json.RawMessage `json:"value"`
}

// ApplyWritesDelete is "com.atproto.space.applyWrites#delete" — an operation that deletes an existing record.
type ApplyWritesDelete struct {
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`
}

// ApplyWritesCreateResult is "com.atproto.space.applyWrites#createResult".
type ApplyWritesCreateResult struct {
	URI string `json:"uri"`
	CID string `json:"cid"`

	// ValidationStatus is "valid", "unknown", or "" (or another value).
	ValidationStatus string `json:"validationStatus,omitempty"`
}

// ApplyWritesUpdateResult is "com.atproto.space.applyWrites#updateResult".
type ApplyWritesUpdateResult struct {
	URI string `json:"uri"`
	CID string `json:"cid"`

	// ValidationStatus is "valid", "unknown", or "" (or another value).
	ValidationStatus string `json:"validationStatus,omitempty"`
}

// ApplyWritesDeleteResult is "com.atproto.space.applyWrites#deleteResult".
type ApplyWritesDeleteResult struct{}

// ApplyWritesWrite is one of [ApplyWritesInput].Writes — a union of [ApplyWritesCreate], [ApplyWritesUpdate], and [ApplyWritesDelete].
//
// Exactly one field is set. The union is closed: a member of any other type is an error.
type ApplyWritesWrite struct {
	Create *ApplyWritesCreate
	Update *ApplyWritesUpdate
	Delete *ApplyWritesDelete
}

var _ json.Marshaler = ApplyWritesWrite{}
var _ json.Unmarshaler = new(ApplyWritesWrite)

func (receiver ApplyWritesWrite) MarshalJSON() ([]byte, error) {
	if count := countSet(nil != receiver.Create, nil != receiver.Update, nil != receiver.Delete); 1 != count {
		return nil, errNotOneMember("ApplyWritesWrite", count)
	}

	switch {
	case nil != receiver.Create:
		return marshalWithType(TypeApplyWritesCreate, receiver.Create)
	case nil != receiver.Update:
		return marshalWithType(TypeApplyWritesUpdate, receiver.Update)
	default:
		return marshalWithType(TypeApplyWritesDelete, receiver.Delete)
	}
}

func (receiver *ApplyWritesWrite) UnmarshalJSON(data []byte) error {
	typ, err := unionType(data)
	if nil != err {
		return err
	}

	var member ApplyWritesWrite
	switch typ {
	case TypeApplyWritesCreate:
		member.Create = new(ApplyWritesCreate)
		err = json.Unmarshal(data, member.Create)
	case TypeApplyWritesUpdate:
		member.Update = new(ApplyWritesUpdate)
		err = json.Unmarshal(data, member.Update)
	case TypeApplyWritesDelete:
		member.Delete = new(ApplyWritesDelete)
		err = json.Unmarshal(data, member.Delete)
	default:
		return erorr.Wrap(ErrUnknownType, "unknown applyWrites write", field.String("type", typ))
	}
	if nil != err {
		return err
	}

	*receiver = member
	return nil
}

// ApplyWritesResult is one of [ApplyWritesOutput].Results — a union of [ApplyWritesCreateResult], [ApplyWritesUpdateResult], and [ApplyWritesDeleteResult].
//
// Exactly one field is set. The union is closed: a member of any other type is an error.
type ApplyWritesResult struct {
	CreateResult *ApplyWritesCreateResult
	UpdateResult *ApplyWritesUpdateResult
	DeleteResult *ApplyWritesDeleteResult
}

var _ json.Marshaler = ApplyWritesResult{}
var _ json.Unmarshaler = new(ApplyWritesResult)

func (receiver ApplyWritesResult) MarshalJSON() ([]byte, error) {
	if count := countSet(nil != receiver.CreateResult, nil != receiver.UpdateResult, nil != receiver.DeleteResult); 1 != count {
		return nil, errNotOneMember("ApplyWritesResult", count)
	}

	switch {
	case nil != receiver.CreateResult:
		return marshalWithType(TypeApplyWritesCreateResult, receiver.CreateResult)
	case nil != receiver.UpdateResult:
		return marshalWithType(TypeApplyWritesUpdateResult, receiver.UpdateResult)
	default:
		return marshalWithType(TypeApplyWritesDeleteResult, receiver.DeleteResult)
	}
}

func (receiver *ApplyWritesResult) UnmarshalJSON(data []byte) error {
	typ, err := unionType(data)
	if nil != err {
		return err
	}

	var member ApplyWritesResult
	switch typ {
	case TypeApplyWritesCreateResult:
		member.CreateResult = new(ApplyWritesCreateResult)
		err = json.Unmarshal(data, member.CreateResult)
	case TypeApplyWritesUpdateResult:
		member.UpdateResult = new(ApplyWritesUpdateResult)
		err = json.Unmarshal(data, member.UpdateResult)
	case TypeApplyWritesDeleteResult:
		member.DeleteResult = new(ApplyWritesDeleteResult)
		err = json.Unmarshal(data, member.DeleteResult)
	default:
		return erorr.Wrap(ErrUnknownType, "unknown applyWrites result", field.String("type", typ))
	}
	if nil != err {
		return err
	}

	*receiver = member
	return nil
}

// ApplyWrites calls "com.atproto.space.applyWrites" — apply a batch of creates, updates, and deletes to a space repo, as one transaction.
//
// Implemented by the PDS. Requires auth.
//
// Errors: SpaceNotFound, RecordNotFound, RecordAlreadyExists.
func ApplyWrites(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input ApplyWritesInput) (ApplyWritesOutput, error) {
	return execute[ApplyWritesOutput](ctx, client, methodURL(endpoint, NSIDApplyWrites, ""), header, input)
}
