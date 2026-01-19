package content_discovery

import (
	"net/http"

	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/http/render"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Post("/content-discovery/books", h.ListBooks)
	r.Post("/content-discovery/download-pdfs", h.DownloadPDFs)
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Post("/content-discovery/books", h.ListBooks)
	r.With(authz.RequireScope("write")).Post("/content-discovery/download-pdfs", h.DownloadPDFs)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Post("/content-discovery/books", h.ListBooks)
	r.With(authz.RequireScope("write")).Post("/content-discovery/download-pdfs", h.DownloadPDFs)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// No learner access to content discovery
}

// ListBooks lists books from selected subjects for test creation.
// @Summary List books from selected subjects
// @Description Fetch books from the provided subject URLs for test creation workflow. Returns book title, detail page URL, and PDF download link when available.
// @Tags Content Discovery
// @Security OAuth2
// @Accept json
// @Produce json
// @Param request body BookListRequest true "Subject IDs and max books"
// @Success 200 {object} BookListResponse "List of books with subject information and PDF links"
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

// DownloadPDFs downloads PDFs from selected books and processes them for content generation.
// @Summary Download PDFs from selected books
// @Description Downloads PDFs from the provided book links, creates document records, and processes them through the file search service for classification generation
// @Tags Content Discovery
// @Security OAuth2
// @Accept json
// @Produce json
// @Param request body PDFDownloadRequest true "Books to download"
// @Success 200 {object} PDFDownloadResponse "PDF download job started"
// @Failure 400 {object} map[string]string "Bad request - invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /content-discovery/download-pdfs [post]
func (h *Handler) DownloadPDFs(w http.ResponseWriter, r *http.Request) {
	var req PDFDownloadRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Books) == 0 {
		render.Error(w, http.StatusBadRequest, "No books provided")
		return
	}

	ctx := r.Context()
	
	// Extract user ID from context
	userID := authz.UserIDFromContext(ctx)
	if userID == "" {
		render.Error(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	
	response, err := h.service.DownloadPDFs(ctx, req, userID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, response)
}
