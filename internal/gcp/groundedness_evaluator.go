package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/genai"

	"learning-core-api/internal/domain/evals"
)

// GroundednessEvaluator implements the GroundednessEvaluator interface using Vertex AI
type GroundednessEvaluator struct {
	client     *genai.Client
	evalPrompt string
}

// NewGroundednessEvaluator creates a new groundedness evaluator with a stored prompt
func NewGroundednessEvaluator(client *genai.Client, evalPrompt string) *GroundednessEvaluator {
	return &GroundednessEvaluator{
		client:     client,
		evalPrompt: evalPrompt,
	}
}

// NewGroundednessEvaluatorFromAPIKey creates a new groundedness evaluator from API key
func NewGroundednessEvaluatorFromAPIKey(ctx context.Context, apiKey string, evalPrompt string) (*GroundednessEvaluator, error) {
	client, err := NewGenAIClient(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	return NewGroundednessEvaluator(client, evalPrompt), nil
}

// GroundingMetadata represents the grounding metadata from file search
type GroundingMetadata struct {
	GroundingChunks   []GroundingChunk   `json:"groundingChunks,omitempty"`
	GroundingSupports []GroundingSupport `json:"groundingSupports,omitempty"`
}

type GroundingChunk struct {
	RetrievedContext RetrievedContext `json:"retrievedContext"`
}

type RetrievedContext struct {
	Text  string `json:"text"`
	Title string `json:"title"`
}

type GroundingSupport struct {
	GroundingChunkIndices []int   `json:"groundingChunkIndices"`
	Segment               Segment `json:"segment"`
}

type Segment struct {
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
	Text       string `json:"text"`
}

// EvaluateGroundednessResult is the structured response from groundedness evaluation
type EvaluateGroundednessResult struct {
	IsGrounded        bool     `json:"is_grounded"`
	UnsupportedClaims []string `json:"unsupported_claims"`
	GroundednessScore float64  `json:"groundedness_score"`
}

// EvaluateGroundedness evaluates if the response is grounded in the provided context
// Uses stored eval prompt template and Vertex AI for evaluation
func (e *GroundednessEvaluator) EvaluateGroundedness(ctx context.Context, question string, expectedAnswer string, groundingMetadata json.RawMessage) (*evals.GroundednessResult, error) {
	if e.client == nil {
		return nil, fmt.Errorf("genai client is required")
	}

	if e.evalPrompt == "" {
		return nil, fmt.Errorf("eval prompt is required")
	}

	// Parse grounding metadata to extract context
	var grounding GroundingMetadata
	if err := json.Unmarshal(groundingMetadata, &grounding); err != nil {
		return nil, fmt.Errorf("failed to parse grounding metadata: %w", err)
	}

	// Extract supporting context from grounding chunks
	referenceContext := ""
	for _, chunk := range grounding.GroundingChunks {
		if chunk.RetrievedContext.Text != "" {
			referenceContext += chunk.RetrievedContext.Text + "\n\n"
		}
	}

	if referenceContext == "" {
		return nil, fmt.Errorf("no supporting context found in grounding metadata")
	}

	// Build evaluation request using stored prompt template
	// Template variables: {{context}} and {{response}}
	evaluationPrompt := e.evalPrompt
	evaluationPrompt = fmt.Sprintf(evaluationPrompt, referenceContext, expectedAnswer)

	// Call Gemini with stored prompt for evaluation
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: evaluationPrompt},
			},
		},
	}

	genConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	resp, err := e.client.Models.GenerateContent(ctx, "gemini-1.5-pro", contents, genConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to call genai for groundedness evaluation: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned from groundedness evaluation")
	}

	// Extract the response text
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
	}

	// Parse the evaluation result
	var evalResult EvaluateGroundednessResult
	if err := json.Unmarshal([]byte(responseText), &evalResult); err != nil {
		return nil, fmt.Errorf("failed to parse groundedness evaluation result: %w", err)
	}

	// Convert to GroundednessResult
	verdict := "WARN"
	if evalResult.IsGrounded {
		verdict = "PASS"
	} else if len(evalResult.UnsupportedClaims) > 0 {
		verdict = "FAIL"
	}

	return &evals.GroundednessResult{
		Score:              evalResult.GroundednessScore,
		Verdict:            verdict,
		Reasoning:          fmt.Sprintf("Grounded: %v, Unsupported claims: %d", evalResult.IsGrounded, len(evalResult.UnsupportedClaims)),
		SupportingSegments: evalResult.UnsupportedClaims,
		CreatedAt:          time.Now().UTC().Format(time.RFC3339),
	}, nil
}
