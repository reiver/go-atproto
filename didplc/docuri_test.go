package didplc_test

import (
	"github.com/reiver/go-atproto/internal/didplc"

	"testing"
)

func TestDocURI(t *testing.T) {

	tests := []struct{
		DID string
		Expected string
	}{
		{
			DID: "",
			Expected: "",
		},



		{
			DID: "apple",
			Expected: "",
		},
		{
			DID: "banana",
			Expected: "",
		},
		{
			DID: "cherry",
			Expected: "",
		},



		{
			DID: "did:once:one",
			Expected: "",
		},
		{
			DID: "did:twice:two",
			Expected: "",
		},
		{
			DID: "did:thrice:three",
			Expected: "",
		},
		{
			DID: "did:fource:four",
			Expected: "",
		},



		{
			DID: "did:plc:",
			Expected: "https://plc.directory/did:plc:",
		},



		{
			DID: "did:plc:abcde12345",
			Expected: "https://plc.directory/did:plc:abcde12345",
		},



		{
			DID: "did:plc:once",
			Expected: "https://plc.directory/did:plc:once",
		},
		{
			DID: "did:plc:twice",
			Expected: "https://plc.directory/did:plc:twice",
		},
		{
			DID: "did:plc:thrice",
			Expected: "https://plc.directory/did:plc:thrice",
		},
		{
			DID: "did:plc:fource",
			Expected: "https://plc.directory/did:plc:fource",
		},
	}

	for testNumber, test := range tests {

		actual := didplc.DocURI(test.DID)
		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual DID Placeholder document URI was not what was expected.", testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			continue
		}
	}
}
