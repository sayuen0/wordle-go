package util

import "testing"

func TestValidateFunc(t *testing.T) {
	test := []struct {
		validateFunc ValidateFunc
		s            string
		want         bool
	}{
		{validateFunc: ValidLowerAlphabet(5), s: "water", want: true},
		{validateFunc: ValidLowerAlphabet(5), s: "orange", want: false},
	}
	for _, tt := range test {
		if got := tt.validateFunc(tt.s); got != tt.want {
			t.Errorf("want %v, got %v", got, tt.want)
		}
	}
}
