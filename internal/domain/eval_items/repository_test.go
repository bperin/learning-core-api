package eval_items_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/eval_items"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
)

func setupTestRepo(t *testing.T) (eval_items.Repository, *store.Queries, func()) {
	tx, cleanup := testutil.NewTestTx(t)
	queries := store.New(tx)
	repo := eval_items.NewRepository(queries)
	return repo, queries, cleanup
}

func createTestUser(t *testing.T, queries *store.Queries) uuid.UUID {
	userID := uuid.New()
	_, err := queries.CreateUser(context.Background(), store.CreateUserParams{
		ID:        userID,
		Email:     "test@example.com",
		Password:  "hashedpassword",
		IsAdmin:   false,
		IsLearner: true,
		IsTeacher: false,
	})
	require.NoError(t, err)
	return userID
}

func createTestEval(t *testing.T, queries *store.Queries, userID uuid.UUID) uuid.UUID {
	eval, err := queries.CreateEval(context.Background(), store.CreateEvalParams{
		Title:  "Test Evaluation",
		Status: "draft",
		UserID: userID,
	})
	require.NoError(t, err)
	return eval.ID
}

func TestEvalItemRepository_Create(t *testing.T) {
	repo, queries, cleanup := setupTestRepo(t)
	defer cleanup()

	userID := createTestUser(t, queries)
	evalID := createTestEval(t, queries, userID)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:      evalID,
			Prompt:      "What is 2 + 2?",
			Options:     []string{"3", "4", "5", "6"},
			CorrectIdx:  1,
			Hint:        stringPtr("Think about basic addition"),
			Explanation: stringPtr("2 + 2 equals 4"),
			Metadata:    json.RawMessage(`{"difficulty":"easy"}`),
		}

		item, err := repo.Create(ctx, req)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, item.ID)
		assert.Equal(t, evalID, item.EvalID)
		assert.Equal(t, "What is 2 + 2?", item.Prompt)
		assert.Equal(t, []string{"3", "4", "5", "6"}, item.Options)
		assert.Equal(t, int32(1), item.CorrectIdx)
		assert.NotNil(t, item.Hint)
		assert.Equal(t, "Think about basic addition", *item.Hint)
		assert.NotNil(t, item.Explanation)
		assert.Equal(t, "2 + 2 equals 4", *item.Explanation)
		assert.NotNil(t, item.Metadata)
		assert.Equal(t, "easy", item.Metadata["difficulty"])
	})

	t.Run("validation error - empty prompt", func(t *testing.T) {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:     evalID,
			Prompt:     "",
			Options:    []string{"A", "B"},
			CorrectIdx: 0,
		}

		_, err := repo.Create(ctx, req)
		assert.Error(t, err)
		assert.True(t, eval_items.IsValidationError(err))
	})

	t.Run("validation error - insufficient options", func(t *testing.T) {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:     evalID,
			Prompt:     "Test question?",
			Options:    []string{"A"},
			CorrectIdx: 0,
		}

		_, err := repo.Create(ctx, req)
		assert.Error(t, err)
		assert.True(t, eval_items.IsValidationError(err))
	})

	t.Run("validation error - invalid correct index", func(t *testing.T) {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:     evalID,
			Prompt:     "Test question?",
			Options:    []string{"A", "B"},
			CorrectIdx: 5,
		}

		_, err := repo.Create(ctx, req)
		assert.Error(t, err)
		assert.True(t, eval_items.IsValidationError(err))
	})
}

func TestEvalItemRepository_GetByID(t *testing.T) {
	repo, queries, cleanup := setupTestRepo(t)
	defer cleanup()

	userID := createTestUser(t, queries)
	evalID := createTestEval(t, queries, userID)

	ctx := context.Background()

	// Create a test item
	req := &eval_items.CreateEvalItemRequest{
		EvalID:     evalID,
		Prompt:     "Test question?",
		Options:    []string{"A", "B", "C"},
		CorrectIdx: 1,
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	t.Run("successful retrieval", func(t *testing.T) {
		item, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, item.ID)
		assert.Equal(t, created.Prompt, item.Prompt)
		assert.Equal(t, created.Options, item.Options)
		assert.Equal(t, created.CorrectIdx, item.CorrectIdx)
	})

	t.Run("not found error", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := repo.GetByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.True(t, eval_items.IsNotFoundError(err))
	})
}

func TestEvalItemRepository_GetByEvalID(t *testing.T) {
	repo, queries, cleanup := setupTestRepo(t)
	defer cleanup()

	userID := createTestUser(t, queries)
	evalID := createTestEval(t, queries, userID)

	ctx := context.Background()

	// Create multiple test items
	for i := 0; i < 3; i++ {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:     evalID,
			Prompt:     fmt.Sprintf("Question %d?", i+1),
			Options:    []string{"A", "B"},
			CorrectIdx: 0,
		}
		_, err := repo.Create(ctx, req)
		require.NoError(t, err)
	}

	t.Run("successful retrieval", func(t *testing.T) {
		items, err := repo.GetByEvalID(ctx, evalID)
		require.NoError(t, err)
		assert.Len(t, items, 3)

		prompts := make(map[string]struct{}, len(items))
		for _, item := range items {
			assert.Equal(t, evalID, item.EvalID)
			prompts[item.Prompt] = struct{}{}
		}
		for i := 0; i < 3; i++ {
			_, ok := prompts[fmt.Sprintf("Question %d?", i+1)]
			assert.True(t, ok)
		}
	})

	t.Run("empty result for non-existent eval", func(t *testing.T) {
		nonExistentEvalID := uuid.New()
		items, err := repo.GetByEvalID(ctx, nonExistentEvalID)
		require.NoError(t, err)
		assert.Len(t, items, 0)
	})
}

func TestEvalItemRepository_List(t *testing.T) {
	repo, queries, cleanup := setupTestRepo(t)
	defer cleanup()

	userID := createTestUser(t, queries)
	evalID := createTestEval(t, queries, userID)

	ctx := context.Background()

	// Create test items
	for i := 0; i < 5; i++ {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:     evalID,
			Prompt:     fmt.Sprintf("Question %d?", i+1),
			Options:    []string{"A", "B"},
			CorrectIdx: 0,
		}
		_, err := repo.Create(ctx, req)
		require.NoError(t, err)
	}

	t.Run("list with pagination", func(t *testing.T) {
		listReq := &eval_items.ListEvalItemsRequest{
			Limit:  3,
			Offset: 0,
		}

		items, err := repo.List(ctx, listReq)
		require.NoError(t, err)
		assert.Len(t, items, 3)
	})

	t.Run("list with eval filter", func(t *testing.T) {
		listReq := &eval_items.ListEvalItemsRequest{
			EvalID: &evalID,
			Limit:  10,
			Offset: 0,
		}

		items, err := repo.List(ctx, listReq)
		require.NoError(t, err)
		assert.Len(t, items, 5)

		for _, item := range items {
			assert.Equal(t, evalID, item.EvalID)
		}
	})
}

func TestEvalItemRepository_Count(t *testing.T) {
	repo, queries, cleanup := setupTestRepo(t)
	defer cleanup()

	userID := createTestUser(t, queries)
	evalID := createTestEval(t, queries, userID)

	ctx := context.Background()

	// Initially should be 0
	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create test items
	for i := 0; i < 3; i++ {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:     evalID,
			Prompt:     fmt.Sprintf("Question %d?", i+1),
			Options:    []string{"A", "B"},
			CorrectIdx: 0,
		}
		_, err := repo.Create(ctx, req)
		require.NoError(t, err)
	}

	// Should now be 3
	count, err = repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestEvalItemRepository_CountByEvalID(t *testing.T) {
	repo, queries, cleanup := setupTestRepo(t)
	defer cleanup()

	userID := createTestUser(t, queries)
	evalID := createTestEval(t, queries, userID)

	ctx := context.Background()

	// Initially should be 0
	count, err := repo.CountByEvalID(ctx, evalID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create test items
	for i := 0; i < 2; i++ {
		req := &eval_items.CreateEvalItemRequest{
			EvalID:     evalID,
			Prompt:     fmt.Sprintf("Question %d?", i+1),
			Options:    []string{"A", "B"},
			CorrectIdx: 0,
		}
		_, err := repo.Create(ctx, req)
		require.NoError(t, err)
	}

	// Should now be 2
	count, err = repo.CountByEvalID(ctx, evalID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
