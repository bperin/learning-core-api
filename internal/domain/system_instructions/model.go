package system_instructions

import (
	"time"

	"github.com/google/uuid"
)

// SystemInstruction represents a system instruction stored in the database.
type SystemInstruction struct {
	ID        uuid.UUID `json:"id"`
	Version   int32     `json:"version"`
	Text      string    `json:"text"`
	IsActive  bool      `json:"is_active"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateSystemInstructionRequest represents data needed to create a system instruction.
type CreateSystemInstructionRequest struct {
	Text      string    `json:"text"`
	IsActive  *bool     `json:"is_active,omitempty"`
	CreatedBy uuid.UUID `json:"created_by"`
}
