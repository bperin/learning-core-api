package document_graph

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/api/documentai/v1"

	"learning-core-api/internal/domain/documents"
	"learning-core-api/internal/gcp"
)

type Service struct {
	repo          *Repository
	documentsRepo documents.Repository
	gcsService    *gcp.GCSService
	docAIService  *gcp.DocumentAIService
}

func NewService(dbRepo *Repository, documentsRepo documents.Repository, gcsService *gcp.GCSService, docAIService *gcp.DocumentAIService) (*Service, error) {
	if dbRepo == nil {
		return nil, fmt.Errorf("graph repository is required")
	}
	if documentsRepo == nil {
		return nil, fmt.Errorf("documents repository is required")
	}
	if gcsService == nil {
		return nil, fmt.Errorf("gcs service is required")
	}
	if docAIService == nil {
		return nil, fmt.Errorf("document ai service is required")
	}

	return &Service{
		repo:          dbRepo,
		documentsRepo: documentsRepo,
		gcsService:    gcsService,
		docAIService:  docAIService,
	}, nil
}

func (s *Service) BuildGraph(ctx context.Context, documentID uuid.UUID) (*BuildResult, error) {
	if documentID == uuid.Nil {
		return nil, fmt.Errorf("document id is required")
	}

	_, _ = s.documentsRepo.UpdateRagStatus(ctx, documentID, documents.RagStatusProcessing)

	doc, err := s.documentsRepo.GetByID(ctx, documentID)
	if err != nil {
		_ = s.markError(ctx, documentID)
		return nil, err
	}
	if doc.StoragePath == nil || *doc.StoragePath == "" {
		_ = s.markError(ctx, documentID)
		return nil, fmt.Errorf("document storage path is required")
	}

	objectName, err := s.resolveObjectName(*doc.StoragePath, doc.StorageBucket)
	if err != nil {
		_ = s.markError(ctx, documentID)
		return nil, err
	}

	content, err := s.gcsService.DownloadFile(ctx, objectName)
	if err != nil {
		_ = s.markError(ctx, documentID)
		return nil, fmt.Errorf("failed to download document: %w", err)
	}

	mimeType := "application/pdf"
	if doc.MimeType != nil && *doc.MimeType != "" {
		mimeType = *doc.MimeType
	}

	processed, err := s.docAIService.ProcessDocument(ctx, content, mimeType)
	if err != nil {
		_ = s.markError(ctx, documentID)
		return nil, err
	}

	nodes, edges := s.buildGraphFromDocument(doc, processed)
	if err := s.repo.ReplaceGraph(ctx, documentID, nodes, edges); err != nil {
		_ = s.markError(ctx, documentID)
		return nil, err
	}

	_, _ = s.documentsRepo.UpdateRagStatus(ctx, documentID, documents.RagStatusReady)

	return &BuildResult{
		DocumentID:   documentID,
		NodesCreated: len(nodes),
		EdgesCreated: len(edges),
	}, nil
}

func (s *Service) QueryGraph(ctx context.Context, documentID uuid.UUID, query string, limit int) (*QueryResponse, error) {
	if documentID == uuid.Nil {
		return nil, fmt.Errorf("document id is required")
	}
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query is required")
	}
	if limit <= 0 {
		limit = 8
	}

	matched, err := s.repo.SearchNodes(ctx, documentID, query, limit)
	if err != nil {
		return nil, err
	}

	baseIDs := make([]uuid.UUID, 0, len(matched))
	for _, node := range matched {
		baseIDs = append(baseIDs, node.ID)
	}

	neighborNodes, edges, err := s.repo.FetchNeighbors(ctx, documentID, baseIDs, limit*4)
	if err != nil {
		return nil, err
	}

	unique := make(map[uuid.UUID]Node)
	for _, node := range matched {
		unique[node.ID] = node
	}
	for _, node := range neighborNodes {
		unique[node.ID] = node
	}

	combined := make([]Node, 0, len(unique))
	for _, node := range unique {
		combined = append(combined, node)
	}

	return &QueryResponse{
		DocumentID: documentID,
		Nodes:      combined,
		Edges:      edges,
	}, nil
}

func (s *Service) resolveObjectName(storagePath string, storageBucket *string) (string, error) {
	path := strings.TrimSpace(storagePath)
	if strings.HasPrefix(path, "gs://") {
		trimmed := strings.TrimPrefix(path, "gs://")
		parts := strings.SplitN(trimmed, "/", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid gcs uri: %s", storagePath)
		}
		bucket := parts[0]
		objectName := parts[1]
		if storageBucket != nil && *storageBucket != "" && bucket != *storageBucket {
			return "", fmt.Errorf("storage bucket mismatch: %s", bucket)
		}
		if bucket != s.gcsService.BucketName() {
			return "", fmt.Errorf("bucket %s is not supported by current gcs service", bucket)
		}
		return objectName, nil
	}

	if storageBucket != nil && *storageBucket != "" && *storageBucket != s.gcsService.BucketName() {
		return "", fmt.Errorf("storage bucket mismatch: %s", *storageBucket)
	}

	return path, nil
}

func (s *Service) buildGraphFromDocument(doc *documents.Document, processed *documentai.GoogleCloudDocumentaiV1Document) ([]Node, []Edge) {
	var nodes []Node
	var edges []Edge

	documentNodeID := uuid.New()
	docMeta, _ := json.Marshal(map[string]string{
		"filename": doc.Filename,
		"title":    safeString(doc.Title),
	})
	nodes = append(nodes, Node{
		ID:         documentNodeID,
		DocumentID: doc.ID,
		NodeType:   "document",
		Text:       safeString(doc.Title),
		Metadata:   docMeta,
	})

	if processed == nil {
		return nodes, edges
	}

	fullText := processed.Text

	for pageIndex, page := range processed.Pages {
		pageNumber := pageIndex + 1
		if page.PageNumber != 0 {
			pageNumber = int(page.PageNumber)
		}
		pageText := extractText(fullText, textAnchorFromLayout(page.Layout))
		pageID := uuid.New()
		pageMeta, _ := json.Marshal(map[string]int{
			"page_number": pageNumber,
		})
		nodes = append(nodes, Node{
			ID:         pageID,
			DocumentID: doc.ID,
			NodeType:   "page",
			Text:       pageText,
			PageNumber: &pageNumber,
			Metadata:   pageMeta,
		})
		edges = append(edges, Edge{
			ID:         uuid.New(),
			DocumentID: doc.ID,
			FromNodeID: documentNodeID,
			ToNodeID:   pageID,
			Relation:   "contains",
		})

		var prevParagraphID uuid.UUID
		for paragraphIndex, paragraph := range page.Paragraphs {
			paragraphText := extractText(fullText, textAnchorFromLayout(paragraph.Layout))
			if strings.TrimSpace(paragraphText) == "" {
				continue
			}
			paragraphID := uuid.New()
			paragraphMeta, _ := json.Marshal(map[string]int{
				"page_number": pageNumber,
				"position":    paragraphIndex + 1,
			})
			nodes = append(nodes, Node{
				ID:         paragraphID,
				DocumentID: doc.ID,
				NodeType:   "paragraph",
				Text:       paragraphText,
				PageNumber: &pageNumber,
				Metadata:   paragraphMeta,
			})
			edges = append(edges, Edge{
				ID:         uuid.New(),
				DocumentID: doc.ID,
				FromNodeID: pageID,
				ToNodeID:   paragraphID,
				Relation:   "contains",
			})
			if prevParagraphID != uuid.Nil {
				edges = append(edges, Edge{
					ID:         uuid.New(),
					DocumentID: doc.ID,
					FromNodeID: prevParagraphID,
					ToNodeID:   paragraphID,
					Relation:   "next",
				})
			}
			prevParagraphID = paragraphID
		}
	}

	return nodes, edges
}

func (s *Service) markError(ctx context.Context, documentID uuid.UUID) error {
	_, err := s.documentsRepo.UpdateRagStatus(ctx, documentID, documents.RagStatusError)
	return err
}

func extractText(fullText string, anchor *documentai.GoogleCloudDocumentaiV1DocumentTextAnchor) string {
	if anchor == nil || len(anchor.TextSegments) == 0 || fullText == "" {
		return ""
	}

	var builder strings.Builder
	for _, segment := range anchor.TextSegments {
		start := int(segment.StartIndex)
		end := int(segment.EndIndex)
		if start < 0 {
			start = 0
		}
		if end > len(fullText) {
			end = len(fullText)
		}
		if start >= len(fullText) || start >= end {
			continue
		}
		builder.WriteString(fullText[start:end])
	}

	return strings.TrimSpace(builder.String())
}

func textAnchorFromLayout(layout *documentai.GoogleCloudDocumentaiV1DocumentPageLayout) *documentai.GoogleCloudDocumentaiV1DocumentTextAnchor {
	if layout == nil {
		return nil
	}
	return layout.TextAnchor
}

func safeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
