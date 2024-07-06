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

func (repo CustomerInMemoryRepository) CreateCustomer(customer Customer) CustomerId {
	id := customer.Id
	repo.Data[id] = customer
	return id
}

func (repo CustomerInMemoryRepository) GetCustomer(id CustomerId) (Customer, bool) {
	value, ok := repo.Data[id]
	if !ok {
		return Customer{}, false
	}
	return value, true
}

func TestCustomerService_CreateCustomer(t *testing.T) {
	customerRepository := newCustomerInMemoryRepository()
	idRepository := &IdMockRepository{
		ReturnedId: "1",
	}
	idService := NewIdService(idRepository)
	service := NewCustomerService(customerRepository, idService)

	command := CreateCustomerCommand{
		Name: "John Doe",
		Age:  30,
	}

	customerId := service.CreateCustomer(command)

	assert.Equal(t, idRepository.ReturnedId, customerId.Raw, "Customer ID should be 'mock-id'")
}

func TestCustomerService_GetExistingCustomer(t *testing.T) {
	customerRepository := newCustomerInMemoryRepository()
	idRepository := &IdMockRepository{
		ReturnedId: "1",
	}
	idService := NewIdService(idRepository)
	service := NewCustomerService(customerRepository, idService)

	command := CreateCustomerCommand{
		Name: "John Doe",
		Age:  30,
	}

	customerId := service.CreateCustomer(command)

	customer, found := service.GetCustomer(customerId)
	expectedCustomer := Customer{
		Id:   customerId,
		Name: "John Doe",
		Age:  30,
	}

	assert.True(t, found, "Customer should be found")
	assert.Equal(t, expectedCustomer, customer, "Returned customer should match mock data")
}

func TestCustomerService_GetNotExistingCustomer(t *testing.T) {
	customerRepository := newCustomerInMemoryRepository()
	idRepository := &IdMockRepository{
		ReturnedId: "1",
	}
	idService := NewIdService(idRepository)
	service := NewCustomerService(customerRepository, idService)

	_, found := service.GetCustomer(CustomerId{Raw: "not-existing"})

	assert.False(t, found, "Customer should not be found")
}
