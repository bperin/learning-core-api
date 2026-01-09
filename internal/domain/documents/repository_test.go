package documents

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
	"learning-core-api/internal/utils"
)

func setupTestDB(t *testing.T) (*sql.DB, *store.Queries, Repository) {
	t.Helper()

	db := testutil.NewTestDB(t)
	queries := store.New(db)
	repo := NewRepository(queries)

	return db, queries, repo
}

func TestDocumentRepository_Create(t *testing.T) {
	db, _, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a test user first
	userID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", 
		userID, "test@example.com", "password123")
	require.NoError(t, err)

	req := CreateDocumentRequest{
		Filename:  "test.pdf",
		Title:     utils.StringPtr("Test Document"),
		MimeType:  utils.StringPtr("application/pdf"),
		Content:   utils.StringPtr("Test content"),
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{"math", "science"},
	}

	doc, err := repo.Create(ctx, req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, doc.ID)
	assert.Equal(t, req.Filename, doc.Filename)
	assert.Equal(t, req.Title, doc.Title)
	assert.Equal(t, req.MimeType, doc.MimeType)
	assert.Equal(t, req.Content, doc.Content)
	assert.Equal(t, req.RagStatus, doc.RagStatus)
	assert.Equal(t, req.UserID, doc.UserID)
	assert.Equal(t, req.Subjects, doc.Subjects)
	assert.False(t, doc.CreatedAt.IsZero())
	assert.False(t, doc.UpdatedAt.IsZero())

	// Clean up
	err = repo.Delete(ctx, doc.ID)
	require.NoError(t, err)
	
	// Clean up user
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	require.NoError(t, err)
}

func TestDocumentRepository_GetByID(t *testing.T) {
	db, _, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	
	// Create a test user first
	userID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", 
		userID, "test@example.com", "password123")
	require.NoError(t, err)

	// Create a document first
	req := CreateDocumentRequest{
		Filename:  "test.pdf",
		Title:     utils.StringPtr("Test Document"),
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{"math"},
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Get the document
	doc, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, doc.ID)
	assert.Equal(t, created.Filename, doc.Filename)
	assert.Equal(t, created.Title, doc.Title)

	// Test not found
	_, err = repo.GetByID(ctx, uuid.New())
	assert.Equal(t, ErrDocumentNotFound, err)

	// Clean up
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)
	
	// Clean up user
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	require.NoError(t, err)
}

func TestDocumentRepository_Update(t *testing.T) {
	db, _, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	
	// Create a test user first
	userID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", 
		userID, "test@example.com", "password123")
	require.NoError(t, err)

	// Create a document first
	req := CreateDocumentRequest{
		Filename:  "test.pdf",
		Title:     utils.StringPtr("Original Title"),
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{"math"},
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Update the document
	updateReq := UpdateDocumentRequest{
		Title:     utils.StringPtr("Updated Title"),
		RagStatus: utils.StringPtr(RagStatusReady),
		Subjects:  []string{"math", "science"},
	}

	updated, err := repo.Update(ctx, created.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Updated Title", *updated.Title)
	assert.Equal(t, RagStatusReady, updated.RagStatus)
	assert.Equal(t, []string{"math", "science"}, updated.Subjects)
	assert.True(t, updated.UpdatedAt.After(created.UpdatedAt))

	// Clean up
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)
	
	// Clean up user
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	require.NoError(t, err)
}

func TestDocumentRepository_UpdateRagStatus(t *testing.T) {
	db, _, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	
	// Create a test user first
	userID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", 
		userID, "test@example.com", "password123")
	require.NoError(t, err)

	// Create a document first
	req := CreateDocumentRequest{
		Filename:  "test.pdf",
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{},
	}

	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Update RAG status
	updated, err := repo.UpdateRagStatus(ctx, created.ID, RagStatusProcessing)
	require.NoError(t, err)
	assert.Equal(t, RagStatusProcessing, updated.RagStatus)
	assert.True(t, updated.UpdatedAt.After(created.UpdatedAt))

	// Test invalid status
	_, err = repo.UpdateRagStatus(ctx, created.ID, "invalid_status")
	assert.Equal(t, ErrInvalidRagStatus, err)

	// Clean up
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)
	
	// Clean up user
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	require.NoError(t, err)
}

func TestDocumentRepository_GetByUser(t *testing.T) {
	db, _, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	
	// Create test users first
	userID := uuid.New()
	otherUserID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", 
		userID, "test1@example.com", "password123")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", 
		otherUserID, "test2@example.com", "password123")
	require.NoError(t, err)

	// Create documents for different users
	req1 := CreateDocumentRequest{
		Filename:  "user1_doc1.pdf",
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{},
	}

	req2 := CreateDocumentRequest{
		Filename:  "user1_doc2.pdf",
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{},
	}

	req3 := CreateDocumentRequest{
		Filename:  "user2_doc1.pdf",
		RagStatus: RagStatusPending,
		UserID:    otherUserID,
		Subjects:  []string{},
	}

	doc1, err := repo.Create(ctx, req1)
	require.NoError(t, err)

	doc2, err := repo.Create(ctx, req2)
	require.NoError(t, err)

	doc3, err := repo.Create(ctx, req3)
	require.NoError(t, err)

	// Get documents for first user
	userDocs, err := repo.GetByUser(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, userDocs, 2)

	// Verify all documents belong to the user
	for _, doc := range userDocs {
		assert.Equal(t, userID, doc.UserID)
	}

	// Get documents for second user
	otherUserDocs, err := repo.GetByUser(ctx, otherUserID)
	require.NoError(t, err)
	assert.Len(t, otherUserDocs, 1)
	assert.Equal(t, otherUserID, otherUserDocs[0].UserID)

	// Clean up
	err = repo.Delete(ctx, doc1.ID)
	require.NoError(t, err)
	err = repo.Delete(ctx, doc2.ID)
	require.NoError(t, err)
	err = repo.Delete(ctx, doc3.ID)
	require.NoError(t, err)
	
	// Clean up users
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", otherUserID)
	require.NoError(t, err)
}

func TestDocumentRepository_Search(t *testing.T) {
	db, _, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	
	// Create a test user first
	userID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", 
		userID, "test@example.com", "password123")
	require.NoError(t, err)

	// Create documents with different titles
	req1 := CreateDocumentRequest{
		Filename:  "math_basics.pdf",
		Title:     utils.StringPtr("Mathematics Fundamentals"),
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{},
	}

	req2 := CreateDocumentRequest{
		Filename:  "science_intro.pdf",
		Title:     utils.StringPtr("Introduction to Science"),
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{},
	}

	req3 := CreateDocumentRequest{
		Filename:  "math_advanced.pdf",
		Title:     utils.StringPtr("Advanced Mathematics"),
		RagStatus: RagStatusPending,
		UserID:    userID,
		Subjects:  []string{},
	}

	doc1, err := repo.Create(ctx, req1)
	require.NoError(t, err)

	doc2, err := repo.Create(ctx, req2)
	require.NoError(t, err)

	doc3, err := repo.Create(ctx, req3)
	require.NoError(t, err)

	// Search for "math" - should return 2 documents
	mathDocs, err := repo.Search(ctx, "math", 10, 0)
	require.NoError(t, err)
	assert.Len(t, mathDocs, 2)

	// Search for "science" - should return 1 document
	scienceDocs, err := repo.Search(ctx, "science", 10, 0)
	require.NoError(t, err)
	assert.Len(t, scienceDocs, 1)

	// Search for non-existent term
	noDocs, err := repo.Search(ctx, "nonexistent", 10, 0)
	require.NoError(t, err)
	assert.Len(t, noDocs, 0)

	// Clean up
	err = repo.Delete(ctx, doc1.ID)
	require.NoError(t, err)
	err = repo.Delete(ctx, doc2.ID)
	require.NoError(t, err)
	err = repo.Delete(ctx, doc3.ID)
	require.NoError(t, err)
	
	// Clean up user
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	require.NoError(t, err)
}
