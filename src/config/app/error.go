package app

import "fmt"

type AppError struct {
	Base  error // The business error (e.g., UserNotFound)
	Cause error // The technical error + stack trace (e.g., sql: no rows)
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v: %v", e.Base, e.Cause)
	}
	return e.Base.Error()
}

func Reject(base error, cause error) error {
	return &AppError{Base: base, Cause: cause}
}
