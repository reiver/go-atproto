package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDClearRecords is the NSID of the procedure "com.atproto.space.clearRecords".
const NSIDClearRecords = "com.atproto.space.clearRecords"

// ClearRecordsInput is the input of "com.atproto.space.clearRecords".
type ClearRecordsInput struct {
	Space string `json:"space"`
	Repo  string `json:"repo"`
}

// ClearRecords calls "com.atproto.space.clearRecords" — delete all the records in the authenticated account's space repo.
//
// In the space authority's repo, the records published from the space's configuration (such as its policy) are kept.
//
// Implemented by the PDS. Requires auth.
func ClearRecords(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input ClearRecordsInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDClearRecords, ""), header, input)
}
