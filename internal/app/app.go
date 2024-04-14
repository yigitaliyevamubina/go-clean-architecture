package app

import (
	"fmt"
	"fourth-exam/post-service-clean-arch/config"
	v1 "fourth-exam/post-service-clean-arch/internal/controller/http/v1"
	"fourth-exam/post-service-clean-arch/internal/usecase"
	"fourth-exam/post-service-clean-arch/internal/usecase/repo"
	"fourth-exam/post-service-clean-arch/pkg/httpserver"
	"fourth-exam/post-service-clean-arch/pkg/logger"
	"fourth-exam/post-service-clean-arch/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	postUseCase := usecase.New(
		repo.New(pg),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, postUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <- interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <- httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - run - httpServer.Shutdown: %w", err))
	}
}
