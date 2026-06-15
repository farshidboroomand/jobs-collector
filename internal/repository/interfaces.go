package repository

import (
	"context"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
)

// BotRepository defines the storage operations for bots.
type BotRepository interface {
	Create(ctx context.Context, Bot *domain.Bot) error
	FindByID(ctx context.Context, id uint) (*domain.Bot, error)
	FindAll(ctx context.Context, page, pageSize int) ([]domain.Bot, int64, error)
}

// ContinentRepository defines the storage operations for continents.
type ContinentRepository interface {
	Create(ctx context.Context, continent *domain.Continent) error
	FindByID(ctx context.Context, id uint) (*domain.Continent, error)
}
