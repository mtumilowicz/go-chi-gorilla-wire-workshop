//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"go-chi-gorilla-wire-workshop/app/domain"
	"go-chi-gorilla-wire-workshop/app/infrastructure"
)

func InitializeApp() *domain.CustomerService {
	wire.Build(
		infrastructure.NewCustomerInMemoryRepository,
		infrastructure.NewIdUuidRepository,
		domain.NewIdService,
		domain.NewCustomerService,
	)
	return &domain.CustomerService{}
}

func InitializeInMemoryApp() *domain.CustomerService {
	wire.Build(
		infrastructure.NewCustomerInMemoryRepository,
		infrastructure.NewIdUuidRepository,
		domain.NewIdService,
		domain.NewCustomerService,
	)
	return &domain.CustomerService{}
}
