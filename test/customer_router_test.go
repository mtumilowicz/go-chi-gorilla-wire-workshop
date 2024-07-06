package test

import (
	"bytes"
	"encoding/json"
	"go-chi-gorilla-wire-workshop/app"
	"go-chi-gorilla-wire-workshop/app/gateway"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCustomerRouter(t *testing.T) {
	t.Run("Create Customer", func(t *testing.T) {
		// given
		customerService := app.InitializeInMemoryApp()
		r := chi.NewRouter()
		gateway.CustomerRouter(*customerService, r)

		// and
		reqBody := gateway.CreateCustomerApiInput{Name: "John Doe"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// when
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusCreated, rr.Code)

		// and
		var apiOutput gateway.CustomerIdApiOutput
		err := json.NewDecoder(rr.Body).Decode(&apiOutput)
		assert.NoError(t, err)
		assert.NotEmpty(t, apiOutput.Id, "Customer ID should not be empty")

		// and
		location := rr.Header().Get("Location")
		assert.NotEmpty(t, location, "Location header should not be empty")
	})

	t.Run("Get Customer", func(t *testing.T) {
		// given
		customerService := app.InitializeInMemoryApp()
		r := chi.NewRouter()
		gateway.CustomerRouter(*customerService, r)

		// and
		reqBody := gateway.CreateCustomerApiInput{Name: "John Doe", Age: 30}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		var idApiOutput gateway.CustomerIdApiOutput
		err := json.NewDecoder(rr.Body).Decode(&idApiOutput)
		customerId := idApiOutput.Id

		// when
		req, _ = http.NewRequest("GET", "/customers/"+customerId, nil)
		rr = httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusOK, rr.Code)

		// and
		var apiOutput gateway.CustomerApiOutput
		err = json.NewDecoder(rr.Body).Decode(&apiOutput)
		assert.NoError(t, err)
		assert.Equal(t, gateway.CustomerApiOutput{
			Name: "John Doe",
			Id:   customerId,
			Age:  30,
		}, apiOutput)
	})

	t.Run("Get Non-Existent Customer", func(t *testing.T) {
		// given
		customerService := app.InitializeInMemoryApp()
		r := chi.NewRouter()
		gateway.CustomerRouter(*customerService, r)

		// and
		req, _ := http.NewRequest("GET", "/customers/NonExistent", nil)
		rr := httptest.NewRecorder()

		// when
		r.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// and
		location := rr.Header().Get("Location")
		assert.Empty(t, location, "Location header should be empty")
	})
}
