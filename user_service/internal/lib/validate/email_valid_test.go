package validate

import (
	"testing"
)

func Test_EmailValid_ShouldReturnFalse_WhenEmailIsInvalid(t *testing.T) {
	testCases := []struct {
		email    string
		expected bool
	}{

		{"user@example.com", true},
		{"john.doe@gmail.com", true},
		{"user+tag@example.co.uk", false},
		{"firstname.lastname@domain.com", true},

		{"", false},
		{"invalid-email", false},
		{"user@domain", false},
		{"@example.com", false},
		{"user@.com", false},
		{"user@domain.", false},

		{"a@b.co", true},
		{"very.long.email.address@very.long.domain.com", false},

		{"user123@example.com", true},
		{"user-name@example.com", true},
		{"user_name@example.com", true},
	}

	for _, tc := range testCases {
		t.Run(tc.email, func(t *testing.T) {
			result := EmailValid(tc.email)

			if result != tc.expected {
				t.Errorf("Email: %s, Expected: %v, Got: %v", tc.email, tc.expected, result)
			}
		})
	}
}
