package chunking_configs

import (
	"time"

	"github.com/google/uuid"
)

// ChunkingConfig represents a chunking configuration stored in the database.
type ChunkingConfig struct {
	ID           uuid.UUID `json:"id"`
	Version      int32     `json:"version"`
	ChunkSize    int32     `json:"chunk_size"`
	ChunkOverlap int32     `json:"chunk_overlap"`
	IsActive     bool      `json:"is_active"`
	CreatedBy    uuid.UUID `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateChunkingConfigRequest represents data needed to create a chunking config.
type CreateChunkingConfigRequest struct {
	ChunkSize    int32     `json:"chunk_size"`
	ChunkOverlap int32     `json:"chunk_overlap"`
	IsActive     *bool     `json:"is_active,omitempty"`
	CreatedBy    uuid.UUID `json:"created_by"`
}
