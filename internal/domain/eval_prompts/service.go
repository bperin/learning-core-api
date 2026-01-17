package eval_prompts

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Service handles eval prompt business logic
type Service struct {
	repo Repository
}

// NewService creates a new eval prompt service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetActive retrieves the active eval prompt for a given eval type
func (s *Service) GetActive(ctx context.Context, evalType string) (*EvalPrompt, error) {
	if evalType == "" {
		return nil, ErrInvalidEvalType
	}

	prompt, err := s.repo.GetActiveByType(ctx, evalType)
	if err != nil {
		return nil, fmt.Errorf("failed to get active prompt for type %q: %w", evalType, err)
	}

	if prompt == nil {
		return nil, fmt.Errorf("%w for type %q", ErrNoActivePrompt, evalType)
	}

	return prompt, nil
}

// GetByVersion retrieves a specific version of an eval prompt
func (s *Service) GetByVersion(ctx context.Context, evalType string, version int32) (*EvalPrompt, error) {
	if evalType == "" {
		return nil, ErrInvalidEvalType
	}

	if version <= 0 {
		return nil, ErrInvalidVersion
	}

	prompt, err := s.repo.GetByTypeAndVersion(ctx, evalType, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt version: %w", err)
	}

	if prompt == nil {
		return nil, fmt.Errorf("%w: type=%q version=%d", ErrPromptNotFound, evalType, version)
	}

	return prompt, nil
}

// Create creates a new eval prompt version
func (s *Service) Create(ctx context.Context, req *CreateEvalPromptRequest) (*EvalPrompt, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	prompt, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create eval prompt: %w", err)
	}

	return prompt, nil
}

// Activate activates a specific eval prompt version
func (s *Service) Activate(ctx context.Context, promptID uuid.UUID) error {
	if promptID == uuid.Nil {
		return fmt.Errorf("prompt id is required")
	}

	if err := s.repo.Activate(ctx, promptID); err != nil {
		return fmt.Errorf("failed to activate prompt: %w", err)
	}

	return nil
}

// Deactivate deactivates a specific eval prompt version
func (s *Service) Deactivate(ctx context.Context, promptID uuid.UUID) error {
	if promptID == uuid.Nil {
		return fmt.Errorf("prompt id is required")
	}

	if err := s.repo.Deactivate(ctx, promptID); err != nil {
		return fmt.Errorf("failed to deactivate prompt: %w", err)
	}

	return nil
}

// List retrieves all versions of an eval prompt type
func (s *Service) List(ctx context.Context, evalType string, limit int32, offset int32) ([]*EvalPrompt, error) {
	if evalType == "" {
		return nil, ErrInvalidEvalType
	}

	prompts, err := s.repo.ListByType(ctx, evalType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	return prompts, nil
}
