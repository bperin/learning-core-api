package evals

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ResultService defines the interface for eval result business logic.
type ResultService interface {
	Create(ctx context.Context, req CreateResultRequest) (*Result, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Result, error)
	ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error)
	ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error)
	DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error
}

// resultService implements the ResultService interface.
type resultService struct {
	repo ResultRepository
}

// NewResultService creates a new eval result service.
func NewResultService(repo ResultRepository) ResultService {
	return &resultService{
		repo: repo,
	}
}

// Create creates a new eval result with business logic validation.
func (s *resultService) Create(ctx context.Context, req CreateResultRequest) (*Result, error) {
	// Business logic validation
	if req.EvalRunID == uuid.Nil {
		return nil, errors.New("eval run ID is required")
	}

	if req.RuleID == uuid.Nil {
		return nil, errors.New("rule ID is required")
	}

	// Create the result
	result := Result{
		EvalRunID: req.EvalRunID,
		RuleID:    req.RuleID,
		Pass:      req.Pass,
		Score:     req.Score,
		Details:   req.Details,
	}

	return s.repo.Create(ctx, result)
}

// GetByID retrieves an eval result by ID.
func (s *resultService) GetByID(ctx context.Context, id uuid.UUID) (*Result, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval result ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// ListByRun retrieves all eval results for an eval run.
func (s *resultService) ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error) {
	if evalRunID == uuid.Nil {
		return nil, errors.New("eval run ID is required")
	}

	return s.repo.ListByRun(ctx, evalRunID)
}

// ListByRule retrieves all eval results for a specific rule.
func (s *resultService) ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error) {
	if ruleID == uuid.Nil {
		return nil, errors.New("rule ID is required")
	}

	return s.repo.ListByRule(ctx, ruleID)
}

// DeleteByRun deletes all eval results for an eval run.
func (s *resultService) DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error {
	if evalRunID == uuid.Nil {
		return errors.New("eval run ID is required")
	}

	return s.repo.DeleteByRun(ctx, evalRunID)
}
