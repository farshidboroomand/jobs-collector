package repository

import (
	"context"
	"errors"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"gorm.io/gorm"
)

type gormBotRepository struct {
	db *gorm.DB
}

// NewBotRepository creates a new Bot repository
func NewBotRepository(db *gorm.DB) BotRepository {
	return &gormBotRepository{db: db}
}

// Create inserts a single Bot into the database
func (r *gormBotRepository) Create(ctx context.Context, bot *domain.Bot) error {
	// Create will insert the record and update the model with generated values
	result := r.db.WithContext(ctx).Create(bot)
	if result.Error != nil {
		return result.Error
	}
	// Bot.ID is now populated with the generated ID
	return nil
}

// FindByID retrieves a bot by their ID
func (r *gormBotRepository) FindByID(ctx context.Context, id uint) (*domain.Bot, error) {
	var bot domain.Bot
	// First finds the first record matching the condition
	result := r.db.WithContext(ctx).First(&bot, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found
		}
		return nil, result.Error
	}
	return &bot, nil
}

// FindAll retrieves all Bots with pagination
func (r *gormBotRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.Bot, int64, error) {
	var bots []domain.Bot
	var total int64

	// Count total records first
	r.db.WithContext(ctx).Model(&domain.Bot{}).Count(&total)

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Retrieve paginated results with ordering
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bots)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return bots, total, nil
}
