package documents_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/documents"
	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"
)

func setupTestDB(t *testing.T) (*store.Queries, func()) {
	t.Helper()

	db := testutil.NewTestDB(t)
	queries := store.New(db)

	cleanup := func() {
		db.Close()
	}

	return queries, cleanup
}

func seedTenant(ctx context.Context, t *testing.T, q *store.Queries) uuid.UUID {
	t.Helper()

	tenant, err := q.CreateTenant(ctx, store.CreateTenantParams{
		Name:     "Test Tenant",
		IsActive: true,
	})
	require.NoError(t, err)

	return tenant.ID
}

func seedUser(ctx context.Context, t *testing.T, q *store.Queries, tenantID uuid.UUID) uuid.UUID {
	t.Helper()

	user, err := q.CreateUser(ctx, store.CreateUserParams{
		TenantID: tenantID,
		Email:    "user@example.com",
	})
	require.NoError(t, err)

	return user.ID
}

func seedSubject(ctx context.Context, t *testing.T, q *store.Queries, userID uuid.UUID) uuid.UUID {
	t.Helper()

	subject, err := q.CreateSubject(ctx, store.CreateSubjectParams{
		UserID:      userID,
		Name:        "Test Subject",
		Description: sql.NullString{String: "Test Description", Valid: true},
	})
	require.NoError(t, err)

	return subject.ID
}

func seedFileSearchStore(ctx context.Context, t *testing.T, q *store.Queries, subjectID uuid.UUID) uuid.UUID {
	t.Helper()

	storeName := "store-" + uuid.NewString()
	created, err := q.CreateFileSearchStore(ctx, store.CreateFileSearchStoreParams{
		SubjectID:      subjectID,
		StoreName:      storeName,
		DisplayName:    sql.NullString{String: "Test Store", Valid: true},
		ChunkingConfig: json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	return created.ID
}

func seedSubjectAndStore(ctx context.Context, t *testing.T, q *store.Queries) (uuid.UUID, uuid.UUID) {
	t.Helper()

	tenantID := seedTenant(ctx, t, q)
	userID := seedUser(ctx, t, q, tenantID)
	subjectID := seedSubject(ctx, t, q, userID)
	storeID := seedFileSearchStore(ctx, t, q, subjectID)

	return subjectID, storeID
}

func TestRepository_CreateAndGetDocument(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := documents.NewRepository(queries)
	ctx := context.Background()

	subjectID, storeID := seedSubjectAndStore(ctx, t, queries)
	metadata := json.RawMessage(`{"source":"test"}`)

	created, err := repo.Create(ctx, documents.Document{
		SubjectID: subjectID,
		StoreID:   storeID,
		Title:     "Test Document",
		SourceURI: "s3://bucket/document.pdf",
		SHA256:    "abc123",
		Metadata:  metadata,
		FileName:  "document.pdf",
		DocName:   "Document",
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, subjectID, created.SubjectID)
	assert.Equal(t, storeID, created.StoreID)
	assert.Equal(t, "Test Document", created.Title)
	assert.Equal(t, "s3://bucket/document.pdf", created.SourceURI)
	assert.Equal(t, "abc123", created.SHA256)
	assert.JSONEq(t, string(metadata), string(created.Metadata))
	assert.Equal(t, "document.pdf", created.FileName)
	assert.Equal(t, "Document", created.DocName)
	assert.Nil(t, created.IndexedAt)

	fetched, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.SourceURI, fetched.SourceURI)
	assert.JSONEq(t, string(metadata), string(fetched.Metadata))
}

func TestRepository_GetBySubjectAndSourceURI(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := documents.NewRepository(queries)
	ctx := context.Background()

	subjectID, storeID := seedSubjectAndStore(ctx, t, queries)
	sourceURI := "s3://bucket/lookup.pdf"

	created, err := repo.Create(ctx, documents.Document{
		SubjectID: subjectID,
		StoreID:   storeID,
		Title:     "Lookup Document",
		SourceURI: sourceURI,
		Metadata:  json.RawMessage(`{"kind":"lookup"}`),
	})
	require.NoError(t, err)

	fetched, err := repo.GetBySubjectAndSourceURI(ctx, subjectID, sourceURI)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.SourceURI, fetched.SourceURI)
}

func TestRepository_ListBySubject(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := documents.NewRepository(queries)
	ctx := context.Background()

	subjectID, storeID := seedSubjectAndStore(ctx, t, queries)

	_, err := repo.Create(ctx, documents.Document{
		SubjectID: subjectID,
		StoreID:   storeID,
		Title:     "Document One",
		SourceURI: "s3://bucket/doc-one.pdf",
		Metadata:  json.RawMessage(`{"rank":1}`),
	})
	require.NoError(t, err)

	_, err = repo.Create(ctx, documents.Document{
		SubjectID: subjectID,
		StoreID:   storeID,
		Title:     "Document Two",
		SourceURI: "s3://bucket/doc-two.pdf",
		Metadata:  json.RawMessage(`{"rank":2}`),
	})
	require.NoError(t, err)

	list, err := repo.ListBySubject(ctx, subjectID)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestRepository_UpdateDocument(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := documents.NewRepository(queries)
	ctx := context.Background()

	subjectID, storeID := seedSubjectAndStore(ctx, t, queries)
	created, err := repo.Create(ctx, documents.Document{
		SubjectID: subjectID,
		StoreID:   storeID,
		Title:     "Original Title",
		SourceURI: "s3://bucket/update.pdf",
		Metadata:  json.RawMessage(`{"status":"draft"}`),
	})
	require.NoError(t, err)

	newTitle := "Updated Title"
	newMetadata := json.RawMessage(`{"status":"final"}`)
	indexedAt := time.Now().UTC()

	err = repo.Update(ctx, created.ID, &newTitle, newMetadata, &indexedAt)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, newTitle, updated.Title)
	assert.JSONEq(t, string(newMetadata), string(updated.Metadata))
	require.NotNil(t, updated.IndexedAt)
	assert.True(t, updated.IndexedAt.Equal(indexedAt))
}

func TestRepository_DeleteDocument(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	repo := documents.NewRepository(queries)
	ctx := context.Background()

	subjectID, storeID := seedSubjectAndStore(ctx, t, queries)
	created, err := repo.Create(ctx, documents.Document{
		SubjectID: subjectID,
		StoreID:   storeID,
		Title:     "Delete Document",
		SourceURI: "s3://bucket/delete.pdf",
		Metadata:  json.RawMessage(`{"delete":true}`),
	})
	require.NoError(t, err)

	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
}
