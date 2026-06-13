package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/farshidboroomand/jobs-collector/config"
	"github.com/farshidboroomand/jobs-collector/internal/platform/application"
	"github.com/farshidboroomand/jobs-collector/internal/platform/database"
	"github.com/farshidboroomand/jobs-collector/internal/repository"
	"github.com/farshidboroomand/jobs-collector/internal/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("failed to load configuration")
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.WithError(err).Fatal("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Fatal("failed to get sql db")
	}
	log.Info("database connected")
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.WithError(err).Error("failed to close sql db")
		}
	}()
	// repositories
	botRepo := repository.NewBotRepository(db)

	// services container
	services := &service.Services{
		Bot: service.NewBotService(botRepo),
	}

	// application
	app, err := application.NewApplication(cfg, services)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize application")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		log.WithError(err).Fatal("application stopped with error")
	}

	log.Info("application stopped")
}
