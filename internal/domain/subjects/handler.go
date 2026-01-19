package subjects

import (
	"strings"

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
	r.With(authz.RequireScope("read")).Get("/subjects/for-selection", h.ListForSelection)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/subjects", h.ListAll)
	r.With(authz.RequireScope("read")).Get("/subjects/for-selection", h.ListForSelection)
}

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

// ListForSelection godoc
// @Summary List subjects for test creation selection
// @Description Retrieve subjects formatted for selection UI (tag cloud style)
// @Tags subjects
// @Accept json
// @Produce json
// @Success 200 {array} subjects.SubjectForSelection
// @Failure 500 {string} string "internal server error"
// @Security OAuth2Auth[read]
// @Router /subjects/for-selection [get]
func (h *Handler) ListForSelection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subjects, err := h.service.ListAll(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Transform subjects for selection UI
	var selectionSubjects []SubjectForSelection
	for _, subject := range subjects {
		// Extract the display name from the full name (e.g., "Education" from "Open Textbook Library - Education Textbooks")
		displayName := extractDisplayName(subject.Name)

		selectionSubjects = append(selectionSubjects, SubjectForSelection{
			ID:          subject.ID,
			DisplayName: displayName,
			FullName:    subject.Name,
			URL:         subject.Url,
		})

		// Also add sub-subjects
		for _, subSubject := range subject.SubSubjects {
			subDisplayName := extractDisplayName(subSubject.Name)
			selectionSubjects = append(selectionSubjects, SubjectForSelection{
				ID:          subSubject.ID,
				DisplayName: subDisplayName,
				FullName:    subSubject.Name,
				URL:         subSubject.Url,
				ParentID:    &subject.ID,
			})
		}
	}

	render.JSON(w, http.StatusOK, selectionSubjects)
}

// extractDisplayName extracts the subject name for display (e.g., "Education" from "Open Textbook Library - Education Textbooks")
func extractDisplayName(fullName string) string {
	// Look for patterns like "Open Textbook Library - Subject Textbooks"
	if idx := strings.Index(fullName, " - "); idx != -1 {
		afterDash := fullName[idx+3:]
		// Remove " Textbooks" suffix if present
		if strings.HasSuffix(afterDash, " Textbooks") {
			return strings.TrimSuffix(afterDash, " Textbooks")
		}
		return afterDash
	}
	return fullName
}
