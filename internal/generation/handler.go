package generation

import (
	"encoding/json"
	"learning-core-api/internal/http/render"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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

// CreateGenerationRun handles POST /generation/runs
func (h *Handler) CreateGenerationRun(w http.ResponseWriter, r *http.Request) {
	var req CreateGenerationRunRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	run, err := h.service.CreateGenerationRun(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, run)
}

// GetGenerationRunByID handles GET /generation/runs/{id}
func (h *Handler) GetGenerationRunByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid generation run ID")
		return
	}

	ctx := r.Context()
	run, err := h.service.GetGenerationRunByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, run)
}

// ListGenerationRunsByModule handles GET /generation/runs/module/{moduleID}
func (h *Handler) ListGenerationRunsByModule(w http.ResponseWriter, r *http.Request) {
	moduleIDStr := chi.URLParam(r, "moduleID")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	ctx := r.Context()
	runs, err := h.service.ListGenerationRunsByModule(ctx, moduleID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, runs)
}

// UpdateGenerationRun handles PUT /generation/runs/{id}
func (h *Handler) UpdateGenerationRun(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid generation run ID")
		return
	}

	var req struct {
		Status        *RunStatus      `json:"status,omitempty"`
		OutputPayload json.RawMessage `json:"outputPayload,omitempty"`
		Error         json.RawMessage `json:"error,omitempty"`
		StartedAt     *time.Time      `json:"startedAt,omitempty"`
		FinishedAt    *time.Time      `json:"finishedAt,omitempty"`
	}
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	err = h.service.UpdateGenerationRun(ctx, id, req.Status, req.OutputPayload, req.Error, req.StartedAt, req.FinishedAt)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Generation run updated successfully"})
}

// DeleteGenerationRun handles DELETE /generation/runs/{id}
func (h *Handler) DeleteGenerationRun(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid generation run ID")
		return
	}

	ctx := r.Context()
	err = h.service.DeleteGenerationRun(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Generation run deleted successfully"})
}

// CreateArtifact handles POST /generation/artifacts
func (h *Handler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	var req CreateArtifactRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	artifact, err := h.service.CreateArtifact(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, artifact)
}

// GetArtifactByID handles GET /generation/artifacts/{id}
func (h *Handler) GetArtifactByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid artifact ID")
		return
	}

	ctx := r.Context()
	artifact, err := h.service.GetArtifactByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artifact)
}

// ListArtifactsByModule handles GET /generation/artifacts/module/{moduleID}
func (h *Handler) ListArtifactsByModule(w http.ResponseWriter, r *http.Request) {
	moduleIDStr := chi.URLParam(r, "moduleID")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	ctx := r.Context()
	artifacts, err := h.service.ListArtifactsByModule(ctx, moduleID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artifacts)
}

// ListArtifactsByModuleAndStatus handles GET /generation/artifacts/module/{moduleID}/status/{status}
func (h *Handler) ListArtifactsByModuleAndStatus(w http.ResponseWriter, r *http.Request) {
	moduleIDStr := chi.URLParam(r, "moduleID")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	statusStr := chi.URLParam(r, "status")
	status := ArtifactStatus(statusStr)

	ctx := r.Context()
	artifacts, err := h.service.ListArtifactsByModuleAndStatus(ctx, moduleID, status)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artifacts)
}

// UpdateArtifactStatus handles PUT /generation/artifacts/{id}/status
func (h *Handler) UpdateArtifactStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid artifact ID")
		return
	}

	var req struct {
		Status ArtifactStatus `json:"status"`
	}
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	err = h.service.UpdateArtifactStatus(ctx, id, req.Status, nil, nil)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Artifact status updated successfully"})
}

// DeleteArtifact handles DELETE /generation/artifacts/{id}
func (h *Handler) DeleteArtifact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid artifact ID")
		return
	}

	ctx := r.Context()
	err = h.service.DeleteArtifact(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Artifact deleted successfully"})
}
