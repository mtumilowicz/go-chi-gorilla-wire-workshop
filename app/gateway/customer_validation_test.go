package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_CreateCustomerApiInput(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateCustomerApiInput
		expectError   bool
		errorMessages []string
	}{
		{
			name: "Valid input",
			input: CreateCustomerApiInput{
				Name: "John Doe",
				Age:  intPtr(30),
			},
			expectError:   false,
			errorMessages: nil,
		},
		{
			name: "Empty Name",
			input: CreateCustomerApiInput{
				Name: "",
				Age:  intPtr(30),
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Name' failed on the 'min' tag"},
		},
		{
			name: "Name too short",
			input: CreateCustomerApiInput{
				Name: "",
				Age:  intPtr(30),
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Name' failed on the 'min' tag"},
		},
		{
			name: "Name too long",
			input: CreateCustomerApiInput{
				Name: "ThisNameIsWayTooLongAndExceedsThirtyCharactersLimit",
				Age:  intPtr(30),
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Name' failed on the 'max' tag"},
		},
		{
			name: "Empty Age",
			input: CreateCustomerApiInput{
				Name: "John Doe",
				Age:  nil,
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Age' failed on the 'required' tag"},
		},
		{
			name: "Age too low",
			input: CreateCustomerApiInput{
				Name: "John Doe",
				Age:  intPtr(0),
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Age' failed on the 'min' tag"},
		},
		{
			name: "Age too high",
			input: CreateCustomerApiInput{
				Name: "John Doe",
				Age:  intPtr(201),
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Age' failed on the 'max' tag"},
		},
		{
			name: "Both fields empty",
			input: CreateCustomerApiInput{
				Name: "",
				Age:  nil,
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Name' failed on the 'min' tag", "Field validation for 'Age' failed on the 'required' tag"},
		},
		{
			name: "Both fields invalid",
			input: CreateCustomerApiInput{
				Name: "",
				Age:  intPtr(0),
			},
			expectError:   true,
			errorMessages: []string{"Field validation for 'Name' failed on the 'min' tag", "Field validation for 'Age' failed on the 'min' tag"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateInput(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				for _, msg := range tc.errorMessages {
					assert.Contains(t, err.Error(), msg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
