package profile

import "errors"

var (
	ErrUnknownTDID          = errors.New("unknown user TDID")
	ErrPasswordTooShort     = errors.New("password too short")
	ErrPhoneNotSpecified    = errors.New("phone number not specified")
	ErrServerProblem        = errors.New("Временные проблемы с сервером")
	ErrRecoveryRequired     = errors.New("recovery required")
	ErrVerificationRequired = errors.New("verification required")
)
