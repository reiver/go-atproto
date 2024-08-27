package sync

// The message that comes back from the Bluesky Firehose websocket is 2 CBOR objects concatenated with each other.
// The first part is called the message-header. The second part is called the message-payload.
//
// SubscriptionMessageHeader represents the message-header.
type SubscriptionMessageHeader struct {
	Operation Operation `cbor:"op,omitempty"`
	Type      string    `cbor:"t"`
}
