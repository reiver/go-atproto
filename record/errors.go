package record

import (
	"github.com/reiver/go-erorr"
)

const (
	errNilSource = erorr.Error("atproto: nil source")
	errNotTyped  = erorr.Error("atproto: not found")
)
