package filesearch

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"learning-core-api/internal/domain/documents"
	"learning-core-api/internal/domain/subjects"

	"github.com/google/uuid"
	"google.golang.org/genai"
)

const defaultPollInterval = 5 * time.Second

// Service coordinates file search store management and uploads.
type Service interface {
	EnsureStore(ctx context.Context, subjectID uuid.UUID, displayName string) (*Store, error)
	UploadToSubject(ctx context.Context, req UploadRequest) (*UploadResult, error)
}

type service struct {
	apiKey       string
	storeRepo    Repository
	subjectRepo  subjects.Repository
	documentRepo documents.Repository
	pollInterval time.Duration
}

// NewService creates a new file search service.
func NewService(apiKey string, storeRepo Repository, subjectRepo subjects.Repository, documentRepo documents.Repository) Service {
	return &service{
		apiKey:       apiKey,
		storeRepo:    storeRepo,
		subjectRepo:  subjectRepo,
		documentRepo: documentRepo,
		pollInterval: defaultPollInterval,
	}
}

// EnsureStore ensures a file search store exists for the subject.
func (s *service) EnsureStore(ctx context.Context, subjectID uuid.UUID, displayName string) (*Store, error) {
	if subjectID == uuid.Nil {
		return nil, errors.New("subject ID is required")
	}

	store, err := s.storeRepo.GetBySubjectID(ctx, subjectID)
	if err == nil {
		return store, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	subject, err := s.subjectRepo.GetByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	resolvedDisplayName := displayName
	if resolvedDisplayName == "" {
		resolvedDisplayName = subject.Name
	}

	client, err := s.newClient(ctx)
	if err != nil {
		return nil, err
	}

	fileSearchStore, err := client.FileSearchStores.Create(ctx, &genai.CreateFileSearchStoreConfig{
		DisplayName: resolvedDisplayName,
	})
	if err != nil {
		return nil, err
	}

	chunkingConfig := json.RawMessage(`{}`)
	created, err := s.storeRepo.Create(ctx, Store{
		SubjectID:      subjectID,
		StoreName:      fileSearchStore.Name,
		DisplayName:    resolvedDisplayName,
		ChunkingConfig: chunkingConfig,
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

// UploadToSubject uploads a file to the subject's file search store and records a document.
func (s *service) UploadToSubject(ctx context.Context, req UploadRequest) (*UploadResult, error) {
	if req.SubjectID == uuid.Nil {
		return nil, errors.New("subject ID is required")
	}
	if req.FilePath == "" {
		return nil, errors.New("file path is required")
	}

	store, err := s.EnsureStore(ctx, req.SubjectID, req.DisplayName)
	if err != nil {
		return nil, err
	}

	client, err := s.newClient(ctx)
	if err != nil {
		return nil, err
	}

	fileName := req.FileName
	if fileName == "" {
		fileName = filepath.Base(req.FilePath)
	}

	sourceURI := req.SourceURI
	if sourceURI == "" {
		sourceURI = req.FilePath
	}

	title := req.Title
	if title == "" {
		title = fileName
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = title
	}

	uploadConfig := &genai.UploadToFileSearchStoreConfig{
		DisplayName: displayName,
		MIMEType:    req.MimeType,
	}

	operation, err := client.FileSearchStores.UploadToFileSearchStoreFromPath(ctx, req.FilePath, store.StoreName, uploadConfig)
	if err != nil {
		return nil, err
	}

	for !operation.Done {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		time.Sleep(s.pollInterval)
		operation, err = client.Operations.GetUploadToFileSearchStoreOperation(ctx, operation, nil)
		if err != nil {
			return nil, err
		}
	}

	if operation.Error != nil {
		return nil, fmt.Errorf("upload failed: %v", operation.Error)
	}

	var documentName string
	if operation.Response != nil {
		documentName = operation.Response.DocumentName
	}

	docName := req.DocName
	if docName == "" {
		docName = documentName
	}

	sha256Value, err := sha256File(req.FilePath)
	if err != nil {
		return nil, err
	}

	indexedAt := time.Now()
	_, err = s.documentRepo.Create(ctx, documents.CreateDocumentRequest{
		SubjectID: req.SubjectID,
		StoreID:   store.ID,
		Title:     title,
		SourceURI: sourceURI,
		SHA256:    sha256Value,
		Metadata:  req.Metadata,
		FileName:  fileName,
		DocName:   docName,
		IndexedAt: &indexedAt,
	})
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Store:        store,
		DocumentName: documentName,
	}, nil
}

func (s *service) newClient(ctx context.Context) (*genai.Client, error) {
	if s.apiKey == "" {
		return nil, errors.New("google API key is required")
	}

	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  s.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
}

func sha256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
