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

type App struct {
	Config   *config.Config
	Router   *gin.Engine
	Services *service.Services
}

func NewApplication(cfg *config.Config, services *service.Services) (*App, error) {
	a := &App{
		Config:   cfg,
		Services: services,
	}

	a.registerRouter()
	a.registerRoutes()

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:              ":" + a.Config.APIPORT,
		Handler:           a.Router,
		ReadHeaderTimeout: time.Second * time.Duration(a.Config.READHEADERTIMEOUT),
	}

	go func() {
		log.Infof("http server listening on %s", a.Config.APIPORT)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("http server error")
		}
	}()

	<-ctx.Done()

	log.Info("shutting down http server")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(a.Config.GRACEFULSHUTDOWNTIMEOUT),
	)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

func (a *App) registerRouter() {
	switch a.Config.ENV {
	case config.EnvTesting:
		gin.SetMode(gin.TestMode)
	case config.EnvStaging:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(
		gin.Recovery(),
	)

	a.Router = router
}

func (a *App) registerRoutes() {
	api := a.Router.Group("/api")

	api.GET("/bots", a.getBots)
}

func (a *App) getBots(c *gin.Context) {
	log.Info("get bots")
	c.JSON(http.StatusOK, nil)
}
