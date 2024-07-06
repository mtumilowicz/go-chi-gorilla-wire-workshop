package domain

type IdRepository interface {
	GetId() string
}

type IdService struct {
	repository IdRepository
}

func NewIdService(repository IdRepository) IdService {
	return IdService{
		repository: repository,
	}
}

func (service *IdService) GenerateId() string {
	return service.repository.GetId()
}
