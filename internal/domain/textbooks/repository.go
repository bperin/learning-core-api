package textbooks

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for textbook data operations
type Repository interface {
	// CreateSubject creates a new subject in the database
	CreateSubject(ctx context.Context, subject *Subject) error

	// CreateSubSubject creates a new sub-subject in the database
	CreateSubSubject(ctx context.Context, subSubject *SubSubject) error

	// GetAllSubjects retrieves all subjects from the database
	GetAllSubjects(ctx context.Context) ([]Subject, error)

	// GetSubjectByID retrieves a subject by ID
	GetSubjectByID(ctx context.Context, id uuid.UUID) (*Subject, error)

	// GetSubjectByName retrieves a subject by name
	GetSubjectByName(ctx context.Context, name string) (*Subject, error)

	// GetSubSubjectsBySubjectID retrieves all sub-subjects for a subject
	GetSubSubjectsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]SubSubject, error)

	// DeleteSubject deletes a subject and its sub-subjects
	DeleteSubject(ctx context.Context, id uuid.UUID) error

	// DeleteAllSubjects deletes all subjects (for re-scraping)
	DeleteAllSubjects(ctx context.Context) error
}

// NewRepository creates a new textbook repository
// For now, this is a placeholder that will be implemented with actual database queries
func NewRepository() Repository {
	// TODO: Implement with actual database queries using SQLC
	return &mockRepository{}
}

// mockRepository is a placeholder implementation
type mockRepository struct{}

func (m *mockRepository) CreateSubject(ctx context.Context, subject *Subject) error {
	return nil
}

func (m *mockRepository) CreateSubSubject(ctx context.Context, subSubject *SubSubject) error {
	return nil
}

func (m *mockRepository) GetAllSubjects(ctx context.Context) ([]Subject, error) {
	return []Subject{}, nil
}

func (m *mockRepository) GetSubjectByID(ctx context.Context, id uuid.UUID) (*Subject, error) {
	return nil, nil
}

func (m *mockRepository) GetSubjectByName(ctx context.Context, name string) (*Subject, error) {
	return nil, nil
}

func (m *mockRepository) GetSubSubjectsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]SubSubject, error) {
	return []SubSubject{}, nil
}

func (m *mockRepository) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockRepository) DeleteAllSubjects(ctx context.Context) error {
	return nil
}
