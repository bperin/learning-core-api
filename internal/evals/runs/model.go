package runs

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
