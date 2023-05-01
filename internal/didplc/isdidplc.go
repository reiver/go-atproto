package didplc

import (
	"strings"
)

func IsDIDPlaceholder(did string) bool {
	return strings.HasPrefix(did, "did:plc:")
}
