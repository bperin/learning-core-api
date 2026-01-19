package model_configs

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
	// No public routes for model configs
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/model-configs", h.ListAll)
	r.With(authz.RequireScope("read")).Get("/model-configs/{id}", h.GetByID)
	r.With(authz.RequireScope("read")).Get("/model-configs/active", h.GetActive)
	r.With(authz.RequireScope("write")).Post("/model-configs", h.Create)
	r.With(authz.RequireScope("write")).Post("/model-configs/{id}/activate", h.Activate)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	// Teachers have no access to model configs
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// Learners have no access to model configs
}

// ListAll lists all model configs.
// @Summary List all model configs
// @Description Get all model configurations
// @Tags Model Configs
// @Security OAuth2[read]
// @Success 200 {array} ModelConfig "List of model configs"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /model-configs [get]
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	configs, err := h.service.ListAll(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, configs)
}

// GetByID retrieves a model config by ID.
// @Summary Get model config by ID
// @Description Retrieve a specific model configuration by its UUID
// @Tags Model Configs
// @Security OAuth2[read]
// @Param id path string true "Config ID (UUID)"
// @Success 200 {object} ModelConfig "Model config details"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 404 {object} map[string]string "Config not found"
// @Router /model-configs/{id} [get]
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

// GetActive retrieves the active model config.
// @Summary Get active model config
// @Description Retrieve the currently active model configuration
// @Tags Model Configs
// @Security OAuth2[read]
// @Success 200 {object} ModelConfig "Active model config"
// @Failure 404 {object} map[string]string "Active config not found"
// @Router /model-configs/active [get]
func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config, err := h.service.GetActive(ctx)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Active config not found")
		return
	}

	render.JSON(w, http.StatusOK, config)
}

// Create creates a new model config.
// @Summary Create model config
// @Description Create a new model configuration
// @Tags Model Configs
// @Security OAuth2[write]
// @Param request body CreateModelConfigRequest true "Model config data"
// @Success 201 {object} ModelConfig "Created model config"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /model-configs [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateModelConfigRequest
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

// Activate activates a model config.
// @Summary Activate model config
// @Description Mark a model configuration as active
// @Tags Model Configs
// @Security OAuth2[write]
// @Param id path string true "Config ID (UUID)"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 404 {object} map[string]string "Config not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /model-configs/{id}/activate [post]
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

	render.JSON(w, http.StatusOK, map[string]string{"message": "Model config activated successfully"})
}
