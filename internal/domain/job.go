package domain

// Job represents a job posting in the system.
type Job struct {
	ID        int
	Title     string
	Company   string
	CountryID int
}
