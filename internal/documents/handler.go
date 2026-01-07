package documents

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"learning-core-api/internal/http/render"
)

// Handler handles HTTP requests for documents
type Handler struct {
	service Service
}

// NewHandler creates a new documents handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the document routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/{id}", h.GetByID)
	r.Post("/", h.Create)
	r.Get("/subject/{subjectID}", h.ListBySubject)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// Create handles POST /documents
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateDocumentRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	document, err := h.service.Create(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, document)
}

// GetByID handles GET /documents/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	ctx := r.Context()
	document, err := h.service.GetByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, document)
}

// ListBySubject handles GET /documents/subject/{subjectID}
func (h *Handler) ListBySubject(w http.ResponseWriter, r *http.Request) {
	subjectIDStr := chi.URLParam(r, "subjectID")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid subject ID")
		return
	}

	ctx := r.Context()
	documents, err := h.service.ListBySubject(ctx, subjectID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, documents)
}

// Update handles PUT /documents/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req struct {
		Title     *string         `json:"title,omitempty"`
		Metadata  json.RawMessage `json:"metadata,omitempty"`
		IndexedAt *time.Time      `json:"indexedAt,omitempty"`
	}
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	err = h.service.Update(ctx, id, req.Title, req.Metadata, req.IndexedAt)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Document updated successfully"})
}

// Delete handles DELETE /documents/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	ctx := r.Context()
	err = h.service.Delete(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}
