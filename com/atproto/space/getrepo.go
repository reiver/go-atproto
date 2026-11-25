package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDGetRepo is the NSID of the query "com.atproto.space.getRepo".
const NSIDGetRepo = "com.atproto.space.getRepo"

// GetRepoParams are the params of "com.atproto.space.getRepo".
type GetRepoParams struct {
	Space string `query:"space"`
	Repo  string `query:"repo"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver GetRepoParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDGetRepo, encodeQuery(receiver))
}

// GetRepo calls "com.atproto.space.getRepo" — download a space repo as a STAR file.
//
// It returns the HTTP response, with its body (the STAR file) unread — the caller must close it.
//
// Implemented by the PDS. Requires auth as the repo's account, or with a space credential.
// The response is authenticated with a MAC (see the package documentation).
//
// Errors: InvalidSpaceMac, SpaceNotFound, RepoNotFound, RepoTakendown, RepoSuspended, RepoDeactivated.
func GetRepo(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params GetRepoParams) (*http.Response, error) {
	return client.Query(ctx, params.URL(endpoint), header)
}
