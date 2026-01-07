package modules

import (
	"learning-core-api/internal/http/render"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for modules
type Handler struct {
	service Service
}

// NewHandler creates a new modules handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the module routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/{id}", h.GetByID)
	r.Post("/", h.Create)
	r.Get("/tenant/{tenantID}", h.ListByTenant)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// Create handles POST /modules
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateModuleRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	module, err := h.service.Create(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, module)
}

// GetByID handles GET /modules/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	ctx := r.Context()
	module, err := h.service.GetByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, module)
}

// ListByTenant handles GET /modules/tenant/{tenantID}
func (h *Handler) ListByTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	ctx := r.Context()
	modules, err := h.service.ListByTenant(ctx, tenantID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, modules)
}

// Update handles PUT /modules/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	var req struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	err = h.service.Update(ctx, id, req.Name, req.Description)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Module updated successfully"})
}

// Delete handles DELETE /modules/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	ctx := r.Context()
	err = h.service.Delete(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Module deleted successfully"})
}
