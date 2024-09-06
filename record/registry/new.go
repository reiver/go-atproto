package registry

import (
	"github.com/reiver/go-reg"
)

var NewFuncs reg.Registry[func()Record]
