package gcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"learning-core-api/internal/config"
)

type GCSService struct {
	client     *storage.Client
	bucketName string
}

func NewGCSService(ctx context.Context, bucketName string) (*GCSService, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcs client: %w", err)
	}

	return &GCSService{client: client, bucketName: bucketName}, nil
}

func NewGCSServiceFromConfig(ctx context.Context, cfg *config.Config) (*GCSService, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	return NewGCSService(ctx, cfg.GCSBucketName)
}

func (s *GCSService) BucketName() string {
	return s.bucketName
}

func (s *GCSService) ObjectURI(objectName string) string {
	return fmt.Sprintf("gs://%s/%s", s.bucketName, objectName)
}

func (s *GCSService) UploadFile(ctx context.Context, objectName string, contentType string, r io.Reader) (string, error) {
	if objectName == "" {
		return "", fmt.Errorf("object name is required")
	}

	w := s.client.Bucket(s.bucketName).Object(objectName).NewWriter(ctx)
	if contentType != "" {
		w.ContentType = contentType
	}

	if _, err := io.Copy(w, r); err != nil {
		_ = w.Close()
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize upload: %w", err)
	}

	return s.ObjectURI(objectName), nil
}

func (s *GCSService) DownloadFile(ctx context.Context, objectName string) ([]byte, error) {
	reader, err := s.OpenReader(ctx, objectName)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *GCSService) OpenReader(ctx context.Context, objectName string) (*storage.Reader, error) {
	if objectName == "" {
		return nil, fmt.Errorf("object name is required")
	}
	reader, err := s.client.Bucket(s.bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open object reader: %w", err)
	}
	return reader, nil
}

func (s *GCSService) BucketAttrs(ctx context.Context) (*storage.BucketAttrs, error) {
	attrs, err := s.client.Bucket(s.bucketName).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bucket attrs: %w", err)
	}
	return attrs, nil
}

func (s *GCSService) DeleteObject(ctx context.Context, objectName string) error {
	if objectName == "" {
		return fmt.Errorf("object name is required")
	}
	if err := s.client.Bucket(s.bucketName).Object(objectName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func (s *GCSService) ListObjects(ctx context.Context, prefix string) ([]*storage.ObjectAttrs, error) {
	query := &storage.Query{Prefix: prefix}
	it := s.client.Bucket(s.bucketName).Objects(ctx, query)
	objects := []*storage.ObjectAttrs{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}
		objects = append(objects, attrs)
	}
	return objects, nil
}

func (s *GCSService) GenerateSignedUploadURL(ctx context.Context, objectName string, contentType string, ttl time.Duration) (string, error) {
	if objectName == "" {
		return "", fmt.Errorf("object name is required")
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	opts := &storage.SignedURLOptions{
		Scheme:      storage.SigningSchemeV4,
		Method:      "PUT",
		Expires:     time.Now().Add(ttl),
		ContentType: contentType,
	}

	url, err := s.client.Bucket(s.bucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed url: %w", err)
	}
	return url, nil
}

func (s *GCSService) Close() error {
	return s.client.Close()
}
