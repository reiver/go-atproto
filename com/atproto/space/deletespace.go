package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDDeleteSpace is the NSID of the procedure "com.atproto.space.deleteSpace".
const NSIDDeleteSpace = "com.atproto.space.deleteSpace"

// DeleteSpaceInput is the input of "com.atproto.space.deleteSpace".
type DeleteSpaceInput struct {
	Space string `json:"space"`
}

// DeleteSpace calls "com.atproto.space.deleteSpace" — delete a space, along with the authority's space repo.
//
// Implemented by the PDS. Requires auth as the space authority.
//
// Errors: SpaceNotFound, NotAuthorized.
func DeleteSpace(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input DeleteSpaceInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDDeleteSpace, ""), header, input)
}
