package evals_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/evals"
	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"
)

func setupTestDB(t *testing.T) (*sql.DB, *store.Queries) {
	db := testutil.NewTestDB(t)
	queries := store.New(db)

	return db, queries
}

func TestSuiteRepository(t *testing.T) {
	db, queries := setupTestDB(t)
	defer db.Close()

	repo := evals.NewSuiteRepository(queries)
	ctx := context.Background()

	t.Run("Create and GetByID", func(t *testing.T) {
		expectedSuite := evals.Suite{
			Name:        "Test Suite",
			Description: "Quality gate for MCQ items",
		}

		createdSuite, err := repo.Create(ctx, expectedSuite)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, createdSuite.ID)

		assert.Equal(t, expectedSuite.Name, createdSuite.Name)
		assert.Equal(t, expectedSuite.Description, createdSuite.Description)

		retrievedSuite, err := repo.GetByID(ctx, createdSuite.ID)
		require.NoError(t, err)
		assert.Equal(t, createdSuite.ID, retrievedSuite.ID)
		assert.Equal(t, createdSuite.Name, retrievedSuite.Name)
		assert.Equal(t, createdSuite.Description, retrievedSuite.Description)
		assert.WithinDuration(t, createdSuite.CreatedAt, retrievedSuite.CreatedAt, time.Second)
	})

	t.Run("GetByName", func(t *testing.T) {
		expectedSuite := evals.Suite{
			Name:        "Named Suite",
			Description: "Gate for flashcards",
		}

		createdSuite, err := repo.Create(ctx, expectedSuite)
		require.NoError(t, err)

		retrievedSuite, err := repo.GetByName(ctx, createdSuite.Name)
		require.NoError(t, err)
		assert.Equal(t, createdSuite.ID, retrievedSuite.ID)
		assert.Equal(t, createdSuite.Name, retrievedSuite.Name)
	})

	t.Run("List", func(t *testing.T) {
		suite1 := evals.Suite{
			Name:        "Suite 1",
			Description: "First suite",
		}
		suite2 := evals.Suite{
			Name:        "Suite 2",
			Description: "Second suite",
		}

		createdSuite1, err := repo.Create(ctx, suite1)
		require.NoError(t, err)
		createdSuite2, err := repo.Create(ctx, suite2)
		require.NoError(t, err)

		allSuites, err := repo.List(ctx)
		require.NoError(t, err)

		var foundSuite1, foundSuite2 bool
		for _, suite := range allSuites {
			if suite.ID == createdSuite1.ID {
				foundSuite1 = true
				assert.Equal(t, suite1.Name, suite.Name)
				assert.Equal(t, suite1.Description, suite.Description)
			}
			if suite.ID == createdSuite2.ID {
				foundSuite2 = true
				assert.Equal(t, suite2.Name, suite.Name)
				assert.Equal(t, suite2.Description, suite.Description)
			}
		}

		assert.True(t, foundSuite1, "Suite 1 should be in the list")
		assert.True(t, foundSuite2, "Suite 2 should be in the list")
	})

	t.Run("Update", func(t *testing.T) {
		originalSuite := evals.Suite{
			Name:        "Update Test Suite",
			Description: "Original description",
		}

		createdSuite, err := repo.Create(ctx, originalSuite)
		require.NoError(t, err)

		updatedSuite, err := repo.Update(ctx, createdSuite.ID, "Updated Name", "Updated description")
		require.NoError(t, err)

		assert.Equal(t, createdSuite.ID, updatedSuite.ID)
		assert.Equal(t, "Updated Name", updatedSuite.Name)
		assert.Equal(t, "Updated description", updatedSuite.Description)
	})

	t.Run("Delete", func(t *testing.T) {
		suiteToDelete := evals.Suite{
			Name:        "Delete Test Suite",
			Description: "Suite to delete",
		}

		createdSuite, err := repo.Create(ctx, suiteToDelete)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, createdSuite.ID)
		require.NoError(t, err)

		err = repo.Delete(ctx, createdSuite.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, createdSuite.ID)
		assert.Error(t, err)
	})

	t.Run("GetNonExistentSuite", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := repo.GetByID(ctx, nonExistentID)
		assert.Error(t, err)
	})

	t.Run("GetSuiteByNameNotFound", func(t *testing.T) {
		_, err := repo.GetByName(ctx, "non-existent-suite-name")
		assert.Error(t, err)
	})
}
