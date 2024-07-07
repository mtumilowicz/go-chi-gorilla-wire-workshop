package app

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate[T any](input T) error {
	if err := validate.Struct(input); err != nil {
		return err
	}
	return nil
}
