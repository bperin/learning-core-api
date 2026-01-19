package subjects

import (
	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/http/render"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for subjects
type Handler struct {
	service Service
}

// NewHandler creates a new subjects handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/subjects", h.ListAll)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {}

// ListAll godoc
// @Summary List all subjects
// @Description Retrieve all academic subjects with their nested sub-subjects
// @Tags subjects
// @Accept json
// @Produce json
// @Success 200 {array} subjects.Subject
// @Failure 500 {string} string "internal server error"
// @Security OAuth2Auth[read]
// @Router /subjects [get]
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subjects, err := h.service.ListAll(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, subjects)
}
