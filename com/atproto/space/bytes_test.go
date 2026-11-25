package space_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/reiver/go-atproto/com/atproto/space"
)

func TestBytes_MarshalJSON(t *testing.T) {
	tests := []struct {
		Bytes    space.Bytes
		Expected string
	}{
		{Bytes: nil, Expected: `{"$bytes":""}`},
		{Bytes: space.Bytes{}, Expected: `{"$bytes":""}`},
		{Bytes: space.Bytes("Hello, world!"), Expected: `{"$bytes":"SGVsbG8sIHdvcmxkIQ"}`},
		{Bytes: space.Bytes{0xfb, 0xff}, Expected: `{"$bytes":"+/8"}`},
	}

	for _, test := range tests {
		actual, err := json.Marshal(test.Bytes)
		if nil != err {
			t.Errorf("For %x, did not expect an error but actually got one: %s", []byte(test.Bytes), err)
			continue
		}
		if test.Expected != string(actual) {
			t.Errorf("For %x, the actual JSON is not what was expected.\nEXPECTED: %s\nACTUAL:   %s", []byte(test.Bytes), test.Expected, actual)
		}
	}
}

func TestBytes_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		JSON     string
		Expected string
	}{
		{JSON: `{"$bytes":"SGVsbG8sIHdvcmxkIQ"}`, Expected: "Hello, world!"},
		{JSON: `{"$bytes":"SGVsbG8sIHdvcmxkIQ=="}`, Expected: "Hello, world!"},
		{JSON: `{"$bytes":"+/8"}`, Expected: "\xfb\xff"},
		{JSON: `{"$bytes":""}`, Expected: ""},
	}

	for _, test := range tests {
		var actual space.Bytes
		if err := json.Unmarshal([]byte(test.JSON), &actual); nil != err {
			t.Errorf("For %s, did not expect an error but actually got one: %s", test.JSON, err)
			continue
		}
		if test.Expected != string(actual) {
			t.Errorf("For %s, expected %q but actually got %q.", test.JSON, test.Expected, actual)
		}
	}
}

func TestBytes_UnmarshalJSON_bad(t *testing.T) {
	tests := []string{
		`"SGVsbG8"`,
		`null`,
		`{}`,
		`{"$bytes":5}`,
		`{"$bytes":"SGVsbG8","extra":true}`,
		`{"$bytes":"not base64!"}`,
		`{"$bytes":"-_8"}`, // base64url, not base64
	}

	for _, test := range tests {
		var actual space.Bytes
		err := json.Unmarshal([]byte(test), &actual)
		if !errors.Is(err, space.ErrBadBytes) {
			t.Errorf("For %s, expected an error wrapping %q but actually got: %v", test, space.ErrBadBytes, err)
		}
	}
}
