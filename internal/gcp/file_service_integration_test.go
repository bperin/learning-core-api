package gcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

func TestFileServiceUploadIntegration(t *testing.T) {
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	storeName := os.Getenv("FILE_STORE_NAME")
	apiKey := os.Getenv("GOOGLE_API_KEY")
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if bucketName == "" || storeName == "" || apiKey == "" {
		_ = godotenv.Load()
		bucketName = os.Getenv("GCS_BUCKET_NAME")
		storeName = os.Getenv("FILE_STORE_NAME")
		apiKey = os.Getenv("GOOGLE_API_KEY")
		credentialsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	}
	if bucketName == "" || storeName == "" || apiKey == "" {
		t.Skipf("missing GCS_BUCKET_NAME, FILE_STORE_NAME, or GOOGLE_API_KEY (bucket=%t store=%t api_key=%t)", bucketName != "", storeName != "", apiKey != "")
	}

	if credentialsPath != "" {
		resolvedPath, err := resolveCredentialsPath(credentialsPath)
		require.NoError(t, err)
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", resolvedPath)
		t.Logf("using GOOGLE_APPLICATION_CREDENTIALS=%s", resolvedPath)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	require.NoError(t, err)

	gcsService, err := NewGCSService(ctx, bucketName)
	require.NoError(t, err)
	defer gcsService.Close()

	_, err = gcsService.BucketAttrs(ctx)
	require.NoError(t, err)

	fileService := NewFileService(gcsService, client, storeName, 10*time.Minute)
	require.NoError(t, fileService.EnsureStore(ctx))
	storeName = fileService.StoreName()

	storesBefore, err := fileService.ListStores(ctx)
	require.NoError(t, err)
	for _, store := range storesBefore {
		t.Logf("store before: %s", store.Name)
	}

	filesBefore, err := fileService.ListFiles(ctx, storeName)
	require.NoError(t, err)
	for _, doc := range filesBefore {
		t.Logf("file before: %s", doc.Name)
	}

	objectName := fmt.Sprintf("test/%s-test.pdf", uuid.NewString())
	localPath := filepath.Join("../../test_docs", "test.pdf")
	fileHandle, err := os.Open(localPath)
	require.NoError(t, err)
	defer fileHandle.Close()

	_, err = gcsService.UploadFile(ctx, objectName, "application/pdf", fileHandle)
	require.NoError(t, err)
	defer func() {
		_ = gcsService.DeleteObject(ctx, objectName)
	}()

	result, err := fileService.UploadToFileSearchStore(ctx, objectName, "synthetic-test", "application/pdf")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Operation)

	operation := result.Operation
	deadline := time.Now().Add(2 * time.Minute)
	for !operation.Done && time.Now().Before(deadline) {
		time.Sleep(3 * time.Second)
		operation, err = client.Operations.GetUploadToFileSearchStoreOperation(ctx, operation, nil)
		require.NoError(t, err)
	}

	require.True(t, operation.Done, "upload operation did not complete in time")
	require.NotNil(t, operation.Response)
	require.NotEmpty(t, operation.Response.DocumentName)

	filesAfter, err := fileService.ListFiles(ctx, storeName)
	require.NoError(t, err)
	for _, doc := range filesAfter {
		t.Logf("file after: %s", doc.Name)
	}

	// No document delete method on genai client; rely on store cleanup policies.
}

func resolveCredentialsPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is empty")
	}
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err != nil {
			return "", fmt.Errorf("credentials file not found: %w", err)
		}
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := wd
	for {
		candidate := filepath.Clean(filepath.Join(current, path))
		if _, err := os.Stat(candidate); err == nil {
			absPath, err := filepath.Abs(candidate)
			if err != nil {
				return "", err
			}
			return absPath, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("credentials file not found from %s", wd)
}
