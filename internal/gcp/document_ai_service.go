package gcp

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/api/documentai/v1"
	"google.golang.org/api/option"
)

type DocumentAIService struct {
	client        *documentai.Service
	processorName string
}

func NewDocumentAIService(ctx context.Context, projectID, location, processorID, apiKey string) (*DocumentAIService, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project id is required")
	}
	if location == "" {
		return nil, fmt.Errorf("document ai location is required")
	}
	if processorID == "" {
		return nil, fmt.Errorf("document ai processor id is required")
	}

	opts := []option.ClientOption{option.WithScopes(documentai.CloudPlatformScope)}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}

	client, err := documentai.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create document ai client: %w", err)
	}

	processorName := fmt.Sprintf("projects/%s/locations/%s/processors/%s", projectID, location, processorID)
	return &DocumentAIService{
		client:        client,
		processorName: processorName,
	}, nil
}

func (s *DocumentAIService) ProcessDocument(ctx context.Context, content []byte, mimeType string) (*documentai.GoogleCloudDocumentaiV1Document, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("document ai client is required")
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("document content is required")
	}
	if mimeType == "" {
		mimeType = "application/pdf"
	}

	req := &documentai.GoogleCloudDocumentaiV1ProcessRequest{
		RawDocument: &documentai.GoogleCloudDocumentaiV1RawDocument{
			Content:  base64.StdEncoding.EncodeToString(content),
			MimeType: mimeType,
		},
	}

	resp, err := s.client.Projects.Locations.Processors.Process(s.processorName, req).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("document ai processing failed: %w", err)
	}
	if resp == nil || resp.Document == nil {
		return nil, fmt.Errorf("document ai response missing document")
	}

	return resp.Document, nil
}
