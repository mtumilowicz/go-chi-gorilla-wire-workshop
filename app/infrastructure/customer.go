package infrastructure

import (
	"github.com/google/uuid"
	"go-chi-gorilla-wire-workshop/app/domain"
	"sync"
)

type CustomerInMemoryRepository struct {
	Data sync.Map
}

func NewCustomerInMemoryRepository() domain.CustomerRepository {
	return &CustomerInMemoryRepository{}
}

func (repo *CustomerInMemoryRepository) CreateCustomer(command domain.CreateCustomerCommand) domain.CustomerId {
	customer := newCustomer(command)
	repo.Data.Store(command.Name, customer)
	return customer.Id
}

func (repo *CustomerInMemoryRepository) GetCustomer(name string) (domain.Customer, bool) {
	value, ok := repo.Data.Load(name)
	if !ok {
		return domain.Customer{}, false
	}
	return value.(domain.Customer), true
}

func newCustomer(command domain.CreateCustomerCommand) domain.Customer {
	id := uuid.New()
	return domain.Customer{
		Id:   domain.CustomerId{Raw: id.String()},
		Name: command.Name,
		Age:  command.Age,
	}
}
