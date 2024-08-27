package sync

import (
	"github.com/reiver/go-xrpc"

	"github.com/reiver/go-atproto/internal/config"
)

func SubscribeRepos() (xrpc.Iterator, error) {
	const nsid string = "com.atproto.sync.subscribeRepos"

        var xrpcURL xrpc.URL = xrpc.URL{
                Host: config.FireHoseHost,
                NSID: nsid,
        }

        return xrpc.Subscribe(xrpcURL.String())
}
