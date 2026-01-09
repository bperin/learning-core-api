package users_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/users"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
)

func setupTestDB(t *testing.T) (*store.Queries, func()) {
	t.Helper()

	tx, cleanup := testutil.NewTestTx(t)
	queries := store.New(tx)

	return queries, cleanup
}

func TestRepository_CreateAndGetUser(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, users.User{
		ID:    uuid.New(),
		Email: "user@example.com",
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "user@example.com", created.Email)
	assert.False(t, created.CreatedAt.IsZero())

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

	created, err := repo.CreateUser(ctx, users.User{
		ID:    uuid.New(),
		Email: "lookup@example.com",
	})
	require.NoError(t, err)

	fetched, err := repo.GetUserByEmail(ctx, created.Email)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Email, fetched.Email)
}

func TestRepository_DeleteUser(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := users.NewRepository(queries)
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, users.User{
		ID:    uuid.New(),
		Email: "delete@example.com",
	})
	require.NoError(t, err)

	err = repo.DeleteUser(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetUserByID(ctx, created.ID)
	assert.Error(t, err)
}
