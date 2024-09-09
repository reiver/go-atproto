package sync

import (
	"bytes"
	"io"

	"github.com/ipld/go-car"
	"github.com/reiver/go-erorr"
)

// The message that comes back from the Bluesky Firehose websocket is 2 CBOR objects concatenated with each other.
// The first part is called the message-header. The second part is called the message-payload.
//
// internalSubscriptionMessagePayload represents the message-payload.
type internalSubscriptionMessagePayload map[string]any

func (receiver internalSubscriptionMessagePayload) Blocks() (*car.CarReader, error) {
	const name string = "blocks"

	if nil == receiver {
		return nil, errNilReceiver
	}

	values, found := receiver[name]
	if !found {
		return nil, errNoBlocks
	}

	var blocks []byte
	{
		var casted bool
		blocks, casted = values.([]byte)
		if !casted {
			return nil, errBlocksNotBytes
		}
	}

	var carreader *car.CarReader
	{
		var reader io.Reader = bytes.NewReader(blocks)
		if nil == reader {
			return nil, errNilReader
		}

		var err error

		carreader, err = car.NewCarReader(reader)
		if nil != err {
			return nil, erorr.Errorf("bsky: problem creating CAR (Content Addressable aRchives) reader: %w", err)
		}
	}

	return carreader, nil
}

func (receiver internalSubscriptionMessagePayload) Rebase() (bool, bool) {
	const name string = "rebase"
	var empty bool

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case bool:
		return casted, true
	default:
		return empty,  false
	}
}

func (receiver internalSubscriptionMessagePayload) Repo() (string, bool) {
	const name string = "repo"
	var empty string

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case string:
		return casted, true
	default:
		return empty,  false
	}
}

func (receiver internalSubscriptionMessagePayload) Rev() (string, bool) {
	const name string = "rev"
	var empty string

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case string:
		return casted, true
	default:
		return empty,  false
	}
}

func (receiver internalSubscriptionMessagePayload) Seq() (uint64, bool) {
	const name string = "seq"
	var empty uint64

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case uint64:
		return casted, true
	default:
		return empty,  false
	}
}

func (receiver internalSubscriptionMessagePayload) Since() (string, bool) {
	const name string = "since"
	var empty string

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case string:
		return casted, true
	default:
		return empty,  false
	}
}

func (receiver internalSubscriptionMessagePayload) TooBig() (bool, bool) {
	const name string = "tooBig"
	var empty bool

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case bool:
		return casted, true
	default:
		return empty,  false
	}
}

func (receiver internalSubscriptionMessagePayload) Time() (string, bool) {
	const name string = "time"
	var empty string

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case string:
		return casted, true
	default:
		return empty,  false
	}
}
