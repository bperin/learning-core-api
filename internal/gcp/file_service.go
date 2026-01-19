package gcp

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
)

type FileService struct {
	gcs             *GCSService
	client          *genai.Client
	fileSearchStore *genai.FileSearchStores
	storeName       string
	chunkingConfig  *genai.ChunkingConfig
}

type UploadURLResponse struct {
	BucketName string
	ObjectName string
	UploadURL  string
}

type FileSearchUploadResult struct {
	StoreName  string
	FileName   string
	Operation  *genai.UploadToFileSearchStoreOperation
	ObjectName string
}

func NewFileService(ctx context.Context, gcs *GCSService, apiKey string, storeName string) (*FileService, error) {
	if gcs == nil {
		return nil, fmt.Errorf("gcs service is required")
	}
	if storeName == "" {
		return nil, fmt.Errorf("file search store name is required")
	}

	client, err := newGenAIClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	s := &FileService{
		gcs:             gcs,
		client:          client,
		fileSearchStore: client.FileSearchStores,
		storeName:       storeName,
		chunkingConfig:  defaultChunkingConfig(),
	}

	if err := s.EnsureStore(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *FileService) StoreName() string {
	return s.storeName
}

func (s *FileService) EnsureStore(ctx context.Context) error {
	if s.storeName == "" {
		return fmt.Errorf("file search store name is required")
	}
	if s.fileSearchStore == nil {
		return fmt.Errorf("file search store client is required")
	}
	store, err := s.fileSearchStore.Get(ctx, s.storeName, nil)
	if err == nil {
		if store != nil && store.Name != "" {
			s.storeName = store.Name
		}
		return nil
	}

	if !isNotFoundError(err) {
		if matched, matchErr := s.matchStoreByDisplayName(ctx); matchErr == nil && matched {
			return nil
		}
		return fmt.Errorf("failed to fetch file search store %q: %w", s.storeName, err)
	}

	if matched, matchErr := s.matchStoreByDisplayName(ctx); matchErr == nil && matched {
		return nil
	}

	created, err := s.fileSearchStore.Create(ctx, &genai.CreateFileSearchStoreConfig{
		DisplayName: s.storeName,
	})
	if err != nil {
		return fmt.Errorf("failed to create file search store %q: %w", s.storeName, err)
	}
	if created != nil && created.Name != "" {
		s.storeName = created.Name
		return nil
	}

	return s.waitForStore(ctx)
}

func (s *FileService) GenerateUploadURL(ctx context.Context, documentID uuid.UUID, filename string, contentType string, ttl time.Duration) (*UploadURLResponse, error) {
	objectName := s.objectName(documentID, filename)
	url, err := s.gcs.GenerateSignedUploadURL(ctx, objectName, contentType, ttl)
	if err != nil {
		return nil, err
	}

	return &UploadURLResponse{
		BucketName: s.gcs.BucketName(),
		ObjectName: objectName,
		UploadURL:  url,
	}, nil
}

func (s *FileService) UploadToFileSearchStore(ctx context.Context, objectName string, displayName string, mimeType string) (*FileSearchUploadResult, error) {
	if s.fileSearchStore == nil {
		return nil, fmt.Errorf("file search store client is required")
	}
	if s.storeName == "" {
		return nil, fmt.Errorf("file search store name is required")
	}

	reader, err := s.gcs.OpenReader(ctx, objectName)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	config := &genai.UploadToFileSearchStoreConfig{
		MIMEType:       mimeType,
		DisplayName:    displayName,
		ChunkingConfig: s.chunkingConfig,
	}

	operation, err := s.fileSearchStore.UploadToFileSearchStore(ctx, reader, s.storeName, config)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to file search store: %w", err)
	}

	return &FileSearchUploadResult{
		StoreName:  s.storeName,
		FileName:   extractDocumentName(operation),
		Operation:  operation,
		ObjectName: objectName,
	}, nil
}

func (s *FileService) ListStores(ctx context.Context) ([]*genai.FileSearchStore, error) {
	if s.client == nil || s.client.FileSearchStores == nil {
		return nil, fmt.Errorf("file search store client is required")
	}
	stores := []*genai.FileSearchStore{}
	for store, err := range s.client.FileSearchStores.All(ctx) {
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (s *FileService) ListFiles(ctx context.Context, storeName string) ([]*genai.Document, error) {
	if s.client == nil || s.client.FileSearchStores == nil || s.client.FileSearchStores.Documents == nil {
		return nil, fmt.Errorf("file search store documents client is required")
	}
	if storeName == "" {
		storeName = s.storeName
	}
	if storeName == "" {
		return nil, fmt.Errorf("file search store name is required")
	}

	documents := []*genai.Document{}
	for doc, err := range s.client.FileSearchStores.Documents.All(ctx, storeName) {
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func (s *FileService) DeleteAllFiles(ctx context.Context, storeName string) error {
	if s.client == nil || s.client.FileSearchStores == nil || s.client.FileSearchStores.Documents == nil {
		return fmt.Errorf("file search store documents client is required")
	}
	if storeName == "" {
		storeName = s.storeName
	}
	if storeName == "" {
		return nil // Nothing to delete
	}

	files, err := s.ListFiles(ctx, storeName)
	if err != nil {
		return fmt.Errorf("failed to list files for deletion in store %q: %w", storeName, err)
	}

	for _, doc := range files {
		if err := s.client.FileSearchStores.Documents.Delete(ctx, doc.Name, nil); err != nil {
			return fmt.Errorf("failed to delete file %q from store %q: %w", doc.Name, storeName, err)
		}
	}
	return nil
}

func (s *FileService) DeleteStore(ctx context.Context, storeName string) error {
	if s.client == nil || s.client.FileSearchStores == nil {
		return fmt.Errorf("file search store client is required")
	}
	if storeName == "" {
		storeName = s.storeName
	}
	if storeName == "" {
		return nil
	}

	force := true
	config := &genai.DeleteFileSearchStoreConfig{
		Force: &force,
	}

	if err := s.client.FileSearchStores.Delete(ctx, storeName, config); err != nil {
		return fmt.Errorf("failed to delete file search store %q: %w", storeName, err)
	}
	return nil
}

func (s *FileService) ClearAllStores(ctx context.Context) error {
	stores, err := s.ListStores(ctx)
	if err != nil {
		return fmt.Errorf("failed to list stores for clearing: %w", err)
	}

	for _, store := range stores {
		if err := s.DeleteStore(ctx, store.Name); err != nil {
			// If we fail to delete the store, try deleting files as fallback
			_ = s.DeleteAllFiles(ctx, store.Name)
			return err
		}
	}
	return nil
}

func (s *FileService) matchStoreByDisplayName(ctx context.Context) (bool, error) {
	if s.client == nil || s.client.FileSearchStores == nil {
		return false, fmt.Errorf("file search store client is required")
	}
	for store, err := range s.client.FileSearchStores.All(ctx) {
		if err != nil {
			return false, err
		}
		if store == nil {
			continue
		}
		if store.Name == s.storeName || store.DisplayName == s.storeName {
			s.storeName = store.Name
			return true, nil
		}
	}
	return false, nil
}

func (s *FileService) waitForStore(ctx context.Context) error {
	deadline, hasDeadline := ctx.Deadline()
	for {
		matched, err := s.matchStoreByDisplayName(ctx)
		if err != nil {
			return err
		}
		if matched {
			return nil
		}
		if hasDeadline && time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for file search store %q", s.storeName)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}
}

func isNotFoundError(err error) bool {
	var apiErr genai.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == 404
	}
	return false
}

func (s *FileService) objectName(documentID uuid.UUID, filename string) string {
	cleanName := strings.TrimSpace(filename)
	if cleanName == "" {
		cleanName = "document"
	}
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	cleanName = filepath.Base(cleanName)
	return fmt.Sprintf("documents/%s/%s", documentID.String(), cleanName)
}

func defaultChunkingConfig() *genai.ChunkingConfig {
	chunkSize := int32(512)
	overlap := int32(50)
	return &genai.ChunkingConfig{
		WhiteSpaceConfig: &genai.WhiteSpaceConfig{
			MaxTokensPerChunk: &chunkSize,
			MaxOverlapTokens:  &overlap,
		},
	}
}

func extractDocumentName(operation *genai.UploadToFileSearchStoreOperation) string {
	if operation == nil || operation.Response == nil {
		return ""
	}
	return operation.Response.DocumentName
}

func newGenAIClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	if apiKey == "" {
		return genai.NewClient(ctx, nil)
	}
	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
}
