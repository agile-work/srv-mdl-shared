package shared

import "fmt"

// ErrorModule defines the struct to the api error
type ErrorModule struct {
	Scope        string
	ErrorMessage string
}

// Error handling error struct to string
func (e *ErrorModule) Error() string {
	return fmt.Sprintf("%s << %s", e.Scope, e.ErrorMessage)
}

// NewError create a new error to the api
func NewError(scope, err string) error {
	return &ErrorModule{
		Scope:        scope,
		ErrorMessage: err,
	}
}

// GetErrorStruct return struct error
func GetErrorStruct(err error) *ErrorModule {
	return err.(*ErrorModule)
}
