package evals

import (
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"learning-core-api/internal/generation"
	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateEvalRun_Cases(t *testing.T) {
	ruleA := uuid.New()
	ruleB := uuid.New()
	score09 := float32(0.9)
	score08 := float32(0.8)
	score06 := float32(0.6)

	tests := []struct {
		name      string
		rulesList []Rule
		results   []Result
		wantPass  bool
		wantScore *float64
	}{
		{
			name: "hard fail rule fails",
			rulesList: []Rule{
				{ID: ruleA, HardFail: true, Weight: 1},
			},
			results: []Result{
				{RuleID: ruleA, Pass: false, Score: &score09},
			},
			wantPass: false,
		},
		{
			name: "hard fail passes, score ok",
			rulesList: []Rule{
				{ID: ruleA, HardFail: true, Weight: 1},
				{ID: ruleB, HardFail: false, Weight: 1},
			},
			results: []Result{
				{RuleID: ruleA, Pass: true, Score: &score09},
				{RuleID: ruleB, Pass: true, Score: &score08},
			},
			wantPass:  true,
			wantScore: floatPtr(0.85),
		},
		{
			name: "missing result fails",
			rulesList: []Rule{
				{ID: ruleA, HardFail: false, Weight: 1},
			},
			results:  []Result{},
			wantPass: false,
		},
		{
			name: "zero weights pass",
			rulesList: []Rule{
				{ID: ruleA, HardFail: false, Weight: 0},
			},
			results: []Result{
				{RuleID: ruleA, Pass: true},
			},
			wantPass: true,
		},
		{
			name: "weighted average",
			rulesList: []Rule{
				{ID: ruleA, Weight: 2},
				{ID: ruleB, Weight: 1},
			},
			results: []Result{
				{RuleID: ruleA, Pass: true, Score: &score09},
				{RuleID: ruleB, Pass: true, Score: &score06},
			},
			wantPass:  true,
			wantScore: floatPtr((2*0.9 + 1*0.6) / 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origRules := cloneRules(tt.rulesList)
			origResults := cloneResults(tt.results)

			agg := AggregateEvalRun(tt.rulesList, tt.results)
			assert.Equal(t, tt.wantPass, agg.OverallPass)

			if tt.wantScore == nil {
				assert.Nil(t, agg.OverallScore)
			} else if assert.NotNil(t, agg.OverallScore) {
				assert.InEpsilon(t, *tt.wantScore, *agg.OverallScore, 0.0001)
				assert.True(t, *agg.OverallScore >= 0 && *agg.OverallScore <= 1)
			}

			assert.True(t, equalRules(origRules, tt.rulesList))
			assert.True(t, equalResults(origResults, tt.results))
		})
	}
}

func TestAggregateAndFinalize_Integration(t *testing.T) {
	ctx := context.Background()
	db := testutil.NewTestDB(t)
	defer db.Close()

	queries := store.New(db)
	runsRepo := NewRunRepository(queries)
	rulesRepo := NewRuleRepository(queries)
	resultsRepo := NewResultRepository(queries)
	artifactRepo := generation.NewRepository(queries)

	svc := NewService(runsRepo, rulesRepo, resultsRepo, artifactRepo)

	tenant, err := queries.CreateTenant(ctx, store.CreateTenantParams{
		Name:     "Eval Tenant",
		IsActive: true,
	})
	require.NoError(t, err)

	module, err := queries.CreateModule(ctx, store.CreateModuleParams{
		TenantID:    tenant.ID,
		Title:       "Eval Module",
		Name:        sql.NullString{String: "eval-module", Valid: true},
		Description: sql.NullString{String: "Eval module", Valid: true},
	})
	require.NoError(t, err)

	genRun, err := queries.CreateGenerationRun(ctx, store.CreateGenerationRunParams{
		ModuleID:       module.ID,
		AgentName:      "test-agent",
		AgentVersion:   "1.0",
		Model:          "gpt-4",
		ModelParams:    json.RawMessage(`{}`),
		PromptID:       uuid.NullUUID{Valid: false},
		StoreName:      "test-store",
		MetadataFilter: json.RawMessage(`{}`),
		Status:         "SUCCEEDED",
		InputPayload:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	artifact, err := queries.CreateArtifact(ctx, store.CreateArtifactParams{
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

	suite, err := queries.CreateEvalSuite(ctx, store.CreateEvalSuiteParams{
		Name:        "suite-" + uuid.New().String()[:8],
		Description: "Integration suite",
	})
	require.NoError(t, err)

	ruleA, err := queries.CreateEvalRule(ctx, store.CreateEvalRuleParams{
		SuiteID:  suite.ID,
		EvalType: "GROUNDEDNESS",
		Weight:   2,
		Params:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	ruleB, err := queries.CreateEvalRule(ctx, store.CreateEvalRuleParams{
		SuiteID:  suite.ID,
		EvalType: "ANSWER_CORRECTNESS",
		Weight:   1,
		Params:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	run, err := queries.CreateEvalRun(ctx, store.CreateEvalRunParams{
		ArtifactID:  artifact.ID,
		SuiteID:     suite.ID,
		JudgeModel:  "gpt-4",
		JudgeParams: json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	_, err = queries.UpsertEvalResult(ctx, store.UpsertEvalResultParams{
		EvalRunID: run.ID,
		RuleID:    ruleA.ID,
		Pass:      true,
		Score:     sql.NullFloat64{Float64: 0.9, Valid: true},
		Details:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	_, err = queries.UpsertEvalResult(ctx, store.UpsertEvalResultParams{
		EvalRunID: run.ID,
		RuleID:    ruleB.ID,
		Pass:      true,
		Score:     sql.NullFloat64{Float64: 0.6, Valid: true},
		Details:   json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	agg, err := svc.AggregateAndFinalize(ctx, run.ID)
	require.NoError(t, err)
	require.NotNil(t, agg.OverallScore)
	assert.True(t, agg.OverallPass)
	assert.InEpsilon(t, 0.8, *agg.OverallScore, 0.0001)

	updatedRun, err := runsRepo.GetByID(ctx, run.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedRun.OverallScore)
	assert.True(t, *updatedRun.OverallPass)
	assert.InEpsilon(t, 0.8, float64(*updatedRun.OverallScore), 0.0001)

	updatedArtifact, err := artifactRepo.GetArtifactByID(ctx, artifact.ID)
	require.NoError(t, err)
	assert.Equal(t, generation.ArtifactStatusApproved, updatedArtifact.Status)
	assert.WithinDuration(t, time.Now(), *updatedArtifact.ApprovedAt, time.Second)
}

func floatPtr(value float64) *float64 {
	return &value
}

func TestAggregateEvalRun_OrderIndependent(t *testing.T) {
	cfg := &quick.Config{
		MaxCount: 50,
		Rand:     rand.New(rand.NewSource(1)),
	}

	err := quick.Check(func(c evalCase) bool {
		agg1 := AggregateEvalRun(c.Rules, c.Results)

		rulesCopy := cloneRules(c.Rules)
		resultsCopy := cloneResults(c.Results)
		rand.Shuffle(len(rulesCopy), func(i, j int) { rulesCopy[i], rulesCopy[j] = rulesCopy[j], rulesCopy[i] })
		rand.Shuffle(len(resultsCopy), func(i, j int) { resultsCopy[i], resultsCopy[j] = resultsCopy[j], resultsCopy[i] })

		agg2 := AggregateEvalRun(rulesCopy, resultsCopy)
		return aggregatesEqual(agg1, agg2)
	}, cfg)
	require.NoError(t, err)
}

func TestAggregateEvalRun_ZeroWeightNoImpact(t *testing.T) {
	cfg := &quick.Config{
		MaxCount: 50,
		Rand:     rand.New(rand.NewSource(2)),
	}

	err := quick.Check(func(c evalCase) bool {
		agg1 := AggregateEvalRun(c.Rules, c.Results)

		ruleID := uuid.New()
		zeroRule := Rule{ID: ruleID, HardFail: false, Weight: 0}
		zeroResult := Result{RuleID: ruleID, Pass: true}

		rulesCopy := append(cloneRules(c.Rules), zeroRule)
		resultsCopy := append(cloneResults(c.Results), zeroResult)

		agg2 := AggregateEvalRun(rulesCopy, resultsCopy)
		return aggregatesEqual(agg1, agg2)
	}, cfg)
	require.NoError(t, err)
}

func TestAggregateEvalRun_HardFailDominates(t *testing.T) {
	cfg := &quick.Config{
		MaxCount: 50,
		Rand:     rand.New(rand.NewSource(3)),
	}

	err := quick.Check(func(c evalCase) bool {
		ruleID := uuid.New()
		hardFailRule := Rule{ID: ruleID, HardFail: true, Weight: 1}
		hardFailResult := Result{RuleID: ruleID, Pass: false}

		rulesCopy := append(cloneRules(c.Rules), hardFailRule)
		resultsCopy := append(cloneResults(c.Results), hardFailResult)

		agg := AggregateEvalRun(rulesCopy, resultsCopy)
		return !agg.OverallPass
	}, cfg)
	require.NoError(t, err)
}

type evalCase struct {
	Rules   []Rule
	Results []Result
}

func (evalCase) Generate(r *rand.Rand, size int) reflect.Value {
	ruleCount := r.Intn(5) + 1
	rulesList := make([]Rule, ruleCount)
	resultsList := make([]Result, ruleCount)

	for i := 0; i < ruleCount; i++ {
		ruleID := uuid.New()
		weightOptions := []float32{0, 0.5, 1, 2}
		weight := weightOptions[r.Intn(len(weightOptions))]

		rulesList[i] = Rule{
			ID:       ruleID,
			Weight:   weight,
			HardFail: r.Intn(2) == 0,
		}

		pass := r.Intn(2) == 0
		var score *float32
		if r.Intn(2) == 0 {
			s := float32(r.Float64())
			score = &s
		}

		resultsList[i] = Result{
			RuleID: ruleID,
			Pass:   pass,
			Score:  score,
		}
	}

	return reflect.ValueOf(evalCase{
		Rules:   rulesList,
		Results: resultsList,
	})
}

func cloneRules(input []Rule) []Rule {
	if input == nil {
		return nil
	}
	out := make([]Rule, len(input))
	copy(out, input)
	return out
}

func cloneResults(input []Result) []Result {
	if input == nil {
		return nil
	}
	out := make([]Result, len(input))
	copy(out, input)
	return out
}

func aggregatesEqual(a, b EvalAggregate) bool {
	if a.OverallPass != b.OverallPass {
		return false
	}
	if a.OverallScore == nil && b.OverallScore == nil {
		return true
	}
	if a.OverallScore == nil || b.OverallScore == nil {
		return false
	}
	return math.Abs(*a.OverallScore-*b.OverallScore) <= 1e-9
}

func equalRules(a, b []Rule) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func equalResults(a, b []Result) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
