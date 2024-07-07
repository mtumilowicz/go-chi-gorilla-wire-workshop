package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type InvalidInput struct {
	Err error
}

func (e InvalidInput) Error() string {
	return fmt.Sprintf("invalid input: %v", e.Err)
}

var validate = validator.New()

func Validate[T any](input T) error {
	if err := validate.Struct(input); err != nil {
		return InvalidInput{Err: err}
	}
	return nil
}
