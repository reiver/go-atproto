package space

import (
	"context"
	"io"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDImportRepo is the NSID of the procedure "com.atproto.space.importRepo".
const NSIDImportRepo = "com.atproto.space.importRepo"

// ImportRepo calls "com.atproto.space.importRepo" — import the authenticated account's space repo from a STAR file ('body'), for account migration.
//
// If the account is the space's authority, the space must first be imported with importSpaceMeta (see [ImportSpaceMeta]).
//
// Implemented by the PDS. Requires auth.
//
// Errors: SpaceNotFound, InvalidRepoFile, RepoMismatch.
func ImportRepo(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, body io.Reader) error {
	return executeSTAR(ctx, client, methodURL(endpoint, NSIDImportRepo, ""), header, body)
}
