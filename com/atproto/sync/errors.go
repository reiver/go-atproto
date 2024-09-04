package sync

import (
	"github.com/reiver/go-erorr"
)

const (
	errNilDecoder                    = erorr.Error("atproto: nil decoder")
	errNilDestination                = erorr.Error("atproto: nil destination")
	errNilPayload                    = erorr.Error("atproto: nil payload")
	errNilReader                     = erorr.Error("atproto: nil reader")
	errNilReceiver                   = erorr.Error("atproto: nil receiver")
	errNilSubscriptionMessageHeader  = erorr.Error("atproto: nil subscription-message-header")
	errNilSubscriptionMessagePayload = erorr.Error("atproto: nil subscription-message-payload")
	errNilXRPCIterator               = erorr.Error("atproto: nil xrpc iterator")
	errNoBlocks                      = erorr.Error("atproto: no blocks")
	errBlocksNotBytes                = erorr.Error("atproto: blocks not bytes")
)
