package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDConfirmRegistration is the NSID of the procedure "com.atproto.space.confirmRegistration".
const NSIDConfirmRegistration = "com.atproto.space.confirmRegistration"

// ConfirmRegistrationInput is the input of "com.atproto.space.confirmRegistration".
type ConfirmRegistrationInput struct {
	Space     string `json:"space"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

// ConfirmRegistrationOutput is the output of "com.atproto.space.confirmRegistration".
type ConfirmRegistrationOutput struct {
	Challenge string `json:"challenge"`
}

// ConfirmRegistration calls "com.atproto.space.confirmRegistration" — confirm a new notification registration.
//
// The space host calls it on the notify endpoint passed to authorizeSpaceSync, before registering it.
// The syncing service checks that it issued the token for the space, and echoes back the challenge.
//
// Implemented by the syncing service. ('endpoint' is the notify endpoint.)
func ConfirmRegistration(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input ConfirmRegistrationInput) (ConfirmRegistrationOutput, error) {
	return execute[ConfirmRegistrationOutput](ctx, client, methodURL(endpoint, NSIDConfirmRegistration, ""), header, input)
}
