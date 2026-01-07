package generation

import (
	"learning-core-api/internal/http/render"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateArtifact handles POST /generation/artifacts.
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

// GetArtifactByID handles GET /generation/artifacts/{id}.
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

// ListArtifactsByModule handles GET /generation/artifacts/module/{moduleID}.
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

// ListArtifactsByModuleAndStatus handles GET /generation/artifacts/module/{moduleID}/status/{status}.
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

// UpdateArtifactStatus handles PUT /generation/artifacts/{id}/status.
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

// DeleteArtifact handles DELETE /generation/artifacts/{id}.
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
