package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDNotifyCredentialsRevoked is the NSID of the procedure "com.atproto.space.notifyCredentialsRevoked".
const NSIDNotifyCredentialsRevoked = "com.atproto.space.notifyCredentialsRevoked"

// NotifyCredentialsRevokedInput is the input of "com.atproto.space.notifyCredentialsRevoked".
type NotifyCredentialsRevokedInput struct {
	Space string `json:"space"`

	// Credentials are the IDs ("jti") of the revoked credentials — 1 to 100 of them.
	Credentials []string `json:"credentials"`
}

// NotifyCredentialsRevoked calls "com.atproto.space.notifyCredentialsRevoked" — notify a repo host that space credentials have been revoked.
//
// The space host sends it.
//
// Implemented by the PDS. Requires service auth.
func NotifyCredentialsRevoked(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input NotifyCredentialsRevokedInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDNotifyCredentialsRevoked, ""), header, input)
}
