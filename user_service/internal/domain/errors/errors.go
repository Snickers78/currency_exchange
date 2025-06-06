package errs

import "errors"

var (
	UserAlreadyExists = errors.New("User already exists")
	WrongCredentials  = errors.New("Wrong email or password")
	BadRequest        = errors.New("Invalid email or password")
	InternalError     = errors.New("Internal error")
	UserNotFound      = errors.New("User not found")
)
