package eval_results

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EvalResult represents the outcome of an evaluation
type EvalResult struct {
	ID                uuid.UUID       `json:"id"`
	EvalItemID        uuid.UUID       `json:"eval_item_id"`
	EvalType          string          `json:"eval_type"`
	EvalPromptID      uuid.UUID       `json:"eval_prompt_id"`
	Score             *float64        `json:"score,omitempty"`
	IsGrounded        *bool           `json:"is_grounded,omitempty"`
	Verdict           string          `json:"verdict"`
	Reasoning         *string         `json:"reasoning,omitempty"`
	UnsupportedClaims json.RawMessage `json:"unsupported_claims,omitempty"`
	GCPEvalID         *string         `json:"gcp_eval_id,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}

// CreateEvalResultRequest represents a request to create an eval result
type CreateEvalResultRequest struct {
	EvalItemID        uuid.UUID       `json:"eval_item_id" validate:"required"`
	EvalType          string          `json:"eval_type" validate:"required,min=1"`
	EvalPromptID      uuid.UUID       `json:"eval_prompt_id" validate:"required"`
	Score             *float64        `json:"score,omitempty"`
	IsGrounded        *bool           `json:"is_grounded,omitempty"`
	Verdict           string          `json:"verdict" validate:"required,oneof=PASS FAIL WARN"`
	Reasoning         *string         `json:"reasoning,omitempty"`
	UnsupportedClaims json.RawMessage `json:"unsupported_claims,omitempty"`
	GCPEvalID         *string         `json:"gcp_eval_id,omitempty"`
}

// Validate validates the CreateEvalResultRequest
func (r *CreateEvalResultRequest) Validate() error {
	if r.EvalItemID == uuid.Nil {
		return ErrInvalidEvalItemID
	}
	if r.EvalType == "" {
		return ErrInvalidEvalType
	}
	if r.EvalPromptID == uuid.Nil {
		return ErrInvalidEvalPromptID
	}
	if r.Verdict == "" {
		return ErrInvalidVerdict
	}
	if r.Verdict != "PASS" && r.Verdict != "FAIL" && r.Verdict != "WARN" {
		return ErrInvalidVerdict
	}
	return nil
}

// EvalResultStats represents aggregate statistics for eval results
type EvalResultStats struct {
	TotalEvals int64   `json:"total_evals"`
	Passed     int64   `json:"passed"`
	Failed     int64   `json:"failed"`
	Warned     int64   `json:"warned"`
	AvgScore   float64 `json:"avg_score"`
	PassRate   float64 `json:"pass_rate"`
}
