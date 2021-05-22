package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateIsPhone(t *testing.T) {
	validate := validator.New()
	validationRules := NewValidationRules(nil)
	_ = validate.RegisterValidation("is_phone", validationRules.ValidateIsPhone, false)

	type TestCase struct {
		Phone              string `validate:"is_phone"`
		ValidationWillFail bool
	}

	testCases := []TestCase{
		{
			Phone:              "111",
			ValidationWillFail: true,
		},
		{
			Phone:              "+9-(701)-234-56-75",
			ValidationWillFail: true,
		},
		{
			Phone:              "777712345671321321231",
			ValidationWillFail: true,
		},
		{
			Phone:              "77771234567",
			ValidationWillFail: false,
		},
		{
			Phone:              "+77771234567",
			ValidationWillFail: false,
		},
		{
			Phone:              "87771234567",
			ValidationWillFail: false,
		},
		{
			Phone:              "+7(701)2345675",
			ValidationWillFail: false,
		},
		{
			Phone:              "+7 (701) 234 56 75",
			ValidationWillFail: false,
		},
		{
			Phone:              "+7-(701)-234-56-75",
			ValidationWillFail: false,
		},
	}

	for _, testCase := range testCases {
		if testCase.ValidationWillFail {
			assert.Error(t, validate.Struct(testCase))
			continue
		}
		assert.NoError(t, validate.Struct(testCase))
	}
}
