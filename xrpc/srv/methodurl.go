package xrpcsrv

// MethodURL return the BlueSky AT-Protocol method-URL based on the service-endpoint and method-ID.
//
// Most likely you usually won't use this function directly, but would use it indirectly, by calling something that calls it.
//
// For example, if the service-endpoint is:
//
//	"https://bsky.social"
//
// And the method-ID is:
//
//	"com.atproto.server.createSession"
//
// Then MethodURL would return:
//
//	"https://bsky.social/xrpc/com.atproto.server.createSession"
//
func MethodURL(serviceEndpoint string, methodID string) string {
	if "" == serviceEndpoint {
		return ""
	}

	var url string = serviceEndpoint

	if '/' != url[len(url)-1] {
		url += "/"
	}

	url += "xrpc/" + methodID

	return url
}
