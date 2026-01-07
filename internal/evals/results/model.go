package results

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
