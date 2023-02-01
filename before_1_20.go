//go:build !go1.20

package errors

import (
	"errors"
)

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
