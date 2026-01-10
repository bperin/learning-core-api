package model_configs

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
)

func setupTestRepo(t *testing.T) (*sql.Tx, *store.Queries, Repository, func()) {
	t.Helper()

	tx, cleanup := testutil.NewTestTx(t)
	queries := store.New(tx)
	repo := NewRepository(queries)

	return tx, queries, repo, cleanup
}

func createTestUser(t *testing.T, db *sql.Tx) uuid.UUID {
	t.Helper()

	ctx := context.Background()
	userID := uuid.New()
	_, err := db.ExecContext(ctx,
		"INSERT INTO users (id, email, password, is_admin, is_learner, is_teacher) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "model-config@example.com", "password123", false, true, false)
	require.NoError(t, err)

	return userID
}

func TestModelConfigRepository_CreateActivate(t *testing.T) {
	db, _, repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	first, err := repo.Create(ctx, CreateModelConfigRequest{
		ModelName:   "gemini-1.5-pro",
		Temperature: 0.7,
		MaxTokens:   512,
		TopP:        0.9,
		TopK:        40,
		MimeType:    "application/json",
		IsActive:    true,
		CreatedBy:   userID,
	})
	require.NoError(t, err)
	assert.True(t, first.IsActive)

	second, err := repo.Create(ctx, CreateModelConfigRequest{
		ModelName:   "gemini-1.5-flash",
		Temperature: 0.3,
		MaxTokens:   256,
		TopP:        0.8,
		TopK:        20,
		MimeType:    "application/json",
		IsActive:    true,
		CreatedBy:   userID,
	})
	require.NoError(t, err)
	assert.True(t, second.IsActive)

	firstReload, err := repo.GetByID(ctx, first.ID)
	require.NoError(t, err)
	assert.False(t, firstReload.IsActive)

	activeConfig, err := repo.GetActive(ctx)
	require.NoError(t, err)
	assert.Equal(t, second.ID, activeConfig.ID)

	err = repo.Activate(ctx, first.ID)
	require.NoError(t, err)

	activeConfig, err = repo.GetActive(ctx)
	require.NoError(t, err)
	assert.Equal(t, first.ID, activeConfig.ID)

	secondReload, err := repo.GetByID(ctx, second.ID)
	require.NoError(t, err)
	assert.False(t, secondReload.IsActive)

	all, err := repo.ListAll(ctx)
	require.NoError(t, err)
	assert.Len(t, all, 3)
}
