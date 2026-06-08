package main

import (
	"database/sql"

	"github.com/farshidboroomand/jobs-collector/config"
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
	if err != nil || db.Ping() != nil {
		log.WithError(err).Fatal("Failed to connect to MySQL")
	}

	// 3. Initialize Repositories (SRP: Each has its own file/struct)
	jobRepo := mysql.NewJobRepository(db)

	// 4. Initialize Services (Injecting Repositories)
	jobService := service.NewJobService(jobRepo)

	log.Info("Application initialized with separate repositories.")
	_ = jobService
}
