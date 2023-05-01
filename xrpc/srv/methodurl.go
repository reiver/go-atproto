package xrpcsrv

func MethodURL(baseurl string, methodID string) string {
	if "" == baseurl {
		return ""
	}

	var url string = baseurl

	if '/' != url[len(url)-1] {
		url += "/"
	}

	url += "xrpc/" + methodID

	return url
}
