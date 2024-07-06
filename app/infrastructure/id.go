package infrastructure

import "github.com/google/uuid"

type IdUuidRepository struct{}

func NewIdUuidRepository() *IdUuidRepository {
	return &IdUuidRepository{}
}

func (repo *IdUuidRepository) GetId() string {
	return uuid.New().String()
}
