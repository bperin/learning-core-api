package evals

import "github.com/go-chi/chi/v5"

// SuiteHandler handles HTTP requests for eval suites.
type SuiteHandler struct {
	service SuiteService
}

// NewSuiteHandler creates a new eval suites handler.
func NewSuiteHandler(service SuiteService) *SuiteHandler {
	return &SuiteHandler{
		service: service,
	}
}

// RegisterRoutes registers the eval suite routes.
func (h *SuiteHandler) RegisterRoutes(r chi.Router) {
	// TODO: define eval suite endpoints.
}
