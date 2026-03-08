package http

import "fmt"

var (
	UserNotFound      = fmt.Errorf("User not found")
	InvalidEmail      = fmt.Errorf("Invalid email")
	InvalidName       = fmt.Errorf("Invalid name: cannot be empty")
	InvalidPassword   = fmt.Errorf("Invalid password: cannot be empty")
	UserAlreadyExists = fmt.Errorf("User already exists")

	AuthEmailNotFound      = fmt.Errorf("Email not found")
	AuthInvalidPassword    = fmt.Errorf("Invalid password")
	AuthEmailAlreadyExists = fmt.Errorf("Email already exists")

	JWTFailedGenerateToken = fmt.Errorf("Token generation failed")
)
