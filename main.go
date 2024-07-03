package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/context"
	"net/http"

	"go-chi-gorilla-wire-workshop/app" // Import the app package
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	service := app.InitializeApp()

	app.Router(service, r)

	http.ListenAndServe(":8080", context.ClearHandler(r))
}
