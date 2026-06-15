package service

import (
	"context"
	"errors"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"github.com/farshidboroomand/jobs-collector/internal/repository"
)

// BotService handles business logic related to Bots.
type BotService struct {
	repo repository.BotRepository
}

// NewBotService creates a new BotService.
func NewBotService(r repository.BotRepository) *BotService {
	return &BotService{repo: r}
}

// ListBots retrieves a list of Bots from the repository.
func (s *BotService) ListBots(ctx context.Context, page, pageSize int) ([]domain.Bot, int64, error) {
	return s.repo.FindAll(ctx, page, pageSize)
}

// PublishBot stores a new Bot in the repository.
func (s *BotService) PublishBot(ctx context.Context, bot *domain.Bot) error {
	if bot.Title == "" {
		return errors.New("title cannot be empty")
	}
	return s.repo.Create(ctx, bot)
}
