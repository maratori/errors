package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maratori/errors"
)

func TestOneError(t *testing.T) {
	t.Parallel()
	const (
		newErr = "new err"
		prefix = "prefix"
		key    = "key"
		value  = "value"
		key2   = "key2"
		value2 = "value2"
	)

	t.Run("new error without fields", func(t *testing.T) {
		t.Parallel()
		err := errors.New(newErr).E()
		require.EqualError(t, err, newErr)

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], newErr)
	})

	t.Run("new error with field", func(t *testing.T) {
		t.Parallel()
		err := errors.New(newErr).WithField(key, value).E()
		require.EqualError(t, err, newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value}, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], newErr)
	})

	t.Run("non-nil wrapped error without fields", func(t *testing.T) {
		t.Parallel()
		err := errors.Wrap(prefix, stderrors.New(newErr)).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix+": "+newErr)
	})

	t.Run("non-nil wrapped error with one field", func(t *testing.T) {
		t.Parallel()
		err := errors.Wrap(prefix, stderrors.New(newErr)).WithField(key, value).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value}, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix+": "+newErr)
	})

	t.Run("non-nil wrapped error with two fields", func(t *testing.T) {
		t.Parallel()
		err := errors.Wrap(prefix, stderrors.New(newErr)).WithField(key, value).WithField(key2, value2).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value, key2: value2}, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix+": "+newErr)
	})

	t.Run("first field has priority", func(t *testing.T) {
		t.Parallel()
		err := errors.Wrap(prefix, stderrors.New(newErr)).WithField(key, value).WithField(key, value2).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value}, fields)
	})

	t.Run("nil wrapped error without fields", func(t *testing.T) {
		t.Parallel()
		err := errors.Wrap(prefix, nil).E()
		require.NoError(t, err)

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Empty(t, errs)
	})

	t.Run("nil wrapped error with fields", func(t *testing.T) {
		t.Parallel()
		err := errors.Wrap(prefix, nil).WithField(key, value).E()
		require.NoError(t, err)

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Empty(t, errs)
	})

	t.Run("errors.Is(): non-nil wrapped error with fields", func(t *testing.T) {
		t.Parallel()
		err := stderrors.New(newErr)
		wrapped := errors.Wrap(prefix, err).WithField(key, value).E()

		require.ErrorIs(t, wrapped, err)
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()
	const (
		newErr  = "new err"
		prefix  = "prefix"
		prefix2 = "prefix2"
		key     = "key"
		value   = "value"
		key2    = "key2"
		value2  = "value2"
	)

	t.Run("wrap error without fields", func(t *testing.T) {
		t.Parallel()
		original := errors.New(newErr).E()
		err := errors.Wrap(prefix, original).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix+": "+newErr)
	})

	t.Run("wrap error without fields and add field", func(t *testing.T) {
		t.Parallel()
		original := errors.Wrap(prefix, stderrors.New(newErr)).E()
		err := errors.Wrap(prefix2, original).WithField(key, value).E()
		require.EqualError(t, err, prefix2+": "+prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value}, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix2+": "+prefix+": "+newErr)
	})

	t.Run("wrap error with one field", func(t *testing.T) {
		t.Parallel()
		original := errors.New(newErr).WithField(key, value).E()
		err := errors.Wrap(prefix, original).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value}, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix+": "+newErr)
	})

	t.Run("wrap error with one field and add new field", func(t *testing.T) {
		t.Parallel()
		original := errors.New(newErr).WithField(key, value).E()
		err := errors.Wrap(prefix, original).WithField(key2, value2).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value, key2: value2}, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix+": "+newErr)
	})

	t.Run("inner field has priority", func(t *testing.T) {
		t.Parallel()
		original := errors.New(newErr).WithField(key, value).E()
		err := errors.Wrap(prefix, original).WithField(key, value2).E()
		require.EqualError(t, err, prefix+": "+newErr)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key: value}, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], prefix+": "+newErr)
	})

	t.Run("fields are not accessible if error is wrapped with fmt", func(t *testing.T) {
		t.Parallel()
		errWithField := errors.New(newErr).WithField(key, value).E()
		wrapped := fmt.Errorf("%w", errWithField)
		fields := errors.FieldsFromError(wrapped)
		require.NotNil(t, fields)
		require.Empty(t, fields)
	})
}

func TestMultipleErrors(t *testing.T) {
	t.Parallel()
	const (
		newErr1 = "new err 1"
		newErr2 = "new err 2"
		newErr3 = "new err 3"
		prefix1 = "prefix1"
		prefix2 = "prefix2"
		key1    = "key1"
		value1  = "value1"
		key2    = "key2"
		value2  = "value2"
		key3    = "key3"
		value3  = "value3"
		key4    = "key4"
		value4  = "value4"
	)

	t.Run("join no errors", func(t *testing.T) {
		t.Parallel()
		err := errors.Join()
		require.NoError(t, err)

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Empty(t, errs)
	})

	t.Run("join returns the single non-nil error", func(t *testing.T) {
		t.Parallel()
		original := stderrors.New(newErr1)
		err := errors.Join(nil, original, nil)
		require.EqualError(t, err, newErr1)
		require.ErrorIs(t, err, original)
		require.NotEqual(t, err, original) // because it's wrapped
		require.True(t, err != original)   //nolint:errorlint,testifylint // use != to validate they are not the same

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 1)
		require.EqualError(t, errs[0], newErr1)
	})

	t.Run("join two errors", func(t *testing.T) {
		t.Parallel()
		original1 := stderrors.New(newErr1)
		original2 := stderrors.New(newErr2)
		err := errors.Join(nil, original1, nil, original2, nil)
		require.EqualError(t, err, newErr1+"\n"+newErr2)
		require.ErrorIs(t, err, original1)
		require.ErrorIs(t, err, original2)

		fields := errors.FieldsFromError(err)
		require.NotNil(t, fields)
		require.Empty(t, fields)

		errs := errors.Errors(err)
		require.Len(t, errs, 2)
		require.EqualError(t, errs[0], newErr1)
		require.EqualError(t, errs[1], newErr2)
	})

	t.Run("join two errors with fields", func(t *testing.T) {
		t.Parallel()
		original1 := errors.New(newErr1).WithField(key1, value1).E()
		original2 := errors.New(newErr2).WithField(key2, value2).E()
		err := errors.Join(nil, original1, nil, original2, nil)
		require.EqualError(t, err, newErr1+"\n"+newErr2)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key1: value1}, fields) // only from the first error

		errs := errors.Errors(err)
		require.Len(t, errs, 2)
		require.EqualError(t, errs[0], newErr1)
		require.EqualError(t, errs[1], newErr2)

		fields = errors.FieldsFromError(errs[0])
		require.Equal(t, errors.Fields{key1: value1}, fields)

		fields = errors.FieldsFromError(errs[1])
		require.Equal(t, errors.Fields{key2: value2}, fields)
	})

	t.Run("append into two errors with fields", func(t *testing.T) {
		t.Parallel()
		err1 := errors.New(newErr1).WithField(key1, value1).E()
		err2 := errors.New(newErr2).WithField(key2, value2).E()
		errors.AppendInto(&err1, err2)
		require.EqualError(t, err1, newErr1+"\n"+newErr2)

		fields := errors.FieldsFromError(err1)
		require.Equal(t, errors.Fields{key1: value1}, fields) // only from the first error

		errs := errors.Errors(err1)
		require.Len(t, errs, 2)
		require.EqualError(t, errs[0], newErr1)
		require.EqualError(t, errs[1], newErr2)

		fields = errors.FieldsFromError(errs[0])
		require.Equal(t, errors.Fields{key1: value1}, fields)

		fields = errors.FieldsFromError(errs[1])
		require.Equal(t, errors.Fields{key2: value2}, fields)
	})

	t.Run("join wrapped errors", func(t *testing.T) {
		t.Parallel()
		original1 := errors.New(newErr1).WithField(key1, value1).E()
		original2 := errors.New(newErr2).WithField(key2, value2).E()
		wrapped1 := errors.Wrap(prefix1, original1).WithField(key3, value3).E()
		wrapped2 := errors.Wrap(prefix2, original2).WithField(key4, value4).E()
		err := errors.Join(nil, wrapped1, nil, wrapped2, nil)
		require.EqualError(t, err, prefix1+": "+newErr1+"\n"+prefix2+": "+newErr2)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key1: value1, key3: value3}, fields) // only from the first chain

		errs := errors.Errors(err)
		require.Len(t, errs, 2)
		require.EqualError(t, errs[0], prefix1+": "+newErr1)
		require.EqualError(t, errs[1], prefix2+": "+newErr2)

		fields = errors.FieldsFromError(errs[0])
		require.Equal(t, errors.Fields{key1: value1, key3: value3}, fields)

		fields = errors.FieldsFromError(errs[1])
		require.Equal(t, errors.Fields{key2: value2, key4: value4}, fields)
	})

	t.Run("wrap joined errors", func(t *testing.T) {
		t.Parallel()
		original1 := errors.New(newErr1).WithField(key1, value1).E()
		original2 := errors.New(newErr2).WithField(key2, value2).E()
		joined := errors.Join(nil, original1, nil, original2, nil)
		err := errors.Wrap(prefix1, joined).WithField(key3, value3).E()
		require.EqualError(t, err, prefix1+": "+newErr1+"\n"+newErr2)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key1: value1, key3: value3}, fields) // only from the first chain

		errs := errors.Errors(err)
		require.Len(t, errs, 2)
		require.EqualError(t, errs[0], prefix1+": "+newErr1)
		require.EqualError(t, errs[1], prefix1+": "+newErr2)

		fields = errors.FieldsFromError(errs[0])
		require.Equal(t, errors.Fields{key1: value1, key3: value3}, fields)

		fields = errors.FieldsFromError(errs[1])
		require.Equal(t, errors.Fields{key2: value2, key3: value3}, fields)
	})

	t.Run("join several times", func(t *testing.T) {
		t.Parallel()
		original1 := errors.New(newErr1).WithField(key1, value1).E()
		original2 := errors.New(newErr2).WithField(key2, value2).E()
		original3 := errors.New(newErr3).WithField(key3, value3).E()
		joined1 := errors.Join(nil, original1, nil, original2, nil)
		joined2 := errors.Join(original3, joined1)
		err := errors.Wrap(prefix1, joined2).WithField(key4, value4).E()
		require.EqualError(t, err, prefix1+": "+newErr3+"\n"+newErr1+"\n"+newErr2)

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key3: value3, key4: value4}, fields) // only from the first chain

		errs := errors.Errors(err)
		require.Len(t, errs, 3)
		require.EqualError(t, errs[0], prefix1+": "+newErr3)
		require.EqualError(t, errs[1], prefix1+": "+newErr1)
		require.EqualError(t, errs[2], prefix1+": "+newErr2)

		fields = errors.FieldsFromError(errs[0])
		require.Equal(t, errors.Fields{key3: value3, key4: value4}, fields)

		fields = errors.FieldsFromError(errs[1])
		require.Equal(t, errors.Fields{key1: value1, key4: value4}, fields)

		fields = errors.FieldsFromError(errs[2])
		require.Equal(t, errors.Fields{key2: value2, key4: value4}, fields)
	})

	t.Run("inner fields may repeat", func(t *testing.T) {
		t.Parallel()
		original1 := errors.New(newErr1).WithField(key1, value1).E()
		original2 := errors.New(newErr2).WithField(key1, value2).E()
		original3 := errors.New(newErr3).WithField(key1, value3).E()
		joined1 := errors.Join(nil, original1, nil, original2, nil)
		joined2 := errors.Join(original3, joined1)
		err := errors.Wrap(prefix1, joined2).WithField(key2, value4).E()

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key1: value3, key2: value4}, fields) // only from the first chain

		errs := errors.Errors(err)
		require.Len(t, errs, 3)

		fields = errors.FieldsFromError(errs[0])
		require.Equal(t, errors.Fields{key1: value3, key2: value4}, fields)

		fields = errors.FieldsFromError(errs[1])
		require.Equal(t, errors.Fields{key1: value1, key2: value4}, fields)

		fields = errors.FieldsFromError(errs[2])
		require.Equal(t, errors.Fields{key1: value2, key2: value4}, fields)
	})

	t.Run("inner field has priority", func(t *testing.T) {
		t.Parallel()
		original1 := errors.New(newErr1).WithField(key1, value1).E()
		original2 := errors.New(newErr2).WithField(key2, value2).E()
		original3 := errors.New(newErr3).WithField(key1, value3).E()
		joined1 := errors.Join(nil, original1, nil, original2, nil)
		joined2 := errors.Join(original3, joined1)
		err := errors.Wrap(prefix1, joined2).WithField(key1, value4).E()

		fields := errors.FieldsFromError(err)
		require.Equal(t, errors.Fields{key1: value3}, fields) // only from the first chain

		errs := errors.Errors(err)
		require.Len(t, errs, 3)

		fields = errors.FieldsFromError(errs[0])
		require.Equal(t, errors.Fields{key1: value3}, fields)

		fields = errors.FieldsFromError(errs[1])
		require.Equal(t, errors.Fields{key1: value1}, fields)

		fields = errors.FieldsFromError(errs[2])
		require.Equal(t, errors.Fields{key1: value4, key2: value2}, fields)
	})
}
