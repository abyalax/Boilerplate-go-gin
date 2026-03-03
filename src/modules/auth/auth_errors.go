package auth

import "fmt"

var (
	ErrEmailNotFound      = fmt.Errorf("email not found")
	ErrInvalidEmail       = fmt.Errorf("invalid email")
	ErrEmailAlreadyExists = fmt.Errorf("email already exists")

	ErrInvalidPassword = fmt.Errorf("invalid password; cannot be empty")
	ErrPasswordToShort = fmt.Errorf("password minimal 6 character")
)
