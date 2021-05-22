package auth

import (
	"errors"
	"fmt"
)

var ErrInternalError = errors.New("internal error")

var ErrUserDoesNotExist = errors.New("user does not exist")
var ErrEmailAlreadyTaken = errors.New("the email address is already taken")
var ErrPhoneNotSpecified = errors.New("phone number not specified")
var ErrEmailNotSpecified = errors.New("email address not specified")
var ErrInvalidEmailAddress = errors.New("invalid email address format")
var ErrPhoneNumNotLinkedToAccount = errors.New("phone number is not linked to account")
var ErrEmailNotLinkedToAccount = errors.New("email is not linked to account")

var ErrInvalidLoginOrPassword = errors.New("login or password is incorrect")
var ErrUserDisabled = errors.New("user disabled")

func ErrInvalidTdIDList(id string) error {
	return errors.New(fmt.Sprintf("invalid tdid in list: %s", id))
}
