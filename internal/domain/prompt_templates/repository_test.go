package prompt_templates

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
	"learning-core-api/internal/utils"
)

func setupTestRepo(t *testing.T) (*sql.Tx, *store.Queries, Repository, func()) {
	t.Helper()

	tx, cleanup := testutil.NewTestTx(t)
	queries := store.New(tx)
	repo := NewRepository(queries)

	return tx, queries, repo, cleanup
}

func TestPromptTemplateRepository_CreateActivate(t *testing.T) {
	_, queries, repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	active := true
	createdBy := "tester"
	generationType := utils.GenerationTypeClassification.String()
	baseVersion := int32(0)
	existing, err := queries.GetPromptTemplatesByGenerationType(ctx, store.GenerationType(generationType))
	require.NoError(t, err)
	for _, tmpl := range existing {
		if tmpl.Version > baseVersion {
			baseVersion = tmpl.Version
		}
	}

	first, err := repo.Create(ctx, CreatePromptTemplateRequest{
		GenerationType: generationType,
		Version:        baseVersion + 1,
		IsActive:       active,
		Title:          "Classification v1",
		Template:       "Classify: {{.text}}",
		CreatedBy:      &createdBy,
	})
	require.NoError(t, err)
	assert.True(t, first.IsActive)

	second, err := repo.CreateVersion(ctx, CreatePromptTemplateVersionRequest{
		GenerationType: generationType,
		IsActive:       active,
		Title:          "Classification v2",
		Template:       "Classify text: {{.text}}",
		CreatedBy:      &createdBy,
	})
	require.NoError(t, err)
	assert.True(t, second.IsActive)

	firstReload, err := repo.GetByID(ctx, first.ID)
	require.NoError(t, err)
	assert.False(t, firstReload.IsActive)

	activeTemplate, err := repo.GetActiveByGenerationType(ctx, generationType)
	require.NoError(t, err)
	assert.Equal(t, second.ID, activeTemplate.ID)

	activated, err := repo.Activate(ctx, first.ID)
	require.NoError(t, err)
	assert.True(t, activated.IsActive)

	activeTemplate, err = repo.GetActiveByGenerationType(ctx, generationType)
	require.NoError(t, err)
	assert.Equal(t, first.ID, activeTemplate.ID)

	secondReload, err := repo.GetByID(ctx, second.ID)
	require.NoError(t, err)
	assert.False(t, secondReload.IsActive)
}
