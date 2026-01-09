package evals

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
	r.With(authz.RequireScope("write")).Post("/evals", h.CreateEval)
	r.With(authz.RequireScope("write")).Post("/evals/{id}/publish", h.PublishEval)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/evals", h.ListEvals)
	r.With(authz.RequireScope("read")).Get("/evals/{id}", h.GetEval)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/evals", h.ListEvals)
	r.With(authz.RequireScope("read")).Get("/evals/{id}", h.GetEval)
}

// CreateEval godoc
// @Summary Create eval
// @Description Admin-only. Create a new eval in draft.
// @Tags evals
// @Accept json
// @Produce json
// @Success 501 {string} string "not implemented"
// @Security OAuth2Auth[write]
// @Router /evals [post]
func (h *Handler) CreateEval(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}

// PublishEval godoc
// @Summary Publish eval
// @Description Admin-only. Publish a draft eval (immutable after publish).
// @Tags evals
// @Accept json
// @Produce json
// @Param id path string true "Eval ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2Auth[write]
// @Router /evals/{id}/publish [post]
func (h *Handler) PublishEval(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}

// ListEvals godoc
// @Summary List evals
// @Description Teacher+Learner. List published evals.
// @Tags evals
// @Accept json
// @Produce json
// @Success 501 {string} string "not implemented"
// @Security OAuth2Auth[read]
// @Router /evals [get]
func (h *Handler) ListEvals(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}

// GetEval godoc
// @Summary Get eval
// @Description Teacher+Learner. Get a published eval by ID.
// @Tags evals
// @Accept json
// @Produce json
// @Param id path string true "Eval ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2Auth[read]
// @Router /evals/{id} [get]
func (h *Handler) GetEval(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}
