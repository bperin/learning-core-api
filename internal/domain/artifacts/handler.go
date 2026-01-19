package artifacts

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	httpPkg "learning-core-api/internal/http"
	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/http/render"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	// No public routes for artifacts
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/artifacts", h.ListArtifacts)
	r.With(authz.RequireScope("read")).Get("/artifacts/{id}", h.GetArtifactByID)
	r.With(authz.RequireScope("read")).Get("/artifacts/stats", h.GetArtifactStats)
	r.With(authz.RequireScope("read")).Get("/artifacts/type/{type}", h.GetArtifactsByType)
	r.With(authz.RequireScope("read")).Get("/artifacts/status/{status}", h.GetArtifactsByStatus)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	// Teachers have limited access to artifacts
	r.With(authz.RequireScope("read")).Get("/artifacts/stats", h.GetArtifactStats)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// Learners have no access to artifacts
}

// ListArtifacts lists all artifacts with pagination.
// @Summary List all artifacts
// @Description Get all artifacts with pagination support
// @Tags Artifacts
// @Security OAuth2Auth[read]
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} httpPkg.PaginatedResponse[store.Artifact] "Paginated list of artifacts"
// @Failure 400 {object} map[string]string "Bad request - invalid pagination parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /artifacts [get]
func (h *Handler) ListArtifacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get pagination parameters
	pagination := httpPkg.GetPaginationParams(r)
	
	// Get artifacts
	artifacts, err := h.service.ListArtifacts(ctx, int32(pagination.Limit), int32(pagination.Offset))
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	// For now, we'll use the count of returned items as total
	// In a real implementation, you'd want a separate count query
	total := int64(len(artifacts))
	
	response := httpPkg.NewPaginatedResponse(artifacts, pagination, total)
	render.JSON(w, http.StatusOK, response)
}

// GetArtifactByID retrieves an artifact by ID.
// @Summary Get artifact by ID
// @Description Retrieve a specific artifact by its UUID
// @Tags Artifacts
// @Security OAuth2Auth[read]
// @Param id path string true "Artifact ID (UUID)"
// @Success 200 {object} store.Artifact "Artifact details"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 404 {object} map[string]string "Artifact not found"
// @Router /artifacts/{id} [get]
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
		render.Error(w, http.StatusNotFound, "Artifact not found")
		return
	}

	render.JSON(w, http.StatusOK, artifact)
}

// GetArtifactsByType retrieves artifacts filtered by type.
// @Summary Get artifacts by type
// @Description Retrieve artifacts filtered by their type
// @Tags Artifacts
// @Security OAuth2Auth[read]
// @Param type path string true "Artifact type"
// @Success 200 {array} store.Artifact "List of artifacts"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /artifacts/type/{type} [get]
func (h *Handler) GetArtifactsByType(w http.ResponseWriter, r *http.Request) {
	artifactType := chi.URLParam(r, "type")
	if artifactType == "" {
		render.Error(w, http.StatusBadRequest, "Artifact type is required")
		return
	}

	ctx := r.Context()
	artifacts, err := h.service.GetArtifactsByType(ctx, artifactType)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artifacts)
}

// GetArtifactsByStatus retrieves artifacts filtered by status.
// @Summary Get artifacts by status
// @Description Retrieve artifacts filtered by their status
// @Tags Artifacts
// @Security OAuth2Auth[read]
// @Param status path string true "Artifact status"
// @Success 200 {array} store.Artifact "List of artifacts"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /artifacts/status/{status} [get]
func (h *Handler) GetArtifactsByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")
	if status == "" {
		render.Error(w, http.StatusBadRequest, "Artifact status is required")
		return
	}

	ctx := r.Context()
	artifacts, err := h.service.GetArtifactsByStatus(ctx, status)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, artifacts)
}

// GetArtifactStats retrieves artifact statistics.
// @Summary Get artifact statistics
// @Description Retrieve statistics about artifacts (counts by status, etc.)
// @Tags Artifacts
// @Security OAuth2Auth[read]
// @Success 200 {object} store.GetArtifactStatsRow "Artifact statistics"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /artifacts/stats [get]
func (h *Handler) GetArtifactStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := h.service.GetArtifactStats(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, stats)
}
