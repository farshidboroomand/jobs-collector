package service

import (
	"errors"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"github.com/farshidboroomand/jobs-collector/internal/repository"
)

// JobService handles business logic related to jobs.
type JobService struct {
	repo repository.JobRepository
}

// NewJobService creates a new JobService.
func NewJobService(r repository.JobRepository) *JobService {
	return &JobService{repo: r}
}

// PublishJob stores a new job in the repository.
func (s *JobService) PublishJob(job *domain.Job) error {
	if job.Title == "" {
		return errors.New("title cannot be empty")
	}
	return s.repo.Create(job)
}
