package chunking_configs

import (
	"net/http"

	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/http/render"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	// No public routes for chunking configs
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/chunking-configs", h.ListAll)
	r.With(authz.RequireScope("read")).Get("/chunking-configs/{id}", h.GetByID)
	r.With(authz.RequireScope("read")).Get("/chunking-configs/active", h.GetActive)
	r.With(authz.RequireScope("write")).Post("/chunking-configs", h.Create)
	r.With(authz.RequireScope("write")).Post("/chunking-configs/{id}/activate", h.Activate)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	// Teachers have no access to chunking configs
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// Learners have no access to chunking configs
}

// ListAll lists all chunking configs.
// @Summary List all chunking configs
// @Description Get all chunking configurations
// @Tags Chunking Configs
// @Security OAuth2Auth[read]
// @Success 200 {array} ChunkingConfig "List of chunking configs"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /chunking-configs [get]
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	configs, err := h.service.ListAll(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, configs)
}

// GetByID retrieves a chunking config by ID.
// @Summary Get chunking config by ID
// @Description Retrieve a specific chunking configuration by its UUID
// @Tags Chunking Configs
// @Security OAuth2Auth[read]
// @Param id path string true "Config ID (UUID)"
// @Success 200 {object} ChunkingConfig "Chunking config details"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 404 {object} map[string]string "Config not found"
// @Router /chunking-configs/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid config ID")
		return
	}

	ctx := r.Context()
	config, err := h.service.GetByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Config not found")
		return
	}

	render.JSON(w, http.StatusOK, config)
}

// GetActive retrieves the active chunking config.
// @Summary Get active chunking config
// @Description Retrieve the currently active chunking configuration
// @Tags Chunking Configs
// @Security OAuth2Auth[read]
// @Success 200 {object} ChunkingConfig "Active chunking config"
// @Failure 404 {object} map[string]string "Active config not found"
// @Router /chunking-configs/active [get]
func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config, err := h.service.GetActive(ctx)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Active config not found")
		return
	}

	render.JSON(w, http.StatusOK, config)
}

// Create creates a new chunking config.
// @Summary Create new chunking config
// @Description Create a new chunking configuration (immutable - cannot edit existing)
// @Tags Chunking Configs
// @Security OAuth2Auth[write]
// @Accept json
// @Param request body CreateChunkingConfigRequest true "Chunking config request"
// @Success 201 {object} ChunkingConfig "Created chunking config"
// @Failure 400 {object} map[string]string "Bad request - invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /chunking-configs [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateChunkingConfigRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := r.Context()
	config, err := h.service.Create(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, config)
}

// Activate marks a chunking config as active.
// @Summary Activate chunking config
// @Description Mark a chunking configuration as active (deactivates other versions)
// @Tags Chunking Configs
// @Security OAuth2Auth[write]
// @Param id path string true "Config ID (UUID)"
// @Success 200 {object} map[string]string "Activation status"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /chunking-configs/{id}/activate [post]
func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid config ID")
		return
	}

	ctx := r.Context()
	if err := h.service.Activate(ctx, id); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"status": "activated"})
}
