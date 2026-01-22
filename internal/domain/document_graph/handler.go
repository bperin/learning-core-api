package document_graph

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"learning-core-api/internal/http/render"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.Post("/documents/{id}/graph/build", h.BuildGraph)
	r.Post("/documents/{id}/graph/query", h.QueryGraph)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.Post("/documents/{id}/graph/build", h.BuildGraph)
	r.Post("/documents/{id}/graph/query", h.QueryGraph)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {}

func (h *Handler) BuildGraph(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		render.Error(w, http.StatusServiceUnavailable, "Graph service unavailable")
		return
	}

	documentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	result, err := h.service.BuildGraph(r.Context(), documentID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, result)
}

func (h *Handler) QueryGraph(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		render.Error(w, http.StatusServiceUnavailable, "Graph service unavailable")
		return
	}

	documentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req QueryRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	response, err := h.service.QueryGraph(r.Context(), documentID, req.Query, req.Limit)
	if err != nil {
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, response)
}
