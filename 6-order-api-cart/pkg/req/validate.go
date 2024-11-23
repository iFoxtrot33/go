package req

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{10,14}$`)
	return phoneRegex.MatchString(phone)
}

func IsValid[T any](payload T) error {
	validate := validator.New()

	if err := validate.RegisterValidation("phone", validatePhone); err != nil {
		return err
	}

	err := validate.Struct(payload)

	return err
}
