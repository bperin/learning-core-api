package auth

import (
	"context"
	"errors"
	"slap-realtime/internal/store"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockAuthRepository struct {
	createOTPFunc               func(ctx context.Context, phone, code, otpType string, expiresAt time.Time) error
	getLastOTPFunc              func(ctx context.Context, phone, otpType string) (*OTP, error)
	deleteOTPsByPhoneFunc       func(ctx context.Context, phone string) error
	getUserByPhoneFunc          func(ctx context.Context, phone string) (uuid.UUID, error)
	getUserByPhoneFullFunc      func(ctx context.Context, phone string) (*store.User, error)
	updateUserVerificationFunc  func(ctx context.Context, userID uuid.UUID, verified bool) error
	createRefreshTokenFunc      func(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error
	getRefreshTokenFunc         func(ctx context.Context, token string) (*store.RefreshToken, error)
	deleteRefreshTokenFunc      func(ctx context.Context, token string) error
	deleteUserRefreshTokensFunc func(ctx context.Context, userID uuid.UUID) error
	getUserByIDFunc             func(ctx context.Context, id uuid.UUID) (*store.User, error)
}

func (m *mockAuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*store.User, error) {
	return m.getUserByIDFunc(ctx, id)
}

func (m *mockAuthRepository) CreateOTP(ctx context.Context, phone, code, otpType string, expiresAt time.Time) error {
	return m.createOTPFunc(ctx, phone, code, otpType, expiresAt)
}

func (m *mockAuthRepository) GetLastOTP(ctx context.Context, phone, otpType string) (*OTP, error) {
	return m.getLastOTPFunc(ctx, phone, otpType)
}

func (m *mockAuthRepository) DeleteOTPsByPhone(ctx context.Context, phone string) error {
	return m.deleteOTPsByPhoneFunc(ctx, phone)
}

func (m *mockAuthRepository) GetUserByPhone(ctx context.Context, phone string) (uuid.UUID, error) {
	return m.getUserByPhoneFunc(ctx, phone)
}

func (m *mockAuthRepository) GetUserByPhoneFull(ctx context.Context, phone string) (*store.User, error) {
	return m.getUserByPhoneFullFunc(ctx, phone)
}

func (m *mockAuthRepository) UpdateUserVerification(ctx context.Context, userID uuid.UUID, verified bool) error {
	return m.updateUserVerificationFunc(ctx, userID, verified)
}

func (m *mockAuthRepository) CreateRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	return m.createRefreshTokenFunc(ctx, token, userID, expiresAt)
}

func (m *mockAuthRepository) GetRefreshToken(ctx context.Context, token string) (*store.RefreshToken, error) {
	return m.getRefreshTokenFunc(ctx, token)
}

func (m *mockAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return m.deleteRefreshTokenFunc(ctx, token)
}

func (m *mockAuthRepository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return m.deleteUserRefreshTokensFunc(ctx, userID)
}

func TestService_SendOtpCode(t *testing.T) {
	phoneNumber := "+1234567890"

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockAuthRepository{
			getUserByPhoneFunc: func(ctx context.Context, phone string) (uuid.UUID, error) {
				return uuid.New(), nil
			},
			createOTPFunc: func(ctx context.Context, phone, code, otpType string, expiresAt time.Time) error {
				return nil
			},
		}
		service := NewService("secret", mockRepo)
		_, err := service.SendOtpCode(context.Background(), phoneNumber)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := &mockAuthRepository{
			getUserByPhoneFunc: func(ctx context.Context, phone string) (uuid.UUID, error) {
				return uuid.Nil, errors.New("not found")
			},
		}
		service := NewService("secret", mockRepo)
		_, err := service.SendOtpCode(context.Background(), phoneNumber)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestService_Verify(t *testing.T) {
	phoneNumber := "+1234567890"
	userID := uuid.New()
	userRole := "user"

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockAuthRepository{
			getUserByPhoneFunc: func(ctx context.Context, phone string) (uuid.UUID, error) {
				return userID, nil
			},
			getUserByPhoneFullFunc: func(ctx context.Context, phone string) (*store.User, error) {
				return &store.User{ID: userID, PhoneNumber: phone, Role: userRole}, nil
			},
			updateUserVerificationFunc: func(ctx context.Context, userID uuid.UUID, verified bool) error {
				return nil
			},
			createRefreshTokenFunc: func(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
				return nil
			},
		}
		service := NewService("secret", mockRepo)

		// 1. Send OTP first
		code, err := service.SendOtpCode(context.Background(), phoneNumber)
		assert.NoError(t, err)

		// 2. Verify
		tokens, user, err := service.Verify(context.Background(), phoneNumber, code)
		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
	})
}

func TestService_RefreshToken(t *testing.T) {
	userID := uuid.New()
	userRole := "user"

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockAuthRepository{
			getRefreshTokenFunc: func(ctx context.Context, token string) (*store.RefreshToken, error) {
				return &store.RefreshToken{Token: token, UserID: userID, ExpiresAt: time.Now().Add(time.Hour)}, nil
			},
			deleteRefreshTokenFunc: func(ctx context.Context, token string) error {
				return nil
			},
			createRefreshTokenFunc: func(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
				return nil
			},
			getUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*store.User, error) {
				return &store.User{ID: userID, Role: userRole}, nil
			},
		}
		service := NewService("secret", mockRepo)

		// 1. Generate a valid refresh token first
		tokens, err := service.GenerateTokenPair(context.Background(), userID, userRole)
		assert.NoError(t, err)

		newTokens, err := service.RefreshToken(context.Background(), tokens.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, newTokens)
		assert.NotEqual(t, tokens.RefreshToken, newTokens.RefreshToken)
	})

	t.Run("token validation error", func(t *testing.T) {
		service := NewService("secret", &mockAuthRepository{})
		_, err := service.RefreshToken(context.Background(), "invalid-token")
		assert.Error(t, err)
	})
}
