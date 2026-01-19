package subjects

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListSubjects(ctx context.Context) ([]Subject, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Subject), args.Error(1)
}

func (m *MockRepository) ListSubSubjectsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]SubSubject, error) {
	args := m.Called(ctx, subjectID)
	return args.Get(0).([]SubSubject), args.Error(1)
}

func TestListAll(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	subjectID := uuid.New()
	subjects := []Subject{
		{
			ID:        subjectID,
			Name:      "Computer Science",
			Url:       "http://cs.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	subSubjects := []SubSubject{
		{
			ID:        uuid.New(),
			SubjectID: subjectID,
			Name:      "Algorithms",
			Url:       "http://algo.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("ListSubjects", ctx).Return(subjects, nil)
	mockRepo.On("ListSubSubjectsBySubjectID", ctx, subjectID).Return(subSubjects, nil)

	result, err := service.ListAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Computer Science", result[0].Name)
	assert.Len(t, result[0].SubSubjects, 1)
	assert.Equal(t, "Algorithms", result[0].SubSubjects[0].Name)

	mockRepo.AssertExpectations(t)
}
