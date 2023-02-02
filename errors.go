// Package errors provides functions to construct errors with custom fields for structured logging.
package errors

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
	"golang.org/x/exp/maps"
)

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

type Fields = map[string]any

func Errors(err error) []error {
	errs := wrapper{err: err}.Errors()
	res := make([]error, 0, len(errs))
	for _, e := range errs {
		res = append(res, e)
	}
	return res
}

func FieldsFromError(err error) Fields {
	errs := wrapper{err: err}.Errors()
	if len(errs) == 0 || errs[0].fields == nil {
		// It may be handy to update map returned by FieldsFromError.
		// Need to return empty map instead of nil.
		return Fields{}
	}
	return maps.Clone(errs[0].fields)
}

type ErrorBuilder struct {
	err treeNode
}

func New(msg string) *ErrorBuilder {
	return Err(errors.New(msg))
}

func Err(err error) *ErrorBuilder {
	// Type switch instead of errors.As() because we don't want to extract wrapped error to not miss wrapper.
	switch e := err.(type) { //nolint:errorlint // see comment above
	case nil:
		return nil
	case treeNode:
		return &ErrorBuilder{
			err: e,
		}
	default:
		return &ErrorBuilder{
			err: wrapper{
				err: err,
			},
		}
	}
}

func (e *ErrorBuilder) E() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *ErrorBuilder) Wrap(prefix string) *ErrorBuilder {
	if e == nil {
		return nil
	}
	e.err = withPrefix{
		err:    e.err,
		prefix: prefix,
	}
	return e
}

func (e *ErrorBuilder) WithFields(fields Fields) *ErrorBuilder {
	if e == nil {
		return nil
	}
	e.err = withFields{
		err:    e.err,
		fields: maps.Clone(fields),
	}
	return e
}

func (e *ErrorBuilder) WithField(key string, value any) *ErrorBuilder {
	return e.WithFields(Fields{
		key: value,
	})
}

func Wrap(prefix string, err error) *ErrorBuilder {
	return Err(err).Wrap(prefix)
}

func WithFields(err error, fields Fields) *ErrorBuilder {
	return Err(err).WithFields(fields)
}

func WithField(err error, key string, value any) *ErrorBuilder {
	return Err(err).WithField(key, value)
}

func Combine(errs ...error) error {
	converted := make([]treeNode, 0, len(errs))
	for _, err := range errs {
		// Type switch instead of errors.As() because we don't want to extract wrapped error to not miss wrapper.
		switch e := err.(type) { //nolint:errorlint // see comment above
		case nil:
			continue
		case many:
			converted = append(converted, e.errors...)
		case treeNode:
			converted = append(converted, e)
		default:
			converted = append(converted, wrapper{
				err: e,
			})
		}
	}
	if len(converted) == 0 {
		return nil
	}
	if len(converted) == 1 {
		return converted[0]
	}
	return many{
		errors: converted,
	}
}

func AppendInto(into *error, err error) {
	if into == nil {
		panic("misuse of errors.AppendInto: into pointer must not be nil")
	}
	if err == nil {
		return
	}
	*into = Combine(*into, err)
}

func joinFields(outer Fields, inner Fields) Fields {
	switch {
	case len(outer) == 0:
		return inner
	case len(inner) == 0:
		return outer
	}
	res := make(Fields, len(outer)+len(inner))
	for k, v := range outer {
		res[k] = v
	}
	for k, v := range inner {
		res[k] = v // inner map has higher priority for duplicated fields
	}
	return res
}

type errorWithFields struct {
	err    error
	fields Fields
}

func (e errorWithFields) Error() string {
	return e.err.Error()
}

func (e errorWithFields) Unwrap() error {
	return e.err
}

type treeNode interface {
	isMyError()
	error
	Errors() []errorWithFields
	// Optional methods:
	//   Unwrap() error
	//   Unwrap() []error
}

//nolint:exhaustruct // false positive
var (
	_ treeNode = wrapper{}
	_ treeNode = withPrefix{}
	_ treeNode = withFields{}
	_ treeNode = many{}
)

func (e wrapper) isMyError()    {}
func (e withPrefix) isMyError() {}
func (e withFields) isMyError() {}
func (e many) isMyError()       {}

type wrapper struct {
	err error
}

func (e wrapper) Errors() []errorWithFields {
	// Type switch instead of errors.As() because we don't want to extract wrapped error to not miss wrapper.
	switch err := e.err.(type) { //nolint:errorlint // see comment above
	case nil:
		return nil
	case treeNode:
		return err.Errors()
	case errorWithFields:
		return []errorWithFields{err}
	default:
		return []errorWithFields{{
			err:    err,
			fields: nil,
		}}
	}
}

func (e wrapper) Error() string {
	return e.err.Error()
}

func (e wrapper) Unwrap() error {
	return e.err
}

type withPrefix struct {
	err    treeNode
	prefix string
}

func (e withPrefix) Errors() []errorWithFields {
	errs := e.err.Errors()
	res := make([]errorWithFields, 0, len(errs))
	for _, err := range errs {
		res = append(res, errorWithFields{
			err:    fmt.Errorf("%s: %w", e.prefix, err.err),
			fields: err.fields,
		})
	}
	return res
}

func (e withPrefix) Error() string {
	return fmt.Sprintf("%s: %s", e.prefix, e.err.Error())
}

func (e withPrefix) Unwrap() error {
	return e.err
}

type withFields struct {
	err    treeNode
	fields Fields
}

func (e withFields) Errors() []errorWithFields {
	errs := e.err.Errors()
	res := make([]errorWithFields, 0, len(errs))
	for _, err := range errs {
		res = append(res, errorWithFields{
			err:    err.err,
			fields: joinFields(e.fields, err.fields),
		})
	}
	return res
}

func (e withFields) Error() string {
	return e.err.Error()
}

func (e withFields) Unwrap() error {
	return e.err
}

type many struct {
	errors []treeNode
}

func (e many) Errors() []errorWithFields {
	res := make([]errorWithFields, 0, len(e.errors))
	for _, err := range e.errors {
		res = append(res, err.Errors()...)
	}
	return res
}

func (e many) Error() string {
	return multierr.Combine(e.Unwrap()...).Error()
}

func (e many) Unwrap() []error {
	res := make([]error, 0, len(e.errors))
	for _, err := range e.errors {
		res = append(res, err)
	}
	return res
}
