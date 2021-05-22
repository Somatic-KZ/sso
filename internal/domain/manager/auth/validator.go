package auth

import (
	"github.com/JetBrainer/sso/utils"
	"github.com/go-playground/validator/v10"
)

const (
	iinMin = 100000000000
	iinMax = 999999999999
)

type ValidationRules struct {
	users *Users
}

func NewValidationRules(users *Users) *ValidationRules {
	return &ValidationRules{users: users}
}

// ValidateIIN валидация ИИН
func (v ValidationRules) ValidateIIN(fl validator.FieldLevel) bool {
	iin := fl.Field().Int()

	if iin > iinMax || iin < iinMin {
		return false
	}

	return true
}

// ValidateUniqueEmail проверяет наличие уже существующего email'а у другого
// пользователя.
func (v ValidationRules) ValidateUniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	existsUser, _ := v.users.ByEmail(email)

	return existsUser == nil
}

// ValidateUniquePhone проверяет наличие уже существующего номера телефона у
// другого пользователя.
func (v ValidationRules) ValidateUniquePhone(fl validator.FieldLevel) bool {
	phone := utils.NormPhoneNum(fl.Field().String())

	// просматриваем первичный телефон
	user, _ := v.users.ByPrimaryPhone(phone)

	return user == nil
}

// ValidateUniquePhones проверяет наличие уже существующего и провалидированного номера телефона у
// другого пользователя.
func (v ValidationRules) ValidateUniquePhones(fl validator.FieldLevel) bool {
	phone := utils.NormPhoneNum(fl.Field().String())

	// просматриваем провалидированные телефоны
	user, _ := v.users.ByPhone(phone)

	return user == nil
}

// ValidateIsPhone проверка на правильность телефона
func (v ValidationRules) ValidateIsPhone(fl validator.FieldLevel) bool {
	return utils.NormPhoneNum(fl.Field().String()) != ""
}
