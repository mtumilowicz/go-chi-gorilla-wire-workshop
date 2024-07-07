package infrastructure

import (
	"go-chi-gorilla-wire-workshop/app/domain"
	"sync"
)

type CustomerInMemoryRepository struct {
	Data sync.Map
}

func NewCustomerInMemoryRepository() domain.CustomerRepository {
	return &CustomerInMemoryRepository{}
}

func (repo *CustomerInMemoryRepository) CreateCustomer(customer domain.Customer) (domain.CustomerId, error) {
	id := customer.Id
	// Use LoadOrStore to check and add atomically
	_, loaded := repo.Data.LoadOrStore(id, customer)
	if loaded {
		return domain.CustomerId{}, domain.CustomerAlreadyExistsError{Id: id}
	}
	return id, nil
}
func (repo *CustomerInMemoryRepository) GetCustomer(id domain.CustomerId) (domain.Customer, bool) {
	value, ok := repo.Data.Load(id)
	if !ok {
		return domain.Customer{}, false
	}
	return value.(domain.Customer), true
}
