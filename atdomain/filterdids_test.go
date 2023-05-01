package atdomain

import (
	"reflect"

	"testing"
)

func TestFilterDIDs(t *testing.T) {

	tests := []struct{
		Responses []string
		Expected []string
	}{
		{
			Responses: []string{},
			Expected : []string(nil),
		},



		{
			Responses: []string{
				"once=apple",
			},
			Expected : []string(nil),
		},
		{
			Responses: []string{
				"once=apple",
				"twice=banana",
			},
			Expected : []string(nil),
		},
		{
			Responses: []string{
				"once=apple",
				"twice=banana",
				"thrice=cherry",
			},
			Expected : []string(nil),
		},
		{
			Responses: []string{
				"once=apple",
				"twice=banana",
				"thrice=cherry",
				"fource=date",
			},
			Expected : []string(nil),
		},



		{
			Responses: []string{
				"did=did:plc:abcde12345",
			},
			Expected : []string{
				"did:plc:abcde12345",
			},
		},
		{
			Responses: []string{
				"did=did:plc:abcde12345",
				"did=did:plc:wxyz6789",
			},
			Expected : []string{
				"did:plc:abcde12345",
				"did:plc:wxyz6789",
			},
		},



		{
			Responses: []string{
				"wow=ok",
				"did=did:plc:abcde12345",
			},
			Expected : []string{
				"did:plc:abcde12345",
			},
		},
		{
			Responses: []string{
				"did=did:plc:abcde12345",
				"wow=ok",
			},
			Expected : []string{
				"did:plc:abcde12345",
			},
		},
		{
			Responses: []string{
				"wow=ok",
				"did=did:plc:abcde12345",
				"did=did:plc:wxyz6789",
			},
			Expected : []string{
				"did:plc:abcde12345",
				"did:plc:wxyz6789",
			},
		},
		{
			Responses: []string{
				"did=did:plc:abcde12345",
				"wow=ok",
				"did=did:plc:wxyz6789",
			},
			Expected : []string{
				"did:plc:abcde12345",
				"did:plc:wxyz6789",
			},
		},
		{
			Responses: []string{
				"did=did:plc:abcde12345",
				"did=did:plc:wxyz6789",
				"wow=ok",
			},
			Expected : []string{
				"did:plc:abcde12345",
				"did:plc:wxyz6789",
			},
		},
	}

	for testNumber, test := range tests {

		actual := filterDIDs(test.Responses)
		expected := test.Expected

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("For test #%d, the actual filtered DIDs is not what was expected.", testNumber)
			t.Logf("EXPECTED: %#v", expected)
			t.Logf("ACTUAL:   %#v", actual)
			continue
		}
	}
}
