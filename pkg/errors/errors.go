package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

func New(status int) Error {
	c := GetCode(status)
	return Error{
		Code: c,
		Err:  errors.New(c.Message),
	}
}

func Newf(status int, msg ...string) Error {
	c := GetCode(status)
	return Error{
		Code: c,
		Err:  errors.New(fmt.Sprintf(c.Message, msg)),
	}
}

type Error struct {
	Err  error
	Code Code
}

func (e Error) Error() string {
	return fmt.Sprintf("status: %d, message: %s", e.Code.Status, e.Code.Message)
}

func (e Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	return e.Code.Status == t.Code.Status
}

func (e Error) As(target interface{}) bool {
	switch target.(type) {
	case Error:
		target = e
		return true
	default:
		return false
	}
}

func (e Error) Unwrap() error {
	return e.Err
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func WithStack(err error) error {
	return errors.WithStack(err)
}

func WithStackByCode(status int) error {
	return errors.WithStack(New(status))
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
