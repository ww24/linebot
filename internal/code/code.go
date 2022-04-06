package code

import "errors"

type Code int

const (
	OK Code = iota
	NotFound
)

type internalError struct {
	source error
	code   Code
}

func (e *internalError) Error() string {
	return e.source.Error()
}

func (e *internalError) Unwrap() error {
	return e.source
}

func With(err error, code Code) error {
	return &internalError{source: err, code: code}
}

func From(err error) Code {
	ie := new(internalError)
	if errors.As(err, &ie) {
		return ie.code
	}
	return OK
}
