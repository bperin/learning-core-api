package auth

import (
	"context"
	"learning-core-api/internal/persistance/store"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*store.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*store.User, error)

	CreateRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*store.Session, error) // Note: Renamed from RefreshToken to Session to match consolidated schema if needed, checking store.Models
	DeleteRefreshToken(ctx context.Context, token string) error
}

type sqlRepository struct {
	q store.Querier
}

func NewRepository(q store.Querier) Repository {
	return &sqlRepository{q: q}
}

func (r *sqlRepository) GetUserByEmail(ctx context.Context, email string) (*store.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqlRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*store.User, error) {
	u, err := r.q.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqlRepository) CreateRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	// In the consolidated schema, sessions are used for refresh tokens
	_, err := r.q.CreateSession(ctx, store.CreateSessionParams{
		ID:        uuid.New().String(), // Session ID
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	return err
}

func (r *sqlRepository) GetRefreshToken(ctx context.Context, token string) (*store.Session, error) {
	rt, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *sqlRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.q.DeleteSession(ctx, token)
}
