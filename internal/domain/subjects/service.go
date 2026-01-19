package subjects

import (
	"context"
)

type Service interface {
	ListAll(ctx context.Context) ([]Subject, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) ListAll(ctx context.Context) ([]Subject, error) {
	subjects, err := s.repo.ListSubjects(ctx)
	if err != nil {
		return nil, err
	}

	for i := range subjects {
		subSubjects, err := s.repo.ListSubSubjectsBySubjectID(ctx, subjects[i].ID)
		if err != nil {
			return nil, err
		}
		subjects[i].SubSubjects = subSubjects
	}

	return subjects, nil
}
