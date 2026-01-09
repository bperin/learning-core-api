package subjects

import (
	"context"
	"database/sql"
	"errors"
	"learning-core-api/internal/persistance/store"

	"github.com/google/uuid"
)

type sqlRepository struct {
	queries *store.Queries
}

// NewRepository creates a new subject repository
func NewRepository(queries *store.Queries) Repository {
	return &sqlRepository{
		queries: queries,
	}
}

func (r *sqlRepository) Create(ctx context.Context, req CreateSubjectRequest) (*Subject, error) {
	params := store.CreateSubjectParams{
		ID:     uuid.New(),
		Name:   req.Name,
		UserID: req.UserID,
	}
	
	if req.Description != nil {
		params.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	dbSubject, err := r.queries.CreateSubject(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainSubject(dbSubject), nil
}

func (r *sqlRepository) GetByID(ctx context.Context, id uuid.UUID) (*Subject, error) {
	dbSubject, err := r.queries.GetSubject(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubjectNotFound
		}
		return nil, err
	}

	return toDomainSubject(dbSubject), nil
}

func (r *sqlRepository) Update(ctx context.Context, id uuid.UUID, req UpdateSubjectRequest) (*Subject, error) {
	// First get the existing subject to preserve unchanged fields
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	params := store.UpdateSubjectParams{
		ID:   id,
		Name: existing.Name,
	}

	// Update name if provided
	if req.Name != nil {
		params.Name = *req.Name
	}

	// Update description if provided
	if req.Description != nil {
		params.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	} else if existing.Description != nil {
		params.Description = sql.NullString{
			String: *existing.Description,
			Valid:  true,
		}
	}

	dbSubject, err := r.queries.UpdateSubject(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainSubject(dbSubject), nil
}

func (r *sqlRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteSubject(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *sqlRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*Subject, error) {
	dbSubjects, err := r.queries.ListSubjectsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	subjects := make([]*Subject, len(dbSubjects))
	for i, dbSubject := range dbSubjects {
		subjects[i] = toDomainSubject(dbSubject)
	}

	return subjects, nil
}

// Helper function to convert store.Subject to domain Subject
func toDomainSubject(dbSubject store.Subject) *Subject {
	subject := &Subject{
		ID:        dbSubject.ID,
		Name:      dbSubject.Name,
		UserID:    dbSubject.UserID,
		CreatedAt: dbSubject.CreatedAt,
		UpdatedAt: dbSubject.UpdatedAt,
	}

	if dbSubject.Description.Valid {
		subject.Description = &dbSubject.Description.String
	}

	return subject
}
