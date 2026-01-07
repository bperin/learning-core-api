package generation

import "github.com/go-chi/chi/v5"

// Handler handles HTTP requests for generation
type Handler struct {
	service Service
}

// NewHandler creates a new generation handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the generation routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	// GenerationRun routes
	r.Get("/runs/{id}", h.GetGenerationRunByID)
	r.Post("/runs", h.CreateGenerationRun)
	r.Get("/runs/module/{moduleID}", h.ListGenerationRunsByModule)
	r.Put("/runs/{id}", h.UpdateGenerationRun)
	r.Delete("/runs/{id}", h.DeleteGenerationRun)

	// Artifact routes
	r.Get("/artifacts/{id}", h.GetArtifactByID)
	r.Post("/artifacts", h.CreateArtifact)
	r.Get("/artifacts/module/{moduleID}", h.ListArtifactsByModule)
	r.Get("/artifacts/module/{moduleID}/status/{status}", h.ListArtifactsByModuleAndStatus)
	r.Put("/artifacts/{id}/status", h.UpdateArtifactStatus)
	r.Delete("/artifacts/{id}", h.DeleteArtifact)
}
