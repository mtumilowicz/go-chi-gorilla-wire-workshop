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

func (repo *CustomerInMemoryRepository) CreateCustomer(command domain.CreateCustomerCommand) {
	customer := newCustomer(command)
	repo.Data.Store(command.Name, customer)
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
		Id:   id.String(),
		Name: command.Name,
		Age:  command.Age,
	}
}
