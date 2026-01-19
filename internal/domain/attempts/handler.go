package attempts

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
	r.With(authz.RequireScope("read")).Get("/attempts/{id}", h.GetAttempt)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/attempts/{id}", h.GetAttempt)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	r.With(authz.RequireScope("write")).Post("/evals/{id}/attempts", h.StartAttempt)
	r.With(authz.RequireScope("read")).Get("/attempts/{id}", h.GetAttempt)
	r.With(authz.RequireScope("write")).Post("/attempts/{id}/answers", h.SubmitAnswer)
	r.With(authz.RequireScope("write")).Post("/attempts/{id}/submit", h.SubmitAttempt)
}

// StartAttempt godoc
// @Summary Start test attempt
// @Description Learner-only. Create a new attempt for an eval.
// @Tags attempts
// @Accept json
// @Produce json
// @Param id path string true "Eval ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2[write]
// @Router /evals/{id}/attempts [post]
func (h *Handler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}

// GetAttempt godoc
// @Summary Get attempt
// @Description Learner+Teacher+Admin. Get a test attempt by ID.
// @Tags attempts
// @Accept json
// @Produce json
// @Param id path string true "Attempt ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2[read]
// @Router /attempts/{id} [get]
func (h *Handler) GetAttempt(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}

// SubmitAnswer godoc
// @Summary Submit answer
// @Description Learner-only. Submit an answer for an attempt.
// @Tags attempts
// @Accept json
// @Produce json
// @Param id path string true "Attempt ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2[write]
// @Router /attempts/{id}/answers [post]
func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}

// SubmitAttempt godoc
// @Summary Submit attempt
// @Description Learner-only. Submit an attempt and finalize score.
// @Tags attempts
// @Accept json
// @Produce json
// @Param id path string true "Attempt ID"
// @Success 501 {string} string "not implemented"
// @Security OAuth2[write]
// @Router /attempts/{id}/submit [post]
func (h *Handler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	render.Error(w, http.StatusNotImplemented, "not implemented")
}
