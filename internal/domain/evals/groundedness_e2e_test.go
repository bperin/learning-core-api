package evals_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/eval_items"
	"learning-core-api/internal/domain/evals"
)

// MockGroundednessEvaluator is a mock implementation for testing
type MockGroundednessEvaluator struct {
	shouldPass bool
}

func (m *MockGroundednessEvaluator) EvaluateGroundedness(ctx context.Context, question string, expectedAnswer string, groundingMetadata json.RawMessage) (*evals.GroundednessResult, error) {
	verdict := "FAIL"
	score := 0.2
	if m.shouldPass {
		verdict = "PASS"
		score = 0.95
	}

	return &evals.GroundednessResult{
		Score:      score,
		Verdict:    verdict,
		Reasoning:  "Mock evaluation result",
		CreatedAt:  "2026-01-15T23:35:00Z",
	}, nil
}

// TestGroundednessEvaluation tests the groundedness evaluation service
func TestGroundednessEvaluation(t *testing.T) {
	ctx := context.Background()

	// Create mock grounding metadata
	groundingMetadata := json.RawMessage(`{
		"groundingChunks": [
			{
				"retrievedContext": {
					"text": "Mitochondria are the powerhouse of the cell, responsible for ATP production.",
					"title": "cell-biology"
				}
			}
		],
		"groundingSupports": [
			{
				"groundingChunkIndices": [0],
				"segment": {
					"startIndex": 0,
					"endIndex": 50,
					"text": "Mitochondria are the powerhouse of the cell"
				}
			}
		]
	}`)

	t.Run("evaluate grounded question", func(t *testing.T) {
		mockEval := &MockGroundednessEvaluator{shouldPass: true}
		service := evals.NewGroundednessService(mockEval)

		item := &eval_items.EvalItem{
			ID:                uuid.New(),
			EvalID:            uuid.New(),
			Prompt:            "What is the primary function of mitochondria?",
			Options:           []string{"Storage", "Energy production", "Protein synthesis", "DNA replication"},
			CorrectIdx:        1,
			GroundingMetadata: groundingMetadata,
		}

		result, err := service.EvaluateEvalItem(ctx, item)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, item.ID, result.EvalItemID)
		assert.Equal(t, "PASS", result.Verdict)
		assert.Greater(t, result.Score, 0.9)
	})

	t.Run("evaluate ungrounded question", func(t *testing.T) {
		mockEval := &MockGroundednessEvaluator{shouldPass: false}
		service := evals.NewGroundednessService(mockEval)

		item := &eval_items.EvalItem{
			ID:                uuid.New(),
			EvalID:            uuid.New(),
			Prompt:            "What is the primary function of mitochondria?",
			Options:           []string{"Storage", "Energy production", "Protein synthesis", "DNA replication"},
			CorrectIdx:        1,
			GroundingMetadata: groundingMetadata,
		}

		result, err := service.EvaluateEvalItem(ctx, item)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "FAIL", result.Verdict)
		assert.Less(t, result.Score, 0.5)
	})

	t.Run("missing grounding metadata", func(t *testing.T) {
		mockEval := &MockGroundednessEvaluator{shouldPass: true}
		service := evals.NewGroundednessService(mockEval)

		item := &eval_items.EvalItem{
			ID:         uuid.New(),
			EvalID:     uuid.New(),
			Prompt:     "What is the primary function of mitochondria?",
			Options:    []string{"Storage", "Energy production", "Protein synthesis", "DNA replication"},
			CorrectIdx: 1,
		}

		result, err := service.EvaluateEvalItem(ctx, item)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "grounding metadata is required")
	})

	t.Run("missing expected answer", func(t *testing.T) {
		mockEval := &MockGroundednessEvaluator{shouldPass: true}
		service := evals.NewGroundednessService(mockEval)

		item := &eval_items.EvalItem{
			ID:                uuid.New(),
			EvalID:            uuid.New(),
			Prompt:            "What is the primary function of mitochondria?",
			Options:           []string{},
			CorrectIdx:        0,
			GroundingMetadata: groundingMetadata,
		}

		result, err := service.EvaluateEvalItem(ctx, item)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "expected answer is required")
	})
}

// TestGroundednessResultStructure tests that the result has the expected fields
func TestGroundednessResultStructure(t *testing.T) {
	result := &evals.GroundednessResult{
		EvalItemID:         uuid.New(),
		Score:              0.95,
		Verdict:            "PASS",
		Reasoning:          "The answer is fully supported by the context",
		SupportingSegments: []string{"segment1", "segment2"},
		CreatedAt:          "2026-01-15T23:35:00Z",
	}

	// Verify all fields are set
	assert.NotEqual(t, uuid.Nil, result.EvalItemID)
	assert.Equal(t, 0.95, result.Score)
	assert.Equal(t, "PASS", result.Verdict)
	assert.NotEmpty(t, result.Reasoning)
	assert.Len(t, result.SupportingSegments, 2)
	assert.NotEmpty(t, result.CreatedAt)

	// Verify it can be marshaled to JSON
	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify it can be unmarshaled back
	var unmarshaled evals.GroundednessResult
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, result.Score, unmarshaled.Score)
	assert.Equal(t, result.Verdict, unmarshaled.Verdict)
}
