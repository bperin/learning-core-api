package auth

import (
	"context"
	"errors"
	"testing"

	"learning-core-api/internal/persistance/store"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("success", func(t *testing.T) {
		repo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (*store.User, error) {
				return &store.User{
					ID:       userID,
					Email:    email,
					Password: "secret",
					IsAdmin:  true,
				}, nil
			},
		}
		service := NewService("secret", repo)

		tokens, user, err := service.LoginWithEmail(context.Background(), "admin@example.com", "secret")
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
					Password: "secret",
				}, nil
			},
		}
		service := NewService("secret", repo)

		_, _, err := service.LoginWithEmail(context.Background(), "admin@example.com", "wrong")
		assert.Error(t, err)
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
