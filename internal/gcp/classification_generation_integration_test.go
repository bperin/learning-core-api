package gcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"

	"learning-core-api/internal/domain/artifacts"
	"learning-core-api/internal/domain/generation"
	"learning-core-api/internal/domain/taxonomy"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
	"learning-core-api/internal/utils"
)

type fileSearchToolSeed struct {
	StoreNames     []string `json:"store_names"`
	MetadataFilter string   `json:"metadata_filter"`
}

func TestClassificationGenerationWithFileSearchIntegration(t *testing.T) {
	storeName := os.Getenv("FILE_STORE_NAME")
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if storeName == "" || apiKey == "" {
		_ = godotenv.Load()
		storeName = os.Getenv("FILE_STORE_NAME")
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if storeName == "" || apiKey == "" {
		t.Skipf("missing FILE_STORE_NAME or GOOGLE_API_KEY (store=%t api_key=%t)", storeName != "", apiKey != "")
	}
	if os.Getenv("TEST_DB_URL") == "" {
		t.Skip("missing TEST_DB_URL for seeded integration test")
	}

	ctx := context.Background()
	db := testutil.NewTestDB(t)
	t.Cleanup(func() {
		_ = db.Close()
	})
	queries := store.New(db)

	userID := uuid.New()
	_, err := queries.CreateUser(ctx, store.CreateUserParams{
		ID:        userID,
		Email:     fmt.Sprintf("classification-test-%s@example.com", userID.String()),
		Password:  "password123",
		IsAdmin:   false,
		IsLearner: true,
		IsTeacher: false,
	})
	require.NoError(t, err)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	require.NoError(t, err)

	storeName, err = ensureFileSearchStore(ctx, client, storeName)
	require.NoError(t, err)

	localPath := filepath.Join("../../test_docs", "test.pdf")
	metadataID := userID.String()

	uploadCfg := &genai.UploadToFileSearchStoreConfig{
		MIMEType:    "application/pdf",
		DisplayName: fmt.Sprintf("classification-%s", metadataID),
		CustomMetadata: []*genai.CustomMetadata{
			{
				Key:         "user_id",
				StringValue: metadataID,
			},
		},
	}

	operation, err := client.FileSearchStores.UploadToFileSearchStoreFromPath(ctx, localPath, storeName, uploadCfg)
	require.NoError(t, err)

	deadline := time.Now().Add(2 * time.Minute)
	for !operation.Done && time.Now().Before(deadline) {
		time.Sleep(3 * time.Second)
		operation, err = client.Operations.GetUploadToFileSearchStoreOperation(ctx, operation, nil)
		require.NoError(t, err)
	}
	require.True(t, operation.Done, "upload operation did not complete in time")
	require.NotNil(t, operation.Response)
	require.NotEmpty(t, operation.Response.DocumentName)

	doc, err := queries.CreateDocument(ctx, store.CreateDocumentParams{
		Filename:          "test.pdf",
		Title:             sql.NullString{String: "Classification Test", Valid: true},
		MimeType:          sql.NullString{String: "application/pdf", Valid: true},
		FileStoreName:     sql.NullString{String: storeName, Valid: true},
		FileStoreFileName: sql.NullString{String: operation.Response.DocumentName, Valid: true},
		RagStatus:         "READY",
		UserID:            userID,
	})
	require.NoError(t, err)

	activeModel, err := queries.GetActiveModelConfig(ctx)
	require.NoError(t, err)
	activeSystemInstruction, err := queries.GetActiveSystemInstruction(ctx)
	require.NoError(t, err)
	activePrompt, err := queries.GetPromptTemplateByGenerationType(ctx, utils.GenerationTypeClassification.DB())
	require.NoError(t, err)
	activeSchema, err := queries.GetActiveSchemaTemplateByGenerationType(ctx, utils.GenerationTypeClassification.DB())
	require.NoError(t, err)

	toolSeed := fileSearchToolSeed{
		StoreNames:     []string{storeName},
		MetadataFilter: fmt.Sprintf(`user_id = "%s"`, metadataID),
	}
	toolConfig, err := json.Marshal(toolSeed)
	require.NoError(t, err)

	schemaPretty, err := json.MarshalIndent(json.RawMessage(activeSchema.SchemaJson), "", "  ")
	require.NoError(t, err)
	t.Logf("system instructions: %s", activeSystemInstruction.Text)
	t.Logf("prompt template (generation_type=%s, version=%d): %s", activePrompt.GenerationType, activePrompt.Version, activePrompt.Template)
	t.Logf("schema template (generation_type=%s, version=%d):\n%s", activeSchema.GenerationType, activeSchema.Version, string(schemaPretty))

	generator, err := NewGenerationServiceFromAPIKey(ctx, apiKey)
	require.NoError(t, err)
	artifactService := artifacts.NewService(db)
	genService, err := generation.NewService(db, artifactService, generator)
	require.NoError(t, err)

	resp, err := genService.Generate(ctx, generation.GenerateRequest{
		UserID: userID,
		Target: generation.Target{
			DocumentID: &doc.ID,
		},
		Instructions: generation.Instructions{
			SystemInstructionID: &activeSystemInstruction.ID,
			GenerationType:      utils.GenerationTypeClassification.String(),
		},
		Output: generation.OutputConfig{
			GenerationType: utils.GenerationTypeClassification.String(),
			Format:         "json",
		},
		Tools: []generation.ToolConfig{
			{
				Type:   "file_search",
				Config: toolConfig,
			},
		},
		ModelConfigID: activeModel.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.OutputText)
	var outputPretty string
	if len(resp.OutputJSON) > 0 {
		pretty, err := json.MarshalIndent(json.RawMessage(resp.OutputJSON), "", "  ")
		require.NoError(t, err)
		outputPretty = string(pretty)
	} else {
		outputPretty = resp.OutputText
	}
	t.Logf("classification output:\n%s", outputPretty)

	taxonomyRepo := taxonomy.NewRepository(queries)
	createdNodes, err := taxonomy.IngestGeneratedTaxonomy(ctx, taxonomyRepo, doc.ID, userID, resp.OutputJSON)
	require.NoError(t, err)
	require.NotEmpty(t, createdNodes)
	for _, node := range createdNodes {
		t.Logf("taxonomy node: path=%s depth=%d id=%s", node.Path, node.Depth, node.ID)
	}

	deepest := createdNodes[0]
	for _, node := range createdNodes[1:] {
		if node.Depth > deepest.Depth {
			deepest = node
		}
	}

	trace, err := traceTaxonomyToRoot(ctx, taxonomyRepo, createdNodes, deepest)
	require.NoError(t, err)
	require.NotEmpty(t, trace)
	t.Logf("taxonomy trace (leaf -> root): %s", strings.Join(trace, " -> "))
}

func traceTaxonomyToRoot(ctx context.Context, repo taxonomy.Repository, created []*taxonomy.TaxonomyNode, leaf *taxonomy.TaxonomyNode) ([]string, error) {
	index := make(map[uuid.UUID]*taxonomy.TaxonomyNode, len(created))
	for _, node := range created {
		index[node.ID] = node
	}

	var trace []string
	current := leaf
	for current != nil {
		trace = append(trace, current.Path)
		if current.ParentID == nil || *current.ParentID == uuid.Nil {
			break
		}
		if next, ok := index[*current.ParentID]; ok {
			current = next
			continue
		}
		fetched, err := repo.GetByID(ctx, *current.ParentID)
		if err != nil {
			return nil, err
		}
		current = fetched
	}

	return trace, nil
}
