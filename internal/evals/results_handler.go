package evals

import "github.com/go-chi/chi/v5"

// ResultHandler handles HTTP requests for eval results.
type ResultHandler struct {
	service ResultService
}

// NewResultHandler creates a new eval results handler.
func NewResultHandler(service ResultService) *ResultHandler {
	return &ResultHandler{
		service: service,
	}
}

// RegisterRoutes registers the eval result routes.
func (h *ResultHandler) RegisterRoutes(r chi.Router) {
	// TODO: define eval result endpoints.
}
