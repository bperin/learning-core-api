package subjects

import (
	"context"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// Repository defines storage operations for subjects.
type Repository interface {
	Create(ctx context.Context, subject Subject) (*Subject, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Subject, error)
	GetByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*Subject, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Subject, error)
	Update(ctx context.Context, id uuid.UUID, name string, description string) (*Subject, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	queries *store.Queries
}

// NewRepository creates a new subjects repository.
func NewRepository(queries *store.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) Create(ctx context.Context, subject Subject) (*Subject, error) {
	dbSubject, err := r.queries.CreateSubject(ctx, store.CreateSubjectParams{
		UserID:      subject.UserID,
		Name:        subject.Name,
		Description: subject.Description,
	})
	if err != nil {
		return nil, err
	}

	return &Subject{
		ID:          dbSubject.ID,
		UserID:      dbSubject.UserID,
		Name:        dbSubject.Name,
		Description: dbSubject.Description,
		CreatedAt:   dbSubject.CreatedAt,
	}, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Subject, error) {
	dbSubject, err := r.queries.GetSubject(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Subject{
		ID:          dbSubject.ID,
		UserID:      dbSubject.UserID,
		Name:        dbSubject.Name,
		Description: dbSubject.Description,
		CreatedAt:   dbSubject.CreatedAt,
	}, nil
}

func (r *repository) GetByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*Subject, error) {
	dbSubject, err := r.queries.GetSubjectByUserAndName(ctx, store.GetSubjectByUserAndNameParams{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	return &Subject{
		ID:          dbSubject.ID,
		UserID:      dbSubject.UserID,
		Name:        dbSubject.Name,
		Description: dbSubject.Description,
		CreatedAt:   dbSubject.CreatedAt,
	}, nil
}

func (r *repository) ListByUser(ctx context.Context, userID uuid.UUID) ([]Subject, error) {
	dbSubjects, err := r.queries.ListSubjectsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	subjects := make([]Subject, len(dbSubjects))
	for i, dbSubject := range dbSubjects {
		subjects[i] = Subject{
			ID:          dbSubject.ID,
			UserID:      dbSubject.UserID,
			Name:        dbSubject.Name,
			Description: dbSubject.Description,
			CreatedAt:   dbSubject.CreatedAt,
		}
	}

	return subjects, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, name string, description string) (*Subject, error) {
	dbSubject, err := r.queries.UpdateSubject(ctx, store.UpdateSubjectParams{
		ID:          id,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return &Subject{
		ID:          dbSubject.ID,
		UserID:      dbSubject.UserID,
		Name:        dbSubject.Name,
		Description: dbSubject.Description,
		CreatedAt:   dbSubject.CreatedAt,
	}, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSubject(ctx, id)
}
