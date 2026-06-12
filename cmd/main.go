package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/farshidboroomand/jobs-collector/config"
	"github.com/farshidboroomand/jobs-collector/internal/platform/application"
	"github.com/farshidboroomand/jobs-collector/internal/repository/mysql"
	"github.com/farshidboroomand/jobs-collector/internal/service"
	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// 1. DSN Construction
	var dsn string
	switch cfg.DBCONNECTION {
	case "mysql":
		dsn = config.DSN(*cfg)
	default:
		log.Fatal("Unsupported driver")
	}

	// 2. Database Connection & Ping
	db, err := sql.Open(cfg.DBCONNECTION, dsn)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)

	if err != nil || db.Ping() != nil {
		log.WithError(err).Fatal("Failed to connect to MySQL")
	}
	defer db.Close()

	// 3. Initialize Repositories (SRP: Each has its own file/struct)
	jobRepo := mysql.NewJobRepository(db)

	// 4. Initialize Services (Injecting Repositories)
	jobService := service.NewJobService(jobRepo)
	_ = jobService

	// 5. Initialize Application
	ctx, cancel := context.WithCancel(context.Background())

	app, err := application.NewApplication(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	app.RunServices(ctx, wg)

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	select {
	case <-closeSignal:
		log.Info("Terminating by os signal")
	case <-ctx.Done():
		log.Info("Terminating by context cancellation")
	}

	cancel()
	wg.Wait()
}
