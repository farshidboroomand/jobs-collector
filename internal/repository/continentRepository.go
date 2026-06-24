package repository

import (
	"context"
	"errors"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"gorm.io/gorm"
)

// gormContinentRepository implements the ContinentRepository interface using GORM for database operations.
type gormContinentRepository struct {
	db *gorm.DB
}

// NewContinentRepository creates a new instance of gormContinentRepository.
func NewContinentRepository(db *gorm.DB) ContinentRepository {
	return &gormContinentRepository{db: db}
}

// Create inserts a new continent record into the database.
func (r *gormContinentRepository) Create(ctx context.Context, continent *domain.Continent) error {
	result := r.db.WithContext(ctx).Create(continent)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByID retrieves a continent record by its ID from the database.
func (r *gormContinentRepository) FindByID(ctx context.Context, id uint) (*domain.Continent, error) {
	var continent domain.Continent
	result := r.db.WithContext(ctx).First(&continent, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &continent, nil
}
