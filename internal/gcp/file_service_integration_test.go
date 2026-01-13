package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	fileService := NewFileService(gcsService, client, storeName)
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

func TestGCSSignedUploadURLIntegration(t *testing.T) {
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	apiKey := os.Getenv("GOOGLE_API_KEY")
	storeName := os.Getenv("FILE_STORE_NAME")
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if bucketName == "" || apiKey == "" || storeName == "" {
		_ = godotenv.Load()
		bucketName = os.Getenv("GCS_BUCKET_NAME")
		apiKey = os.Getenv("GOOGLE_API_KEY")
		storeName = os.Getenv("FILE_STORE_NAME")
		credentialsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	}
	if bucketName == "" || apiKey == "" || storeName == "" {
		t.Skip("missing GCS_BUCKET_NAME, GOOGLE_API_KEY, or FILE_STORE_NAME")
	}
	t.Logf("using bucket: %s", bucketName)

	if credentialsPath != "" {
		resolvedPath, err := resolveCredentialsPath(credentialsPath)
		require.NoError(t, err)
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", resolvedPath)
		t.Logf("using GOOGLE_APPLICATION_CREDENTIALS=%s", resolvedPath)
	}

	ctx := context.Background()
	gcsService, err := NewGCSService(ctx, bucketName)
	require.NoError(t, err)
	defer gcsService.Close()

	_, err = gcsService.BucketAttrs(ctx)
	require.NoError(t, err)
	t.Log("connected to GCS bucket")

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	require.NoError(t, err)
	fileService := NewFileService(gcsService, client, storeName)
	require.NoError(t, fileService.EnsureStore(ctx))
	storeName = fileService.StoreName()

	// List files before
	filesBefore, err := fileService.ListFiles(ctx, storeName)
	require.NoError(t, err)
	for _, doc := range filesBefore {
		t.Logf("file store file before: %s", doc.Name)
	}

	objectName := fmt.Sprintf("signed/%s-test.pdf", uuid.NewString())
	localPath := filepath.Join("../../test_docs", "test.pdf")
	t.Logf("uploading %s to %s", localPath, objectName)

	fileHandle, err := os.Open(localPath)
	require.NoError(t, err)
	defer fileHandle.Close()

	t.Log("generating signed upload URL...")
	signedURL, err := gcsService.GenerateSignedUploadURL(ctx, objectName, "application/pdf", 10*time.Minute)
	require.NoError(t, err)
	t.Logf("signed URL generated (length: %d)", len(signedURL))

	t.Log("executing HTTP PUT to signed URL...")
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, signedURL, fileHandle)
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/pdf")

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	t.Logf("HTTP response status: %s", response.Status)
	require.Equal(t, http.StatusOK, response.StatusCode)
	_, _ = io.Copy(io.Discard, response.Body)

	t.Log("verifying upload by downloading file...")
	data, err := gcsService.DownloadFile(ctx, objectName)
	require.NoError(t, err)
	t.Logf("successfully downloaded %d bytes", len(data))
	require.NotEmpty(t, data)

	// Now upload from GCS to File Search Store
	t.Log("uploading from GCS to Gemini File Search Store...")
	result, err := fileService.UploadToFileSearchStore(ctx, objectName, "signed-url-test", "application/pdf")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Operation)

	operation := result.Operation
	deadline := time.Now().Add(2 * time.Minute)
	for !operation.Done && time.Now().Before(deadline) {
		t.Log("waiting for File Search Store operation...")
		time.Sleep(3 * time.Second)
		operation, err = client.Operations.GetUploadToFileSearchStoreOperation(ctx, operation, nil)
		require.NoError(t, err)
	}

	require.True(t, operation.Done, "upload operation did not complete in time")
	require.NotNil(t, operation.Response)
	t.Logf("File Search Store document created: %s", operation.Response.DocumentName)

	// List files after
	filesAfter, err := fileService.ListFiles(ctx, storeName)
	require.NoError(t, err)
	for _, doc := range filesAfter {
		t.Logf("file store file after: %s", doc.Name)
	}
}

func TestFileSearchMetadataFilterIntegration(t *testing.T) {
	storeName := os.Getenv("FILE_STORE_NAME")
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if storeName == "" || apiKey == "" {
		_ = godotenv.Load()
		storeName = os.Getenv("FILE_STORE_NAME")
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if storeName == "" || apiKey == "" {
		t.Skipf("missing FILE_STORE_NAME or GOOGLE_API_KEY (store=%t api_key=%t)", storeName != "", apiKey != "")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	require.NoError(t, err)

	storeName, err = ensureFileSearchStore(ctx, client, storeName)
	require.NoError(t, err)
	t.Logf("using file search store: %s", storeName)

	localPath := filepath.Join("../../test_docs", "test.pdf")
	metadataID := uuid.NewString()

	uploadCfg := &genai.UploadToFileSearchStoreConfig{
		MIMEType:    "application/pdf",
		DisplayName: fmt.Sprintf("rag-test-%s", metadataID),
		CustomMetadata: []*genai.CustomMetadata{
			{
				Key:         "id",
				StringValue: metadataID,
			},
		},
	}

	operation, err := client.FileSearchStores.UploadToFileSearchStoreFromPath(ctx, localPath, storeName, uploadCfg)
	require.NoError(t, err)

	deadline := time.Now().Add(2 * time.Minute)
	for !operation.Done && time.Now().Before(deadline) {
		time.Sleep(3 * time.Second)
		operation, err = client.Operations.GetUploadToFileSearchStoreOperation(ctx, operation, nil)
		require.NoError(t, err)
	}

	require.True(t, operation.Done, "upload operation did not complete in time")
	require.NotNil(t, operation.Response)
	require.NotEmpty(t, operation.Response.DocumentName)
	t.Logf("uploaded document: %s", operation.Response.DocumentName)

	filter := fmt.Sprintf("id = \"%s\"", metadataID)
	tools := []*genai.Tool{
		{
			FileSearch: &genai.FileSearch{
				FileSearchStoreNames: []string{storeName},
				MetadataFilter:       filter,
			},
		},
	}
	genConfig := &genai.GenerateContentConfig{Tools: tools}

	resp, err := client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text("Summarize the document in a few paragraphs highlight key topics and components"), genConfig)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Candidates)

	var outputText string
	for _, part := range resp.Candidates[0].Content.Parts {
		outputText += part.Text
	}
	t.Logf("model response: %s", outputText)

	outputPath, err := writeFileSearchOutput(metadataID, outputText, resp.Candidates[0].GroundingMetadata)
	require.NoError(t, err)
	t.Logf("saved output to %s", outputPath)

	if resp.Candidates[0].GroundingMetadata != nil {
		data, err := json.MarshalIndent(resp.Candidates[0].GroundingMetadata, "", "  ")
		require.NoError(t, err)
		t.Logf("grounding metadata: %s", string(data))
	} else {
		t.Log("grounding metadata: <nil>")
	}
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

func ensureFileSearchStore(ctx context.Context, client *genai.Client, storeName string) (string, error) {
	if client == nil || client.FileSearchStores == nil {
		return "", fmt.Errorf("file search store client is required")
	}
	if storeName == "" {
		return "", fmt.Errorf("file search store name is required")
	}

	store, err := client.FileSearchStores.Get(ctx, storeName, nil)
	if err == nil && store != nil && store.Name != "" {
		return store.Name, nil
	}
	if err != nil && !isNotFoundError(err) {
		if matched, matchErr := matchStoreByDisplayName(ctx, client, storeName); matchErr == nil && matched != nil {
			return matched.Name, nil
		}
		return "", err
	}

	if matched, matchErr := matchStoreByDisplayName(ctx, client, storeName); matchErr == nil && matched != nil {
		return matched.Name, nil
	}

	created, err := client.FileSearchStores.Create(ctx, &genai.CreateFileSearchStoreConfig{
		DisplayName: storeName,
	})
	if err != nil {
		return "", err
	}
	if created != nil && created.Name != "" {
		return created.Name, nil
	}

	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		if matched, matchErr := matchStoreByDisplayName(ctx, client, storeName); matchErr == nil && matched != nil {
			return matched.Name, nil
		}
		time.Sleep(3 * time.Second)
	}

	return "", fmt.Errorf("timed out waiting for file search store %q", storeName)
}

func matchStoreByDisplayName(ctx context.Context, client *genai.Client, storeName string) (*genai.FileSearchStore, error) {
	for store, err := range client.FileSearchStores.All(ctx) {
		if err != nil {
			return nil, err
		}
		if store == nil {
			continue
		}
		if store.Name == storeName || store.DisplayName == storeName {
			return store, nil
		}
	}
	return nil, nil
}

func writeFileSearchOutput(metadataID string, outputText string, grounding *genai.GroundingMetadata) (string, error) {
	dir := filepath.Join("tmp", "filesearch_outputs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	payload := struct {
		MetadataID string                   `json:"metadata_id"`
		OutputText string                   `json:"output_text"`
		Grounding  *genai.GroundingMetadata `json:"grounding_metadata,omitempty"`
		WrittenAt  time.Time                `json:"written_at"`
	}{
		MetadataID: metadataID,
		OutputText: outputText,
		Grounding:  grounding,
		WrittenAt:  time.Now().UTC(),
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, fmt.Sprintf("filesearch_%s.json", metadataID))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
