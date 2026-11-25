package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDNotifyRepoUpdate is the NSID of the procedure "com.atproto.space.notifyRepoUpdate".
const NSIDNotifyRepoUpdate = "com.atproto.space.notifyRepoUpdate"

// NotifyRepoUpdateInput is the input of "com.atproto.space.notifyRepoUpdate".
type NotifyRepoUpdateInput struct {
	Space   string `json:"space"`
	Repo    string `json:"repo"`
	RepoRev string `json:"repoRev"`
	Data    Bytes  `json:"data"`
}

// NotifyRepoUpdate calls "com.atproto.space.notifyRepoUpdate" — notify the space host that a space repo has been updated.
//
// Repo hosts send it, with service auth.
//
// Implemented by the space host.
//
// Errors: SpaceNotFound, FutureRev (the repo revision is too far in the future).
func NotifyRepoUpdate(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input NotifyRepoUpdateInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDNotifyRepoUpdate, ""), header, input)
}
