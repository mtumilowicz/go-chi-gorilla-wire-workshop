package domain

type Customer struct {
	Id   string
	Name string
	Age  int
}

type CreateCustomerCommand struct {
	Name string
	Age  int
}

type CustomerRepository interface {
	CreateCustomer(command CreateCustomerCommand)
	GetCustomer(name string) (Customer, bool)
}

type CustomerService struct {
	repository CustomerRepository
	idService  IdService
}

func NewCustomerService(repository CustomerRepository, idService IdService) CustomerService {
	return CustomerService{repository: repository, idService: idService}
}

func (service CustomerService) CreateCustomer(command CreateCustomerCommand) {
	service.repository.CreateCustomer(command)
}

func (service CustomerService) GetCustomer(name string) (Customer, bool) {
	return service.repository.GetCustomer(name)
}
