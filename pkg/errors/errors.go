package errors

import (
	"errors"
)

var (
	ErrInternal        = errors.New("Internal server error")
	ErrNotFound        = errors.New("NotFound")
	ErrBadRequest      = errors.New("ErrBadRequest")
	ErrInvalidPassword = errors.New("InvalidPassword")
)
