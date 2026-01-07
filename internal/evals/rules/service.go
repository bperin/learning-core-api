package rules

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines the interface for eval rule business logic
type Service interface {
	Create(ctx context.Context, req CreateRuleRequest) (*Rule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Rule, error)
	GetBySuiteAndType(ctx context.Context, suiteID uuid.UUID, evalType string) (*Rule, error)
	ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]Rule, error)
	ListByEvalType(ctx context.Context, evalType string) ([]Rule, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRuleRequest) (*Rule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBySuite(ctx context.Context, suiteID uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new eval rule service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new eval rule with business logic validation
func (s *service) Create(ctx context.Context, req CreateRuleRequest) (*Rule, error) {
	// Business logic validation
	if req.SuiteID == uuid.Nil {
		return nil, errors.New("suite ID is required")
	}

	if req.EvalType == "" {
		return nil, errors.New("eval type is required")
	}

	if req.Weight <= 0 {
		return nil, errors.New("weight must be greater than 0")
	}

	// Check if rule with same suite and type already exists
	existing, err := s.repo.GetBySuiteAndType(ctx, req.SuiteID, req.EvalType)
	if err == nil && existing != nil {
		return nil, errors.New("eval rule with this suite and type already exists")
	}

	// Validate min/max score relationship
	if req.MinScore != nil && req.MaxScore != nil && *req.MinScore > *req.MaxScore {
		return nil, errors.New("min score cannot be greater than max score")
	}

	// Create the rule
	rule := Rule{
		SuiteID:  req.SuiteID,
		EvalType: req.EvalType,
		MinScore: req.MinScore,
		MaxScore: req.MaxScore,
		Weight:   req.Weight,
		HardFail: req.HardFail,
		Params:   req.Params,
	}

	return s.repo.Create(ctx, rule)
}

// GetByID retrieves an eval rule by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval rule ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetBySuiteAndType retrieves an eval rule by suite ID and eval type
func (s *service) GetBySuiteAndType(ctx context.Context, suiteID uuid.UUID, evalType string) (*Rule, error) {
	if suiteID == uuid.Nil {
		return nil, errors.New("suite ID is required")
	}

	if evalType == "" {
		return nil, errors.New("eval type is required")
	}

	return s.repo.GetBySuiteAndType(ctx, suiteID, evalType)
}

// ListBySuite retrieves all eval rules for a suite
func (s *service) ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]Rule, error) {
	if suiteID == uuid.Nil {
		return nil, errors.New("suite ID is required")
	}

	return s.repo.ListBySuite(ctx, suiteID)
}

// ListByEvalType retrieves all eval rules for an eval type
func (s *service) ListByEvalType(ctx context.Context, evalType string) ([]Rule, error) {
	if evalType == "" {
		return nil, errors.New("eval type is required")
	}

	return s.repo.ListByEvalType(ctx, evalType)
}

// Update updates an eval rule with business logic validation
func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateRuleRequest) (*Rule, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval rule ID is required")
	}

	// Check if rule exists
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate min/max score relationship if both are being updated
	if req.MinScore != nil && req.MaxScore != nil && *req.MinScore > *req.MaxScore {
		return nil, errors.New("min score cannot be greater than max score")
	}

	// If only one of min/max scores is being updated, validate against the other
	if req.MinScore != nil && current.MaxScore != nil && *req.MinScore > *current.MaxScore {
		return nil, errors.New("min score cannot be greater than current max score")
	}
	if req.MaxScore != nil && current.MinScore != nil && *current.MinScore > *req.MaxScore {
		return nil, errors.New("current min score cannot be greater than max score")
	}

	return s.repo.Update(ctx, id, req)
}

// Delete deletes an eval rule
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("eval rule ID is required")
	}

	// Check if rule exists before deleting
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// DeleteBySuite deletes all eval rules for a suite
func (s *service) DeleteBySuite(ctx context.Context, suiteID uuid.UUID) error {
	if suiteID == uuid.Nil {
		return errors.New("suite ID is required")
	}

	return s.repo.DeleteBySuite(ctx, suiteID)
}
