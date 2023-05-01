package atdomain

import (
	"fmt"
	"net"
)

// lookupTXTError wraps any error we get back from net.LookupTXT().
func lookupTXTError(err error, txtdomainname string) error {
	if nil == err {
		return nil
	}

	switch casted := err.(type) {
	case *net.DNSError:
		switch {
		case casted.IsNotFound:
			return DIDNotFound{
				txtdomainname:txtdomainname,
				err:err,
			}
		}
	}

	return fmt.Errorf("atproto: problem looking-up DNS TXT record(s) for domnain-name %q: %w", txtdomainname, err)
}
