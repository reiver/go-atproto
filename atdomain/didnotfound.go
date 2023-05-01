package atdomain

import (
	"fmt"
)

var _ error = DIDNotFound{}

// DIDNotFound is an error that can be returned by LookupDID if the "_atproto." sub-domain doesn't exist
// or there are no "did=" DNS TXT records for that sub-domain.
type DIDNotFound struct {
	txtdomainname string
	err error
}

func (receiver DIDNotFound) Error() string {
	return fmt.Sprintf("atproto: DID not found — DNS TXT domain-name=%q", receiver.txtdomainname)
}

func (receiver DIDNotFound) Unwrap() error {
	return receiver.err
}
