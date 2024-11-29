package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	res := fmt.Sprintf("%d errors occured:\n", len(e.errs))
	for _, err := range e.errs {
		res += "\t* " + err.Error()
	}
	res += "\n"
	return res
}

func (e *MultiError) Unwrap() error {
	if len(e.errs) == 0 {
		return nil
	}
	return &MultiError{errs: e.errs[1:]}
}

func (e *MultiError) Is(target error) bool {
	for _, err := range e.errs {
		if err.Error() == target.Error() {
			return true
		}
	}
	return false
}

func (e *MultiError) As(target any) bool {
	for _, err := range e.errs {
		if reflect.TypeOf(err) == reflect.TypeOf(target).Elem() {
			return true
		}
	}
	return false
}

func Append(err error, errs ...error) *MultiError {
	switch err := err.(type) {
	case *MultiError:
		for _, e := range errs {
			switch e := e.(type) {
			case *MultiError:
				err.errs = append(err.errs, e.errs...)
			default:
				err.errs = append(err.errs, e)
			}
		}
		return err
	default:
		newErrs := make([]error, 0, len(errs)+1)
		if err != nil {
			newErrs = append(newErrs, err)
		}
		newErrs = append(newErrs, errs...)
		return Append(&MultiError{}, newErrs...)
	}
}

type RandomError struct {
	code int
}

func (re RandomError) Error() string {
	return fmt.Sprintf("error with code %d", re.code)
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
	expectedMessage = "1 errors occured:\n\t* error 2\n"
	assert.EqualError(t, errors.Unwrap(err), expectedMessage)
}
