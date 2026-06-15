package application

import (
	"context"
	"net/http"
	"time"

	"github.com/farshidboroomand/jobs-collector/config"
	"github.com/farshidboroomand/jobs-collector/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// App encapsulates the application configuration, router, services, and handlers.
type App struct {
	Config     *config.Config
	Router     *gin.Engine
	Services   *service.Services
	BotHandler *BotHandler
}

// NewApplication initializes a new App instance with the provided configuration and services.
func NewApplication(cfg *config.Config, services *service.Services) (*App, error) {
	a := &App{
		Config:     cfg,
		Services:   services,
		BotHandler: NewBotHandler(services),
	}

	a.registerRouter()
	a.registerRoutes()

	return a, nil
}

// registerRoutes sets up the API routes for the application.
func (a *App) registerRoutes() {
	api := a.Router.Group("/api")
	{
		// Cleaner routing: all Bot routes handled by BotHandler
		api.GET("/bots", a.BotHandler.List)
		api.POST("/bots", a.BotHandler.CreateBot) // Assuming CreateBot is implemented in BotHandler
	}
}

// registerRouter configures the Gin router based on the application environment.
func (a *App) registerRouter() {
	switch a.Config.ENV {
	case config.EnvTesting:
		gin.SetMode(gin.TestMode)
	case config.EnvStaging:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger()) // Added Logger for better visibility
	a.Router = r
}

// Run starts the HTTP server and listens for incoming requests. It also handles graceful shutdown on context cancellation.
func (a *App) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:              ":" + a.Config.APIPORT,
		Handler:           a.Router,
		ReadHeaderTimeout: time.Second * time.Duration(a.Config.READHEADERTIMEOUT),
	}

	go func() {
		log.Infof("HTTP server listening on %s", a.Config.APIPORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("HTTP server error")
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(a.Config.GRACEFULSHUTDOWNTIMEOUT),
	)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
