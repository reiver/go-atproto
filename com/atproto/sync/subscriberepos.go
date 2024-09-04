package sync

import (
	reiver_iter "github.com/reiver/go-iter"
	"github.com/reiver/go-xrpc"

	"github.com/reiver/go-atproto/internal/car"
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

	var xrpcIterator xrpc.Iterator
	{
		var err error

		xrpcIterator, err = xrpc.Subscribe(xrpcURL.String())
		if nil != err {
			return nil, err
		}
		if nil == xrpcIterator {
			return nil, errNilXRPCIterator
		}
	}

	var iterators reiver_iter.Iterators
	{
		fn := func(value []byte) (reiver_iter.Iterator, error) {
			iterator, err := car.NewIteratorFromBytes(value)
			return iterator, err
		}

		iterators = &reiver_iter.SplitIterators[[]byte]{
			Iterator: xrpcIterator,
			Func: fn,
		}
	}

	var iterator iter.Iterator
	{
		var flattenedIterator reiver_iter.FlattenedIterator
		flattenedIterator.Iterators = iterators
	}

        return iterator, nil
}
