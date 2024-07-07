package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chi-gorilla-wire-workshop/app/validation"
)

func TestCreateCustomerCommand_Validation(t *testing.T) {
	tests := []struct {
		name      string
		command   CreateCustomerCommand
		expectErr bool
	}{
		{
			name:      "valid command",
			command:   CreateCustomerCommand{Name: "John Doe", Age: 25},
			expectErr: false,
		},
		{
			name:      "empty name",
			command:   CreateCustomerCommand{Name: "", Age: 25},
			expectErr: true,
		},
		{
			name:      "name too long",
			command:   CreateCustomerCommand{Name: "John Doe with a very long name that exceeds 30 characters", Age: 25},
			expectErr: true,
		},
		{
			name:      "age too low",
			command:   CreateCustomerCommand{Name: "John Doe", Age: 0},
			expectErr: true,
		},
		{
			name:      "age too high",
			command:   CreateCustomerCommand{Name: "John Doe", Age: 201},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.command)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func FuzzCreateCustomerCommand_Validation(f *testing.F) {
	// Seed with a few valid examples
	f.Add("John Doe", 25)
	f.Add("Jane Smith", 30)
	f.Add("Alice Johnson", 45)

	f.Fuzz(func(t *testing.T, name string, age int) {
		if len(name) < 1 || len(name) > 30 {
			return
		}
		if age < 1 || age > 200 {
			return
		}

		command := CreateCustomerCommand{Name: name, Age: age}
		err := validation.Validate(command)
		assert.NoError(t, err)
	})
}

func TestCustomer_Validation(t *testing.T) {
	tests := []struct {
		name      string
		customer  Customer
		expectErr bool
	}{
		{
			name:      "valid customer",
			customer:  Customer{Id: CustomerId{Raw: "123"}, Name: "John Doe", Age: 25},
			expectErr: false,
		},
		{
			name:      "empty id",
			customer:  Customer{Id: CustomerId{Raw: ""}, Name: "John Doe", Age: 25},
			expectErr: true,
		},
		{
			name:      "empty name",
			customer:  Customer{Id: CustomerId{Raw: "123"}, Name: "", Age: 25},
			expectErr: true,
		},
		{
			name:      "name too long",
			customer:  Customer{Id: CustomerId{Raw: "123"}, Name: "John Doe with a very long name that exceeds 30 characters", Age: 25},
			expectErr: true,
		},
		{
			name:      "age too low",
			customer:  Customer{Id: CustomerId{Raw: "123"}, Name: "John Doe", Age: 0},
			expectErr: true,
		},
		{
			name:      "age too high",
			customer:  Customer{Id: CustomerId{Raw: "123"}, Name: "John Doe", Age: 201},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.customer)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateCustomerCommand_toCustomer(t *testing.T) {
	tests := []struct {
		name       string
		command    CreateCustomerCommand
		customerId CustomerId
		expectErr  bool
	}{
		{
			name:       "valid command",
			command:    CreateCustomerCommand{Name: "John Doe", Age: 25},
			customerId: CustomerId{Raw: "123"},
			expectErr:  false,
		},
		{
			name:       "invalid command",
			command:    CreateCustomerCommand{Name: "", Age: 25},
			customerId: CustomerId{Raw: "123"},
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := tt.command.toCustomer(tt.customerId)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, Customer{}, customer)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.customerId, customer.Id)
				assert.Equal(t, tt.command.Name, customer.Name)
				assert.Equal(t, tt.command.Age, customer.Age)
			}
		})
	}
}
