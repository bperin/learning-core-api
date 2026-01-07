package evals

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func SeedTenant(ctx context.Context, q *store.Queries) uuid.UUID {
	t, err := q.CreateTenant(ctx, store.CreateTenantParams{
		Name:     "Demo",
		IsActive: true,
	})
	if err != nil {
		panic(err)
	}
	return t.ID
}

func SeedEvalSuite(ctx context.Context, q *store.Queries) uuid.UUID {
	s, err := q.CreateEvalSuite(ctx, store.CreateEvalSuiteParams{
		Name:        "gate_fast_" + uuid.New().String()[:8],
		Description: "Fast gate",
	})
	if err != nil {
		panic(err)
	}
	return s.ID
}

func SeedEvalRule(ctx context.Context, q *store.Queries, suiteID uuid.UUID) uuid.UUID {
	r, err := q.CreateEvalRule(ctx, store.CreateEvalRuleParams{
		SuiteID:  suiteID,
		EvalType: "GROUNDEDNESS",
		MinScore: sql.NullFloat64{Float64: 0.75, Valid: true},
		Weight:   1.0,
		HardFail: true,
		Params:   json.RawMessage(`{}`),
	})
	if err != nil {
		panic(err)
	}
	return r.ID
}

func TestEvalResults_Upsert(t *testing.T) {
	ctx := context.Background()
	dbConn := testutil.NewTestDB(t)
	defer dbConn.Close()

	q := store.New(dbConn)

	tenantID := SeedTenant(ctx, q)

	// Need a module and generation run to create an artifact
	module, err := q.CreateModule(ctx, store.CreateModuleParams{
		TenantID:    tenantID,
		Description: sql.NullString{String: "Test Description", Valid: true},
	})
	require.NoError(t, err)

	genRun, err := q.CreateGenerationRun(ctx, store.CreateGenerationRunParams{
		ModuleID:       module.ID,
		AgentName:      "test-agent",
		AgentVersion:   "1.0",
		Model:          "gpt-4",
		ModelParams:    json.RawMessage(`{}`),
		PromptID:       uuid.NullUUID{Valid: false},
		StoreName:      "test-store",
		MetadataFilter: json.RawMessage(`{}`),
		Status:         "PENDING",
		InputPayload:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

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

	suiteID := SeedEvalSuite(ctx, q)
	ruleID := SeedEvalRule(ctx, q, suiteID)

	// create eval_run
	run, err := q.CreateEvalRun(ctx, store.CreateEvalRunParams{
		ArtifactID:  artifact.ID,
		SuiteID:     suiteID,
		JudgeModel:  "gemini",
		JudgeParams: json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	// UPSERT #1
	res1, err := q.UpsertEvalResult(ctx, store.UpsertEvalResultParams{
		EvalRunID: run.ID,
		RuleID:    ruleID,
		Pass:      false,
		Score:     sql.NullFloat64{Float64: 0.6, Valid: true},
		Details:   json.RawMessage(`{"reason":"bad grounding"}`),
	})
	require.NoError(t, err)
	require.False(t, res1.Pass)

	// UPSERT #2 (same key)
	res2, err := q.UpsertEvalResult(ctx, store.UpsertEvalResultParams{
		EvalRunID: run.ID,
		RuleID:    ruleID,
		Pass:      true,
		Score:     sql.NullFloat64{Float64: 0.9, Valid: true},
		Details:   json.RawMessage(`{"fixed":true}`),
	})
	require.NoError(t, err)
	require.True(t, res2.Pass)

	// ensure only one row exists
	results, err := q.ListEvalResultsByRun(ctx, run.ID)
	require.NoError(t, err)
	require.Len(t, results, 1)
}
