package record

import (
	"github.com/reiver/go-atproto/record/registry"
)

func New(name string) (Record, bool) {
	var fn func()Record
	var found bool

	fn, found = registry.NewFuncs.Get(name)
	if !found {
		return nil, false
	}

	return fn(), true
}
