package validate

import "regexp"

func EmailValid(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9][a-zA-Z0-9\-]*\.[a-zA-Z]{2,}$`)
	isValid := regex.MatchString(email)
	return isValid
}
