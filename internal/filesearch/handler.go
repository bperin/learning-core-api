package filesearch

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"learning-core-api/internal/http/render"
)

// Handler handles HTTP requests for file search operations.
type Handler struct {
	service Service
}

// NewHandler creates a new file search handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers file search routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/subject/{subjectID}/upload", h.UploadToSubject)
}

// UploadToSubject handles POST /file-search/subject/{subjectID}/upload.
func (h *Handler) UploadToSubject(w http.ResponseWriter, r *http.Request) {
	subjectIDStr := chi.URLParam(r, "subjectID")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid subject ID")
		return
	}

	var req UploadRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	req.SubjectID = subjectID

	result, err := h.service.UploadToSubject(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, result)
}
