package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-chi-gorilla-wire-workshop/app/domain"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CreateCustomerApiInput struct {
	Name string `json:"name" validate:"required,min=1,max=30"`
	Age  *int   `json:"age" validate:"required,min=1,max=200"`
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

func (apiInput CreateCustomerApiInput) toCommand() (domain.CreateCustomerCommand, error) {
	if err := ValidateInput(apiInput); err != nil {
		return domain.CreateCustomerCommand{}, err
	}
	return domain.CreateCustomerCommand{
		Name: apiInput.Name,
		Age:  *apiInput.Age,
	}, nil
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
			command, err := apiInput.toCommand()
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			}
			customerId, err := service.CreateCustomer(command)
			if err != nil {
				message, status := customerErrorToHttp(err)
				http.Error(w, message, status)
			}
			location := fmt.Sprintf("%s/%s", baseUrl, customerId)
			w.Header().Set("Location", location)
			w.WriteHeader(http.StatusCreated)
			apiOutput := newCustomerIdApiOutput(customerId)
			json.NewEncoder(w).Encode(apiOutput)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			customer, found := service.GetCustomer(domain.CustomerId{Raw: id})
			if !found {
				http.Error(w, "customer not found", http.StatusNotFound)
				return
			}
			apiOutput := newCustomerApiOutput(customer)
			json.NewEncoder(w).Encode(apiOutput)
		})
	})
}

func customerErrorToHttp(err error) (message string, httpCode int) {
	switch {
	case errors.As(err, &domain.CustomerAlreadyExistsError{}):
		var e domain.CustomerAlreadyExistsError
		errors.As(err, &e)
		return err.Error(), http.StatusBadRequest
	default:
		return err.Error(), http.StatusInternalServerError
	}
}
