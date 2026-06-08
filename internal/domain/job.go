package domain

import "github.com/google/uuid"

// Job represents a job posting in the system.
type Job struct {
	ID        uuid.UUID
	Title     string
	Company   string
	CountryID int
}
