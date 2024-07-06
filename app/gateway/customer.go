package gateway

import (
	"encoding/json"
	"go-chi-gorilla-wire-workshop/app/domain"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CreateCustomerApiInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CustomerApiOutput struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func newCustomerApiOutput(customer domain.Customer) CustomerApiOutput {
	return CustomerApiOutput{
		Id:   customer.Id,
		Name: customer.Name,
		Age:  customer.Age,
	}
}

func (apiInput *CreateCustomerApiInput) toCommand() domain.CreateCustomerCommand {
	return domain.CreateCustomerCommand{
		Name: apiInput.Name,
		Age:  apiInput.Age,
	}
}

func CustomerRouter(service domain.CustomerService, r *chi.Mux) {
	r.Route("/customers", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var apiInput CreateCustomerApiInput
			if err := json.NewDecoder(r.Body).Decode(&apiInput); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			command := apiInput.toCommand()
			service.CreateCustomer(command)
			w.WriteHeader(http.StatusCreated)
		})

		r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")
			customer, found := service.GetCustomer(name)
			if !found {
				http.Error(w, "customer not found", http.StatusNotFound)
				return
			}
			apiOutput := newCustomerApiOutput(customer)
			json.NewEncoder(w).Encode(apiOutput)
		})
	})
}
