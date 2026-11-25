package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDExportSpaceMeta is the NSID of the query "com.atproto.space.exportSpaceMeta".
const NSIDExportSpaceMeta = "com.atproto.space.exportSpaceMeta"

// ExportSpaceMetaParams are the params of "com.atproto.space.exportSpaceMeta".
type ExportSpaceMetaParams struct {
	Space string `query:"space"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver ExportSpaceMetaParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDExportSpaceMeta, encodeQuery(receiver))
}

// ExportSpaceMeta calls "com.atproto.space.exportSpaceMeta" — export a space's state (such as its policy and member list) as a STAR file, for account migration.
//
// It returns the HTTP response, with its body (the STAR file) unread — the caller must close it.
//
// Implemented by the PDS. Requires auth as the space authority.
//
// Errors: SpaceNotFound, NotAuthorized.
func ExportSpaceMeta(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params ExportSpaceMetaParams) (*http.Response, error) {
	return client.Query(ctx, params.URL(endpoint), header)
}
