package identity

import (
	"net/url"

	"github.com/reiver/go-xrpc"

	"github.com/reiver/go-atproto/internal/config"
)

func ResolveHandle(dst any, handle string) error {

	const nsid string = "com.atproto.identity.resolveHandle"
	const host string = config.DefaultHost

	var query string = "handle=" + url.QueryEscape(handle)

	var url string = xrpc.ConstructURL(host, nsid, query).String()

	return xrpc.Query(dst, url)
}

type ResolveHandleResponse struct {
	DID string `json:"did"`
}
