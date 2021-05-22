package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/dongri/phonenumber"
)

// NormPhoneNum нормализует телефонные номера.
func NormPhoneNum(num string) string {
	if len(num) > 10 && num[0] == '8' {
		return phonenumber.Parse(num[1:], "KZ")
	}

	return phonenumber.Parse(num, "KZ")
}

// MaskPhoneNum накладывает скрывающую маску на номер телефона
func MaskPhoneNum(num string) string {
	normPhone := NormPhoneNum(num)
	if len(normPhone) < 11 {
		return normPhone
	}

	maskedPhone := ""
	for i, c := range normPhone {
		if i >= 4 && i <= 8 {
			maskedPhone += "*"
			continue
		}
		maskedPhone += string(c)
	}

	return maskedPhone
}

const maxUserLenInEmailAddr = 64

var (
	ErrInvalidFormat = errors.New("email: invalid format")
	userRegexp       = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp       = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
)

// метод для валидации email'ов
func ValidateEmail(email string) error {
	if len(email) < 6 || len(email) > 254 {
		return ErrInvalidFormat
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return ErrInvalidFormat
	}

	user := email[:at]
	host := email[at+1:]

	if len(user) > maxUserLenInEmailAddr {
		return ErrInvalidFormat
	}

	if !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return ErrInvalidFormat
	}

	return nil
}
