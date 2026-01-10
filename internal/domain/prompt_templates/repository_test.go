package prompt_templates

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

func TestPromptTemplateRepository_CreateActivate(t *testing.T) {
	_, _, repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	active := true
	createdBy := "tester"
	key := "intent_extraction_" + uuid.NewString()

	first, err := repo.Create(ctx, CreatePromptTemplateRequest{
		Key:       key,
		Version:   1,
		IsActive:  active,
		Title:     "Intent Extraction v1",
		Template:  "Extract intent from: {{.text}}",
		CreatedBy: &createdBy,
	})
	require.NoError(t, err)
	assert.True(t, first.IsActive)

	second, err := repo.CreateVersion(ctx, CreatePromptTemplateVersionRequest{
		Key:       key,
		IsActive:  active,
		Title:     "Intent Extraction v2",
		Template:  "Extract intents: {{.text}}",
		CreatedBy: &createdBy,
	})
	require.NoError(t, err)
	assert.True(t, second.IsActive)

	firstReload, err := repo.GetByID(ctx, first.ID)
	require.NoError(t, err)
	assert.False(t, firstReload.IsActive)

	activeTemplate, err := repo.GetActiveByKey(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, second.ID, activeTemplate.ID)

	activated, err := repo.Activate(ctx, first.ID)
	require.NoError(t, err)
	assert.True(t, activated.IsActive)

	activeTemplate, err = repo.GetActiveByKey(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, first.ID, activeTemplate.ID)

	secondReload, err := repo.GetByID(ctx, second.ID)
	require.NoError(t, err)
	assert.False(t, secondReload.IsActive)
}
