package auth

import (
	"context"
	"errors"
	"testing"

	"learning-core-api/internal/persistance/store"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepository struct {
	getUserByEmailFunc func(ctx context.Context, email string) (*store.User, error)
	getUserByIDFunc    func(ctx context.Context, id uuid.UUID) (*store.User, error)
}

func (m *mockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*store.User, error) {
	return m.getUserByEmailFunc(ctx, email)
}

func (m *mockAuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*store.User, error) {
	return m.getUserByIDFunc(ctx, id)
}

func (m *mockAuthRepository) CreateRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt interface{}) error {
	return nil
}

func (m *mockAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return nil
}

func TestService_LoginWithEmail(t *testing.T) {
	userID := uuid.New()
	plainPassword := "testpassword123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		repo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (*store.User, error) {
				return &store.User{
					ID:       userID,
					Email:    email,
					Password: string(hashedPassword),
					IsAdmin:  true,
				}, nil
			},
		}
		service := NewService("secret", repo)

		tokens, user, err := service.LoginWithEmail(context.Background(), "admin@example.com", plainPassword)
		require.NoError(t, err)
		require.NotNil(t, tokens)
		require.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
	})

	t.Run("invalid password", func(t *testing.T) {
		repo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (*store.User, error) {
				return &store.User{
					ID:       userID,
					Email:    email,
					Password: string(hashedPassword),
				}, nil
			},
		}
		service := NewService("secret", repo)

		_, _, err := service.LoginWithEmail(context.Background(), "admin@example.com", "wrongpassword")
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		repo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (*store.User, error) {
				return nil, errors.New("not found")
			},
		}
		service := NewService("secret", repo)

		_, _, err := service.LoginWithEmail(context.Background(), "missing@example.com", "secret")
		assert.Error(t, err)
	})
}

func TestService_RefreshToken(t *testing.T) {
	service := NewService("secret", &mockAuthRepository{})
	roles := []string{"admin"}

	tokens, err := service.GenerateTokenPair(context.Background(), uuid.New(), roles)
	require.NoError(t, err)

	refreshed, err := service.RefreshToken(context.Background(), tokens.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshed.AccessToken)
	assert.NotEmpty(t, refreshed.RefreshToken)
	assert.Equal(t, roles[0], refreshed.Role)
}

func TestService_RefreshToken_Invalid(t *testing.T) {
	service := NewService("secret", &mockAuthRepository{})
	_, err := service.RefreshToken(context.Background(), "invalid-token")
	assert.Error(t, err)
}

func TestService_GenerateTokenPair_ScopesInAccessToken(t *testing.T) {
	service := NewService("secret", &mockAuthRepository{})
	userID := uuid.New()
	roles := []string{"admin"}

	tokens, err := service.GenerateTokenPair(context.Background(), userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, tokens.AccessToken)

	token, err := jwt.Parse(tokens.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	scopes, ok := claims["scopes"].([]interface{})
	require.True(t, ok, "scopes claim should exist and be an array")
	require.Len(t, scopes, 2, "admin role should have read and write scopes")

	scopeStrs := make([]string, len(scopes))
	for i, s := range scopes {
		scopeStrs[i] = s.(string)
	}
	assert.Contains(t, scopeStrs, "read")
	assert.Contains(t, scopeStrs, "write")
}

func TestService_GenerateTokenPair_ScopesInRefreshToken(t *testing.T) {
	service := NewService("secret", &mockAuthRepository{})
	userID := uuid.New()
	roles := []string{"learner"}

	tokens, err := service.GenerateTokenPair(context.Background(), userID, roles)
	require.NoError(t, err)
	require.NotEmpty(t, tokens.RefreshToken)

	token, err := jwt.Parse(tokens.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	scopes, ok := claims["scopes"].([]interface{})
	require.True(t, ok, "scopes claim should exist and be an array")
	require.Len(t, scopes, 2, "learner role should have read and write scopes")

	scopeStrs := make([]string, len(scopes))
	for i, s := range scopes {
		scopeStrs[i] = s.(string)
	}
	assert.Contains(t, scopeStrs, "read")
	assert.Contains(t, scopeStrs, "write")
}

func TestService_GenerateTokenPair_RolesInToken(t *testing.T) {
	service := NewService("secret", &mockAuthRepository{})
	userID := uuid.New()
	roles := []string{"admin", "teacher"}

	tokens, err := service.GenerateTokenPair(context.Background(), userID, roles)
	require.NoError(t, err)

	token, err := jwt.Parse(tokens.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	rolesInterface, ok := claims["roles"].([]interface{})
	require.True(t, ok, "roles claim should exist and be an array")
	require.Len(t, rolesInterface, 2)

	roleStrs := make([]string, len(rolesInterface))
	for i, r := range rolesInterface {
		roleStrs[i] = r.(string)
	}
	assert.Contains(t, roleStrs, "admin")
	assert.Contains(t, roleStrs, "teacher")
}

func TestService_RefreshToken_PreservesScopesAndRoles(t *testing.T) {
	service := NewService("secret", &mockAuthRepository{})
	userID := uuid.New()
	originalRoles := []string{"admin"}

	tokens, err := service.GenerateTokenPair(context.Background(), userID, originalRoles)
	require.NoError(t, err)

	refreshed, err := service.RefreshToken(context.Background(), tokens.RefreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, refreshed.AccessToken)

	token, err := jwt.Parse(refreshed.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	scopes, ok := claims["scopes"].([]interface{})
	require.True(t, ok, "refreshed token should have scopes")
	require.Len(t, scopes, 2)

	roles, ok := claims["roles"].([]interface{})
	require.True(t, ok, "refreshed token should have roles")
	require.Len(t, roles, 1)
	assert.Equal(t, "admin", roles[0].(string))
}

func TestService_ScopesFromRoles_AllRoles(t *testing.T) {
	tests := []struct {
		name     string
		roles    []string
		expected []string
	}{
		{
			name:     "admin role",
			roles:    []string{"admin"},
			expected: []string{"read", "write"},
		},
		{
			name:     "teacher role",
			roles:    []string{"teacher"},
			expected: []string{"read", "write"},
		},
		{
			name:     "learner role",
			roles:    []string{"learner"},
			expected: []string{"read", "write"},
		},
		{
			name:     "multiple roles",
			roles:    []string{"admin", "teacher"},
			expected: []string{"read", "write"},
		},
		{
			name:     "empty roles",
			roles:    []string{},
			expected: []string{"read"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopes := scopesFromRoles(tt.roles)
			assert.ElementsMatch(t, tt.expected, scopes)
		})
	}
}
