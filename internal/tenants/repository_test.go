package tenants_test

import (
	"context"
	"learning-core-api/internal/store"
	"learning-core-api/internal/tenants"
	"learning-core-api/internal/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (tenants.Repository, context.Context, func()) {
	ctx := context.Background()
	db, cleanup := testutil.StartPostgres(ctx)

	err := testutil.Migrate(db)
	require.NoError(t, err)

	queries := store.New(db)
	repo := tenants.NewRepository(queries)

	return repo, ctx, cleanup
}

func TestTenantRepository_Create(t *testing.T) {
	repo, ctx, cleanup := setupTest(t)
	defer cleanup()

	tenant := tenants.Tenant{
		Name: "Test Tenant",
	}

	created, err := repo.Create(ctx, tenant)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, tenant.Name, created.Name)
	assert.True(t, created.IsActive)
}

func TestTenantRepository_GetByID(t *testing.T) {
	repo, ctx, cleanup := setupTest(t)
	defer cleanup()

	tenant := tenants.Tenant{
		Name: "Test Get",
	}

	created, _ := repo.Create(ctx, tenant)

	retrieved, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
}
