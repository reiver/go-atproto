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
// Decode into a map[string]any or an appropriate struct.
func SubscribeRepos() (iter.Iterator, error) {
	const nsid string = "com.atproto.sync.subscribeRepos"

        var xrpcURL xrpc.URL = xrpc.URL{
                Host: config.FireHoseHost,
                NSID: nsid,
        }

	var url string = xrpcURL.String()

	var xrpcIterator xrpc.Iterator
	{
		var err error

		xrpcIterator, err = xrpc.Subscribe(url)
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
			var msg internalSubscriptionMessage = internalSubscriptionMessage(value)

			var header  internalSubscriptionMessageHeader
			var payload internalSubscriptionMessagePayload

			err := msg.Decode(&header, &payload)
			if nil != err {
				return reiver_iter.Empty, nil
			}

			var blocks []byte
			{
				const name string = "blocks"

				if nil == payload {
					return reiver_iter.Empty, nil
				}

				values, found := payload[name]
				if !found {
					return reiver_iter.Empty, nil
				}

				var casted bool
				blocks, casted = values.([]byte)
				if !casted {
					return reiver_iter.Empty, nil
				}
			}

			iterator, err := car.NewIteratorFromBytes(blocks)
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

		iterator = &flattenedIterator
	}

        return iterator, nil
}
