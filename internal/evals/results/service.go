package results

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines the interface for eval result business logic
type Service interface {
	Create(ctx context.Context, req CreateResultRequest) (*Result, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Result, error)
	ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error)
	ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error)
	DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new eval result service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new eval result with business logic validation
func (s *service) Create(ctx context.Context, req CreateResultRequest) (*Result, error) {
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

// GetByID retrieves an eval result by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Result, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval result ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// ListByRun retrieves all eval results for an eval run
func (s *service) ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error) {
	if evalRunID == uuid.Nil {
		return nil, errors.New("eval run ID is required")
	}

	return s.repo.ListByRun(ctx, evalRunID)
}

// ListByRule retrieves all eval results for a specific rule
func (s *service) ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error) {
	if ruleID == uuid.Nil {
		return nil, errors.New("rule ID is required")
	}

	return s.repo.ListByRule(ctx, ruleID)
}

// DeleteByRun deletes all eval results for an eval run
func (s *service) DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error {
	if evalRunID == uuid.Nil {
		return errors.New("eval run ID is required")
	}

	return s.repo.DeleteByRun(ctx, evalRunID)
}
