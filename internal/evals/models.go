package evals

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Run represents an evaluation run
type Run struct {
	ID              uuid.UUID       `json:"id"`
	ArtifactID      uuid.UUID       `json:"artifact_id"`
	GenerationRunID *uuid.UUID      `json:"generation_run_id,omitempty"`
	SuiteID         uuid.UUID       `json:"suite_id"`
	JudgeModel      string          `json:"judge_model"`
	JudgeParams     json.RawMessage `json:"judge_params"`
	Status          string          `json:"status"`
	OverallPass     *bool           `json:"overall_pass,omitempty"`
	OverallScore    *float32        `json:"overall_score,omitempty"`
	Error           json.RawMessage `json:"error"`
	StartedAt       time.Time       `json:"started_at"`
	FinishedAt      *time.Time      `json:"finished_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// CreateRunRequest represents the request to create an evaluation run
type CreateRunRequest struct {
	ArtifactID      uuid.UUID       `json:"artifact_id"`
	GenerationRunID *uuid.UUID      `json:"generation_run_id,omitempty"`
	SuiteID         uuid.UUID       `json:"suite_id"`
	JudgeModel      string          `json:"judge_model"`
	JudgeParams     json.RawMessage `json:"judge_params"`
}

// UpdateRunResultRequest represents the request to update run results
type UpdateRunResultRequest struct {
	Status       string          `json:"status"`
	OverallPass  *bool           `json:"overall_pass,omitempty"`
	OverallScore *float32        `json:"overall_score,omitempty"`
	Error        json.RawMessage `json:"error"`
}

// Result represents an evaluation result
type Result struct {
	ID        uuid.UUID       `json:"id"`
	EvalRunID uuid.UUID       `json:"eval_run_id"`
	RuleID    uuid.UUID       `json:"rule_id"`
	Pass      bool            `json:"pass"`
	Score     *float32        `json:"score,omitempty"`
	Details   json.RawMessage `json:"details"`
	CreatedAt time.Time       `json:"created_at"`
}

// CreateResultRequest represents the request to create an evaluation result
type CreateResultRequest struct {
	EvalRunID uuid.UUID       `json:"eval_run_id"`
	RuleID    uuid.UUID       `json:"rule_id"`
	Pass      bool            `json:"pass"`
	Score     *float32        `json:"score,omitempty"`
	Details   json.RawMessage `json:"details"`
}

// Rule represents an evaluation rule that defines thresholds for an eval type within a suite
type Rule struct {
	ID        uuid.UUID       `json:"id"`
	SuiteID   uuid.UUID       `json:"suite_id"`
	EvalType  string          `json:"eval_type"`
	MinScore  *float32        `json:"min_score,omitempty"`
	MaxScore  *float32        `json:"max_score,omitempty"`
	Weight    float32         `json:"weight"`
	HardFail  bool            `json:"hard_fail"`
	Params    json.RawMessage `json:"params"`
	CreatedAt time.Time       `json:"created_at"`
}

// CreateRuleRequest represents the request to create an evaluation rule
type CreateRuleRequest struct {
	SuiteID  uuid.UUID       `json:"suite_id"`
	EvalType string          `json:"eval_type"`
	MinScore *float32        `json:"min_score,omitempty"`
	MaxScore *float32        `json:"max_score,omitempty"`
	Weight   float32         `json:"weight"`
	HardFail bool            `json:"hard_fail"`
	Params   json.RawMessage `json:"params"`
}

// UpdateRuleRequest represents the request to update an evaluation rule
type UpdateRuleRequest struct {
	MinScore *float32        `json:"min_score,omitempty"`
	MaxScore *float32        `json:"max_score,omitempty"`
	Weight   *float32        `json:"weight,omitempty"`
	HardFail *bool           `json:"hard_fail,omitempty"`
	Params   json.RawMessage `json:"params,omitempty"`
}

// Suite represents an evaluation suite
type Suite struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateSuiteRequest represents the request to create an evaluation suite
type CreateSuiteRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
