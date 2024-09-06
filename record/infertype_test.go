package record_test

import (
	"testing"

        "github.com/reiver/go-atproto/record"
)

func TestInferType(t *testing.T) {

	tests := []struct{
		Src any
		Expected string
	}{
		{
			Src: map[string]string{
				"$type":"once.twice.thrice.fource",
			},
			Expected: "once.twice.thrice.fource",
		},
		{
			Src: map[string]any{
				"$type":"once.twice.thrice.fource",
			},
			Expected: "once.twice.thrice.fource",
		},
		{
			Src: map[any]string{
				"$type":"once.twice.thrice.fource",
			},
			Expected: "once.twice.thrice.fource",
		},
		{
			Src: map[any]any{
				"$type":"once.twice.thrice.fource",
			},
			Expected: "once.twice.thrice.fource",
		},



		{
			Src: map[string]string{
				"$type":"once.twice.thrice.fource",
				"apple":"1",
				"banana":"TWO",
				"cherry":"3",
			},
			Expected: "once.twice.thrice.fource",
		},
		{
			Src: map[string]any{
				"$type":"once.twice.thrice.fource",
				"apple":1,
				"banana":"TWO",
				"cherry":"3",
			},
			Expected: "once.twice.thrice.fource",
		},
		{
			Src: map[any]string{
				"$type":"once.twice.thrice.fource",
				"apple":"1",
				"banana":"TWO",
				"cherry":"3",
			},
			Expected: "once.twice.thrice.fource",
		},
		{
			Src: map[any]any{
				"$type":"once.twice.thrice.fource",
				"apple":1,
				"banana":"TWO",
				"cherry":"3",
			},
			Expected: "once.twice.thrice.fource",
		},
	}

	for testNumber, test := range tests {

		actual, err := record.InferType(test.Src)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error but actually got one.", testNumber)
			t.Logf("ERROR: (%T) %s", err, err)
			t.Logf("SRC: (%T) %#v", test.Src, test.Src)
			continue
		}

		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual $type is not what was expected.", testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			t.Logf("SRC: (%T) %#v", test.Src, test.Src)
			continue
		}
	}
}
