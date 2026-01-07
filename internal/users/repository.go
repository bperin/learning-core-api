package users

import (
	"context"
	"database/sql"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// Repository defines the interface for user operations
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error)
	ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, displayName *string, isActive *bool) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// UserRole operations
	CreateUserRole(ctx context.Context, userRole UserRole) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error)
	DeleteUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new user repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// CreateUser creates a new user
func (r *repository) CreateUser(ctx context.Context, user User) (*User, error) {
	params := store.CreateUserParams{
		TenantID: user.TenantID,
		Email:    user.Email,
	}

	if user.DisplayName != nil {
		params.DisplayName = sql.NullString{String: *user.DisplayName, Valid: true}
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	user := toDomainUser(dbUser)
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	dbUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	user := toDomainUser(dbUser)
	return &user, nil
}

// GetUserByEmail retrieves a user by email and tenant ID
func (r *repository) GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, store.GetUserByEmailParams{
		TenantID: tenantID,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}

	user := toDomainUser(dbUser)
	return &user, nil
}

// ListUsersByTenant retrieves all users for a tenant
func (r *repository) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]User, error) {
	dbUsers, err := r.queries.ListUsersByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return toDomainUsers(dbUsers), nil
}

// UpdateUser updates a user
func (r *repository) UpdateUser(ctx context.Context, id uuid.UUID, displayName *string, isActive *bool) error {
	params := store.UpdateUserParams{
		ID: id,
	}

	current, err := r.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if displayName != nil {
		params.DisplayName = sql.NullString{String: *displayName, Valid: true}
	} else if current.DisplayName != nil {
		params.DisplayName = sql.NullString{String: *current.DisplayName, Valid: true}
	} else {
		params.DisplayName = sql.NullString{Valid: false}
	}

	if isActive != nil {
		params.IsActive = *isActive
	} else {
		params.IsActive = current.IsActive
	}

	return r.queries.UpdateUser(ctx, params)
}

// DeleteUser deletes a user
func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

// CreateUserRole creates a new user role
func (r *repository) CreateUserRole(ctx context.Context, userRole UserRole) error {
	params := store.CreateUserRoleParams{
		UserID: userRole.UserID,
		Role:   string(userRole.Role),
	}

	return r.queries.CreateUserRole(ctx, params)
}

// GetUserRoles retrieves all roles for a user
func (r *repository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error) {
	dbUserRoles, err := r.queries.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toDomainUserRoles(dbUserRoles), nil
}

// DeleteUserRole deletes a user role
func (r *repository) DeleteUserRole(ctx context.Context, userID uuid.UUID, role UserRoleType) error {
	return r.queries.DeleteUserRole(ctx, store.DeleteUserRoleParams{
		UserID: userID,
		Role:   string(role),
	})
}

func toDomainUser(dbUser store.User) User {
	return User{
		ID:          dbUser.ID,
		TenantID:    dbUser.TenantID,
		Email:       dbUser.Email,
		DisplayName: nullStringToPtr(dbUser.DisplayName),
		IsActive:    dbUser.IsActive,
		CreatedAt:   dbUser.CreatedAt,
	}
}

func toDomainUsers(dbUsers []store.User) []User {
	users := make([]User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = toDomainUser(dbUser)
	}
	return users
}

func toDomainUserRole(dbUserRole store.UserRole) UserRole {
	return UserRole{
		UserID:    dbUserRole.UserID,
		Role:      UserRoleType(dbUserRole.Role),
		GrantedAt: dbUserRole.GrantedAt,
	}
}

func toDomainUserRoles(dbUserRoles []store.UserRole) []UserRole {
	userRoles := make([]UserRole, len(dbUserRoles))
	for i, dbUserRole := range dbUserRoles {
		userRoles[i] = toDomainUserRole(dbUserRole)
	}
	return userRoles
}

func nullStringToPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}
