package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDAuthorizeSpaceSync is the NSID of the procedure "com.atproto.space.authorizeSpaceSync".
const NSIDAuthorizeSpaceSync = "com.atproto.space.authorizeSpaceSync"

// AuthorizeSpaceSyncInput is the input of "com.atproto.space.authorizeSpaceSync".
type AuthorizeSpaceSyncInput struct {
	Space string `json:"space"`

	// CredentialKey is the P-256 did:key that the credential is bound to.
	CredentialKey string `json:"credentialKey"`

	// DelegationToken is from getDelegationToken (see [GetDelegationToken]).
	DelegationToken   string `json:"delegationToken,omitempty"`
	ClientAttestation string `json:"clientAttestation,omitempty"`

	// Notify, if not nil, registers an endpoint to be notified of updates to the space (with notifySpaceUpdate).
	// The registration stays active while the credential is valid.
	// After it expires, the endpoint receives only a prompt to reauthorize, or notice of the space's deletion.
	Notify *AuthorizeSpaceSyncNotify `json:"notify,omitempty"`
}

// AuthorizeSpaceSyncNotify is "com.atproto.space.authorizeSpaceSync#notify".
type AuthorizeSpaceSyncNotify struct {
	// Endpoint is the base HTTPS URL that the notification methods are called on.
	Endpoint string `json:"endpoint"`

	// Token is an opaque secret (at most 512 characters), included in the body of each request to the endpoint.
	// The endpoint and token together identify the registration, so the token should stay the same across authorizations.
	Token string `json:"token"`
}

// AuthorizeSpaceSyncOutput is the output of "com.atproto.space.authorizeSpaceSync".
type AuthorizeSpaceSyncOutput struct {
	Credential string `json:"credential"`

	// SpaceRev is the space's current revision. It is left out before the space's first update.
	SpaceRev string `json:"spaceRev,omitempty"`
}

// AuthorizeSpaceSync calls "com.atproto.space.authorizeSpaceSync" — authorize syncing a space.
// It returns a new credential for reading data across the space, and the space's current revision;
// and (optionally) registers for notifications of updates to the space.
//
// Depending on the space's policy, a delegation token, a client attestation, or both may be required.
//
// Implemented by the space host.
//
// Errors: SpaceNotFound, SpaceDeleted, AccountNotAuthorized, AppNotAuthorized, NotAuthorized,
// InvalidDelegationToken, InvalidClientAttestation, RegistrationNotConfirmed,
// RepoTakendown (the space authority), RepoDeactivated (the space authority).
func AuthorizeSpaceSync(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, input AuthorizeSpaceSyncInput) (AuthorizeSpaceSyncOutput, error) {
	return execute[AuthorizeSpaceSyncOutput](ctx, client, methodURL(endpoint, NSIDAuthorizeSpaceSync, ""), header, input)
}
