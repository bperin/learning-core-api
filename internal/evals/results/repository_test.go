package results_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/evals/results"
	"learning-core-api/internal/modules"
	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (results.Repository, *store.Queries, context.Context, func()) {
	ctx := context.Background()
	db := testutil.NewTestDB(t)

	queries := store.New(db)
	repo := results.NewRepository(queries)

	cleanup := func() {
		db.Close()
	}

	return repo, queries, ctx, cleanup
}

func seedDependencies(ctx context.Context, t *testing.T, q *store.Queries) (uuid.UUID, uuid.UUID, uuid.UUID) {
	// 1. Create Tenant
	tenant, err := q.CreateTenant(ctx, store.CreateTenantParams{
		Name:     "Test Tenant",
		IsActive: true,
	})
	require.NoError(t, err)

	// 2. Create Module using modules Repository
	modRepo := modules.NewRepository(q)
	module, err := modRepo.Create(ctx, modules.Module{
		TenantID:    tenant.ID,
		Title:       "Test Title",
		Name:        "Test Module " + uuid.New().String()[:8],
		Description: "Test Description",
	})
	require.NoError(t, err)

	// 3. Create Generation Run
	genRun, err := q.CreateGenerationRun(ctx, store.CreateGenerationRunParams{
		ModuleID:       module.ID,
		AgentName:      "test-agent",
		AgentVersion:   "1.0",
		Model:          "gpt-4",
		ModelParams:    json.RawMessage(`{}`),
		StoreName:      "test-store",
		MetadataFilter: json.RawMessage(`{}`),
		Status:         "SUCCEEDED",
		InputPayload:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	// 4. Create Artifact
	artifact, err := q.CreateArtifact(ctx, store.CreateArtifactParams{
		ModuleID:        module.ID,
		GenerationRunID: genRun.ID,
		Type:            "MCQ_ITEM",
		Status:          "PENDING_EVAL",
		SchemaVersion:   "1.0",
		ArtifactPayload: json.RawMessage(`{}`),
		Grounding:       json.RawMessage(`{}`),
		Tags:            []string{},
	})
	require.NoError(t, err)

	// 5. Create Eval Suite
	suite, err := q.CreateEvalSuite(ctx, store.CreateEvalSuiteParams{
		Name:        "test-suite-" + uuid.New().String()[:8],
		Description: "Test suite",
	})
	require.NoError(t, err)

	// 6. Create Eval Rule
	rule, err := q.CreateEvalRule(ctx, store.CreateEvalRuleParams{
		SuiteID:  suite.ID,
		EvalType: "GROUNDEDNESS",
		MinScore: sql.NullFloat64{Float64: 0.7, Valid: true},
		Weight:   1.0,
		Params:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	// 7. Create Eval Run
	run, err := q.CreateEvalRun(ctx, store.CreateEvalRunParams{
		ArtifactID:  artifact.ID,
		SuiteID:     suite.ID,
		JudgeModel:  "gpt-4",
		JudgeParams: json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	return run.ID, rule.ID, suite.ID
}

func TestRepository_Create(t *testing.T) {
	repo, queries, ctx, cleanup := setupTest(t)
	defer cleanup()

	runID, ruleID, _ := seedDependencies(ctx, t, queries)

	score := float32(0.85)
	res := results.Result{
		EvalRunID: runID,
		RuleID:    ruleID,
		Pass:      true,
		Score:     &score,
		Details:   json.RawMessage(`{"reason": "good"}`),
	}

	created, err := repo.Create(ctx, res)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, res.EvalRunID, created.EvalRunID)
	assert.Equal(t, res.RuleID, created.RuleID)
	assert.Equal(t, res.Pass, created.Pass)
	assert.Equal(t, *res.Score, *created.Score)
	assert.JSONEq(t, string(res.Details), string(created.Details))
}

func TestRepository_GetByID(t *testing.T) {
	repo, queries, ctx, cleanup := setupTest(t)
	defer cleanup()

	runID, ruleID, _ := seedDependencies(ctx, t, queries)

	score := float32(0.85)
	res := results.Result{
		EvalRunID: runID,
		RuleID:    ruleID,
		Pass:      true,
		Score:     &score,
		Details:   json.RawMessage(`{"reason": "good"}`),
	}

	created, _ := repo.Create(ctx, res)

	retrieved, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.EvalRunID, retrieved.EvalRunID)
	assert.Equal(t, created.RuleID, retrieved.RuleID)
	assert.Equal(t, created.Pass, retrieved.Pass)
	assert.Equal(t, *created.Score, *retrieved.Score)
}

func TestRepository_ListByRun(t *testing.T) {
	repo, queries, ctx, cleanup := setupTest(t)
	defer cleanup()

	runID, ruleID, suiteID := seedDependencies(ctx, t, queries)

	// Create another rule in the same suite
	rule2, _ := queries.CreateEvalRule(ctx, store.CreateEvalRuleParams{
		SuiteID:  suiteID,
		EvalType: "ANSWER_CORRECTNESS",
		Weight:   1.0,
		Params:   []byte(`{}`),
	})

	repo.Create(ctx, results.Result{
		EvalRunID: runID,
		RuleID:    ruleID,
		Pass:      true,
		Details:   json.RawMessage(`{}`),
	})

	repo.Create(ctx, results.Result{
		EvalRunID: runID,
		RuleID:    rule2.ID,
		Pass:      false,
		Details:   json.RawMessage(`{}`),
	})

	list, err := repo.ListByRun(ctx, runID)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestRepository_ListByRule(t *testing.T) {
	repo, queries, ctx, cleanup := setupTest(t)
	defer cleanup()

	runID, ruleID, _ := seedDependencies(ctx, t, queries)

	repo.Create(ctx, results.Result{
		EvalRunID: runID,
		RuleID:    ruleID,
		Pass:      true,
		Details:   json.RawMessage(`{}`),
	})

	list, err := repo.ListByRule(ctx, ruleID)
	require.NoError(t, err)
	assert.NotEmpty(t, list)
	assert.Equal(t, ruleID, list[0].RuleID)
}

func TestRepository_DeleteByRun(t *testing.T) {
	repo, queries, ctx, cleanup := setupTest(t)
	defer cleanup()

	runID, ruleID, _ := seedDependencies(ctx, t, queries)

	repo.Create(ctx, results.Result{
		EvalRunID: runID,
		RuleID:    ruleID,
		Pass:      true,
		Details:   json.RawMessage(`{}`),
	})

	err := repo.DeleteByRun(ctx, runID)
	require.NoError(t, err)

	list, _ := repo.ListByRun(ctx, runID)
	assert.Empty(t, list)
}
