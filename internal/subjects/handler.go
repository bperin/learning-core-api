package subjects

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"learning-core-api/internal/http/render"
)

// Handler handles HTTP requests for subjects.
type Handler struct {
	service Service
}

// NewHandler creates a new subjects handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the subject routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/{id}", h.GetByID)
	r.Post("/", h.Create)
	r.Get("/user/{userID}", h.ListByUser)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// Create handles POST /subjects.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSubjectRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	subject, err := h.service.Create(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, subject)
}

// GetByID handles GET /subjects/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid subject ID")
		return
	}

	subject, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, subject)
}

// ListByUser handles GET /subjects/user/{userID}.
func (h *Handler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	subjects, err := h.service.ListByUser(r.Context(), userID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, subjects)
}

// Update handles PUT /subjects/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid subject ID")
		return
	}

	var req UpdateSubjectRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	subject, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, subject)
}

// Delete handles DELETE /subjects/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid subject ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Subject deleted successfully"})
}
