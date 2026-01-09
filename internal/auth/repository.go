package auth

import (
	"context"
	"learning-core-api/internal/persistance/store"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*store.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*store.User, error)
	CreateRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt interface{}) error
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

func (r *sqlRepository) CreateRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt interface{}) error {
	// Sessions are removed as per user feedback
	return nil
}

func (r *sqlRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	// Sessions are removed as per user feedback
	return nil
}
