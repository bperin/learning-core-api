package model_configs_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/model_configs"
	"learning-core-api/internal/domain/users"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
)

func setupHandler(t *testing.T) (*model_configs.Handler, uuid.UUID, func()) {
	t.Helper()

	tx, cleanup := testutil.NewTestTx(t)
	queries := store.New(tx)

	// Create admin user for testing
	adminUser := users.User{
		ID:    uuid.New(),
		Email: "admin@test.com",
	}

	userRepo := users.NewRepository(queries)
	createdUser, err := userRepo.CreateUser(context.Background(), adminUser, "hashedpassword", users.UserRoleAdmin)
	require.NoError(t, err)

	repo := model_configs.NewRepository(queries)
	service := model_configs.NewService(repo)
	handler := model_configs.NewHandler(service)

	return handler, createdUser.ID, cleanup
}

func TestHandler_Create(t *testing.T) {
	handler, adminUserID, cleanup := setupHandler(t)
	defer cleanup()

	payload := model_configs.CreateModelConfigRequest{
		ModelName:   "gemini-1.5-pro",
		Temperature: 0.7,
		MaxTokens:   2048,
		TopP:        0.9,
		TopK:        40.0,
		MimeType:    "application/json",
		IsActive:    true,
		CreatedBy:   adminUserID,
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/model-configs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result model_configs.ModelConfig
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "gemini-1.5-pro", result.ModelName)
	assert.Equal(t, 0.7, result.Temperature)
	assert.True(t, result.IsActive)
}

func TestHandler_ListAll(t *testing.T) {
	handler, _, cleanup := setupHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/model-configs", nil)
	w := httptest.NewRecorder()

	handler.ListAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []*model_configs.ModelConfig
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	// Should have at least the seeded model config
	assert.GreaterOrEqual(t, len(result), 1)
}

func TestHandler_GetActive_Found(t *testing.T) {
	handler, _, cleanup := setupHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/model-configs/active", nil)
	w := httptest.NewRecorder()

	handler.GetActive(w, req)

	// Should find the seeded active config
	assert.Equal(t, http.StatusOK, w.Code)

	var result model_configs.ModelConfig
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.True(t, result.IsActive)
}

func TestHandler_Activate(t *testing.T) {
	handler, adminUserID, cleanup := setupHandler(t)
	defer cleanup()

	// First create a config
	payload := model_configs.CreateModelConfigRequest{
		ModelName:   "gemini-1.5-pro",
		Temperature: 0.3,
		MaxTokens:   4096,
		TopP:        0.95,
		TopK:        20.0,
		MimeType:    "application/json",
		IsActive:    false,
		CreatedBy:   adminUserID,
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	createReq := httptest.NewRequest(http.MethodPost, "/model-configs", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()

	handler.Create(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	var created model_configs.ModelConfig
	err = json.Unmarshal(createW.Body.Bytes(), &created)
	require.NoError(t, err)

	// Now activate it using chi context
	activateReq := httptest.NewRequest(http.MethodPost, "/model-configs/"+created.ID.String()+"/activate", nil)

	// Set up chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", created.ID.String())
	activateReq = activateReq.WithContext(context.WithValue(activateReq.Context(), chi.RouteCtxKey, rctx))

	activateW := httptest.NewRecorder()
	handler.Activate(activateW, activateReq)

	assert.Equal(t, http.StatusOK, activateW.Code)
}
