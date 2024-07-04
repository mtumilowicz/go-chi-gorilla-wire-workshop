package domain

type Customer struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CustomerRepository interface {
	CreateCustomer(customer Customer)
	GetCustomer(name string) (Customer, bool)
}

type CustomerService struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return CustomerService{repo: repo}
}

func (service CustomerService) CreateCustomer(customer Customer) {
	service.repo.CreateCustomer(customer)
}

func (service CustomerService) GetCustomer(name string) (Customer, bool) {
	return service.repo.GetCustomer(name)
}
