package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDGetLatestCommit is the NSID of the query "com.atproto.space.getLatestCommit".
const NSIDGetLatestCommit = "com.atproto.space.getLatestCommit"

// GetLatestCommitParams are the params of "com.atproto.space.getLatestCommit".
type GetLatestCommitParams struct {
	Space string `query:"space"`
	Repo  string `query:"repo"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver GetLatestCommitParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDGetLatestCommit, encodeQuery(receiver))
}

// GetLatestCommitOutput is the output of "com.atproto.space.getLatestCommit".
type GetLatestCommitOutput struct {
	Commit Commit `json:"commit"`
}

// GetLatestCommit calls "com.atproto.space.getLatestCommit" — get the current commit of a space repo.
//
// Implemented by the PDS. Requires auth as the repo's account, or with a space credential.
// The response is authenticated with a MAC (see the package documentation).
//
// Errors: InvalidSpaceMac, SpaceNotFound, RepoNotFound, RepoTakendown, RepoSuspended, RepoDeactivated.
func GetLatestCommit(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params GetLatestCommitParams) (GetLatestCommitOutput, error) {
	return query[GetLatestCommitOutput](ctx, client, params.URL(endpoint), header)
}
