package sessions

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines the interface for session business logic
type Service interface {
	// Session operations
	CreateSession(ctx context.Context, req CreateSessionRequest) (*Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
	ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]Session, error)
	ListSessionsByModule(ctx context.Context, moduleID uuid.UUID) ([]Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error

	// Attempt operations
	CreateAttempt(ctx context.Context, req CreateAttemptRequest) (*Attempt, error)
	GetAttemptByID(ctx context.Context, id uuid.UUID) (*Attempt, error)
	ListAttemptsBySession(ctx context.Context, sessionID uuid.UUID) ([]Attempt, error)
	ListAttemptsByTenant(ctx context.Context, tenantID uuid.UUID) ([]Attempt, error)
	DeleteAttempt(ctx context.Context, id uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new session service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateSession creates a new session with business logic validation
func (s *service) CreateSession(ctx context.Context, req CreateSessionRequest) (*Session, error) {
	// Business logic validation
	if req.ModuleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	if req.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	if req.TenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	// Create the session
	session := Session{
		TenantID: req.TenantID,
		ModuleID: req.ModuleID,
		UserID:   req.UserID,
	}

	return s.repo.CreateSession(ctx, session)
}

// GetSessionByID retrieves a session by ID
func (s *service) GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	if id == uuid.Nil {
		return nil, errors.New("session ID is required")
	}

	return s.repo.GetSessionByID(ctx, id)
}

// ListSessionsByUser retrieves all sessions for a user
func (s *service) ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	return s.repo.ListSessionsByUser(ctx, userID)
}

// ListSessionsByModule retrieves all sessions for a module
func (s *service) ListSessionsByModule(ctx context.Context, moduleID uuid.UUID) ([]Session, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListSessionsByModule(ctx, moduleID)
}

// DeleteSession deletes a session
func (s *service) DeleteSession(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("session ID is required")
	}

	// Check if session exists before deleting
	_, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteSession(ctx, id)
}

// CreateAttempt creates a new attempt with business logic validation
func (s *service) CreateAttempt(ctx context.Context, req CreateAttemptRequest) (*Attempt, error) {
	// Business logic validation
	if req.SessionID == uuid.Nil {
		return nil, errors.New("session ID is required")
	}

	if req.TenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	if req.ArtifactID == uuid.Nil {
		return nil, errors.New("artifact ID is required")
	}

	// Create the attempt
	attempt := Attempt{
		SessionID:  req.SessionID,
		TenantID:   req.TenantID,
		ArtifactID: req.ArtifactID,
		IsCorrect:  req.IsCorrect,
		UserAnswer: req.UserAnswer,
	}

	return s.repo.CreateAttempt(ctx, attempt)
}

// GetAttemptByID retrieves an attempt by ID
func (s *service) GetAttemptByID(ctx context.Context, id uuid.UUID) (*Attempt, error) {
	if id == uuid.Nil {
		return nil, errors.New("attempt ID is required")
	}

	return s.repo.GetAttemptByID(ctx, id)
}

// ListAttemptsBySession retrieves all attempts for a session
func (s *service) ListAttemptsBySession(ctx context.Context, sessionID uuid.UUID) ([]Attempt, error) {
	if sessionID == uuid.Nil {
		return nil, errors.New("session ID is required")
	}

	return s.repo.ListAttemptsBySession(ctx, sessionID)
}

// ListAttemptsByTenant retrieves all attempts for a tenant
func (s *service) ListAttemptsByTenant(ctx context.Context, tenantID uuid.UUID) ([]Attempt, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.ListAttemptsByTenant(ctx, tenantID)
}

// DeleteAttempt deletes an attempt
func (s *service) DeleteAttempt(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("attempt ID is required")
	}

	// Check if attempt exists before deleting
	_, err := s.repo.GetAttemptByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteAttempt(ctx, id)
}
