package sessions

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Session represents a learning session
type Session struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	ModuleID  uuid.UUID `json:"module_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Attempt represents an attempt within a session
type Attempt struct {
	ID         uuid.UUID       `json:"id"`
	SessionID  uuid.UUID       `json:"session_id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	ArtifactID uuid.UUID       `json:"artifact_id"`
	IsCorrect  bool            `json:"is_correct"`
	UserAnswer json.RawMessage `json:"user_answer,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

// CreateSessionRequest represents the request to create a session
type CreateSessionRequest struct {
	TenantID uuid.UUID `json:"tenant_id"`
	ModuleID uuid.UUID `json:"module_id"`
	UserID   uuid.UUID `json:"user_id"`
}

// CreateAttemptRequest represents the request to create an attempt
type CreateAttemptRequest struct {
	SessionID  uuid.UUID       `json:"session_id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	ArtifactID uuid.UUID       `json:"artifact_id"`
	IsCorrect  bool            `json:"is_correct"`
	UserAnswer json.RawMessage `json:"user_answer,omitempty"`
}
