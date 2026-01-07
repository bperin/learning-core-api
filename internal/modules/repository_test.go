package modules_test

import (
	"context"
	"learning-core-api/internal/modules"
	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (modules.Repository, *store.Queries, context.Context, func()) {
	ctx := context.Background()
	db, cleanup := testutil.StartPostgres(ctx)

	err := testutil.Migrate(db)
	require.NoError(t, err)

	queries := store.New(db)
	repo := modules.NewRepository(queries)

	return repo, queries, ctx, cleanup
}

func createTestTenant(t *testing.T, repo modules.Repository, ctx context.Context, db *store.Queries) uuid.UUID {
	params := store.CreateTenantParams{
		Name:     "Test Tenant",
		IsActive: true,
	}
	tenant, err := db.CreateTenant(ctx, params)
	require.NoError(t, err)
	return tenant.ID
}

func TestModuleRepository_Lifecycle(t *testing.T) {
	repo, queries, ctx, cleanup := setupTest(t)
	defer cleanup()

	// 1. Need a tenant first
	tenantID := createTestTenant(t, repo, ctx, queries)

	// 2. Test Create
	module := modules.Module{
		TenantID:    tenantID,
		Title:       "Test Title",
		Name:        "Test Module",
		Description: "A test module",
	}

	created, err := repo.Create(ctx, module)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, module.Title, created.Title)
	assert.Equal(t, module.Name, created.Name)
	assert.Equal(t, module.Description, created.Description)
	assert.Equal(t, tenantID, created.TenantID)

	// 3. Test GetByID
	retrieved, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)

	// 4. Test GetByName
	byName, err := repo.GetByName(ctx, tenantID, created.Name)
	require.NoError(t, err)
	assert.Equal(t, created.ID, byName.ID)

	// 5. Test ListByTenant
	list, err := repo.ListByTenant(ctx, tenantID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, created.ID, list[0].ID)

	// 6. Test Update
	newName := "Updated Name"
	newDesc := "Updated Description"
	err = repo.Update(ctx, created.ID, &newName, &newDesc)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, newDesc, updated.Description)

	// 7. Test Delete
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
}
