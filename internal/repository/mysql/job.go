package mysql

import (
	"database/sql"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"github.com/farshidboroomand/jobs-collector/internal/repository"
)

type jobRepo struct{ db *sql.DB }

// NewJobRepository creates a new JobRepository backed by SQL database.
func NewJobRepository(db *sql.DB) repository.JobRepository {
	return &jobRepo{
		db: db,
	}
}

func (r *jobRepo) Create(j *domain.Job) error {
	_, err := r.db.Exec(
		"INSERT INTO jobs (id, title, company, country_id) VALUES (?, ?, ?, ?)",
		j.ID,
		j.Title,
		j.Company,
		j.CountryID,
	)
	return err
}

func (r *jobRepo) GetByID(id int) (*domain.Job, error) {
	var j domain.Job
	err := r.db.QueryRow("SELECT id, title, company FROM jobs WHERE id = ?", id).Scan(&j.ID, &j.Title, &j.Company)
	return &j, err
}
