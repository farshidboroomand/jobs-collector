package service

import (
	"testing"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Object
type MockJobRepo struct{ mock.Mock }

func (m *MockJobRepo) Create(j *domain.Job) error {
	return m.Called(j).Error(0)
}

func (m *MockJobRepo) GetByID(id int) (*domain.Job, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Job), args.Error(1)
}

// Test Case
func TestPublishJob_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockJobRepo)
	job := &domain.Job{Title: "Go Developer", Company: "Google"}
	mockRepo.On("Create", job).Return(nil)

	// Act
	svc := NewJobService(mockRepo)
	err := svc.PublishJob(job)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
