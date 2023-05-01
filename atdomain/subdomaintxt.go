package atdomain

// With the AT-Protocol, the DID is stored in a DNS TXT record on a particular sub-domain of the domain.
//
// If, for example, the domain is "example.com", then the sub-domain is "_atproto.example.com".
// And if, for example, the domain is "changelog.ca", then the sub-domain is "_atproto.changelog.ca".
// And if, for example, the domain is "once.twice.thrice.fource.xyz", then the sub-domain is "_atproto.once.twice.thrice.fource.xyz".
//
// Here we construct the sub-domain that will be used in the DNS TXT record lookup.
func subdomainTXT(domainname string) string {
	return "_atproto." + domainname
}
