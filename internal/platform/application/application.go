package application

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/farshidboroomand/jobs-collector/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// App represents application configs in the system.
type App struct {
	Config *config.Config
	Router *gin.Engine
}

// NewApplication creates an instance of router and routes.
func NewApplication(_ context.Context, cfg *config.Config) (*App, error) {
	a := &App{Config: cfg}

	a.registerRouter()
	a.registerRoutes()

	return a, nil
}

// RunServices run routers.
func (a *App) RunServices(ctx context.Context, wg *sync.WaitGroup) {
	a.runRouter(ctx, wg)
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
		// log.Logger(log.StandardLogger()),
		gin.Recovery(),
	)
	a.Router = router
}

func (a *App) runRouter(ctx context.Context, wg *sync.WaitGroup) {
	srv := &http.Server{
		Addr:              ":" + a.Config.APIPORT,
		Handler:           a.Router,
		ReadHeaderTimeout: time.Second * time.Duration(a.Config.READHEADERTIMEOUT),
	}

	runAsync(
		ctx, wg, func(ctx context.Context) {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.WithError(err).Fatal("error on running router")
				}
			}()

			<-ctx.Done()

			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				time.Second*time.Duration(a.Config.GRACEFULSHUTDOWNTIMEOUT),
			)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.WithContext(ctx).WithError(err).Error("could not gracefully shutdown the server")
			}
			log.Debug("router successfully closed")
		},
	)
}

func runAsync(ctx context.Context, wg *sync.WaitGroup, fn func(ctx context.Context)) {
	wg.Go(func() {
		fn(ctx)
	})
}
