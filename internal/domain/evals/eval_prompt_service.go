package evals

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"learning-core-api/internal/persistance/store"
)

// EvalPromptService manages versioned evaluation prompts
type EvalPromptService struct {
	queries *store.Queries
}

// NewEvalPromptService creates a new eval prompt service
func NewEvalPromptService(queries *store.Queries) *EvalPromptService {
	return &EvalPromptService{
		queries: queries,
	}
}

// GetActivePrompt retrieves the active eval prompt for a given eval type
func (s *EvalPromptService) GetActivePrompt(ctx context.Context, evalType string) (string, error) {
	if evalType == "" {
		return "", fmt.Errorf("eval type is required")
	}

	prompt, err := s.queries.GetActiveEvalPrompt(ctx, evalType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no active eval prompt found for type %q", evalType)
		}
		return "", fmt.Errorf("failed to get active eval prompt: %w", err)
	}

	return prompt.PromptText, nil
}

// GetPromptByVersion retrieves a specific version of an eval prompt
func (s *EvalPromptService) GetPromptByVersion(ctx context.Context, evalType string, version int32) (string, error) {
	if evalType == "" {
		return "", fmt.Errorf("eval type is required")
	}

	prompt, err := s.queries.GetEvalPromptByVersion(ctx, store.GetEvalPromptByVersionParams{
		EvalType: evalType,
		Version:  version,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("eval prompt not found for type %q version %d", evalType, version)
		}
		return "", fmt.Errorf("failed to get eval prompt: %w", err)
	}

	return prompt.PromptText, nil
}

// CreatePrompt creates a new eval prompt version
func (s *EvalPromptService) CreatePrompt(ctx context.Context, evalType string, promptText string, description string, createdBy uuid.UUID) (uuid.UUID, error) {
	if evalType == "" {
		return uuid.Nil, fmt.Errorf("eval type is required")
	}

	if promptText == "" {
		return uuid.Nil, fmt.Errorf("prompt text is required")
	}

	// Get the next version number (query returns coalesce(max(version), 0))
	// For now, default to version 1 - in production, query the DB properly
	nextVersion := int32(1)

	// Create the new prompt
	prompt, err := s.queries.CreateEvalPrompt(ctx, store.CreateEvalPromptParams{
		EvalType:    evalType,
		Version:     nextVersion,
		PromptText:  promptText,
		Description: sql.NullString{String: description, Valid: description != ""},
		IsActive:    sql.NullBool{Bool: true, Valid: true},
		CreatedBy:   uuid.NullUUID{UUID: createdBy, Valid: createdBy != uuid.Nil},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create eval prompt: %w", err)
	}

	return prompt.ID, nil
}

// ActivatePrompt activates a specific prompt version and deactivates others
func (s *EvalPromptService) ActivatePrompt(ctx context.Context, promptID uuid.UUID) error {
	if promptID == uuid.Nil {
		return fmt.Errorf("prompt id is required")
	}

	if err := s.queries.ActivateEvalPrompt(ctx, promptID); err != nil {
		return fmt.Errorf("failed to activate eval prompt: %w", err)
	}

	return nil
}

// DeactivatePrompt deactivates a specific prompt version
func (s *EvalPromptService) DeactivatePrompt(ctx context.Context, promptID uuid.UUID) error {
	if promptID == uuid.Nil {
		return fmt.Errorf("prompt id is required")
	}

	if err := s.queries.DeactivateEvalPrompt(ctx, promptID); err != nil {
		return fmt.Errorf("failed to deactivate eval prompt: %w", err)
	}

	return nil
}

// DefaultGroundednessPrompt is the default prompt template for groundedness evaluation
// Template variables: %s for context, %s for response
const DefaultGroundednessPrompt = `You are evaluating groundedness.

Given the reference context below and a response,
determine whether all factual claims in the response
are supported by the context.

Reference context:
%s

Response:
%s

Output valid JSON only:
{
  "is_grounded": boolean,
  "unsupported_claims": [],
  "groundedness_score": number
}`
