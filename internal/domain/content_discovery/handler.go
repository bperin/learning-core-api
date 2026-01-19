package content_discovery

import (
	"net/http"

	"github.com/go-chi/chi/v5"

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
	// No public routes for content discovery
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Post("/content-discovery/books", h.ListBooks)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Post("/content-discovery/books", h.ListBooks)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// No learner access to content discovery
}

// ListBooks lists books from selected subjects for test creation.
// @Summary List books from selected subjects
// @Description Fetch books from the provided subject URLs for test creation workflow
// @Tags Content Discovery
// @Security OAuth2Auth[read]
// @Accept json
// @Produce json
// @Param request body BookListRequest true "Subject IDs and max books"
// @Success 200 {object} BookListResponse "List of books with subject information"
// @Failure 400 {object} map[string]string "Bad request - invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /content-discovery/books [post]
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	var req BookListRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := r.Context()
	response, err := h.service.ListBooksFromSubjects(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, response)
}
