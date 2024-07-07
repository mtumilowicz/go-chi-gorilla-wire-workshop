package domain

import (
	"fmt"
	"go-chi-gorilla-wire-workshop/app/validation"
)

type CustomerAlreadyExistsError struct {
	Id CustomerId
}

func (e CustomerAlreadyExistsError) Error() string {
	return fmt.Sprintf("customer with ID %s already exists", e.Id)
}

type Customer struct {
	Id   CustomerId
	Name string `validate:"min=1,max=30"`
	Age  int    `validate:"min=1,max=200"`
}

type CreateCustomerCommand struct {
	Name string `validate:"min=1,max=30"`
	Age  int    `validate:"min=1,max=200"`
}

func (c CreateCustomerCommand) toCustomer(id CustomerId) (Customer, error) {
	customer := Customer{
		Id:   id,
		Name: c.Name,
		Age:  c.Age,
	}
	if err := validation.Validate(customer); err != nil {
		return Customer{}, err
	}
	return customer, nil
}

type CustomerId struct {
	Raw string `validate:"min=1"`
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
	customer, err := command.toCustomer(customerId)
	if err != nil {
		return CustomerId{}, err
	}
	return service.repository.CreateCustomer(customer)
}

func (service CustomerService) GetCustomer(id CustomerId) (Customer, bool) {
	return service.repository.GetCustomer(id)
}
