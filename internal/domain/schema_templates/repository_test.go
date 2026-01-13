package schema_templates

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
		userID, "schema-template@example.com", "password123", false, true, false)
	require.NoError(t, err)

	return userID
}

func TestSchemaTemplateRepository_CreateActivate(t *testing.T) {
	db, queries, repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)
	active := true
	generationType := "CLASSIFICATION"
	baseVersion := int32(0)
	existing, err := queries.ListSchemaTemplatesByGenerationType(ctx, store.GenerationType(generationType))
	require.NoError(t, err)
	for _, tmpl := range existing {
		if tmpl.Version > baseVersion {
			baseVersion = tmpl.Version
		}
	}

	first, err := repo.Create(ctx, CreateSchemaTemplateRequest{
		GenerationType: generationType,
		SchemaJSON:     []byte(`{"type":"object"}`),
		IsActive:       &active,
		CreatedBy:      userID,
	})
	require.NoError(t, err)
	assert.True(t, first.IsActive)
	assert.Equal(t, baseVersion+1, first.Version)

	second, err := repo.Create(ctx, CreateSchemaTemplateRequest{
		GenerationType: generationType,
		SchemaJSON:     []byte(`{"type":"object","properties":{"x":{"type":"string"}}}`),
		IsActive:       &active,
		CreatedBy:      userID,
	})
	require.NoError(t, err)
	assert.True(t, second.IsActive)
	assert.Equal(t, first.Version+1, second.Version)

	firstReload, err := repo.GetByID(ctx, first.ID)
	require.NoError(t, err)
	assert.False(t, firstReload.IsActive)

	activeSchema, err := repo.GetActiveByGenerationType(ctx, generationType)
	require.NoError(t, err)
	assert.Equal(t, second.ID, activeSchema.ID)

	activated, err := repo.Activate(ctx, first.ID)
	require.NoError(t, err)
	assert.True(t, activated.IsActive)

	activeSchema, err = repo.GetActiveByGenerationType(ctx, generationType)
	require.NoError(t, err)
	assert.Equal(t, first.ID, activeSchema.ID)

	secondReload, err := repo.GetByID(ctx, second.ID)
	require.NoError(t, err)
	assert.False(t, secondReload.IsActive)
}
