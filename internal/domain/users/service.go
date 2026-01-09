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
	GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error)
	ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]User, error)
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
	// Business logic validation
	if req.TenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	// Validate email format
	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Check if user with same email already exists in the tenant
	existing, err := s.repo.GetUserByEmail(ctx, req.TenantID, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists in the tenant")
	}

	// Create the user
	user := User{
		TenantID:    req.TenantID,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		IsActive:    true, // Default to active
	}

	return s.repo.CreateUser(ctx, user)
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	// Simple email validation - in a real application, use a proper email validation library
	if len(email) < 5 {
		return false
	}

	// Check for @ and domain
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

// GetUserByEmail retrieves a user by email and tenant ID
func (s *service) GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	return s.repo.GetUserByEmail(ctx, tenantID, email)
}

// ListUsersByTenant retrieves all users for a tenant
func (s *service) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]User, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.ListUsersByTenant(ctx, tenantID)
}

// UpdateUser updates a user with business logic validation
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, displayName *string, isActive *bool) error {
	if id == uuid.Nil {
		return errors.New("user ID is required")
	}

	// Check if user exists
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

	// Check if user exists before deleting
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

	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Create the user role
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
