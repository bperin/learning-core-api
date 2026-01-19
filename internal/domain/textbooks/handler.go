package textbooks

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"learning-core-api/internal/http/render"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/textbooks/subjects", h.GetAllSubjects)
	r.Get("/textbooks/subjects/{name}", h.GetSubjectByName)
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.Post("/textbooks/scrape", h.ScrapeAndStore)
	r.Get("/admin/textbooks/subjects", h.AdminListSubjects)
	r.Get("/admin/textbooks/subjects/{slug}/books", h.AdminGetBooksBySubject)
	r.Post("/admin/textbooks/download", h.AdminDownloadBooks)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	// Teachers can view subjects but not scrape
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// Learners can view subjects
}

// GetAllSubjects retrieves all subjects
func (h *Handler) GetAllSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subjects, err := h.service.GetAllSubjects(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to retrieve subjects")
		return
	}

	render.JSON(w, http.StatusOK, map[string]interface{}{
		"subjects": subjects,
		"count":    len(subjects),
	})
}

// GetSubjectByName retrieves a subject by name
func (h *Handler) GetSubjectByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		render.Error(w, http.StatusBadRequest, "Subject name is required")
		return
	}

	ctx := r.Context()
	subject, err := h.service.GetSubjectByName(ctx, name)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to retrieve subject")
		return
	}

	if subject == nil {
		render.Error(w, http.StatusNotFound, "Subject not found")
		return
	}

	render.JSON(w, http.StatusOK, subject)
}

// ScrapeAndStore scrapes subjects from Open Textbook Library and stores them
// Only accessible to admins
func (h *Handler) ScrapeAndStore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.service.ScrapeAndStoreSubjects(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to scrape and store subjects")
		return
	}

	render.JSON(w, http.StatusOK, map[string]interface{}{
		"message":    "Subjects scraped and stored successfully",
		"count":      len(result.Subjects),
		"scraped_at": result.ScrapedAt,
	})
}

// AdminListSubjects returns a list of all available subjects for admin workflow
func (h *Handler) AdminListSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subjects, err := h.service.GetAllSubjects(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to retrieve subjects")
		return
	}

	render.JSON(w, http.StatusOK, map[string]interface{}{
		"subjects": subjects,
		"count":    len(subjects),
	})
}

// AdminGetBooksBySubject returns books for a given subject slug
func (h *Handler) AdminGetBooksBySubject(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		render.Error(w, http.StatusBadRequest, "Subject slug is required")
		return
	}

	ctx := r.Context()
	books, err := h.service.GetBooksBySubject(ctx, slug)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	render.JSON(w, http.StatusOK, map[string]interface{}{
		"subject": slug,
		"books":   books,
		"count":   len(books),
	})
}

// AdminDownloadRequest represents a request to download selected books
type AdminDownloadRequest struct {
	SubjectSlug string   `json:"subject_slug"`
	BookURLs    []string `json:"book_urls"`
	MaxBooks    int      `json:"max_books,omitempty"`
}

// AdminDownloadBooks downloads selected books from a subject
func (h *Handler) AdminDownloadBooks(w http.ResponseWriter, r *http.Request) {
	var req AdminDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SubjectSlug == "" {
		render.Error(w, http.StatusBadRequest, "Subject slug is required")
		return
	}

	if len(req.BookURLs) == 0 {
		render.Error(w, http.StatusBadRequest, "At least one book URL is required")
		return
	}

	ctx := r.Context()
	result, err := h.service.DownloadBooks(ctx, req.SubjectSlug, req.BookURLs)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to download books")
		return
	}

	render.JSON(w, http.StatusOK, map[string]interface{}{
		"message":         "Books download initiated",
		"subject":         req.SubjectSlug,
		"requested":       len(req.BookURLs),
		"downloaded":      result.DownloadedCount,
		"failed":          result.FailedCount,
		"download_path":   result.DownloadPath,
	})
}
