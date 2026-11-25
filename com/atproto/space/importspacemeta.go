package space

import (
	"context"
	"io"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDImportSpaceMeta is the NSID of the procedure "com.atproto.space.importSpaceMeta".
const NSIDImportSpaceMeta = "com.atproto.space.importSpaceMeta"

// ImportSpaceMeta calls "com.atproto.space.importSpaceMeta" — import a space's state from a STAR file ('body') made by exportSpaceMeta, for account migration.
//
// It must be called before importing the authority's space repo with importRepo (see [ImportRepo]).
//
// Implemented by the PDS. Requires auth as the space authority.
//
// Errors: NotAuthorized, UnsupportedPolicy, InvalidSpaceMeta.
func ImportSpaceMeta(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, body io.Reader) error {
	return executeSTAR(ctx, client, methodURL(endpoint, NSIDImportSpaceMeta, ""), header, body)
}
