package users

import (
	"context"
	"learning-core-api/internal/persistance/store"

	"github.com/google/uuid"
)

// Repository defines the interface for user persistence
type Repository interface {
	CreateUser(ctx context.Context, user User, password string, role UserRoleType) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, displayName *string, isActive *bool) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Role methods
	CreateUserRole(ctx context.Context, userRole UserRole) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error)
	DeleteUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error
}

type sqlRepository struct {
	queries *store.Queries
}

// NewRepository creates a new user repository
func NewRepository(queries *store.Queries) Repository {
	return &sqlRepository{
		queries: queries,
	}
}

func (r *sqlRepository) CreateUser(ctx context.Context, user User, password string, role UserRoleType) (*User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	isAdmin := role == UserRoleAdmin
	isLearner := role == UserRoleLearner
	isTeacher := role == UserRoleInstructor

	params := store.CreateUserParams{
		ID:        user.ID,
		Email:     user.Email,
		Password:  password,
		IsAdmin:   isAdmin,
		IsLearner: isLearner,
		IsTeacher: isTeacher,
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

func (r *sqlRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	dbUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

func (r *sqlRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

func (r *sqlRepository) UpdateUser(ctx context.Context, id uuid.UUID, displayName *string, isActive *bool) error {
	// Implementation depends on available UpdateUser query
	return nil
}

func (r *sqlRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *sqlRepository) CreateUserRole(ctx context.Context, userRole UserRole) error {
	// Implementation depends on available CreateUserRole query
	return nil
}

func (r *sqlRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error) {
	// Implementation depends on available GetUserRoles query
	return nil, nil
}

func (r *sqlRepository) DeleteUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error {
	// Implementation depends on available DeleteUserRole query
	return nil
}

// Helper functions for conversion
func toDomainUser(dbUser store.User) *User {
	return &User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
		// DisplayName: dbUser.DisplayName, // Add if present in store.User
		CreatedAt: dbUser.CreatedAt,
	}
}
