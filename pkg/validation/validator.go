package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrTypeCastToValErrs tells that incoming error is not convertible to validator.ValidationErrors
const ErrTypeCastToValErrs = "could not type cast validation errors from"

// Validator is a composition struct that uses external validator as a base and expands
// it by additional methods.
type Validator struct {
	*validator.Validate
}

// New returns new instance of our custom validator.
func New() *Validator {
	v := &Validator{
		validator.New(),
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

// MapFromValidationErrors returns map of fields and validation violations on those fields.
func (v Validator) MapFromValidationErrors(vErr error) map[string]string {
	validationErrors := vErr.(validator.ValidationErrors)
	errorsMap := make(map[string]string)
	for _, fieldError := range validationErrors {
		errorsMap[fieldError.Field()] = fieldError.ActualTag()
	}

	return errorsMap
}
