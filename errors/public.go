package errors

import (
	"errors"
)

// Public wraps the original error with a new error that has a
// `Public() string` method that will return a message that is
// acceptable to display to the public. This error can also be
// unwrapped using the traditional `errors` package approach.
func Public(err error, msg string, code int) error {
	return ResponseCodeError{err, msg, code}
}

var (
	As = errors.As
	Is = errors.Is
)

type ResponseCodeError struct {
	err  error
	msg  string
	code int
}

func (pe ResponseCodeError) Error() string {
	return pe.err.Error()
}
func (pe ResponseCodeError) Public() string {
	return pe.msg
}
func (pe ResponseCodeError) Unwrap() error {
	return pe.err
}
