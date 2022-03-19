package util

type ValidateFunc func(s string) bool

func ValidLowerAlphabet(length int) ValidateFunc {
	return func(s string) bool {
		// check length
		if len(s) != length {
			return false
		}
		for _, r := range s {
			if r < 'a' || r > 'z' {
				return false
			}
		}
		return true
	}
}

func ValidAlphabet(length int) ValidateFunc {
	return func(s string) bool {
		if len(s) != length {
			return false
		}
		for _, r := range s {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
				return false
			}
		}
		return true
	}
}
