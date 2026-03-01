package users

import "fmt"

// Domain errors
var (
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrInvalidEmail      = fmt.Errorf("invalid email")
	ErrInvalidName       = fmt.Errorf("invalid name: cannot be empty")
	ErrInvalidPassword   = fmt.Errorf("invalid password: cannot be empty")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
)
