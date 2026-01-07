package generation

import (
	"encoding/json"
	"learning-core-api/internal/http/render"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateGenerationRun handles POST /generation/runs.
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

// GetGenerationRunByID handles GET /generation/runs/{id}.
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

// ListGenerationRunsByModule handles GET /generation/runs/module/{moduleID}.
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

// UpdateGenerationRun handles PUT /generation/runs/{id}.
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

// DeleteGenerationRun handles DELETE /generation/runs/{id}.
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
