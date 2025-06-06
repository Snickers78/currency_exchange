package validate

import (
	"testing"
)

func Test_PasswordValid_ShouldReturnFalse_WhenPasswordIsEmptyOrHasLessThan8Symbols(t *testing.T) {
	testCases := []struct {
		password string
		expected bool
	}{
		{"", false},
		{"12345", false},
		{"fffghhs", false},
		{"qwerty1hh", true},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			result := PasswordValid(tc.password)

			if result != tc.expected {
				t.Errorf("Password: %s = %v, want %v", tc.password, result, tc.expected)
			}
		})
	}
}
