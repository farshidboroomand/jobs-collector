package repository

import "github.com/farshidboroomand/jobs-collector/internal/domain"

// JobRepository defines the storage operations for jobs.
type JobRepository interface {
	Create(job *domain.Job) error
	GetByID(id int) (*domain.Job, error)
}
