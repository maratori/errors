//go:build go1.20

package errors

import (
	"errors"
)

func (e many) Error() string {
	return errors.Join(e.Unwrap()...).Error()
}
