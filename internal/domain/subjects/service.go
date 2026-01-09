package subjects

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

// Service defines the business logic interface for subjects
type Service interface {
	// CreateSubject creates a new subject with validation
	CreateSubject(ctx context.Context, req CreateSubjectRequest) (*Subject, error)
	
	// GetSubject retrieves a subject by ID
	GetSubject(ctx context.Context, id uuid.UUID) (*Subject, error)
	
	// UpdateSubject updates an existing subject with validation
	UpdateSubject(ctx context.Context, id uuid.UUID, req UpdateSubjectRequest) (*Subject, error)
	
	// DeleteSubject deletes a subject
	DeleteSubject(ctx context.Context, id uuid.UUID) error
	
	// ListUserSubjects retrieves all subjects for a user
	ListUserSubjects(ctx context.Context, userID uuid.UUID) ([]*Subject, error)
}

type service struct {
	repo Repository
}

// NewService creates a new subject service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateSubject(ctx context.Context, req CreateSubjectRequest) (*Subject, error) {
	// Validate the request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Create the subject
	return s.repo.Create(ctx, req)
}

func (s *service) GetSubject(ctx context.Context, id uuid.UUID) (*Subject, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateSubject(ctx context.Context, id uuid.UUID, req UpdateSubjectRequest) (*Subject, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Validate the request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, req)
}

func (s *service) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidUserID
	}

	// Check if subject exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) ListUserSubjects(ctx context.Context, userID uuid.UUID) ([]*Subject, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return s.repo.ListByUser(ctx, userID)
}

// validateCreateRequest validates the create subject request
func (s *service) validateCreateRequest(req CreateSubjectRequest) error {
	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidSubjectName
	}
	
	if len(req.Name) > 255 {
		return ErrSubjectNameTooLong
	}

	// Validate user ID
	if req.UserID == uuid.Nil {
		return ErrInvalidUserID
	}

	// Validate description if provided
	if req.Description != nil && len(*req.Description) > 1000 {
		return ErrInvalidDescription
	}

	return nil
}

// validateUpdateRequest validates the update subject request
func (s *service) validateUpdateRequest(req UpdateSubjectRequest) error {
	// Validate name if provided
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return ErrInvalidSubjectName
		}
		
		if len(*req.Name) > 255 {
			return ErrSubjectNameTooLong
		}
	}

	// Validate description if provided
	if req.Description != nil && len(*req.Description) > 1000 {
		return ErrInvalidDescription
	}

	return nil
}
