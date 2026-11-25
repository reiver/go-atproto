package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDCreateSpace is the NSID of the procedure "com.atproto.space.createSpace".
const NSIDCreateSpace = "com.atproto.space.createSpace"

// CreateSpaceInput is the input of "com.atproto.space.createSpace".
type CreateSpaceInput struct {
	// SpaceType is the type of the space (an NSID).
	SpaceType string `json:"spaceType"`

	// Skey is the space key. If "", the PDS makes one (a TID).
	Skey string `json:"skey,omitempty"`

	Policy Policy `json:"policy"`
}

// CreateSpaceOutput is the output of "com.atproto.space.createSpace".
type CreateSpaceOutput struct {
	Space string `json:"space"`
}

// CreateSpace calls "com.atproto.space.createSpace" — create a new space, with the authenticated account as its authority.
//
// Implemented by the PDS. Requires auth.
//
// Errors: SpaceAlreadyExists, UnsupportedPolicy.
func CreateSpace(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input CreateSpaceInput) (CreateSpaceOutput, error) {
	return execute[CreateSpaceOutput](ctx, client, methodURL(endpoint, NSIDCreateSpace, ""), header, input)
}
