package taxonomy

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/documents"
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

func createTestUser(t *testing.T, db *sql.Tx) uuid.UUID {
	t.Helper()

	ctx := context.Background()
	userID := uuid.New()
	_, err := db.ExecContext(ctx,
		"INSERT INTO users (id, email, password, is_admin, is_learner, is_teacher) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "taxonomy@example.com", "password123", false, true, false)
	require.NoError(t, err)

	return userID
}

func createTestDocument(t *testing.T, repo documents.Repository, userID uuid.UUID, filename string) *documents.Document {
	t.Helper()

	ctx := context.Background()
	doc, err := repo.Create(ctx, documents.CreateDocumentRequest{
		Filename:  filename,
		RagStatus: documents.RagStatusPending,
		UserID:    userID,
	})
	require.NoError(t, err)

	return doc
}

func TestTaxonomyRepository_CreateActivate(t *testing.T) {
	db, _, repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	createTestUser(t, db)

	first, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Biology",
		Path:     "biology",
		Depth:    0,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
	})
	require.NoError(t, err)
	assert.True(t, first.IsActive)
	assert.Equal(t, int32(1), first.Version)

	second, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Biology",
		Path:     "biology",
		Depth:    0,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
	})
	require.NoError(t, err)
	assert.True(t, second.IsActive)
	assert.Equal(t, int32(2), second.Version)

	firstReload, err := repo.GetByID(ctx, first.ID)
	require.NoError(t, err)
	assert.False(t, firstReload.IsActive)

	active, err := repo.GetActiveByPath(ctx, "biology")
	require.NoError(t, err)
	assert.Equal(t, second.ID, active.ID)

	activated, err := repo.Activate(ctx, first.ID)
	require.NoError(t, err)
	assert.True(t, activated.IsActive)

	active, err = repo.GetActiveByPath(ctx, "biology")
	require.NoError(t, err)
	assert.Equal(t, first.ID, active.ID)
}

func TestTaxonomyRepository_ListDocumentsByPrefix(t *testing.T) {
	db, queries, repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	userID := createTestUser(t, db)
	docRepo := documents.NewRepository(queries)

	docBioRoot := createTestDocument(t, docRepo, userID, "biology-root.pdf")
	docBioDNA := createTestDocument(t, docRepo, userID, "biology-dna.pdf")
	docBioRNA := createTestDocument(t, docRepo, userID, "biology-rna.pdf")
	docBioEcology := createTestDocument(t, docRepo, userID, "biology-ecology.pdf")
	docBioPending := createTestDocument(t, docRepo, userID, "biology-pending.pdf")
	docHistory1 := createTestDocument(t, docRepo, userID, "history-ancient.pdf")
	docHistory2 := createTestDocument(t, docRepo, userID, "history-modern.pdf")

	root, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Biology",
		Path:     "biology",
		Depth:    0,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
	})
	require.NoError(t, err)

	molecular, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Molecular Biology",
		Path:     "biology/molecular-biology",
		Depth:    1,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
		ParentID: &root.ID,
	})
	require.NoError(t, err)

	dna, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "DNA",
		Path:     "biology/molecular-biology/dna",
		Depth:    2,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
		ParentID: &molecular.ID,
	})
	require.NoError(t, err)

	rna, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "RNA",
		Path:     "biology/molecular-biology/rna",
		Depth:    2,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
		ParentID: &molecular.ID,
	})
	require.NoError(t, err)

	ecology, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Ecology",
		Path:     "biology/ecology",
		Depth:    1,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
		ParentID: &root.ID,
	})
	require.NoError(t, err)

	conservation, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Conservation",
		Path:     "biology/ecology/conservation",
		Depth:    2,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
		ParentID: &ecology.ID,
	})
	require.NoError(t, err)

	pending, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Biology Pending",
		Path:     "biology/pending",
		Depth:    1,
		State:    string(TaxonomyStateAIGenerated),
		IsActive: true,
		ParentID: &root.ID,
	})
	require.NoError(t, err)

	historyRoot, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "History",
		Path:     "history",
		Depth:    0,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
	})
	require.NoError(t, err)

	ancient, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Ancient History",
		Path:     "history/ancient",
		Depth:    1,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
		ParentID: &historyRoot.ID,
	})
	require.NoError(t, err)

	modern, err := repo.CreateNode(ctx, CreateTaxonomyNodeRequest{
		Name:     "Modern History",
		Path:     "history/modern",
		Depth:    1,
		State:    string(TaxonomyStateApproved),
		IsActive: true,
		ParentID: &historyRoot.ID,
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docBioRoot.ID,
		TaxonomyNodeID: root.ID,
		State:          string(TaxonomyStateApproved),
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docBioDNA.ID,
		TaxonomyNodeID: dna.ID,
		State:          string(TaxonomyStateApproved),
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docBioRNA.ID,
		TaxonomyNodeID: rna.ID,
		State:          string(TaxonomyStateApproved),
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docBioEcology.ID,
		TaxonomyNodeID: conservation.ID,
		State:          string(TaxonomyStateApproved),
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docBioPending.ID,
		TaxonomyNodeID: root.ID,
		State:          string(TaxonomyStateAIGenerated),
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docBioPending.ID,
		TaxonomyNodeID: pending.ID,
		State:          string(TaxonomyStateApproved),
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docHistory1.ID,
		TaxonomyNodeID: ancient.ID,
		State:          string(TaxonomyStateApproved),
	})
	require.NoError(t, err)

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     docHistory2.ID,
		TaxonomyNodeID: modern.ID,
		State:          string(TaxonomyStateApproved),
	})
	require.NoError(t, err)

	biologyDocs, err := repo.ListDocumentsByPrefix(ctx, "biology")
	require.NoError(t, err)
	if assert.Len(t, biologyDocs, 4) {
		ids := map[uuid.UUID]struct{}{
			biologyDocs[0].ID: {},
			biologyDocs[1].ID: {},
			biologyDocs[2].ID: {},
			biologyDocs[3].ID: {},
		}
		assert.Contains(t, ids, docBioRoot.ID)
		assert.Contains(t, ids, docBioDNA.ID)
		assert.Contains(t, ids, docBioRNA.ID)
		assert.Contains(t, ids, docBioEcology.ID)
	}

	rnaDocs, err := repo.ListDocumentsByPrefix(ctx, "biology/molecular-biology/rna")
	require.NoError(t, err)
	if assert.Len(t, rnaDocs, 1) {
		assert.Equal(t, docBioRNA.ID, rnaDocs[0].ID)
	}

	molecularDocs, err := repo.ListDocumentsByPrefix(ctx, "biology/molecular-biology")
	require.NoError(t, err)
	if assert.Len(t, molecularDocs, 2) {
		ids := map[uuid.UUID]struct{}{
			molecularDocs[0].ID: {},
			molecularDocs[1].ID: {},
		}
		assert.Contains(t, ids, docBioDNA.ID)
		assert.Contains(t, ids, docBioRNA.ID)
	}

	historyDocs, err := repo.ListDocumentsByPrefix(ctx, "history")
	require.NoError(t, err)
	if assert.Len(t, historyDocs, 2) {
		ids := map[uuid.UUID]struct{}{
			historyDocs[0].ID: {},
			historyDocs[1].ID: {},
		}
		assert.Contains(t, ids, docHistory1.ID)
		assert.Contains(t, ids, docHistory2.ID)
	}
}
