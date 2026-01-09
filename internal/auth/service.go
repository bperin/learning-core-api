package auth

import (
	"context"
	"errors"
	authdto "learning-core-api/internal/auth/dto"
	"learning-core-api/internal/persistance/store"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secret []byte
	repo   Repository
}

func NewService(secret string, repo Repository) *Service {
	return &Service{
		secret: []byte(secret),
		repo:   repo,
	}
}

func (s *Service) LoginWithEmail(ctx context.Context, email, password string) (*authdto.TokenResponse, *store.User, error) {
	// Check if user exists
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Login attempt for non-existent user: %s", email)
		return nil, nil, errors.New("invalid credentials")
	}

	// Verify password
	// NOTE: In a real app, use bcrypt.CompareHashAndPassword
	if user.Password != password {
		log.Printf("Invalid password for user: %s", email)
		return nil, nil, errors.New("invalid credentials")
	}

	tokens, err := s.GenerateTokenPair(ctx, user.ID, "user") // Defaulting to "user" for now, or use user.Role if available
	if err != nil {
		return nil, nil, err
	}

	return &authdto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Role:         tokens.Role,
	}, user, nil
}

func (s *Service) GenerateTokenPair(ctx context.Context, userID uuid.UUID, role string) (*authdto.TokenPair, error) {
	// Access Token
	accessClaims := &authdto.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(120 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	// Refresh Token (longer lived)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &authdto.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	// Save refresh token in database (optional, can be implemented in repo)
	err = s.repo.CreateRefreshToken(ctx, refreshStr, userID, refreshExpiresAt)
	if err != nil {
		log.Printf("warning: failed to save refresh token: %v", err)
		// We can still return the tokens even if saving fails, or return error
	}

	return &authdto.TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		Role:         role,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*authdto.TokenResponse, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// rotation logic here...
	s.repo.DeleteRefreshToken(ctx, refreshToken)

	tokens, err := s.GenerateTokenPair(ctx, claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	return &authdto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Role:         tokens.Role,
	}, nil
}

func (s *Service) ValidateToken(tokenStr string) (*authdto.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &authdto.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authdto.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
