package validate

import "regexp"

func PasswordValid(password string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\/-]{5,}$`)
	isValid := regex.MatchString(password)
	return isValid
}
