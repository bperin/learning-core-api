package users

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Email       string    `json:"email"`
	DisplayName *string   `json:"display_name,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRole represents a user's role
type UserRole struct {
	UserID    uuid.UUID    `json:"user_id"`
	Role      UserRoleType `json:"role"`
	GrantedAt time.Time    `json:"granted_at"`
}

// UserRoleType represents the type of user role
type UserRoleType string

const (
	UserRoleAdmin      UserRoleType = "ADMIN"
	UserRoleInstructor UserRoleType = "INSTRUCTOR"
	UserRoleLearner    UserRoleType = "LEARNER"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	Email       string    `json:"email"`
	DisplayName *string   `json:"display_name,omitempty"`
}

// CreateTenantRequest represents the request to create a tenant
type CreateTenantRequest struct {
	Name string `json:"name"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
