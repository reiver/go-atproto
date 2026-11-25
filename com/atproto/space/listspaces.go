package space

import (
	"context"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDListSpaces is the NSID of the query "com.atproto.space.listSpaces".
const NSIDListSpaces = "com.atproto.space.listSpaces"

// ListSpacesParams are the params of "com.atproto.space.listSpaces".
type ListSpacesParams struct {
	// SpaceDID, if not "", lists only the spaces whose authority is this DID.
	SpaceDID string `query:"spaceDid"`

	// SpaceType, if not "", lists only the spaces of this type (an NSID).
	SpaceType string `query:"spaceType"`

	// Limit is from 1 to 100. If 0, it is left out (the default is 50).
	Limit  int64  `query:"limit"`
	Cursor string `query:"cursor"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver ListSpacesParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDListSpaces, encodeQuery(receiver))
}

// ListSpacesOutput is the output of "com.atproto.space.listSpaces".
type ListSpacesOutput struct {
	Cursor string                `json:"cursor,omitempty"`
	Spaces []ListSpacesSpaceView `json:"spaces"`
}

// ListSpacesSpaceView is "com.atproto.space.listSpaces#spaceView".
type ListSpacesSpaceView struct {
	Space string `json:"space"`

	// RepoRev is the current revision of the account's repo in the space, if it has one.
	RepoRev string `json:"repoRev,omitempty"`

	// Data is the current commit data of the account's repo in the space, if it has one.
	Data Bytes `json:"data,omitempty"`

	// SpaceRev is the current revision of the space. It is only included if the account is the space's authority.
	SpaceRev string `json:"spaceRev,omitempty"`
}

// ListSpaces calls "com.atproto.space.listSpaces" — list the spaces that the authenticated account has written to, or is the authority for.
//
// Implemented by the PDS. Requires auth.
//
// Errors: InvalidCursor.
func ListSpaces(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params ListSpacesParams) (ListSpacesOutput, error) {
	return query[ListSpacesOutput](ctx, client, params.URL(endpoint), header)
}
