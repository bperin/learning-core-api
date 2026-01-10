package gcp

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"

	"learning-core-api/internal/domain/generation"
)

type GenerationService struct {
	client *genai.Client
}

func NewGenerationService(client *genai.Client) (*GenerationService, error) {
	if client == nil {
		return nil, fmt.Errorf("genai client is required")
	}
	return &GenerationService{
		client: client,
	}, nil
}

func NewGenerationServiceFromAPIKey(ctx context.Context, apiKey string) (*GenerationService, error) {
	client, err := NewGenAIClient(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	return NewGenerationService(client)
}

func (s *GenerationService) Generate(ctx context.Context, req generation.GeneratorRequest) (*generation.GeneratorResponse, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("genai client is required")
	}
	if req.Model == nil {
		return nil, fmt.Errorf("model config is required")
	}
	if req.Model.Name == "" {
		return nil, fmt.Errorf("model name is required")
	}

	modelName := req.Model.Name
	genConfig := &genai.GenerateContentConfig{}

	if req.Model != nil {
		if req.Model.Temperature != nil {
			genConfig.Temperature = req.Model.Temperature
		}
		if req.Model.MaxTokens != nil {
			genConfig.MaxOutputTokens = *req.Model.MaxTokens
		}
		if req.Model.TopP != nil {
			genConfig.TopP = req.Model.TopP
		}
		if req.Model.TopK != nil {
			tk := float32(*req.Model.TopK)
			genConfig.TopK = &tk
		}
		if req.Model.MimeType != "" {
			genConfig.ResponseMIMEType = req.Model.MimeType
		}
	}

	if req.SystemInstruction != "" {
		genConfig.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{
				{Text: req.SystemInstruction},
			},
		}
	}

	if len(req.OutputSchema) > 0 {
		schema := &genai.Schema{}
		if err := json.Unmarshal(req.OutputSchema, schema); err != nil {
			return nil, fmt.Errorf("failed to parse response schema: %w", err)
		}
		genConfig.ResponseMIMEType = "application/json"
		genConfig.ResponseSchema = schema
	}

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: req.Prompt},
			},
		},
	}

	for _, tool := range req.Tools {
		if tool.Type == "file_search" {
			// TODO: Configure file search tool integration.
		}
	}

	resp, err := s.client.Models.GenerateContent(ctx, modelName, contents, genConfig)
	if err != nil {
		return nil, fmt.Errorf("genai call failed: %w", err)
	}
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned")
	}

	var outputText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			outputText += part.Text
		}
	}

	return &generation.GeneratorResponse{
		OutputText:   outputText,
		FinishReason: string(resp.Candidates[0].FinishReason),
		ModelUsed:    modelName,
	}, nil
}
