package didplc

import (
	"strings"
)

// IsDIDPlaceholder returns whether the passed value is a DID Placeholder.
func IsDIDPlaceholder(did string) bool {
	return strings.HasPrefix(did, "did:plc:")
}
