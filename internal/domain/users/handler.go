package users

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/http/render"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    = 3
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32
	saltLen       = 16
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

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Post("/signup", h.Signup)
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/users/{id}", h.GetUserByID)
	r.With(authz.RequireScope("write")).Post("/users", h.CreateUser)
	r.With(authz.RequireScope("read")).Get("/users/email/{email}", h.GetUserByEmail)
	r.With(authz.RequireScope("write")).Delete("/users/{id}", h.DeleteUser)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param request body users.CreateUserRequest true "Create User Request"
// @Success 201 {object} users.User
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Security OAuth2Auth[write]
// @Router /users [post]
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

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} users.User
// @Failure 400 {string} string "invalid request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Security OAuth2Auth[read]
// @Router /users/{id} [get]
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

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Retrieve a specific user by their email address
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} users.User
// @Failure 400 {string} string "invalid request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Security OAuth2Auth[read]
// @Router /users/email/{email} [get]
func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		render.Error(w, http.StatusBadRequest, "Email is required")
		return
	}

	ctx := r.Context()
	user, err := h.service.GetUserByEmail(ctx, email)
	if err != nil {
		render.Error(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, user)
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

// DeleteUser godoc
// @Summary Delete user
// @Description Remove a user from the system
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Security OAuth2Auth[write]
// @Router /users/{id} [delete]
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

// Signup godoc
// @Summary Sign up a new user
// @Description Create a new user account with email, password, and role
// @Tags users
// @Accept json
// @Produce json
// @Param request body users.SignupRequest true "Signup Request"
// @Success 201 {object} users.User
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /signup [post]
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string       `json:"email"`
		Password    string       `json:"password"`
		DisplayName *string      `json:"display_name,omitempty"`
		Role        UserRoleType `json:"role"`
	}
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Password == "" {
		render.Error(w, http.StatusBadRequest, "Password is required")
		return
	}

	if req.Role != UserRoleAdmin && req.Role != UserRoleInstructor && req.Role != UserRoleLearner {
		render.Error(w, http.StatusBadRequest, "Invalid role. Must be ADMIN, INSTRUCTOR, or LEARNER")
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	ctx := r.Context()
	createReq := CreateUserRequest{
		Email:       req.Email,
		Password:    hashedPassword,
		DisplayName: req.DisplayName,
		Role:        req.Role,
	}

	user, err := h.service.CreateUser(ctx, createReq)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, user)
}

// hashPassword hashes a password using Argon2id
func hashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	hashStr := base64.StdEncoding.EncodeToString(hash)
	saltStr := base64.StdEncoding.EncodeToString(salt)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argon2Memory, argon2Time, argon2Threads, saltStr, hashStr), nil
}
