package space

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDPutRecord is the NSID of the procedure "com.atproto.space.putRecord".
const NSIDPutRecord = "com.atproto.space.putRecord"

// PutRecordInput is the input of "com.atproto.space.putRecord".
type PutRecordInput struct {
	Space      string `json:"space"`
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`

	// Validate: false skips Lexicon schema validation of the record; true requires it; nil validates only records of known Lexicons.
	Validate *bool `json:"validate,omitempty"`

	Record json.RawMessage `json:"record"`
}

// PutRecordOutput is the output of "com.atproto.space.putRecord".
type PutRecordOutput struct {
	URI string `json:"uri"`
	CID string `json:"cid"`

	// ValidationStatus is "valid", "unknown", or "" (or another value).
	ValidationStatus string `json:"validationStatus,omitempty"`
}

// PutRecord calls "com.atproto.space.putRecord" — write a record in a space repo, creating or updating it as needed.
//
// Implemented by the PDS. Requires auth.
//
// Errors: SpaceNotFound.
func PutRecord(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input PutRecordInput) (PutRecordOutput, error) {
	return execute[PutRecordOutput](ctx, client, methodURL(endpoint, NSIDPutRecord, ""), header, input)
}
