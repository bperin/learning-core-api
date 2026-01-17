package eval_results

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Service handles eval result business logic
type Service struct {
	repo Repository
}

// NewService creates a new eval result service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new eval result
func (s *Service) Create(ctx context.Context, req *CreateEvalResultRequest) (*EvalResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	result, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create eval result: %w", err)
	}

	return result, nil
}

// GetByID retrieves an eval result by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*EvalResult, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("eval result id is required")
	}

	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get eval result: %w", err)
	}

	if result == nil {
		return nil, ErrResultNotFound
	}

	return result, nil
}

// GetByEvalItem retrieves all eval results for a specific eval item
func (s *Service) GetByEvalItem(ctx context.Context, evalItemID uuid.UUID) ([]*EvalResult, error) {
	if evalItemID == uuid.Nil {
		return nil, fmt.Errorf("eval item id is required")
	}

	results, err := s.repo.GetByEvalItem(ctx, evalItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get eval results for item: %w", err)
	}

	return results, nil
}

// GetLatestByEvalItem retrieves the latest eval result for a specific eval item and type
func (s *Service) GetLatestByEvalItem(ctx context.Context, evalItemID uuid.UUID, evalType string) (*EvalResult, error) {
	if evalItemID == uuid.Nil {
		return nil, fmt.Errorf("eval item id is required")
	}

	if evalType == "" {
		return nil, ErrInvalidEvalType
	}

	result, err := s.repo.GetLatestByEvalItem(ctx, evalItemID, evalType)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest eval result: %w", err)
	}

	return result, nil
}

// List retrieves eval results with pagination
func (s *Service) List(ctx context.Context, limit int32, offset int32) ([]*EvalResult, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	results, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list eval results: %w", err)
	}

	return results, nil
}

// ListByType retrieves eval results of a specific type with pagination
func (s *Service) ListByType(ctx context.Context, evalType string, limit int32, offset int32) ([]*EvalResult, error) {
	if evalType == "" {
		return nil, ErrInvalidEvalType
	}

	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	results, err := s.repo.ListByType(ctx, evalType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list eval results by type: %w", err)
	}

	return results, nil
}

// GetStats retrieves aggregate statistics for eval results of a type
func (s *Service) GetStats(ctx context.Context, evalType string) (*EvalResultStats, error) {
	if evalType == "" {
		return nil, ErrInvalidEvalType
	}

	stats, err := s.repo.GetStats(ctx, evalType)
	if err != nil {
		return nil, fmt.Errorf("failed to get eval result stats: %w", err)
	}

	return stats, nil
}

// Count returns the total number of eval results
func (s *Service) Count(ctx context.Context) (int64, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count eval results: %w", err)
	}

	return count, nil
}

// CountByType returns the count of eval results for a specific type
func (s *Service) CountByType(ctx context.Context, evalType string) (int64, error) {
	if evalType == "" {
		return 0, ErrInvalidEvalType
	}

	count, err := s.repo.CountByType(ctx, evalType)
	if err != nil {
		return 0, fmt.Errorf("failed to count eval results by type: %w", err)
	}

	return count, nil
}
