package evals_test

import (
	"context"
	"database/sql"
	"testing"

	"learning-core-api/internal/store"
	"learning-core-api/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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

func TestEvalRules_CRUD(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := testutil.StartPostgres(ctx)
	defer cleanup()

	require.NoError(t, testutil.Migrate(dbConn))

	q := store.New(dbConn)

	suiteID := SeedEvalSuite(ctx, q)

	types := []string{
		"SCHEMA_VALIDATION",
		"GROUNDEDNESS",
		"ANSWER_CORRECTNESS",
		"DISTRACTOR_QUALITY",
		"DIFFICULTY_CALIBRATION",
		"CONCEPT_ALIGNMENT",
	}

	for _, et := range types {
		t.Run("Type_"+et, func(t *testing.T) {
			// CREATE
			rule, err := q.CreateEvalRule(ctx, store.CreateEvalRuleParams{
				SuiteID:  suiteID,
				EvalType: et,
				MinScore: sql.NullFloat64{Float64: 0.8, Valid: true},
				Weight:   1.0,
				HardFail: false,
				Params:   []byte(`{"type":"` + et + `"}`),
			})
			require.NoError(t, err)
			require.Equal(t, et, rule.EvalType)

			// GET
			got, err := q.GetEvalRule(ctx, rule.ID)
			require.NoError(t, err)
			require.Equal(t, rule.ID, got.ID)

			// UPDATE
			updated, err := q.UpdateEvalRule(ctx, store.UpdateEvalRuleParams{
				ID:       rule.ID,
				MinScore: sql.NullFloat64{Float64: 0.9, Valid: true},
				MaxScore: sql.NullFloat64{},
				Weight:   2.0,
				HardFail: true,
				Params:   []byte(`{"updated":true}`),
			})
			require.NoError(t, err)
			require.True(t, updated.HardFail)
		})
	}

	// LIST
	rules, err := q.ListEvalRulesBySuite(ctx, suiteID)
	require.NoError(t, err)
	require.Len(t, rules, len(types))

	// DELETE ALL
	for _, r := range rules {
		err = q.DeleteEvalRule(ctx, r.ID)
		require.NoError(t, err)
	}
}
