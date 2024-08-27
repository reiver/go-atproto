package sync

import (
	"bytes"
	"io"
	"strings"

	"github.com/brianolson/cbor_go" // because it can handle extra data after the (first) CBOR
	cbor2 "github.com/fxamacker/cbor/v2" // because it can decode into a map
	"github.com/reiver/go-erorr"
)

// SubscriptionMessage represents a message from the Bluesky Firehose websocket.
type SubscriptionMessage []byte

// Decode decodes the message from the Bluesky Firehose websocket into the message-header and the message-payload.
func (receiver SubscriptionMessage) Decode(header *SubscriptionMessageHeader, payload *SubscriptionMessagePayload) error {

	//@TODO: This is NOT an ideal way of decoding this.
	//       Maybe implement a CBOR codec later to do this better.

	if nil == header {
		return errNilSubscriptionMessageHeader
	}
	if nil == payload {
		return errNilPayload
	}

	{
		var reader io.Reader = bytes.NewReader(receiver)
		if nil == reader {
			return errNilReader
		}

		var decoder *cbor.Decoder = cbor.NewDecoder(reader)
		if nil == decoder {
			return errNilDecoder
		}

		err := decoder.Decode(header)
		if nil != err {
			return erorr.Errorf("atproto: problem decoding cbor header: %w", err)
		}
	}

	var headersize int
	{
		var builder strings.Builder

		err := cbor.Encode(&builder, header)
		if nil != err {
			return erorr.Errorf("atproto: problem encoding cbor header (after decoding it, to get length): %w", err)
		}

		headersize = builder.Len()
	}

	{
		var p []byte = receiver[headersize:]

		err := cbor2.Unmarshal(p, &payload)
		if nil != err {
			return erorr.Errorf("atproto: problem decoding cbor payload: %w", err)
		}
	}

	return nil
}

