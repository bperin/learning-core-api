package subjects

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines subject business logic.
type Service interface {
	Create(ctx context.Context, req CreateSubjectRequest) (*Subject, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Subject, error)
	GetByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*Subject, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Subject, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateSubjectRequest) (*Subject, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

// NewService creates a new subjects service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateSubjectRequest) (*Subject, error) {
	if req.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	existing, err := s.repo.GetByUserAndName(ctx, req.UserID, req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("subject with this name already exists")
	}

	subject := Subject{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
	}

	return s.repo.Create(ctx, subject)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Subject, error) {
	if id == uuid.Nil {
		return nil, errors.New("subject ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*Subject, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	return s.repo.GetByUserAndName(ctx, userID, name)
}

func (s *service) ListByUser(ctx context.Context, userID uuid.UUID) ([]Subject, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateSubjectRequest) (*Subject, error) {
	if id == uuid.Nil {
		return nil, errors.New("subject ID is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, req.Name, req.Description)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("subject ID is required")
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
