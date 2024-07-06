package gateway

import (
	"encoding/json"
	"fmt"
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

type CustomerIdApiOutput struct {
	Id string `json:"id"`
}

func newCustomerIdApiOutput(customerId domain.CustomerId) CustomerIdApiOutput {
	return CustomerIdApiOutput{
		Id: customerId.Raw,
	}
}

func newCustomerApiOutput(customer domain.Customer) CustomerApiOutput {
	return CustomerApiOutput{
		Id:   customer.Id.Raw,
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
	baseUrl := "/customers"
	r.Route(baseUrl, func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var apiInput CreateCustomerApiInput
			if err := json.NewDecoder(r.Body).Decode(&apiInput); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			command := apiInput.toCommand()
			customerId := service.CreateCustomer(command)
			location := fmt.Sprintf("%s/%s", baseUrl, customerId)
			w.Header().Set("Location", location)
			w.WriteHeader(http.StatusCreated)
			apiOutput := newCustomerIdApiOutput(customerId)
			json.NewEncoder(w).Encode(apiOutput)
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
