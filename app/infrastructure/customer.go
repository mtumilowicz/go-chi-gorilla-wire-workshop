package infrastructure

import (
	"go-chi-gorilla-wire-workshop/app/domain"
	"sync"
)

type CustomerInMemoryRepository struct {
	Data sync.Map
}

func NewCustomerInMemoryRepository() *CustomerInMemoryRepository {
	return &CustomerInMemoryRepository{}
}

func (repo *CustomerInMemoryRepository) CreateCustomer(customer domain.Customer) {
	repo.Data.Store(customer.Name, customer)
}

func (repo *CustomerInMemoryRepository) GetCustomer(name string) (domain.Customer, bool) {
	value, ok := repo.Data.Load(name)
	if !ok {
		return domain.Customer{}, false
	}
	return value.(domain.Customer), true
}
