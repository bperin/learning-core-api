package users

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID `json:"id"`
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

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email       string       `json:"email"`
	Password    string       `json:"password"`
	DisplayName *string      `json:"display_name,omitempty"`
	Role        UserRoleType `json:"role"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
