package gateway

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateInput[T any](input T) error {
	if err := validate.Struct(input); err != nil {
		return err
	}
	return nil
}
