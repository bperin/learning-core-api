package evals

import "github.com/go-chi/chi/v5"

// RuleHandler handles HTTP requests for eval rules.
type RuleHandler struct {
	service RuleService
}

// NewRuleHandler creates a new eval rules handler.
func NewRuleHandler(service RuleService) *RuleHandler {
	return &RuleHandler{
		service: service,
	}
}

// RegisterRoutes registers the eval rule routes.
func (h *RuleHandler) RegisterRoutes(r chi.Router) {
	// TODO: define eval rule endpoints.
}
