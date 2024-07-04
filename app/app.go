package app

import (
	"encoding/json"
	"go-chi-gorilla-wire-workshop/app/domain"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router(service domain.CustomerService, r *chi.Mux) {
	r.Route("/customers", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var customer domain.Customer
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
