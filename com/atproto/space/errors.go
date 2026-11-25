package space

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	ErrBadBytes    = erorr.Error("space: bad bytes")
	ErrBadUnion    = erorr.Error("space: bad union")
	ErrUnknownType = erorr.Error("space: unknown $type in a closed union")
	ErrMissingType = erorr.Error("space: missing $type in a union")
)
