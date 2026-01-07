package generation

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// GenerationRun represents a generation run
type GenerationRun struct {
	ID             uuid.UUID       `json:"id"`
	ModuleID       uuid.UUID       `json:"module_id"`
	AgentName      string          `json:"agent_name"`
	AgentVersion   string          `json:"agent_version"`
	Model          string          `json:"model"`
	ModelParams    json.RawMessage `json:"model_params"`
	PromptID       *uuid.UUID      `json:"prompt_id,omitempty"`
	StoreName      string          `json:"store_name"`
	MetadataFilter json.RawMessage `json:"metadata_filter"`
	Status         RunStatus       `json:"status"`
	InputPayload   json.RawMessage `json:"input_payload"`
	OutputPayload  json.RawMessage `json:"output_payload"`
	Error          json.RawMessage `json:"error"`
	StartedAt      *time.Time      `json:"started_at,omitempty"`
	FinishedAt     *time.Time      `json:"finished_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

// RunStatus represents the status of a run
type RunStatus string

const (
	RunStatusPending   RunStatus = "PENDING"
	RunStatusRunning   RunStatus = "RUNNING"
	RunStatusCompleted RunStatus = "COMPLETED"
	RunStatusFailed    RunStatus = "FAILED"
)

// Artifact represents an artifact generated during a run
type Artifact struct {
	ID              uuid.UUID       `json:"id"`
	ModuleID        uuid.UUID       `json:"module_id"`
	GenerationRunID uuid.UUID       `json:"generation_run_id"`
	Type            ArtifactType    `json:"type"`
	Status          ArtifactStatus  `json:"status"`
	SchemaVersion   string          `json:"schema_version"`
	Difficulty      *string         `json:"difficulty,omitempty"`
	Tags            []string        `json:"tags"`
	ArtifactPayload json.RawMessage `json:"artifact_payload"`
	Grounding       json.RawMessage `json:"grounding"`
	EvidenceVersion *string         `json:"evidence_version,omitempty"`
	ApprovedAt      *time.Time      `json:"approved_at,omitempty"`
	RejectedAt      *time.Time      `json:"rejected_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ArtifactType represents the type of artifact
type ArtifactType string

const (
	ArtifactTypeQuestion    ArtifactType = "QUESTION"
	ArtifactTypeExplanation ArtifactType = "EXPLANATION"
	ArtifactTypeSummary     ArtifactType = "SUMMARY"
	ArtifactTypeDocument    ArtifactType = "DOCUMENT"
)

// ArtifactStatus represents the status of an artifact
type ArtifactStatus string

const (
	ArtifactStatusPendingEval ArtifactStatus = "PENDING_EVAL"
	ArtifactStatusApproved    ArtifactStatus = "APPROVED"
	ArtifactStatusRejected    ArtifactStatus = "REJECTED"
)

// CreateGenerationRunRequest represents the request to create a generation run
type CreateGenerationRunRequest struct {
	ModuleID       uuid.UUID       `json:"module_id"`
	AgentName      string          `json:"agent_name"`
	AgentVersion   string          `json:"agent_version"`
	Model          string          `json:"model"`
	ModelParams    json.RawMessage `json:"model_params"`
	PromptID       *uuid.UUID      `json:"prompt_id,omitempty"`
	StoreName      string          `json:"store_name"`
	MetadataFilter json.RawMessage `json:"metadata_filter"`
	InputPayload   json.RawMessage `json:"input_payload"`
}

// CreateArtifactRequest represents the request to create an artifact
type CreateArtifactRequest struct {
	ModuleID        uuid.UUID       `json:"module_id"`
	GenerationRunID uuid.UUID       `json:"generation_run_id"`
	Type            ArtifactType    `json:"type"`
	SchemaVersion   string          `json:"schema_version"`
	Difficulty      *string         `json:"difficulty,omitempty"`
	Tags            []string        `json:"tags"`
	ArtifactPayload json.RawMessage `json:"artifact_payload"`
	Grounding       json.RawMessage `json:"grounding"`
	EvidenceVersion *string         `json:"evidence_version,omitempty"`
}
