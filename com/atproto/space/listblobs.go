package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDListBlobs is the NSID of the query "com.atproto.space.listBlobs".
const NSIDListBlobs = "com.atproto.space.listBlobs"

// ListBlobsParams are the params of "com.atproto.space.listBlobs".
type ListBlobsParams struct {
	Space string `query:"space"`
	Repo  string `query:"repo"`

	// Limit is from 1 to 1000. If 0, it is left out (the default is 500).
	Limit  int64  `query:"limit"`
	Cursor string `query:"cursor"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver ListBlobsParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDListBlobs, encodeQuery(receiver))
}

// ListBlobsOutput is the output of "com.atproto.space.listBlobs".
type ListBlobsOutput struct {
	Cursor string   `json:"cursor,omitempty"`
	CIDs   []string `json:"cids"`
}

// ListBlobs calls "com.atproto.space.listBlobs" — list the CIDs of the blobs referenced by records in a space repo.
//
// Implemented by the PDS. Requires auth as the repo's account, or with a space credential.
//
// Errors: SpaceNotFound, RepoNotFound, RepoTakendown, RepoSuspended, RepoDeactivated.
func ListBlobs(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params ListBlobsParams) (ListBlobsOutput, error) {
	return query[ListBlobsOutput](ctx, client, params.URL(endpoint), header)
}
