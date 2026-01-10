package gcp

import (
	"context"
	
	"google.golang.org/genai"
)

// NewGenAIClient creates a new GenAI client configured for the Gemini API.
func NewGenAIClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	if apiKey == "" {
		return genai.NewClient(ctx, nil)
	}
	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
}
