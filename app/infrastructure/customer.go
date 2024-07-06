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
	id := customer.Id
	repo.Data.Store(id, customer)
	return id
}

func (repo *CustomerInMemoryRepository) GetCustomer(id domain.CustomerId) (domain.Customer, bool) {
	value, ok := repo.Data.Load(id)
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
