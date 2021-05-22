package errors

import "errors"

var (
	ErrEmptyStruct      = errors.New("empty structure")
	ErrInvalidID        = errors.New("invalid id")
	ErrDoesNotExist     = errors.New("does not exist")
	ErrValidationFailed = errors.New("validation failed")
	ErrTooManyRequests  = errors.New("too many requests")
)
