package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDUpdateSpace is the NSID of the procedure "com.atproto.space.updateSpace".
const NSIDUpdateSpace = "com.atproto.space.updateSpace"

// UpdateSpaceInput is the input of "com.atproto.space.updateSpace".
type UpdateSpaceInput struct {
	Space  string `json:"space"`
	Policy Policy `json:"policy"`
}

// UpdateSpace calls "com.atproto.space.updateSpace" — replace the policy of a space.
//
// Changing to a different type of policy discards the state held by the previous policy (such as a member list).
//
// Implemented by the PDS. Requires auth as the space authority.
//
// Errors: SpaceNotFound, NotAuthorized, UnsupportedPolicy.
func UpdateSpace(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input UpdateSpaceInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDUpdateSpace, ""), header, input)
}
