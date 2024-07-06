package infrastructure

import (
	"github.com/google/uuid"
	"go-chi-gorilla-wire-workshop/app/domain"
)

type IdUuidRepository struct{}

func NewIdUuidRepository() domain.IdRepository {
	return &IdUuidRepository{}
}

func (repo *IdUuidRepository) GetId() string {
	return uuid.New().String()
}
