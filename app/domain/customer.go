package domain

import (
	"fmt"
)

type CustomerAlreadyExistsError struct {
	Id CustomerId
}

func (e CustomerAlreadyExistsError) Error() string {
	return fmt.Sprintf("customer with ID %s already exists", e.Id)
}

type Customer struct {
	Id   CustomerId
	Name string
	Age  int
}

type CreateCustomerCommand struct {
	Name string
	Age  int
}

func (c CreateCustomerCommand) toCustomer(id CustomerId) Customer {
	return Customer{
		Id:   id,
		Name: c.Name,
		Age:  c.Age,
	}
}

type CustomerId struct {
	Raw string
}

type CustomerRepository interface {
	CreateCustomer(customer Customer) (CustomerId, error)
	GetCustomer(id CustomerId) (Customer, bool)
}

type CustomerService struct {
	repository CustomerRepository
	idService  IdService
}

func NewCustomerService(repository CustomerRepository, idService IdService) CustomerService {
	return CustomerService{repository: repository, idService: idService}
}

func (service CustomerService) CreateCustomer(command CreateCustomerCommand) (CustomerId, error) {
	id := service.idService.GenerateId()
	customerId := CustomerId{Raw: id}
	customer := command.toCustomer(customerId)
	return service.repository.CreateCustomer(customer)
}

func (service CustomerService) GetCustomer(id CustomerId) (Customer, bool) {
	return service.repository.GetCustomer(id)
}
