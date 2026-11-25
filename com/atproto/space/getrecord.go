package space

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDGetRecord is the NSID of the query "com.atproto.space.getRecord".
const NSIDGetRecord = "com.atproto.space.getRecord"

// GetRecordParams are the params of "com.atproto.space.getRecord".
type GetRecordParams struct {
	Space      string `query:"space"`
	Repo       string `query:"repo"`
	Collection string `query:"collection"`
	Rkey       string `query:"rkey"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver GetRecordParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDGetRecord, encodeQuery(receiver))
}

// GetRecordOutput is the output of "com.atproto.space.getRecord".
type GetRecordOutput struct {
	URI   string          `json:"uri"`
	CID   string          `json:"cid"`
	Value json.RawMessage `json:"value"`
}

// GetRecord calls "com.atproto.space.getRecord" — get a single record from a space repo.
//
// Implemented by the PDS. Requires auth as the repo's account, or with a space credential.
// The response is authenticated with a MAC (see the package documentation).
//
// Errors: InvalidSpaceMac, RecordNotFound, SpaceNotFound, RepoNotFound, RepoTakendown, RepoSuspended, RepoDeactivated.
func GetRecord(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params GetRecordParams) (GetRecordOutput, error) {
	return query[GetRecordOutput](ctx, client, params.URL(endpoint), header)
}
