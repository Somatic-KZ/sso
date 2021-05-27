package validation

import (
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/pkg/validation"
)

func New(users *auth.Users) *validation.Validator {
	v := validation.New()

	validationRules := auth.NewValidationRules(users)
	_ = v.RegisterValidation("iin", validationRules.ValidateIIN, false)
	_ = v.RegisterValidation("unique_email", validationRules.ValidateUniqueEmail, false)
	_ = v.RegisterValidation("is_phone", validationRules.ValidateIsPhone, false)

	return v
}
