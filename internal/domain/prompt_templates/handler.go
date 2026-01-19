package prompt_templates

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
	// No public routes for prompt templates
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/prompt-templates", h.ListByGenerationType)
	r.With(authz.RequireScope("read")).Get("/prompt-templates/{id}", h.GetByID)
	r.With(authz.RequireScope("read")).Get("/prompt-templates/generation-type/{generationType}", h.GetActiveByGenerationType)
	r.With(authz.RequireScope("write")).Post("/prompt-templates", h.Create)
	r.With(authz.RequireScope("write")).Post("/prompt-templates/{id}/activate", h.Activate)
	r.With(authz.RequireScope("write")).Post("/prompt-templates/{id}/deactivate", h.Deactivate)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	// Teachers can read classification and generation schemas only
	r.With(authz.RequireScope("read")).Get("/prompt-templates/generation-type/{generationType}", h.GetActiveByGenerationType)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// Learners can read active classification and generation schemas only
	r.With(authz.RequireScope("read")).Get("/prompt-templates/generation-type/{generationType}", h.GetActiveByGenerationType)
}

// ListByGenerationType lists all prompt templates for a generation type.
// @Summary List prompt templates by generation type
// @Description Get all prompt templates for a specific generation type
// @Tags Prompt Templates
// @Security OAuth2[read]
// @Param generation_type query string true "Generation type (CLASSIFICATION or QUESTIONS)"
// @Success 200 {array} PromptTemplate "List of prompt templates"
// @Failure 400 {object} map[string]string "Bad request - missing generation_type"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /prompt-templates [get]
func (h *Handler) ListByGenerationType(w http.ResponseWriter, r *http.Request) {
	generationType := r.URL.Query().Get("generation_type")
	if generationType == "" {
		render.Error(w, http.StatusBadRequest, "generation_type query parameter is required")
		return
	}

	ctx := r.Context()
	templates, err := h.service.ListByGenerationType(ctx, generationType)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, templates)
}

// GetByID retrieves a prompt template by ID.
// @Summary Get prompt template by ID
// @Description Retrieve a specific prompt template by its UUID
// @Tags Prompt Templates
// @Security OAuth2[read]
// @Param id path string true "Template ID (UUID)"
// @Success 200 {object} PromptTemplate "Prompt template details"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 404 {object} map[string]string "Template not found"
// @Router /prompt-templates/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	ctx := r.Context()
	template, err := h.service.GetByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Template not found")
		return
	}

	render.JSON(w, http.StatusOK, template)
}

// GetActiveByGenerationType retrieves the active prompt template for a generation type.
// @Summary Get active prompt template by generation type
// @Description Retrieve the currently active prompt template for a specific generation type
// @Tags Prompt Templates
// @Security OAuth2[read]
// @Param generationType path string true "Generation type (CLASSIFICATION or QUESTIONS)"
// @Success 200 {object} PromptTemplate "Active prompt template"
// @Failure 400 {object} map[string]string "Bad request - missing generation_type"
// @Failure 404 {object} map[string]string "Active template not found"
// @Router /prompt-templates/generation-type/{generationType} [get]
func (h *Handler) GetActiveByGenerationType(w http.ResponseWriter, r *http.Request) {
	generationType := chi.URLParam(r, "generationType")
	if generationType == "" {
		render.Error(w, http.StatusBadRequest, "generation_type is required")
		return
	}

	ctx := r.Context()
	template, err := h.service.GetActiveByGenerationType(ctx, generationType)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Active template not found")
		return
	}

	render.JSON(w, http.StatusOK, template)
}

// Create creates a new version of a prompt template.
// @Summary Create new prompt template version
// @Description Create a new version of a prompt template (immutable - cannot edit existing)
// @Tags Prompt Templates
// @Security OAuth2[write]
// @Accept json
// @Param request body CreatePromptTemplateVersionRequest true "Template version request"
// @Success 201 {object} PromptTemplate "Created prompt template"
// @Failure 400 {object} map[string]string "Bad request - invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /prompt-templates [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePromptTemplateVersionRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := r.Context()
	template, err := h.service.CreateVersion(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, template)
}

// Activate marks a prompt template as active.
// @Summary Activate prompt template
// @Description Mark a prompt template version as active (deactivates other versions)
// @Tags Prompt Templates
// @Security OAuth2[write]
// @Param id path string true "Template ID (UUID)"
// @Success 200 {object} PromptTemplate "Activated prompt template"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /prompt-templates/{id}/activate [post]
func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	ctx := r.Context()
	template, err := h.service.Activate(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, template)
}

// Deactivate marks a prompt template as inactive.
// @Summary Deactivate prompt template
// @Description Mark a prompt template version as inactive
// @Tags Prompt Templates
// @Security OAuth2[write]
// @Param id path string true "Template ID (UUID)"
// @Success 200 {object} PromptTemplate "Deactivated prompt template"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /prompt-templates/{id}/deactivate [post]
func (h *Handler) Deactivate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	ctx := r.Context()
	template, err := h.service.Deactivate(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, template)
}
