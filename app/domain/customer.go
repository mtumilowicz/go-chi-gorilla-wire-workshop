package domain

type Customer struct {
	Id   CustomerId
	Name string
	Age  int
}

type CreateCustomerCommand struct {
	Name string
	Age  int
}

type CustomerId struct {
	Raw string
}

type CustomerRepository interface {
	CreateCustomer(command CreateCustomerCommand) CustomerId
	GetCustomer(id CustomerId) (Customer, bool)
}

type CustomerService struct {
	repository CustomerRepository
	idService  *IdService
}

func NewCustomerService(repository CustomerRepository, idService *IdService) *CustomerService {
	return &CustomerService{repository: repository, idService: idService}
}

func (service CustomerService) CreateCustomer(command CreateCustomerCommand) CustomerId {
	return service.repository.CreateCustomer(command)
}

func (service CustomerService) GetCustomer(id CustomerId) (Customer, bool) {
	return service.repository.GetCustomer(id)
}
