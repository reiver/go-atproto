package car

import (
	"github.com/reiver/go-erorr"
)

const (
	errNilBlock        = erorr.Error("atproto: nil block")
	errNilBlockRawData = erorr.Error("atproto: nil block raw-data")
	errNilCarReader    = erorr.Error("atproto: nil car-reader")
	errNilReceiver     = erorr.Error("atproto: nil receiver")
)
