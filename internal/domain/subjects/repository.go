package subjects

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	ListSubjects(ctx context.Context) ([]Subject, error)
	ListSubSubjectsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]SubSubject, error)
}
