package sync

import (
	"github.com/reiver/go-erorr"
)

const (
	errNilDecoder                   = erorr.Error("atproto: nil decoder")
	errNilPayload                   = erorr.Error("atproto: nil payload")
	errNilReader                    = erorr.Error("atproto: nil reader")
	errNilReceiver                  = erorr.Error("atproto: nil receiver")
	errNilSubscriptionMessageHeader = erorr.Error("atproto: nil subscription-message-header")
	errNoBlocks                     = erorr.Error("atproto: no blocks")
	errBlocksNotBytes               = erorr.Error("atproto: blocks not bytes")
)
