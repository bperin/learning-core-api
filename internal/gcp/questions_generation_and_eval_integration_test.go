package gcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"

	"learning-core-api/internal/domain/artifacts"
	"learning-core-api/internal/domain/eval_items"
	"learning-core-api/internal/domain/eval_results"
	"learning-core-api/internal/domain/evals"
	"learning-core-api/internal/domain/generation"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
	"learning-core-api/internal/utils"
)

// TestQuestionsGenerationAndEvalIntegration tests the full pipeline:
// 1. Generate questions from document with file search
// 2. Save questions as eval items with grounding metadata
// 3. Evaluate questions using groundedness evaluator with RAG context
func TestQuestionsGenerationAndEvalIntegration(t *testing.T) {
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

	// 1. Setup: Create user and document
	userID := uuid.New()
	_, err := queries.CreateUser(ctx, store.CreateUserParams{
		ID:        userID,
		Email:     fmt.Sprintf("questions-eval-test-%s@example.com", userID.String()),
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
		DisplayName: fmt.Sprintf("questions-eval-%s", metadataID),
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
		Title:             sql.NullString{String: "Questions Generation Test", Valid: true},
		MimeType:          sql.NullString{String: "application/pdf", Valid: true},
		FileStoreName:     sql.NullString{String: storeName, Valid: true},
		FileStoreFileName: sql.NullString{String: operation.Response.DocumentName, Valid: true},
		RagStatus:         "READY",
		UserID:            userID,
	})
	require.NoError(t, err)

	// 2. Generate questions from document
	activeModel, err := queries.GetActiveModelConfig(ctx)
	require.NoError(t, err)
	activeSystemInstruction, err := queries.GetActiveSystemInstruction(ctx)
	require.NoError(t, err)

	toolSeed := fileSearchToolSeed{
		StoreNames:     []string{storeName},
		MetadataFilter: fmt.Sprintf(`user_id = "%s"`, metadataID),
	}
	toolConfig, err := json.Marshal(toolSeed)
	require.NoError(t, err)

	t.Logf("generating questions from document %s", doc.ID)
	generator, err := NewGenerationServiceFromAPIKey(ctx, apiKey)
	require.NoError(t, err)
	artifactService := artifacts.NewService(db)
	genService, err := generation.NewService(db, artifactService, generator, nil)
	require.NoError(t, err)

	resp, err := genService.Generate(ctx, generation.GenerateRequest{
		UserID: userID,
		Target: generation.Target{
			DocumentID: &doc.ID,
		},
		Instructions: generation.Instructions{
			SystemInstructionID: &activeSystemInstruction.ID,
			GenerationType:      utils.GenerationTypeQuestions.String(),
			Variables: map[string]interface{}{
				"question_count": 5,
			},
		},
		Output: generation.OutputConfig{
			GenerationType: utils.GenerationTypeQuestions.String(),
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
	t.Logf("questions output:\n%s", outputPretty)
	t.Logf("grounding metadata captured: %d bytes", len(resp.GroundingMetadata))

	if len(resp.GroundingMetadata) == 0 {
		t.Logf("WARNING: grounding metadata is empty, checking if it's in artifact instead")
	}

	// 3. Create eval and save questions as eval items
	eval, err := queries.CreateEval(ctx, store.CreateEvalParams{
		Title:       "Generated Questions",
		Status:      "draft",
		UserID:      userID,
		Description: sql.NullString{String: "Questions generated from test document", Valid: true},
	})
	require.NoError(t, err)

	// Parse questions from response
	var questionsOutput struct {
		Questions []struct {
			Question       string   `json:"question"`
			ExpectedAnswer string   `json:"expected_answer"`
			Options        []string `json:"options,omitempty"`
			CorrectIdx     int      `json:"correct_idx,omitempty"`
		} `json:"questions"`
	}
	err = json.Unmarshal(resp.OutputJSON, &questionsOutput)
	require.NoError(t, err)
	require.NotEmpty(t, questionsOutput.Questions, "should have generated questions")

	t.Logf("saving %d questions as eval items", len(questionsOutput.Questions))
	evalItemRepo := eval_items.NewRepository(queries)
	var savedItems []*eval_items.EvalItem

	for i, q := range questionsOutput.Questions {
		// Create options if not provided
		options := q.Options
		correctIdx := q.CorrectIdx
		if len(options) == 0 {
			options = []string{q.ExpectedAnswer, "Incorrect option 1", "Incorrect option 2", "Incorrect option 3"}
			correctIdx = 0
		}

		item, err := evalItemRepo.Create(ctx, &eval_items.CreateEvalItemRequest{
			EvalID:            eval.ID,
			Prompt:            q.Question,
			Options:           options,
			CorrectIdx:        int32(correctIdx),
			GroundingMetadata: resp.GroundingMetadata,
			SourceDocumentID:  &doc.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, item)
		savedItems = append(savedItems, item)
		t.Logf("saved eval item %d: %s (id=%s)", i+1, q.Question[:50], item.ID)
	}

	// 4. Evaluate questions using groundedness evaluator
	t.Logf("evaluating %d questions for groundedness", len(savedItems))
	groundednessEvaluator := NewGroundednessEvaluator(client, evals.DefaultGroundednessPrompt)
	groundednessService := evals.NewGroundednessService(groundednessEvaluator)
	evalResultsService := eval_results.NewService(eval_results.NewRepository(queries))

	// Get or create eval prompt
	var evalPromptID uuid.UUID
	activeEvalPrompt, err := queries.GetActiveEvalPrompt(ctx, "groundedness")
	if err == nil && activeEvalPrompt.ID != uuid.Nil {
		evalPromptID = activeEvalPrompt.ID
	} else {
		// Create default eval prompt if not exists
		evalPrompt, err := queries.CreateEvalPrompt(ctx, store.CreateEvalPromptParams{
			EvalType:   "groundedness",
			Version:    1,
			PromptText: evals.DefaultGroundednessPrompt,
			IsActive:   sql.NullBool{Bool: true, Valid: true},
		})
		require.NoError(t, err)
		evalPromptID = evalPrompt.ID
	}

	for i, item := range savedItems {
		// Evaluate groundedness
		result, err := groundednessService.EvaluateEvalItem(ctx, item)
		require.NoError(t, err)
		require.NotNil(t, result)

		t.Logf("eval result %d: verdict=%s score=%.2f", i+1, result.Verdict, result.Score)

		// Save eval result
		unsupportedClaimsJSON, _ := json.Marshal(result.SupportingSegments)
		isGrounded := result.Verdict == "PASS"
		_, err = evalResultsService.Create(ctx, &eval_results.CreateEvalResultRequest{
			EvalItemID:        item.ID,
			EvalType:          "groundedness",
			EvalPromptID:      evalPromptID,
			Score:             &result.Score,
			IsGrounded:        &isGrounded,
			Verdict:           result.Verdict,
			Reasoning:         &result.Reasoning,
			UnsupportedClaims: unsupportedClaimsJSON,
		})
		require.NoError(t, err)
	}

	// 5. Verify results
	stats, err := evalResultsService.GetStats(ctx, "groundedness")
	require.NoError(t, err)
	require.NotNil(t, stats)

	t.Logf("evaluation summary:")
	t.Logf("  total evals: %d", stats.TotalEvals)
	t.Logf("  passed: %d", stats.Passed)
	t.Logf("  failed: %d", stats.Failed)
	t.Logf("  warned: %d", stats.Warned)
	t.Logf("  pass rate: %.2f%%", stats.PassRate)
	t.Logf("  avg score: %.2f", stats.AvgScore)

	require.Equal(t, int64(len(savedItems)), stats.TotalEvals)
	require.Greater(t, stats.TotalEvals, int64(0))
}
