package rules

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"
	"learning-core-api/internal/utils"

	"github.com/google/uuid"
)

// Repository defines the interface for eval rule operations
type Repository interface {
	Create(ctx context.Context, rule Rule) (*Rule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Rule, error)
	GetBySuiteAndType(ctx context.Context, suiteID uuid.UUID, evalType string) (*Rule, error)
	ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]Rule, error)
	ListByEvalType(ctx context.Context, evalType string) ([]Rule, error)
	Update(ctx context.Context, id uuid.UUID, rule UpdateRuleRequest) (*Rule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBySuite(ctx context.Context, suiteID uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new eval rule repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Create creates a new eval rule
func (r *repository) Create(ctx context.Context, rule Rule) (*Rule, error) {
	params := store.CreateEvalRuleParams{
		SuiteID:  rule.SuiteID,
		EvalType: rule.EvalType,
		Weight:   rule.Weight,
		HardFail: rule.HardFail,
		Params:   rule.Params,
	}

	if rule.MinScore != nil {
		params.MinScore = sql.NullFloat64{Float64: float64(*rule.MinScore), Valid: true}
	} else {
		params.MinScore = sql.NullFloat64{Valid: false}
	}

	if rule.MaxScore != nil {
		params.MaxScore = sql.NullFloat64{Float64: float64(*rule.MaxScore), Valid: true}
	} else {
		params.MaxScore = sql.NullFloat64{Valid: false}
	}

	dbRule, err := r.queries.CreateEvalRule(ctx, params)
	if err != nil {
		return nil, err
	}

	var minScore, maxScore *float32
	if dbRule.MinScore.Valid {
		s := float32(dbRule.MinScore.Float64)
		minScore = &s
	}
	if dbRule.MaxScore.Valid {
		s := float32(dbRule.MaxScore.Float64)
		maxScore = &s
	}

	return &Rule{
		ID:        dbRule.ID,
		SuiteID:   dbRule.SuiteID,
		EvalType:  utils.InterfaceToString(dbRule.EvalType), // Convert interface{} to string
		MinScore:  minScore,
		MaxScore:  maxScore,
		Weight:    dbRule.Weight,
		HardFail:  dbRule.HardFail,
		Params:    dbRule.Params,
		CreatedAt: dbRule.CreatedAt,
	}, nil
}

// GetByID retrieves an eval rule by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	dbRule, err := r.queries.GetEvalRule(ctx, id)
	if err != nil {
		return nil, err
	}

	var minScore, maxScore *float32
	if dbRule.MinScore.Valid {
		s := float32(dbRule.MinScore.Float64)
		minScore = &s
	}
	if dbRule.MaxScore.Valid {
		s := float32(dbRule.MaxScore.Float64)
		maxScore = &s
	}

	return &Rule{
		ID:        dbRule.ID,
		SuiteID:   dbRule.SuiteID,
		EvalType:  utils.InterfaceToString(dbRule.EvalType), // Convert interface{} to string
		MinScore:  minScore,
		MaxScore:  maxScore,
		Weight:    dbRule.Weight,
		HardFail:  dbRule.HardFail,
		Params:    dbRule.Params,
		CreatedAt: dbRule.CreatedAt,
	}, nil
}

// GetBySuiteAndType retrieves an eval rule by suite ID and eval type
func (r *repository) GetBySuiteAndType(ctx context.Context, suiteID uuid.UUID, evalType string) (*Rule, error) {
	params := store.GetEvalRuleBySuiteAndTypeParams{
		SuiteID:  suiteID,
		EvalType: evalType,
	}

	dbRule, err := r.queries.GetEvalRuleBySuiteAndType(ctx, params)
	if err != nil {
		return nil, err
	}

	var minScore, maxScore *float32
	if dbRule.MinScore.Valid {
		s := float32(dbRule.MinScore.Float64)
		minScore = &s
	}
	if dbRule.MaxScore.Valid {
		s := float32(dbRule.MaxScore.Float64)
		maxScore = &s
	}

	return &Rule{
		ID:        dbRule.ID,
		SuiteID:   dbRule.SuiteID,
		EvalType:  utils.InterfaceToString(dbRule.EvalType), // Convert interface{} to string
		MinScore:  minScore,
		MaxScore:  maxScore,
		Weight:    dbRule.Weight,
		HardFail:  dbRule.HardFail,
		Params:    dbRule.Params,
		CreatedAt: dbRule.CreatedAt,
	}, nil
}

// ListBySuite retrieves all eval rules for a suite
func (r *repository) ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]Rule, error) {
	dbRules, err := r.queries.ListEvalRulesBySuite(ctx, suiteID)
	if err != nil {
		return nil, err
	}

	rules := make([]Rule, len(dbRules))
	for i, dbRule := range dbRules {
		var minScore, maxScore *float32
		if dbRule.MinScore.Valid {
			s := float32(dbRule.MinScore.Float64)
			minScore = &s
		}
		if dbRule.MaxScore.Valid {
			s := float32(dbRule.MaxScore.Float64)
			maxScore = &s
		}

		rules[i] = Rule{
			ID:        dbRule.ID,
			SuiteID:   dbRule.SuiteID,
			EvalType:  utils.InterfaceToString(dbRule.EvalType), // Convert interface{} to string
			MinScore:  minScore,
			MaxScore:  maxScore,
			Weight:    dbRule.Weight,
			HardFail:  dbRule.HardFail,
			Params:    dbRule.Params,
			CreatedAt: dbRule.CreatedAt,
		}
	}

	return rules, nil
}

// ListByEvalType retrieves all eval rules for an eval type
func (r *repository) ListByEvalType(ctx context.Context, evalType string) ([]Rule, error) {
	dbRules, err := r.queries.ListEvalRulesByEvalType(ctx, evalType)
	if err != nil {
		return nil, err
	}

	rules := make([]Rule, len(dbRules))
	for i, dbRule := range dbRules {
		var minScore, maxScore *float32
		if dbRule.MinScore.Valid {
			s := float32(dbRule.MinScore.Float64)
			minScore = &s
		}
		if dbRule.MaxScore.Valid {
			s := float32(dbRule.MaxScore.Float64)
			maxScore = &s
		}

		rules[i] = Rule{
			ID:        dbRule.ID,
			SuiteID:   dbRule.SuiteID,
			EvalType:  utils.InterfaceToString(dbRule.EvalType), // Convert interface{} to string
			MinScore:  minScore,
			MaxScore:  maxScore,
			Weight:    dbRule.Weight,
			HardFail:  dbRule.HardFail,
			Params:    dbRule.Params,
			CreatedAt: dbRule.CreatedAt,
		}
	}

	return rules, nil
}

// Update updates an eval rule
func (r *repository) Update(ctx context.Context, id uuid.UUID, rule UpdateRuleRequest) (*Rule, error) {
	// Get current rule to preserve values if not updating
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Prepare params
	var paramsBytes json.RawMessage
	if rule.Params != nil {
		paramsBytes = rule.Params
	} else {
		paramsBytes = current.Params
	}

	params := store.UpdateEvalRuleParams{
		ID:       id,
		Weight:   current.Weight,
		HardFail: current.HardFail,
		Params:   paramsBytes,
	}

	if rule.MinScore != nil {
		params.MinScore = sql.NullFloat64{Float64: float64(*rule.MinScore), Valid: true}
	} else {
		if current.MinScore != nil {
			params.MinScore = sql.NullFloat64{Float64: float64(*current.MinScore), Valid: true}
		} else {
			params.MinScore = sql.NullFloat64{Valid: false}
		}
	}

	if rule.MaxScore != nil {
		params.MaxScore = sql.NullFloat64{Float64: float64(*rule.MaxScore), Valid: true}
	} else {
		if current.MaxScore != nil {
			params.MaxScore = sql.NullFloat64{Float64: float64(*current.MaxScore), Valid: true}
		} else {
			params.MaxScore = sql.NullFloat64{Valid: false}
		}
	}

	if rule.Weight != nil {
		params.Weight = *rule.Weight
	} else {
		params.Weight = current.Weight
	}

	if rule.HardFail != nil {
		params.HardFail = *rule.HardFail
	} else {
		params.HardFail = current.HardFail
	}

	dbRule, err := r.queries.UpdateEvalRule(ctx, params)
	if err != nil {
		return nil, err
	}

	var minScore, maxScore *float32
	if dbRule.MinScore.Valid {
		s := float32(dbRule.MinScore.Float64)
		minScore = &s
	}
	if dbRule.MaxScore.Valid {
		s := float32(dbRule.MaxScore.Float64)
		maxScore = &s
	}

	return &Rule{
		ID:        dbRule.ID,
		SuiteID:   dbRule.SuiteID,
		EvalType:  utils.InterfaceToString(dbRule.EvalType), // Convert interface{} to string
		MinScore:  minScore,
		MaxScore:  maxScore,
		Weight:    dbRule.Weight,
		HardFail:  dbRule.HardFail,
		Params:    dbRule.Params,
		CreatedAt: dbRule.CreatedAt,
	}, nil
}

// Delete deletes an eval rule
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEvalRule(ctx, id)
}

// DeleteBySuite deletes all eval rules for a suite
func (r *repository) DeleteBySuite(ctx context.Context, suiteID uuid.UUID) error {
	return r.queries.DeleteEvalRulesBySuite(ctx, suiteID)
}
