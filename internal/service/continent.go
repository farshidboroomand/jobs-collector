package service

import (
	"context"
	"errors"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"github.com/farshidboroomand/jobs-collector/internal/repository"
)

// ContinentService handles business logic relate d to Continents.
type ContinentService struct {
	repo repository.ContinentRepository
}

// NewContinentService creates a new ContinentService.
func NewContinentService(r repository.ContinentRepository) *ContinentService {
	return &ContinentService{repo: r}
}

// GetContinentByID retrieves a Continent by its ID.
func (s *ContinentService) GetContinentByID(ctx context.Context, id uint) (*domain.Continent, error) {
	return s.repo.FindByID(ctx, id)
}

// PublishContinent stores a new Continent in the repository.
func (s *ContinentService) PublishContinent(ctx context.Context, continent *domain.Continent) error {
	if continent.Name == "" {
		return errors.New("name cannot be empty")
	}
	return s.repo.Create(ctx, continent)
}
