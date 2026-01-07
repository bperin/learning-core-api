package infra

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterDeps struct {
	JWTSecret string
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware(deps.JWTSecret))

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(UserIDKey)
			w.Write([]byte(fmt.Sprintf("User ID: %v", userID)))
		})
	})

	return r
}
