package auth

import "errors"

var (
	ErrInvalidNumberFormat  = errors.New("invalid phone number format")
	ErrInvalidToken         = errors.New("token is incorrect or expired")
	ErrUnknownTDID          = errors.New("unknown user TDID")
	ErrPasswordTooShort     = errors.New("password too short")
	ErrPhoneNotSpecified    = errors.New("phone number not specified")
	ErrServerProblem        = errors.New("Временные проблемы с сервером")
	ErrRecoveryRequired     = errors.New("recovery required")
	ErrVerificationRequired = errors.New("verification required")
)
