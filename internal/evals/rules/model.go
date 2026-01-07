package rules

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
