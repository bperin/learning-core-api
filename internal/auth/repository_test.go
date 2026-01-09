package auth

import (
	"context"
	"testing"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) (*store.Queries, Repository, func()) {
	t.Helper()

	db := testutil.NewTestDB(t)
	queries := store.New(db)
	repo := NewRepository(queries)

	cleanup := func() {
		db.Close()
	}

	return queries, repo, cleanup
}

func seedUser(ctx context.Context, t *testing.T, q *store.Queries, email string) uuid.UUID {
	t.Helper()

	id := uuid.New()
	_, err := q.CreateUser(ctx, store.CreateUserParams{
		ID:       id,
		Email:    email,
		Password: "password123",
		IsAdmin:  false,
	})
	require.NoError(t, err)

	return id
}

func TestRepository_GetUserByEmail(t *testing.T) {
	queries, repo, cleanup := setupRepo(t)
	defer cleanup()

	ctx := context.Background()
	userID := seedUser(ctx, t, queries, "user@example.com")

	found, err := repo.GetUserByEmail(ctx, "user@example.com")
	require.NoError(t, err)
	assert.Equal(t, userID, found.ID)
	assert.Equal(t, "user@example.com", found.Email)
}

func TestRepository_GetUserByID(t *testing.T) {
	queries, repo, cleanup := setupRepo(t)
	defer cleanup()

	ctx := context.Background()
	userID := seedUser(ctx, t, queries, "lookup@example.com")

	found, err := repo.GetUserByID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, found.ID)
	assert.Equal(t, "lookup@example.com", found.Email)
}

func TestRepository_GetUserByEmail_NotFound(t *testing.T) {
	_, repo, cleanup := setupRepo(t)
	defer cleanup()

	_, err := repo.GetUserByEmail(context.Background(), "missing@example.com")
	assert.Error(t, err)
}
