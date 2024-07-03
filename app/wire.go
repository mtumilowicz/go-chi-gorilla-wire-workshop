//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
)

// InitializeApp initializes the application with the necessary dependencies.
func InitializeApp() *CustomerService {
	wire.Build(NewCustomerService, NewRepository)
	return &CustomerService{}
}

// NewRepository creates a new Repository instance.
func NewRepository() *Repository {
	return &Repository{}
}
