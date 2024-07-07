package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type IdMockRepository struct {
	ReturnedId string
}

func (m *IdMockRepository) GetId() string {
	return m.ReturnedId
}

type CustomerInMemoryRepository struct {
	Data map[CustomerId]Customer
}

func newCustomerInMemoryRepository() CustomerRepository {
	return &CustomerInMemoryRepository{
		Data: map[CustomerId]Customer{},
	}
}

func (repo CustomerInMemoryRepository) CreateCustomer(customer Customer) (CustomerId, error) {
	id := customer.Id
	if _, ok := repo.Data[id]; ok {
		return CustomerId{}, CustomerAlreadyExistsError{Id: id}
	}
	repo.Data[id] = customer
	return id, nil
}

func (repo CustomerInMemoryRepository) GetCustomer(id CustomerId) (Customer, bool) {
	value, ok := repo.Data[id]
	if !ok {
		return Customer{}, false
	}
	return value, true
}

func TestCustomerService_CreateCustomer(t *testing.T) {
	// given
	customerRepository := newCustomerInMemoryRepository()
	idRepository := &IdMockRepository{
		ReturnedId: "1",
	}
	idService := NewIdService(idRepository)
	service := NewCustomerService(customerRepository, idService)

	// and
	command := CreateCustomerCommand{
		Name: "John Doe",
		Age:  30,
	}

	// when
	customerId, _ := service.CreateCustomer(command)

	// then
	assert.Equal(t, idRepository.ReturnedId, customerId.Raw, "Customer ID should be 'mock-id'")
}

func TestCustomerService_CreateExistingCustomer(t *testing.T) {
	// given
	customerRepository := newCustomerInMemoryRepository()
	idRepository := &IdMockRepository{
		ReturnedId: "1",
	}
	idService := NewIdService(idRepository)
	service := NewCustomerService(customerRepository, idService)

	// and
	command := CreateCustomerCommand{
		Name: "John Doe",
		Age:  30,
	}

	// and
	customerId, _ := service.CreateCustomer(command)

	// when
	_, err := service.CreateCustomer(command)

	// then
	assert.Error(t, err, "Creating customer with same id should produce an error")
	assert.Equal(t, CustomerAlreadyExistsError{Id: customerId}, err, "Error should be of type CustomerAlreadyExistsError")
}

func TestCustomerService_GetExistingCustomer(t *testing.T) {
	// given
	customerRepository := newCustomerInMemoryRepository()
	idRepository := &IdMockRepository{
		ReturnedId: "1",
	}
	idService := NewIdService(idRepository)
	service := NewCustomerService(customerRepository, idService)

	// and
	command := CreateCustomerCommand{
		Name: "John Doe",
		Age:  30,
	}

	// and
	customerId, _ := service.CreateCustomer(command)

	// when
	customer, found := service.GetCustomer(customerId)

	// then
	assert.True(t, found, "Customer should be found")
	expectedCustomer := Customer{
		Id:   customerId,
		Name: command.Name,
		Age:  command.Age,
	}
	assert.Equal(t, expectedCustomer, customer, "Returned customer should match mock data")
}

func TestCustomerService_GetNotExistingCustomer(t *testing.T) {
	// given
	customerRepository := newCustomerInMemoryRepository()
	idRepository := &IdMockRepository{
		ReturnedId: "1",
	}
	idService := NewIdService(idRepository)
	service := NewCustomerService(customerRepository, idService)

	// when
	_, found := service.GetCustomer(CustomerId{Raw: "not-existing"})

	// then
	assert.False(t, found, "Customer should not be found")
}
