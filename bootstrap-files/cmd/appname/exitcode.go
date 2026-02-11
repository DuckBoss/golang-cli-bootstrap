package main

import (
	"fmt"
)

// exitCodeError is an error that carries an exit code for the process.
type exitCodeError struct {
	err  error
	code int
}

func (e *exitCodeError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return fmt.Sprintf("exit code %d", e.code)
}

func (e *exitCodeError) ExitCode() int {
	return e.code
}

// ExitCode wraps err with an exit code. If err is nil, ExitCode returns nil.
func ExitCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &exitCodeError{err: err, code: code}
}

// GetExitCode returns the exit code from err if it implements ExitCode() int,
// otherwise returns ExitCodeError (1).
func GetExitCode(err error) int {
	type exitCoder interface {
		ExitCode() int
	}
	if err == nil {
		return ExitCodeSuccess
	}
	if e, ok := err.(exitCoder); ok {
		return e.ExitCode()
	}
	return ExitCodeError
}
