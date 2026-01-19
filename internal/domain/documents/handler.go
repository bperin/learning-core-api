package documents

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"learning-core-api/internal/gcp"
	"learning-core-api/internal/http/render"
)

type Handler struct {
	service    Service
	gcsService *gcp.GCSService
}

func NewHandler(service Service, gcsService *gcp.GCSService) *Handler {
	return &Handler{
		service:    service,
		gcsService: gcsService,
	}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	// Public routes for documents (if any)
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.Post("/documents/signed-url", h.GetSignedUploadURL)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.Post("/documents/signed-url", h.GetSignedUploadURL)
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// Learners cannot upload documents
}

// GetSignedUploadURLRequest represents a request for a signed upload URL
type GetSignedUploadURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	TTLSeconds  int    `json:"ttl_seconds,omitempty"`
}

// GetSignedUploadURLResponse represents the response with a signed URL
type GetSignedUploadURLResponse struct {
	SignedURL string `json:"signed_url"`
	Filename  string `json:"filename"`
	ExpiresAt string `json:"expires_at"`
}

// GetSignedUploadURL generates a signed URL for uploading a document to GCS
// Only admin and instructor (teacher) roles can access this endpoint
func (h *Handler) GetSignedUploadURL(w http.ResponseWriter, r *http.Request) {
	var req GetSignedUploadURLRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate request
	if req.Filename == "" {
		render.Error(w, http.StatusBadRequest, "Filename is required")
		return
	}

	if req.ContentType == "" {
		render.Error(w, http.StatusBadRequest, "Content type is required")
		return
	}

	// Validate content type
	if !isAllowedContentType(req.ContentType) {
		render.Error(w, http.StatusBadRequest, "Content type not allowed")
		return
	}

	// Get TTL from request or use default (15 minutes)
	ttl := 15 * time.Minute
	if req.TTLSeconds > 0 {
		ttl = time.Duration(req.TTLSeconds) * time.Second
	}

	// Generate signed URL
	ctx := r.Context()
	signedURL, err := h.gcsService.GenerateSignedUploadURL(ctx, req.Filename, req.ContentType, ttl)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to generate signed URL")
		return
	}

	response := GetSignedUploadURLResponse{
		SignedURL: signedURL,
		Filename:  req.Filename,
		ExpiresAt: time.Now().Add(ttl).Format(time.RFC3339),
	}

	render.JSON(w, http.StatusOK, response)
}

// isAllowedContentType checks if the content type is allowed for upload
func isAllowedContentType(contentType string) bool {
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"text/plain":      true,
		"text/markdown":   true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
		"application/msword": true, // .doc
		"text/html":          true,
		"application/json":   true,
	}
	return allowedTypes[contentType]
}
