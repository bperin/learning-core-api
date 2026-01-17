package evals

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"learning-core-api/internal/domain/eval_items"
)

// GroundednessResult represents the result of a groundedness evaluation
type GroundednessResult struct {
	EvalItemID         uuid.UUID   `json:"eval_item_id"`
	Score              float64     `json:"score"`
	Verdict            string      `json:"verdict"` // PASS, FAIL, WARN
	Reasoning          string      `json:"reasoning,omitempty"`
	SupportingSegments []string    `json:"supporting_segments,omitempty"`
	CreatedAt          string      `json:"created_at"`
}

// GroundednessEvaluator defines the interface for groundedness evaluation
type GroundednessEvaluator interface {
	// EvaluateGroundedness evaluates if the expected answer is grounded in the provided context
	// Returns a score (0-1) and verdict (PASS/FAIL/WARN)
	EvaluateGroundedness(ctx context.Context, question string, expectedAnswer string, groundingMetadata json.RawMessage) (*GroundednessResult, error)
}

// GroundednessService handles groundedness evaluation for eval items
type GroundednessService struct {
	evaluator      GroundednessEvaluator
	promptService  *EvalPromptService
}

// NewGroundednessService creates a new groundedness service
func NewGroundednessService(evaluator GroundednessEvaluator) *GroundednessService {
	return &GroundednessService{
		evaluator: evaluator,
	}
}

// NewGroundednessServiceWithPrompts creates a new groundedness service with prompt management
func NewGroundednessServiceWithPrompts(evaluator GroundednessEvaluator, promptService *EvalPromptService) *GroundednessService {
	return &GroundednessService{
		evaluator:     evaluator,
		promptService: promptService,
	}
}

// EvaluateEvalItem evaluates the groundedness of an eval item's expected answer
func (s *GroundednessService) EvaluateEvalItem(ctx context.Context, item *eval_items.EvalItem) (*GroundednessResult, error) {
	if item == nil {
		return nil, fmt.Errorf("eval item is required")
	}

	if item.Prompt == "" {
		return nil, fmt.Errorf("eval item prompt is required")
	}

	if len(item.GroundingMetadata) == 0 {
		return nil, fmt.Errorf("eval item grounding metadata is required for groundedness evaluation")
	}

	// Get the correct answer from options (for multiple choice)
	expectedAnswer := ""
	if len(item.Options) > 0 {
		if int(item.CorrectIdx) < len(item.Options) {
			expectedAnswer = item.Options[item.CorrectIdx]
		}
	}

	if expectedAnswer == "" {
		return nil, fmt.Errorf("expected answer is required for groundedness evaluation")
	}

	// Call the evaluator
	result, err := s.evaluator.EvaluateGroundedness(ctx, item.Prompt, expectedAnswer, item.GroundingMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate groundedness: %w", err)
	}

	// Set the eval item ID
	result.EvalItemID = item.ID

	return result, nil
}
