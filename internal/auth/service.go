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
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Login attempt for non-existent user: %s", email)
		return nil, nil, errors.New("invalid credentials")
	}

	if user.Password != password {
		log.Printf("Invalid password for user: %s", email)
		return nil, nil, errors.New("invalid credentials")
	}

	var roles []string
	if user.IsAdmin {
		roles = append(roles, "admin")
	}
	// Default to learner for non-admin users (schema only has is_admin today).
	if !user.IsAdmin {
		roles = append(roles, "learner")
	}

	tokens, err := s.GenerateTokenPair(ctx, user.ID, roles)
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

func (s *Service) GenerateTokenPair(ctx context.Context, userID uuid.UUID, roles []string) (*authdto.TokenPair, error) {
	rolesJoined := ""
	if len(roles) > 0 {
		rolesJoined = roles[0]
	}
	scopes := scopesFromRoles(roles)

	// Access Token
	accessClaims := jwt.MapClaims{
		"sub":    userID.String(),
		"roles":  roles,
		"scopes": scopes,
		"exp":    time.Now().Add(120 * time.Minute).Unix(),
		"iat":    time.Now().Unix(),
		"jti":    uuid.New().String(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"sub":    userID.String(),
		"roles":  roles,
		"scopes": scopes,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
		"jti":    uuid.New().String(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &authdto.TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		Role:         rolesJoined,
	}, nil
}

func scopesFromRoles(roles []string) []string {
	hasWrite := false
	for _, role := range roles {
		switch role {
		case "admin", "teacher", "learner":
			hasWrite = true
		}
	}
	if hasWrite {
		return []string{"read", "write"}
	}
	return []string{"read"}
}

func (s *Service) RefreshToken(ctx context.Context, refreshTokenStr string) (*authdto.TokenResponse, error) {
	// Parse and validate the refresh token
	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, _ := claims["sub"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	var roles []string
	if rClaim, ok := claims["roles"]; ok {
		if v, ok := rClaim.([]interface{}); ok {
			for _, r := range v {
				if s, ok := r.(string); ok {
					roles = append(roles, s)
				}
			}
		}
	}

	tokens, err := s.GenerateTokenPair(ctx, userID, roles)
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
