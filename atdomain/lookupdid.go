package atdomain

import (
	"fmt"
	"net"
)

// LookupDID returns any DIDs registered for the domain-name.
func LookupDID(domainname string) ([]string, error) {

	var txtdomainname string = subdomainTXT(domainname)

	// Return all the DNS TXT records for the "_atproto" sub-domain.
	//
	// Note that these could contain records that are NOT DIDs.
	// For example:
	//
	//	did=did:plc:m2jwplpernhxkzbo4ev5ljwj
	//	once=apple
	//	twice=banana
	//	thrice=cherry
	//	fource=date
	//	did=did:something:abcde12345
	//
	// We will filter out the non-DIDs later.
	var txtresponses []string
	{
		var err error

		txtresponses, err = net.LookupTXT(txtdomainname)
		if nil != err {
			err = lookupTXTError(err, txtdomainname)
			return nil, err
		}
	}

	var dids []string = filterDIDs(txtresponses)
	if len(dids) < 1 {
		return nil, DIDNotFound{
			txtdomainname:txtdomainname,
			err:fmt.Errorf("atproto: no \"did=\" DNS TXT records for %q", txtdomainname),
		}
	}

	return dids, nil
}
