package drivers

import "errors"

var ErrUserIDNotSpec = errors.New("user TDID not specified")
var ErrUserLoginNotSpec = errors.New("login not specified")
var ErrEmptyUserStruct = errors.New("empty user structure")
var ErrUserDoesNotExist = errors.New("the user does not exist")
var ErrUserPhoneNotSpec = errors.New("phone number not specified")
var ErrUserEmailNotSpec = errors.New("email not specified")

var ErrTokenNotSpec = errors.New("token not specified")
var ErrTokenNotFound = errors.New("token not found")

var ErrEmptyRoleStruct = errors.New("empty role structure")
var ErrRoleDoesNotExist = errors.New("role does not exist")

var ErrEmptyStruct = errors.New("empty structure")

var ErrReceiverDoesNotExist = errors.New("the receiver does not exist")
