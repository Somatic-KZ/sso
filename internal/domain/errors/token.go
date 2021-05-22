package errors

import "errors"

var (
	ErrTokenHasExpired   = errors.New("token has expired")
	ErrTokenDoesNotExist = errors.New("token does not exist")
	ErrTokenNotSpec      = errors.New("token not specified")
	ErrOTPNotSpec        = errors.New("one time password not specified")
	ErrTokenTriesExpired = errors.New("token tries expired")
)
