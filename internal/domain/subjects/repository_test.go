package subjects

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

func setupTestDB(t *testing.T) (*sql.Tx, *store.Queries, Repository, func()) {
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
		userID, "test@example.com", "password123", false, true, false)
	require.NoError(t, err)

	return userID
}

func TestSubjectRepository_Create(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	description := "A comprehensive mathematics course"
	req := CreateSubjectRequest{
		Name:        "Mathematics",
		Description: &description,
		UserID:      userID,
	}

	subject, err := repo.Create(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, subject)
	assert.NotEqual(t, uuid.Nil, subject.ID)
	assert.Equal(t, "Mathematics", subject.Name)
	assert.NotNil(t, subject.Description)
	assert.Equal(t, description, *subject.Description)
	assert.Equal(t, userID, subject.UserID)
	assert.False(t, subject.CreatedAt.IsZero())
	assert.False(t, subject.UpdatedAt.IsZero())
}

func TestSubjectRepository_CreateWithoutDescription(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	req := CreateSubjectRequest{
		Name:   "Science",
		UserID: userID,
	}

	subject, err := repo.Create(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, subject)
	assert.Equal(t, "Science", subject.Name)
	assert.Nil(t, subject.Description)
	assert.Equal(t, userID, subject.UserID)
}

func TestSubjectRepository_GetByID(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	// Create a subject first
	description := "History course"
	req := CreateSubjectRequest{
		Name:        "History",
		Description: &description,
		UserID:      userID,
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Get the subject by ID
	subject, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.NotNil(t, subject)
	assert.Equal(t, created.ID, subject.ID)
	assert.Equal(t, "History", subject.Name)
	assert.NotNil(t, subject.Description)
	assert.Equal(t, description, *subject.Description)
	assert.Equal(t, userID, subject.UserID)
}

func TestSubjectRepository_GetByID_NotFound(t *testing.T) {
	_, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	nonExistentID := uuid.New()

	subject, err := repo.GetByID(ctx, nonExistentID)
	assert.Error(t, err)
	assert.Equal(t, ErrSubjectNotFound, err)
	assert.Nil(t, subject)
}

func TestSubjectRepository_Update(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	// Create a subject first
	description := "Original description"
	req := CreateSubjectRequest{
		Name:        "Original Name",
		Description: &description,
		UserID:      userID,
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Update the subject
	newName := "Updated Name"
	newDescription := "Updated description"
	updateReq := UpdateSubjectRequest{
		Name:        &newName,
		Description: &newDescription,
	}

	updated, err := repo.Update(ctx, created.ID, updateReq)
	require.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, newName, updated.Name)
	assert.NotNil(t, updated.Description)
	assert.Equal(t, newDescription, *updated.Description)
	assert.Equal(t, userID, updated.UserID)
	assert.True(t, updated.UpdatedAt.After(created.UpdatedAt) || updated.UpdatedAt.Equal(created.UpdatedAt))
}

func TestSubjectRepository_UpdatePartial(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	// Create a subject first
	description := "Original description"
	req := CreateSubjectRequest{
		Name:        "Original Name",
		Description: &description,
		UserID:      userID,
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Update only the name
	newName := "Updated Name Only"
	updateReq := UpdateSubjectRequest{
		Name: &newName,
	}

	updated, err := repo.Update(ctx, created.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.NotNil(t, updated.Description)
	assert.Equal(t, description, *updated.Description) // Description should remain unchanged
}

func TestSubjectRepository_Delete(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	// Create a subject first
	req := CreateSubjectRequest{
		Name:   "To Be Deleted",
		UserID: userID,
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Delete the subject
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify it's deleted
	subject, err := repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrSubjectNotFound, err)
	assert.Nil(t, subject)
}

func TestSubjectRepository_ListByUser(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID1 := createTestUser(t, db)

	// Create second user
	userID2 := uuid.New()
	_, err := db.ExecContext(ctx,
		"INSERT INTO users (id, email, password, is_admin, is_learner, is_teacher) VALUES ($1, $2, $3, $4, $5, $6)",
		userID2, "test2@example.com", "password123", false, true, false)
	require.NoError(t, err)

	// Create subjects for user1
	subjects1 := []CreateSubjectRequest{
		{Name: "Math", UserID: userID1},
		{Name: "Science", UserID: userID1},
		{Name: "History", UserID: userID1},
	}

	for _, req := range subjects1 {
		_, err := repo.Create(ctx, req)
		require.NoError(t, err)
	}

	// Create subject for user2
	_, err = repo.Create(ctx, CreateSubjectRequest{
		Name:   "Art",
		UserID: userID2,
	})
	require.NoError(t, err)

	// List subjects for user1
	userSubjects, err := repo.ListByUser(ctx, userID1)
	require.NoError(t, err)
	assert.Len(t, userSubjects, 3)

	// Verify all subjects belong to user1
	for _, subject := range userSubjects {
		assert.Equal(t, userID1, subject.UserID)
	}

	// List subjects for user2
	user2Subjects, err := repo.ListByUser(ctx, userID2)
	require.NoError(t, err)
	assert.Len(t, user2Subjects, 1)
	assert.Equal(t, "Art", user2Subjects[0].Name)
	assert.Equal(t, userID2, user2Subjects[0].UserID)
}

func TestSubjectRepository_ListByUser_Empty(t *testing.T) {
	db, _, repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)

	// List subjects for user with no subjects
	subjects, err := repo.ListByUser(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, subjects)
}
