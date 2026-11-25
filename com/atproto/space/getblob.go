package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDGetBlob is the NSID of the query "com.atproto.space.getBlob".
const NSIDGetBlob = "com.atproto.space.getBlob"

// GetBlobParams are the params of "com.atproto.space.getBlob".
type GetBlobParams struct {
	Space string `query:"space"`
	Repo  string `query:"repo"`
	CID   string `query:"cid"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver GetBlobParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDGetBlob, encodeQuery(receiver))
}

// GetBlob calls "com.atproto.space.getBlob" — get a blob referenced by a record in a space repo, as originally uploaded.
//
// It returns the HTTP response, with its body (the blob) unread — the caller must close it.
//
// Implemented by the PDS. Requires auth as the repo's account, or with a space credential.
//
// Errors: BlobNotFound, SpaceNotFound, RepoNotFound, RepoTakendown, RepoSuspended, RepoDeactivated.
func GetBlob(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params GetBlobParams) (*http.Response, error) {
	return client.Query(ctx, params.URL(endpoint), header)
}
