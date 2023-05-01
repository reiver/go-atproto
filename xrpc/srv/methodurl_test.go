package xrpcsrv_test

import (
	"github.com/reiver/go-atproto/xrpc/srv"

	"testing"
)

func TestMethodURL(t *testing.T) {

	tests := []struct{
		BaseURL  string
		MethodID string
		Expected string
	}{
		{
			BaseURL:  "https://bsky.social",
			MethodID:                          "com.atproto.server.createSession",
			Expected: "https://bsky.social/xrpc/com.atproto.server.createSession",
		},
		{
			BaseURL:  "https://bsky.social/",
			MethodID:                          "com.atproto.server.createSession",
			Expected: "https://bsky.social/xrpc/com.atproto.server.createSession",
		},



		{
			BaseURL:  "https://bsky.social",
			MethodID:                          "app.bsky.actor.getProfile",
			Expected: "https://bsky.social/xrpc/app.bsky.actor.getProfile",
		},
		{
			BaseURL:  "https://bsky.social/",
			MethodID:                          "app.bsky.actor.getProfile",
			Expected: "https://bsky.social/xrpc/app.bsky.actor.getProfile",
		},



		{
			BaseURL:  "https://example.com",
			MethodID:                          "once.twice.thrice.fource",
			Expected: "https://example.com/xrpc/once.twice.thrice.fource",
		},
		{
			BaseURL:  "https://example.com/",
			MethodID:                          "once.twice.thrice.fource",
			Expected: "https://example.com/xrpc/once.twice.thrice.fource",
		},
	}

	for testNumber, test := range tests {

		actual := xrpcsrv.MethodURL(test.BaseURL, test.MethodID)
		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual Method-URL is not what was expected.", testNumber)
			t.Logf("EXPECTED %q", expected)
			t.Logf("ACTUAL:  %q", actual)
			t.Logf("BASE-URL: %q", test.BaseURL)
			t.Logf("METHOD-ID: %q", test.MethodID)
			continue
		}
	}
}
