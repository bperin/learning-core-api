package reviews

import (
	"net/http"

	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/http/render"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("write")).Post("/eval-items/{id}/reviews", h.CreateEvalItemReview)
	r.With(authz.RequireScope("read")).Get("/eval-items/{id}/reviews", h.ListEvalItemReviews)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.With(authz.RequireScope("write")).Post("/eval-items/{id}/reviews", h.CreateEvalItemReview)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {}

// CreateEvalItemReview godoc
// @Summary Create eval item review
// @Description Teacher+Admin. Submit a review verdict for an eval item.
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path string true "Eval Item ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2[write]
// @Router /eval-items/{id}/reviews [post]
func (h *Handler) CreateEvalItemReview(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}

// ListEvalItemReviews godoc
// @Summary List eval item reviews
// @Description Admin-only. List review history for an eval item.
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path string true "Eval Item ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2[read]
// @Router /eval-items/{id}/reviews [get]
func (h *Handler) ListEvalItemReviews(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}
