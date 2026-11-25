package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDDeleteRecord is the NSID of the procedure "com.atproto.space.deleteRecord".
const NSIDDeleteRecord = "com.atproto.space.deleteRecord"

// DeleteRecordInput is the input of "com.atproto.space.deleteRecord".
type DeleteRecordInput struct {
	Space      string `json:"space"`
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`
}

// DeleteRecord calls "com.atproto.space.deleteRecord" — delete a record from a space repo (or make sure it does not exist).
//
// Implemented by the PDS. Requires auth.
//
// Errors: SpaceNotFound.
func DeleteRecord(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input DeleteRecordInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDDeleteRecord, ""), header, input)
}
