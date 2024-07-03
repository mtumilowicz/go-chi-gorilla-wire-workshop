package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/context"

	"go-chi-gorilla-wire-workshop/app" // Import the app package
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	repo := &app.Repository{
		Data: sync.Map{},
	}

	service := app.NewCustomerService(repo)

	app.Router(service, r)

	http.ListenAndServe(":8080", context.ClearHandler(r))
}
