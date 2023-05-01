package didplc

// DocURI returns the URI used to resolve a DID Placeholder (ex: "did:plc:abcde12345") into a DID document.
//
// DocURI returns an empty string it what it is given is not an DID Placeholder.
//
// For example:
//
//	var did string = "did:plc:abcde12345"
//
//	uri := didplc.DocURI(did)
//
//	// uri == "https://plc.directory/did:plc:abcde12345"
//
// I'm assuming the AT-Protocol people will eventually change how this works.
//
// Having a single gateway to resolve DID Placeholder URIs into DID Documents
// is very much NOT decentralized!
func DocURI(did string) string {
	if !IsDIDPlaceholder(did) {
		return ""
	}

	const baseuri string = "https://plc.directory/"

	var uri string = baseuri + did

	return uri
}
