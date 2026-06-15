package repository

import (
	"context"
	"errors"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"gorm.io/gorm"
)

// gromBotRepository implements the BotRepository interface using GORM for database operations.
type gormBotRepository struct {
	db *gorm.DB
}

// NewBotRepository creates a new Bot repository.
func NewBotRepository(db *gorm.DB) BotRepository {
	return &gormBotRepository{db: db}
}

// Create inserts a single Bot into the database.
func (r *gormBotRepository) Create(ctx context.Context, bot *domain.Bot) error {
	result := r.db.WithContext(ctx).Create(bot)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByID retrieves a bot by their ID.
func (r *gormBotRepository) FindByID(ctx context.Context, id uint) (*domain.Bot, error) {
	var bot domain.Bot
	result := r.db.WithContext(ctx).First(&bot, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &bot, nil
}

// FindAll retrieves all Bots with pagination.
func (r *gormBotRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.Bot, int64, error) {
	var bots []domain.Bot
	var total int64

	r.db.WithContext(ctx).Model(&domain.Bot{}).Count(&total)

	offset := (page - 1) * pageSize

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
