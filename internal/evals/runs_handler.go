package evals

import "github.com/go-chi/chi/v5"

// RunHandler handles HTTP requests for eval runs.
type RunHandler struct {
	service RunService
}

// NewRunHandler creates a new eval runs handler.
func NewRunHandler(service RunService) *RunHandler {
	return &RunHandler{
		service: service,
	}
}

// RegisterRoutes registers the eval run routes.
func (h *RunHandler) RegisterRoutes(r chi.Router) {
	// TODO: define eval run endpoints.
}
