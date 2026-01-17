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

type fileSearchToolConfig struct {
	StoreNames     []string `json:"store_names"`
	MetadataFilter string   `json:"metadata_filter,omitempty"`
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
			var cfg fileSearchToolConfig
			if len(tool.Config) == 0 {
				return nil, fmt.Errorf("file_search tool config is required")
			}
			if err := json.Unmarshal(tool.Config, &cfg); err != nil {
				return nil, fmt.Errorf("failed to parse file_search config: %w", err)
			}
			if len(cfg.StoreNames) == 0 {
				return nil, fmt.Errorf("file_search store_names is required")
			}
			genConfig.Tools = append(genConfig.Tools, &genai.Tool{
				FileSearch: &genai.FileSearch{
					FileSearchStoreNames: cfg.StoreNames,
					MetadataFilter:       cfg.MetadataFilter,
				},
			})
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

	// Extract grounding metadata if available
	var groundingMetadata json.RawMessage
	if resp.Candidates[0].GroundingMetadata != nil {
		if groundingBytes, err := json.Marshal(resp.Candidates[0].GroundingMetadata); err == nil {
			groundingMetadata = groundingBytes
		}
	}

	return &generation.GeneratorResponse{
		OutputText:        outputText,
		FinishReason:      string(resp.Candidates[0].FinishReason),
		ModelUsed:         modelName,
		GroundingMetadata: groundingMetadata,
	}, nil
}
