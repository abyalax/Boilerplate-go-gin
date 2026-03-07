package reject

import "fmt"

var (
	UserNotFound      = fmt.Errorf("user not found")
	InvalidEmail      = fmt.Errorf("invalid email")
	InvalidName       = fmt.Errorf("invalid name: cannot be empty")
	InvalidPassword   = fmt.Errorf("invalid password: cannot be empty")
	UserAlreadyExists = fmt.Errorf("user already exists")

	AuthEmailNotFound      = fmt.Errorf("email not found")
	AuthInvalidPassword    = fmt.Errorf("invalid password")
	AuthEmailAlreadyExists = fmt.Errorf("email already exists")
)
