package customerror

import "fmt"

// Error defines the struct to the api error
type Error struct {
	Code         int
	Scope        string
	ErrorMessage string
}

// Error handling error struct to string
func (e *Error) Error() string {
	return fmt.Sprintf("%s << %s", e.Scope, e.ErrorMessage)
}

// New create a new error with scope and the error message
func New(code int, scope, message string) error {
	return &Error{
		Code:         code,
		Scope:        scope,
		ErrorMessage: message,
	}
}

// Cast return struct error
func Cast(err error) *Error {
	return err.(*Error)
}
