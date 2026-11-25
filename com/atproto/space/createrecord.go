package space

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDCreateRecord is the NSID of the procedure "com.atproto.space.createRecord".
const NSIDCreateRecord = "com.atproto.space.createRecord"

// CreateRecordInput is the input of "com.atproto.space.createRecord".
type CreateRecordInput struct {
	Space      string `json:"space"`
	Repo       string `json:"repo"`
	Collection string `json:"collection"`

	// Rkey is the record key. If "", the PDS makes one.
	Rkey string `json:"rkey,omitempty"`

	// Validate: false skips Lexicon schema validation of the record; true requires it; nil validates only records of known Lexicons.
	Validate *bool `json:"validate,omitempty"`

	// Record is the record itself. It must have a "$type" field.
	Record json.RawMessage `json:"record"`
}

// CreateRecordOutput is the output of "com.atproto.space.createRecord".
type CreateRecordOutput struct {
	URI string `json:"uri"`
	CID string `json:"cid"`

	// ValidationStatus is "valid", "unknown", or "" (or another value).
	ValidationStatus string `json:"validationStatus,omitempty"`
}

// CreateRecord calls "com.atproto.space.createRecord" — create a single new record in a space repo.
//
// Implemented by the PDS. Requires auth.
//
// Errors: SpaceNotFound, RecordAlreadyExists.
func CreateRecord(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input CreateRecordInput) (CreateRecordOutput, error) {
	return execute[CreateRecordOutput](ctx, client, methodURL(endpoint, NSIDCreateRecord, ""), header, input)
}
