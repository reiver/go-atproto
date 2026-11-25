package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDRemoveWriter is the NSID of the procedure "com.atproto.space.removeWriter".
const NSIDRemoveWriter = "com.atproto.space.removeWriter"

// RemoveWriterInput is the input of "com.atproto.space.removeWriter".
type RemoveWriterInput struct {
	Space string `json:"space"`
	Repo  string `json:"repo"`
}

// RemoveWriter calls "com.atproto.space.removeWriter" — remove an account from a space's writer set (as listed by listWriters).
//
// The removal is sent to the registered syncers.
// The account is added back if it makes another authorized write to the space.
// It does nothing if the account is not in the writer set.
//
// Implemented by the PDS. Requires auth as the space authority.
//
// Errors: SpaceNotFound, NotAuthorized.
func RemoveWriter(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input RemoveWriterInput) error {
	return executeNoOutput(ctx, client, methodURL(endpoint, NSIDRemoveWriter, ""), header, input)
}
