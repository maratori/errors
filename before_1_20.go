//go:build !go1.20

package errors

import (
	"errors"
	"strings"
)

func (e many) Error() string {
	errs := e.Unwrap()
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "\n")
}

func (e many) Is(target error) bool { // need to implement because multi-error is not supported before go1.20
	for _, err := range e.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e many) As(target any) bool { // need to implement because multi-error is not supported before go1.20
	for _, err := range e.errors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}
