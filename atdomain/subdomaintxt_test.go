package atdomain

import (
	"testing"
)

func TestSubdomainTXT(t *testing.T) {

	tests := []struct{
		DomainName string
		Expected   string
	}{
		{
			DomainName:        "example.com",
			Expected: "_atproto.example.com",
		},
		{
			DomainName:        "changelog.ca",
			Expected: "_atproto.changelog.ca",
		},
		{
			DomainName:        "once.twice.thrice.fource.xyz",
			Expected: "_atproto.once.twice.thrice.fource.xyz",
		},
	}

	for testNumber, test := range tests {

		actual := subdomainTXT(test.DomainName)
		expected := test.Expected

		if expected != actual {
			t.Errorf("For test %d, the actual at-protocol DNS TXT sub-domain is not what was expected.", testNumber)
			t.Logf("EXPECTED %q: ", expected)
			t.Logf("ACTUAL   %q: ", actual)
			continue
		}
	}
}
