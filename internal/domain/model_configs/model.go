package model_configs

import (
	"time"

	"github.com/google/uuid"
)

// ModelConfig represents a versioned model configuration stored in the database.
type ModelConfig struct {
	ID          uuid.UUID `json:"id"`
	Version     int32     `json:"version"`
	ModelName   string    `json:"model_name"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int32     `json:"max_tokens"`
	TopP        float64   `json:"top_p"`
	TopK        float64   `json:"top_k"`
	MimeType    string    `json:"mime_type"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateModelConfigRequest represents the data needed to create a model config.
type CreateModelConfigRequest struct {
	ModelName   string    `json:"model_name"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int32     `json:"max_tokens"`
	TopP        float64   `json:"top_p"`
	TopK        float64   `json:"top_k"`
	MimeType    string    `json:"mime_type"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   uuid.UUID `json:"created_by"`
}
