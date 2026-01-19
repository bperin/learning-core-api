package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/users"
	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
)

func setupHandler(t *testing.T) (*users.Handler, func()) {
	t.Helper()

	tx, cleanup := testutil.NewTestTx(t)
	queries := store.New(tx)
	repo := users.NewRepository(queries)
	service := users.NewService(repo)
	handler := users.NewHandler(service)

	return handler, cleanup
}

func TestHandler_SignupLearner(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	payload := map[string]interface{}{
		"email":    "learner@example.com",
		"password": "securepassword123",
		"role":     "LEARNER",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result users.User
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "learner@example.com", result.Email)
}

func TestHandler_SignupInstructor(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	payload := map[string]interface{}{
		"email":    "instructor@example.com",
		"password": "securepassword456",
		"role":     "INSTRUCTOR",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result users.User
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "instructor@example.com", result.Email)
}

func TestHandler_SignupAdmin(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	payload := map[string]interface{}{
		"email":    "admin@example.com",
		"password": "securepassword789",
		"role":     "ADMIN",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result users.User
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "admin@example.com", result.Email)
}

func TestHandler_SignupMissingPassword(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	payload := map[string]interface{}{
		"email": "nopass@example.com",
		"role":  "LEARNER",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_SignupInvalidRole(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	payload := map[string]interface{}{
		"email":    "invalidrole@example.com",
		"password": "securepassword",
		"role":     "SUPERUSER",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_SignupDuplicateEmail(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	payload := map[string]interface{}{
		"email":    "duplicate@example.com",
		"password": "securepassword",
		"role":     "LEARNER",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	// First signup
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Signup(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Second signup with same email should fail
	body, err = json.Marshal(payload)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Signup(w, req)
	// Should get 500 error for duplicate email
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_GetCurrentUser(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	// First, create a user via signup
	payload := map[string]interface{}{
		"email":    "currentuser@example.com",
		"password": "securepassword",
		"role":     "ADMIN",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Signup(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Extract user ID from signup response
	var signupResult users.User
	err = json.Unmarshal(w.Body.Bytes(), &signupResult)
	require.NoError(t, err)

	// Now test GetCurrentUser with the user ID in context
	req = httptest.NewRequest(http.MethodGet, "/users/me", nil)
	ctx := authz.WithAuth(req.Context(), signupResult.ID.String(), []string{"admin"}, []string{"read", "write"})
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	handler.GetCurrentUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result users.User
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, signupResult.ID, result.ID)
	assert.Equal(t, "currentuser@example.com", result.Email)
}

func TestHandler_GetCurrentUser_NoUserID(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	// Request without user ID in context
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	w := httptest.NewRecorder()

	handler.GetCurrentUser(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
