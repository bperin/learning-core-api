package evals

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"
	"learning-core-api/internal/utils"

	"github.com/google/uuid"
)

// RuleRepository defines the interface for eval rule operations.
type RuleRepository interface {
	Create(ctx context.Context, rule Rule) (*Rule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Rule, error)
	GetBySuiteAndType(ctx context.Context, suiteID uuid.UUID, evalType string) (*Rule, error)
	ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]Rule, error)
	ListByEvalType(ctx context.Context, evalType string) ([]Rule, error)
	Update(ctx context.Context, id uuid.UUID, rule UpdateRuleRequest) (*Rule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBySuite(ctx context.Context, suiteID uuid.UUID) error
}

// ruleRepository implements the RuleRepository interface.
type ruleRepository struct {
	queries *store.Queries
}

// NewRuleRepository creates a new eval rule repository.
func NewRuleRepository(queries *store.Queries) RuleRepository {
	return &ruleRepository{
		queries: queries,
	}
}

// Create creates a new eval rule.
func (r *ruleRepository) Create(ctx context.Context, rule Rule) (*Rule, error) {
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

	rule := toDomainRule(dbRule)
	return &rule, nil
}

// GetByID retrieves an eval rule by ID.
func (r *ruleRepository) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	dbRule, err := r.queries.GetEvalRule(ctx, id)
	if err != nil {
		return nil, err
	}

	rule := toDomainRule(dbRule)
	return &rule, nil
}

// GetBySuiteAndType retrieves an eval rule by suite ID and eval type.
func (r *ruleRepository) GetBySuiteAndType(ctx context.Context, suiteID uuid.UUID, evalType string) (*Rule, error) {
	params := store.GetEvalRuleBySuiteAndTypeParams{
		SuiteID:  suiteID,
		EvalType: evalType,
	}

	dbRule, err := r.queries.GetEvalRuleBySuiteAndType(ctx, params)
	if err != nil {
		return nil, err
	}

	rule := toDomainRule(dbRule)
	return &rule, nil
}

// ListBySuite retrieves all eval rules for a suite.
func (r *ruleRepository) ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]Rule, error) {
	dbRules, err := r.queries.ListEvalRulesBySuite(ctx, suiteID)
	if err != nil {
		return nil, err
	}

	return toDomainRules(dbRules), nil
}

// ListByEvalType retrieves all eval rules for an eval type.
func (r *ruleRepository) ListByEvalType(ctx context.Context, evalType string) ([]Rule, error) {
	dbRules, err := r.queries.ListEvalRulesByEvalType(ctx, evalType)
	if err != nil {
		return nil, err
	}

	return toDomainRules(dbRules), nil
}

// Update updates an eval rule.
func (r *ruleRepository) Update(ctx context.Context, id uuid.UUID, rule UpdateRuleRequest) (*Rule, error) {
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

	rule := toDomainRule(dbRule)
	return &rule, nil
}

// Delete deletes an eval rule.
func (r *ruleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEvalRule(ctx, id)
}

// DeleteBySuite deletes all eval rules for a suite.
func (r *ruleRepository) DeleteBySuite(ctx context.Context, suiteID uuid.UUID) error {
	return r.queries.DeleteEvalRulesBySuite(ctx, suiteID)
}

func toDomainRule(dbRule store.EvalRule) Rule {
	return Rule{
		ID:        dbRule.ID,
		SuiteID:   dbRule.SuiteID,
		EvalType:  utils.InterfaceToString(dbRule.EvalType),
		MinScore:  utils.NullFloat64ToFloat32Ptr(dbRule.MinScore),
		MaxScore:  utils.NullFloat64ToFloat32Ptr(dbRule.MaxScore),
		Weight:    dbRule.Weight,
		HardFail:  dbRule.HardFail,
		Params:    dbRule.Params,
		CreatedAt: dbRule.CreatedAt,
	}
}

func toDomainRules(dbRules []store.EvalRule) []Rule {
	rules := make([]Rule, len(dbRules))
	for i, dbRule := range dbRules {
		rules[i] = toDomainRule(dbRule)
	}
	return rules
}
