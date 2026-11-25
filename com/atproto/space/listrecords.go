package space

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDListRecords is the NSID of the query "com.atproto.space.listRecords".
const NSIDListRecords = "com.atproto.space.listRecords"

// ListRecordsParams are the params of "com.atproto.space.listRecords".
type ListRecordsParams struct {
	Space string `query:"space"`
	Repo  string `query:"repo"`

	// Collection, if not "", lists only the records in that collection.
	Collection string `query:"collection"`

	// Limit is from 1 to 1000. If 0, it is left out (the default is 50).
	Limit         int64  `query:"limit"`
	Cursor        string `query:"cursor"`
	ExcludeValues bool   `query:"excludeValues"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver ListRecordsParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDListRecords, encodeQuery(receiver))
}

// ListRecordsOutput is the output of "com.atproto.space.listRecords".
type ListRecordsOutput struct {
	Cursor  string              `json:"cursor,omitempty"`
	Records []ListRecordsRecord `json:"records"`
}

// ListRecordsRecord is "com.atproto.space.listRecords#record".
type ListRecordsRecord struct {
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`
	CID        string `json:"cid"`

	// Value is the record. It is left out if ExcludeValues was set.
	Value json.RawMessage `json:"value,omitempty"`
}

// ListRecords calls "com.atproto.space.listRecords" — list the records in a space repo (optionally, in only one collection).
//
// Implemented by the PDS. Requires auth as the repo's account, or with a space credential.
//
// Errors: SpaceNotFound, RepoNotFound, RepoTakendown, RepoSuspended, RepoDeactivated, InvalidCursor.
func ListRecords(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params ListRecordsParams) (ListRecordsOutput, error) {
	return query[ListRecordsOutput](ctx, client, params.URL(endpoint), header)
}
