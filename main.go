package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/context"
	"go-chi-gorilla-wire-workshop/app"
	"go-chi-gorilla-wire-workshop/app/gateway"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Use(middleware.Recoverer)

	customerService := app.InitializeApp()

	gateway.CustomerRouter(customerService, r)

	http.ListenAndServe(":8080", context.ClearHandler(r))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		token := r.Header.Get("Authorization")

		// Validate the token (you can add more validation logic here if needed)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}
