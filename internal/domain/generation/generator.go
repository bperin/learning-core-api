package generation

import (
	"context"
	"encoding/json"
)

// Generator defines the boundary for provider-specific generation calls.
type Generator interface {
	Generate(ctx context.Context, req GeneratorRequest) (*GeneratorResponse, error)
}

type GeneratorRequest struct {
	Prompt            string
	SystemInstruction string
	OutputSchema      json.RawMessage
	Tools             []ToolConfig
	Model             *ModelConfig
}

type GeneratorResponse struct {
	OutputText        string
	FinishReason      string
	ModelUsed         string
	GroundingMetadata json.RawMessage
}
