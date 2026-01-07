package sessions

import (
	"context"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// Repository defines the interface for session operations
type Repository interface {
	// Session operations
	CreateSession(ctx context.Context, session Session) (*Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
	ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]Session, error)
	ListSessionsByModule(ctx context.Context, moduleID uuid.UUID) ([]Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error

	// Attempt operations
	CreateAttempt(ctx context.Context, attempt Attempt) (*Attempt, error)
	GetAttemptByID(ctx context.Context, id uuid.UUID) (*Attempt, error)
	ListAttemptsBySession(ctx context.Context, sessionID uuid.UUID) ([]Attempt, error)
	ListAttemptsByTenant(ctx context.Context, tenantID uuid.UUID) ([]Attempt, error)
	DeleteAttempt(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new session repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// CreateSession creates a new session
func (r *repository) CreateSession(ctx context.Context, session Session) (*Session, error) {
	params := store.CreateSessionParams{
		TenantID: session.TenantID,
		ModuleID: session.ModuleID,
		UserID:   session.UserID,
	}

	dbSession, err := r.queries.CreateSession(ctx, params)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        dbSession.ID,
		TenantID:  dbSession.TenantID,
		ModuleID:  dbSession.ModuleID,
		UserID:    dbSession.UserID,
		CreatedAt: dbSession.CreatedAt,
	}, nil
}

// GetSessionByID retrieves a session by ID
func (r *repository) GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	dbSession, err := r.queries.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        dbSession.ID,
		TenantID:  dbSession.TenantID,
		ModuleID:  dbSession.ModuleID,
		UserID:    dbSession.UserID,
		CreatedAt: dbSession.CreatedAt,
	}, nil
}

// ListSessionsByUser retrieves all sessions for a user
func (r *repository) ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	dbSessions, err := r.queries.ListSessionsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = Session{
			ID:        dbSession.ID,
			TenantID:  dbSession.TenantID,
			ModuleID:  dbSession.ModuleID,
			UserID:    dbSession.UserID,
			CreatedAt: dbSession.CreatedAt,
		}
	}

	return sessions, nil
}

// ListSessionsByModule retrieves all sessions for a module
func (r *repository) ListSessionsByModule(ctx context.Context, moduleID uuid.UUID) ([]Session, error) {
	dbSessions, err := r.queries.ListSessionsByModule(ctx, moduleID)
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = Session{
			ID:        dbSession.ID,
			TenantID:  dbSession.TenantID,
			ModuleID:  dbSession.ModuleID,
			UserID:    dbSession.UserID,
			CreatedAt: dbSession.CreatedAt,
		}
	}

	return sessions, nil
}

// DeleteSession deletes a session
func (r *repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSession(ctx, id)
}

// CreateAttempt creates a new attempt
func (r *repository) CreateAttempt(ctx context.Context, attempt Attempt) (*Attempt, error) {
	params := store.CreateAttemptParams{
		SessionID:  attempt.SessionID,
		TenantID:   attempt.TenantID,
		ArtifactID: attempt.ArtifactID,
		IsCorrect:  attempt.IsCorrect,
		UserAnswer: pqtype.NullRawMessage{RawMessage: attempt.UserAnswer, Valid: attempt.UserAnswer != nil},
	}

	dbAttempt, err := r.queries.CreateAttempt(ctx, params)
	if err != nil {
		return nil, err
	}

	return &Attempt{
		ID:         dbAttempt.ID,
		SessionID:  dbAttempt.SessionID,
		TenantID:   dbAttempt.TenantID,
		ArtifactID: dbAttempt.ArtifactID,
		IsCorrect:  dbAttempt.IsCorrect,
		UserAnswer: dbAttempt.UserAnswer.RawMessage,
		CreatedAt:  dbAttempt.CreatedAt,
	}, nil
}

// GetAttemptByID retrieves an attempt by ID
func (r *repository) GetAttemptByID(ctx context.Context, id uuid.UUID) (*Attempt, error) {
	dbAttempt, err := r.queries.GetAttempt(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Attempt{
		ID:         dbAttempt.ID,
		SessionID:  dbAttempt.SessionID,
		TenantID:   dbAttempt.TenantID,
		ArtifactID: dbAttempt.ArtifactID,
		IsCorrect:  dbAttempt.IsCorrect,
		UserAnswer: dbAttempt.UserAnswer.RawMessage,
		CreatedAt:  dbAttempt.CreatedAt,
	}, nil
}

// ListAttemptsBySession retrieves all attempts for a session
func (r *repository) ListAttemptsBySession(ctx context.Context, sessionID uuid.UUID) ([]Attempt, error) {
	dbAttempts, err := r.queries.ListAttemptsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	attempts := make([]Attempt, len(dbAttempts))
	for i, dbAttempt := range dbAttempts {
		attempts[i] = Attempt{
			ID:         dbAttempt.ID,
			SessionID:  dbAttempt.SessionID,
			TenantID:   dbAttempt.TenantID,
			ArtifactID: dbAttempt.ArtifactID,
			IsCorrect:  dbAttempt.IsCorrect,
			UserAnswer: dbAttempt.UserAnswer.RawMessage,
			CreatedAt:  dbAttempt.CreatedAt,
		}
	}

	return attempts, nil
}

// ListAttemptsByTenant retrieves all attempts for a tenant
func (r *repository) ListAttemptsByTenant(ctx context.Context, tenantID uuid.UUID) ([]Attempt, error) {
	dbAttempts, err := r.queries.ListAttemptsByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	attempts := make([]Attempt, len(dbAttempts))
	for i, dbAttempt := range dbAttempts {
		attempts[i] = Attempt{
			ID:         dbAttempt.ID,
			SessionID:  dbAttempt.SessionID,
			TenantID:   dbAttempt.TenantID,
			ArtifactID: dbAttempt.ArtifactID,
			IsCorrect:  dbAttempt.IsCorrect,
			UserAnswer: dbAttempt.UserAnswer.RawMessage,
			CreatedAt:  dbAttempt.CreatedAt,
		}
	}

	return attempts, nil
}

// DeleteAttempt deletes an attempt
func (r *repository) DeleteAttempt(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAttempt(ctx, id)
}
