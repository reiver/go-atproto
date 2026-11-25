package space

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// NSIDListRepoOps is the NSID of the query "com.atproto.space.listRepoOps".
const NSIDListRepoOps = "com.atproto.space.listRepoOps"

// ListRepoOpsParams are the params of "com.atproto.space.listRepoOps".
type ListRepoOpsParams struct {
	Space string `query:"space"`
	Repo  string `query:"repo"`

	// Since, if not "", lists only the operations after this repo revision (a TID).
	Since string `query:"since"`

	// Limit is from 1 to 1000. If 0, it is left out (the default is 100).
	Limit         int64  `query:"limit"`
	Cursor        string `query:"cursor"`
	ExcludeValues bool   `query:"excludeValues"`
}

// URL returns the URL of the query at 'endpoint' (ex: "https://pds.example.com").
func (receiver ListRepoOpsParams) URL(endpoint string) string {
	return methodURL(endpoint, NSIDListRepoOps, encodeQuery(receiver))
}

// ListRepoOpsOutput is the output of "com.atproto.space.listRepoOps".
type ListRepoOpsOutput struct {
	Ops []ListRepoOpsOpEntry `json:"ops"`

	// Commit is the repo's current commit. Once the last page (the one with no cursor) has been applied, compare it with the result.
	Commit Commit `json:"commit"`

	Cursor string `json:"cursor,omitempty"`
}

// ListRepoOpsOpEntry is "com.atproto.space.listRepoOps#opEntry" — one operation on a space repo.
// Operations that share a rev were applied in the same commit.
type ListRepoOpsOpEntry struct {
	Rev        string `json:"rev"`
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`

	// CID is the CID of the record after the operation. It is nil (JSON null) for a delete.
	CID *string `json:"cid"`

	// Prev is the CID of the record before the operation. It is nil (JSON null) for a create.
	Prev *string `json:"prev"`

	// Value is the record after the operation.
	// It is left out for a delete, if ExcludeValues was set, or if the record has changed since.
	Value json.RawMessage `json:"value,omitempty"`
}

// ListRepoOps calls "com.atproto.space.listRepoOps" — list the operations on a space repo after a given revision (for incremental sync).
//
// Hosts may truncate the log; use [GetRepo] for a full sync.
//
// Implemented by the PDS. Requires auth as the repo's account, or with a space credential.
// The response is authenticated with a MAC (see the package documentation).
//
// Errors: InvalidSpaceMac, SpaceNotFound, RepoNotFound, RepoTakendown, RepoSuspended, RepoDeactivated, InvalidCursor.
func ListRepoOps(ctx context.Context, client xrpc.Client, endpoint string, header http.Header, params ListRepoOpsParams) (ListRepoOpsOutput, error) {
	return query[ListRepoOpsOutput](ctx, client, params.URL(endpoint), header)
}
