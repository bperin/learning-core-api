package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
)

func TestBasicSQLCOperations(t *testing.T) {
	tx, cleanup := testutil.NewTestTx(t)
	defer cleanup()

	queries := store.New(tx)
	ctx := context.Background()

	// Test basic user operations
	userID := uuid.New()
	user, err := queries.CreateUser(ctx, store.CreateUserParams{
		ID:        userID,
		Email:     "test@example.com",
		Password:  "hashedpassword",
		IsAdmin:   false,
		IsLearner: true,
		IsTeacher: false,
	})
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)

	// Test get user
	fetchedUser, err := queries.GetUser(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, fetchedUser.ID)

	// Test count users
	count, err := queries.CountUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Test basic eval operations
	eval, err := queries.CreateEval(ctx, store.CreateEvalParams{
		Title:  "Test Eval",
		Status: "draft",
		UserID: userID,
	})
	require.NoError(t, err)
	assert.Equal(t, "Test Eval", eval.Title)
	assert.Equal(t, "draft", eval.Status)

	// Test get eval
	fetchedEval, err := queries.GetEval(ctx, eval.ID)
	require.NoError(t, err)
	assert.Equal(t, eval.ID, fetchedEval.ID)

	t.Log("✅ All basic SQLC operations working correctly!")
}
