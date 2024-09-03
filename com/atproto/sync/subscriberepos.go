package sync

import (
	"github.com/reiver/go-xrpc"

	"github.com/reiver/go-atproto/internal/config"
	"github.com/reiver/go-atproto/iter"
)

// SubscribeRepos calls XRPC "com.atproto.sync.subscribeRepos".
//
// Decode into [SubscriptionMessage].
func SubscribeRepos() (iter.Iterator, error) {
	const nsid string = "com.atproto.sync.subscribeRepos"

        var xrpcURL xrpc.URL = xrpc.URL{
                Host: config.FireHoseHost,
                NSID: nsid,
        }

        return xrpc.Subscribe(xrpcURL.String())
}
