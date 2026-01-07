package users

import (
	"learning-core-api/internal/http/render"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for users
type Handler struct {
	service Service
}

// NewHandler creates a new users handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the user routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	// User routes
	r.Get("/users/{id}", h.GetUserByID)
	r.Post("/users", h.CreateUser)
	r.Get("/users/tenant/{tenantID}", h.ListUsersByTenant)
	r.Get("/users/tenant/{tenantID}/email/{email}", h.GetUserByEmail)
	r.Put("/users/{id}", h.UpdateUser)
	r.Delete("/users/{id}", h.DeleteUser)

	// UserRole routes
	r.Post("/users/{userID}/roles/{role}", h.CreateUserRole)
	r.Get("/users/{userID}/roles", h.GetUserRoles)
	r.Delete("/users/{userID}/roles/{role}", h.DeleteUserRole)

}

// CreateUser handles POST /users/users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, user)
}

// GetUserByID handles GET /users/users/{id}
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx := r.Context()
	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, user)
}

// GetUserByEmail handles GET /users/users/tenant/{tenantID}/email/{email}
func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	email := chi.URLParam(r, "email")
	if email == "" {
		render.Error(w, http.StatusBadRequest, "Email is required")
		return
	}

	ctx := r.Context()
	user, err := h.service.GetUserByEmail(ctx, tenantID, email)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, user)
}

// ListUsersByTenant handles GET /users/users/tenant/{tenantID}
func (h *Handler) ListUsersByTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	ctx := r.Context()
	users, err := h.service.ListUsersByTenant(ctx, tenantID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, users)
}

// UpdateUser handles PUT /users/users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		DisplayName *string `json:"displayName,omitempty"`
		IsActive    *bool   `json:"isActive,omitempty"`
	}
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := r.Context()
	err = h.service.UpdateUser(ctx, id, req.DisplayName, req.IsActive)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// DeleteUser handles DELETE /users/users/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx := r.Context()
	err = h.service.DeleteUser(ctx, id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// CreateUserRole handles POST /users/users/{userID}/roles/{role}
func (h *Handler) CreateUserRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roleStr := chi.URLParam(r, "role")
	role := UserRoleType(roleStr)

	ctx := r.Context()
	err = h.service.CreateUserRole(ctx, userID, role)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "User role created successfully"})
}

// GetUserRoles handles GET /users/users/{userID}/roles
func (h *Handler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx := r.Context()
	roles, err := h.service.GetUserRoles(ctx, userID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, roles)
}

// DeleteUserRole handles DELETE /users/users/{userID}/roles/{role}
func (h *Handler) DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roleStr := chi.URLParam(r, "role")
	role := UserRoleType(roleStr)

	ctx := r.Context()
	err = h.service.DeleteUserRole(ctx, userID, role)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "User role deleted successfully"})
}
