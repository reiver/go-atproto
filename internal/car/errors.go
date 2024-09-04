package car

import (
	"github.com/reiver/go-erorr"
)

const (
	errNilBlock        = erorr.Error("atproto: nil block")
	errNilBlockRawData = erorr.Error("atproto: nil block raw-data")
	errNilBytes        = erorr.Error("atproto: nil bytes")
	errNilCarReader    = erorr.Error("atproto: nil car-reader")
	errNilReader       = erorr.Error("atproto: nil reader")
	errNilReceiver     = erorr.Error("atproto: nil receiver")
)
