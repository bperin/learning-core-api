package users

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for user business logic
type Service interface {
	// User operations
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, displayName *string, isActive *bool) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// UserRole operations
	CreateUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error)
	DeleteUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateUser creates a new user with business logic validation
func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	// Validate email format
	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Check if user with same email already exists
	existing, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	if req.Role == "" {
		return nil, errors.New("role is required")
	}

	// Create the user
	user := User{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		IsActive:    true, // Default to active
	}

	return s.repo.CreateUser(ctx, user, req.Password, req.Role)
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	if len(email) < 5 {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	local, domain := parts[0], parts[1]
	if len(local) == 0 || len(domain) < 3 {
		return false
	}
	return true
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	return s.repo.GetUserByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.repo.GetUserByEmail(ctx, email)
}

// UpdateUser updates a user with business logic validation
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, displayName *string, isActive *bool) error {
	if id == uuid.Nil {
		return errors.New("user ID is required")
	}
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.UpdateUser(ctx, id, displayName, isActive)
}

// DeleteUser deletes a user
func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("user ID is required")
	}
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteUser(ctx, id)
}

// CreateUserRole creates a new user role with business logic validation
func (s *service) CreateUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error {
	if userID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if role == "" {
		return errors.New("role is required")
	}
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	userRole := UserRole{
		UserID:    userID,
		Role:      role,
		GrantedAt: time.Now(),
	}
	return s.repo.CreateUserRole(ctx, userRole)
}

// GetUserRoles retrieves all roles for a user
func (s *service) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	return s.repo.GetUserRoles(ctx, userID)
}

// DeleteUserRole deletes a user role
func (s *service) DeleteUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error {
	if userID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if role == "" {
		return errors.New("role is required")
	}
	return s.repo.DeleteUserRole(ctx, userID, role)
}
