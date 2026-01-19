package subjects

import (
	"context"

	"learning-core-api/internal/persistance/store"

	"github.com/google/uuid"
)

type repositoryImpl struct {
	store *store.Queries
}

func NewRepository(s *store.Queries) Repository {
	return &repositoryImpl{store: s}
}

func (r *repositoryImpl) ListSubjects(ctx context.Context) ([]Subject, error) {
	rows, err := r.store.ListSubjects(ctx)
	if err != nil {
		return nil, err
	}

	subjects := make([]Subject, len(rows))
	for i, row := range rows {
		subjects[i] = Subject{
			ID:        row.ID,
			Name:      row.Name,
			Url:       row.Url,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}
	return subjects, nil
}

func (r *repositoryImpl) ListSubSubjectsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]SubSubject, error) {
	rows, err := r.store.ListSubSubjectsBySubjectID(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	subSubjects := make([]SubSubject, len(rows))
	for i, row := range rows {
		subSubjects[i] = SubSubject{
			ID:        row.ID,
			SubjectID: row.SubjectID,
			Name:      row.Name,
			Url:       row.Url,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}
	return subSubjects, nil
}
