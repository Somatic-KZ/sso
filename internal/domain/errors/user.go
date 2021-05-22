package errors

import "errors"

var (
	ErrUserIDNotSpec        = errors.New("user TDID not specified")
	ErrPhoneAlreadyVerified = errors.New("phone already verified")
)
