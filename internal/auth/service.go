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
	if user.Password != password {
		log.Printf("Invalid password for user: %s", email)
		return nil, nil, errors.New("invalid credentials")
	}

	// Determine roles string for the JWT
	var roles []string
	if user.IsAdmin {
		roles = append(roles, "admin")
	}
	if user.IsTeacher {
		roles = append(roles, "teacher")
	}
	if user.IsLearner {
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
		Role:         tokens.Role, // Note: This might need to be joined or changed to Roles in DTO
	}, user, nil
}

func (s *Service) GenerateTokenPair(ctx context.Context, userID uuid.UUID, roles []string) (*authdto.TokenPair, error) {
	rolesJoined := ""
	if len(roles) > 0 {
		rolesJoined = roles[0] // For backward compatibility in Role field if needed
	}

	// Access Token
	accessClaims := jwt.MapClaims{
		"sub":   userID.String(),
		"roles": roles,
		"exp":   time.Now().Add(120 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.New().String(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"sub":   userID.String(),
		"roles": roles,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.New().String(),
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

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*authdto.TokenResponse, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// For simplicity, we just pass the roles from the old token
	tokens, err := s.GenerateTokenPair(ctx, claims.UserID, []string{claims.Role}) // This logic is simplified
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
