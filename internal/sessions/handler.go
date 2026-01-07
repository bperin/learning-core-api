package sessions

import (
	"learning-core-api/internal/http/render"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for sessions
type Handler struct {
	service Service
}

// NewHandler creates a new sessions handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the session routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	// Session routes
	r.Get("/sessions/{id}", h.GetSessionByID)
	r.Post("/sessions", h.CreateSession)
	r.Get("/sessions/user/{userID}", h.ListSessionsByUser)
	r.Get("/sessions/module/{moduleID}", h.ListSessionsByModule)
	r.Delete("/sessions/{id}", h.DeleteSession)

	// Attempt routes
	r.Get("/attempts/{id}", h.GetAttemptByID)
	r.Post("/attempts", h.CreateAttempt)
	r.Get("/attempts/session/{sessionID}", h.ListAttemptsBySession)
	r.Get("/attempts/tenant/{tenantID}", h.ListAttemptsByTenant)
	r.Delete("/attempts/{id}", h.DeleteAttempt)
}

// CreateSession handles POST /sessions/sessions
func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	session, err := h.service.CreateSession(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, session)
}

// GetSessionByID handles GET /sessions/sessions/{id}
func (h *Handler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	ctx := r.Context()
	session, err := h.service.GetSessionByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, session)
}

// ListSessionsByUser handles GET /sessions/sessions/user/{userID}
func (h *Handler) ListSessionsByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx := r.Context()
	sessions, err := h.service.ListSessionsByUser(ctx, userID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, sessions)
}

// ListSessionsByModule handles GET /sessions/sessions/module/{moduleID}
func (h *Handler) ListSessionsByModule(w http.ResponseWriter, r *http.Request) {
	moduleIDStr := chi.URLParam(r, "moduleID")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	ctx := r.Context()
	sessions, err := h.service.ListSessionsByModule(ctx, moduleID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, sessions)
}

// DeleteSession handles DELETE /sessions/sessions/{id}
func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	ctx := r.Context()
	err = h.service.DeleteSession(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Session deleted successfully"})
}

// CreateAttempt handles POST /sessions/attempts
func (h *Handler) CreateAttempt(w http.ResponseWriter, r *http.Request) {
	var req CreateAttemptRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	attempt, err := h.service.CreateAttempt(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, attempt)
}

// GetAttemptByID handles GET /sessions/attempts/{id}
func (h *Handler) GetAttemptByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid attempt ID")
		return
	}

	ctx := r.Context()
	attempt, err := h.service.GetAttemptByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, attempt)
}

// ListAttemptsBySession handles GET /sessions/attempts/session/{sessionID}
func (h *Handler) ListAttemptsBySession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "sessionID")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	ctx := r.Context()
	attempts, err := h.service.ListAttemptsBySession(ctx, sessionID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, attempts)
}

// ListAttemptsByTenant handles GET /sessions/attempts/tenant/{tenantID}
func (h *Handler) ListAttemptsByTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	ctx := r.Context()
	attempts, err := h.service.ListAttemptsByTenant(ctx, tenantID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, attempts)
}

// DeleteAttempt handles DELETE /sessions/attempts/{id}
func (h *Handler) DeleteAttempt(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid attempt ID")
		return
	}

	ctx := r.Context()
	err = h.service.DeleteAttempt(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Attempt deleted successfully"})
}
