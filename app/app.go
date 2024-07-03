package app

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

// Customer represents a customer entity.
type Customer struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Repository acts as an in-memory store for customers.
type Repository struct {
	Data sync.Map
}

// CreateCustomer adds a new customer to the repository.
func (repo *Repository) CreateCustomer(customer Customer) {
	repo.Data.Store(customer.Name, customer)
}

// GetCustomer retrieves a customer from the repository by name.
func (repo *Repository) GetCustomer(name string) (Customer, bool) {
	value, ok := repo.Data.Load(name)
	if !ok {
		return Customer{}, false
	}
	return value.(Customer), true
}

// CustomerService provides customer-related operations.
type CustomerService struct {
	repo *Repository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(repo *Repository) *CustomerService {
	return &CustomerService{repo: repo}
}

// CreateCustomer adds a new customer using the repository.
func (service *CustomerService) CreateCustomer(customer Customer) {
	service.repo.CreateCustomer(customer)
}

// GetCustomer retrieves a customer by name using the repository.
func (service *CustomerService) GetCustomer(name string) (Customer, bool) {
	return service.repo.GetCustomer(name)
}

// Router sets up the HTTP routes.
func Router(service *CustomerService, r *chi.Mux) {
	r.Route("/customers", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var customer Customer
			if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			service.CreateCustomer(customer)
			w.WriteHeader(http.StatusCreated)
		})

		r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")
			customer, found := service.GetCustomer(name)
			if !found {
				http.Error(w, "customer not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(customer)
		})
	})
}
