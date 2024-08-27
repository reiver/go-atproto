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
// SubscriptionMessagePayload represents the message-payload.
type SubscriptionMessagePayload map[string]any

func (receiver SubscriptionMessagePayload) Blocks() (*car.CarReader, error) {
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

func (receiver SubscriptionMessagePayload) Rebase() (bool, bool) {
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

func (receiver SubscriptionMessagePayload) Rev() (string, bool) {
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

func (receiver SubscriptionMessagePayload) Seq() (int, bool) {
	const name string = "seq"
	var empty int

	if nil == receiver {
		return empty, false
	}

	value, found := receiver[name]
	if !found {
		return empty, false
	}

	switch casted := value.(type) {
	case int:
		return casted, true
	default:
		return empty,  false
	}
}

func (receiver SubscriptionMessagePayload) Since() (string, bool) {
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

func (receiver SubscriptionMessagePayload) TooBig() (bool, bool) {
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

func (receiver SubscriptionMessagePayload) Time() (string, bool) {
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
