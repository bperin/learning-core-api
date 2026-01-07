package users_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"
	"learning-core-api/internal/users"
)

func setupTestDB(t *testing.T) (*store.Queries, func()) {
	t.Helper()

	db := testutil.NewTestDB(t)
	queries := store.New(db)

	cleanup := func() {
		db.Close()
	}

	return queries, cleanup
}

func seedTenant(ctx context.Context, t *testing.T, q *store.Queries) uuid.UUID {
	t.Helper()

	tenant, err := q.CreateTenant(ctx, store.CreateTenantParams{
		Name:     "Test Tenant",
		IsActive: true,
	})
	require.NoError(t, err)

	return tenant.ID
}

func TestRepository_CreateAndGetUser(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	tenantID := seedTenant(ctx, t, queries)
	displayName := "Test User"

	created, err := repo.CreateUser(ctx, users.User{
		TenantID:    tenantID,
		Email:       "user@example.com",
		DisplayName: &displayName,
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, tenantID, created.TenantID)
	assert.Equal(t, "user@example.com", created.Email)
	require.NotNil(t, created.DisplayName)
	assert.Equal(t, displayName, *created.DisplayName)
	assert.True(t, created.IsActive)

	fetched, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Email, fetched.Email)
}

func TestRepository_GetUserByEmail(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	tenantID := seedTenant(ctx, t, queries)

	created, err := repo.CreateUser(ctx, users.User{
		TenantID: tenantID,
		Email:    "lookup@example.com",
	})
	require.NoError(t, err)

	fetched, err := repo.GetUserByEmail(ctx, tenantID, created.Email)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Email, fetched.Email)
}

func TestRepository_ListUsersByTenant(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	tenantID := seedTenant(ctx, t, queries)

	_, err := repo.CreateUser(ctx, users.User{
		TenantID: tenantID,
		Email:    "user1@example.com",
	})
	require.NoError(t, err)

	_, err = repo.CreateUser(ctx, users.User{
		TenantID: tenantID,
		Email:    "user2@example.com",
	})
	require.NoError(t, err)

	list, err := repo.ListUsersByTenant(ctx, tenantID)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestRepository_UpdateUser(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	tenantID := seedTenant(ctx, t, queries)
	created, err := repo.CreateUser(ctx, users.User{
		TenantID: tenantID,
		Email:    "update@example.com",
	})
	require.NoError(t, err)

	newName := "Updated Name"
	isActive := false

	err = repo.UpdateUser(ctx, created.ID, &newName, &isActive)
	require.NoError(t, err)

	updated, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, updated.DisplayName)
	assert.Equal(t, newName, *updated.DisplayName)
	assert.False(t, updated.IsActive)
}

func TestRepository_DeleteUser(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	tenantID := seedTenant(ctx, t, queries)
	created, err := repo.CreateUser(ctx, users.User{
		TenantID: tenantID,
		Email:    "delete@example.com",
	})
	require.NoError(t, err)

	err = repo.DeleteUser(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetUserByID(ctx, created.ID)
	assert.Error(t, err)
}

func TestRepository_UserRoles(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	tenantID := seedTenant(ctx, t, queries)
	created, err := repo.CreateUser(ctx, users.User{
		TenantID: tenantID,
		Email:    "roles@example.com",
	})
	require.NoError(t, err)

	err = repo.CreateUserRole(ctx, users.UserRole{
		UserID: created.ID,
		Role:   users.UserRoleInstructor,
	})
	require.NoError(t, err)

	roles, err := repo.GetUserRoles(ctx, created.ID)
	require.NoError(t, err)
	require.Len(t, roles, 1)
	assert.Equal(t, created.ID, roles[0].UserID)
	assert.Equal(t, users.UserRoleInstructor, roles[0].Role)
	assert.False(t, roles[0].GrantedAt.IsZero())

	err = repo.DeleteUserRole(ctx, created.ID, users.UserRoleInstructor)
	require.NoError(t, err)

	roles, err = repo.GetUserRoles(ctx, created.ID)
	require.NoError(t, err)
	assert.Empty(t, roles)
}
