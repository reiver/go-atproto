package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDGetDelegationToken is the NSID of the query "com.atproto.space.getDelegationToken".
const NSIDGetDelegationToken = "com.atproto.space.getDelegationToken"

// GetDelegationTokenParams are the params of "com.atproto.space.getDelegationToken".
type GetDelegationTokenParams struct {
	Space string `query:"space"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver GetDelegationTokenParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDGetDelegationToken, encodeQuery(receiver))
}

// GetDelegationTokenOutput is the output of "com.atproto.space.getDelegationToken".
type GetDelegationTokenOutput struct {
	Token string `json:"token"`
}

// GetDelegationToken calls "com.atproto.space.getDelegationToken" — get a delegation token for a space,
// to exchange with the space host for a space credential (see [AuthorizeSpaceSync]).
//
// Implemented by the PDS. Requires auth.
func GetDelegationToken(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params GetDelegationTokenParams) (GetDelegationTokenOutput, error) {
	return query[GetDelegationTokenOutput](ctx, client, params.URL(endpoint), header)
}
