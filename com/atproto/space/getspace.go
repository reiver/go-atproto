package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDGetSpace is the NSID of the query "com.atproto.space.getSpace".
const NSIDGetSpace = "com.atproto.space.getSpace"

// GetSpaceParams are the params of "com.atproto.space.getSpace".
type GetSpaceParams struct {
	Space string `query:"space"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver GetSpaceParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDGetSpace, encodeQuery(receiver))
}

// GetSpaceOutput is the output of "com.atproto.space.getSpace".
type GetSpaceOutput struct {
	Space  string `json:"space"`
	Policy Policy `json:"policy"`
}

// GetSpace calls "com.atproto.space.getSpace" — get a space and its policy.
//
// Implemented by the PDS. Requires auth as the space authority.
//
// Errors: SpaceNotFound, NotAuthorized.
func GetSpace(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params GetSpaceParams) (GetSpaceOutput, error) {
	return query[GetSpaceOutput](ctx, client, params.URL(endpoint), header)
}
