package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/ajs/currency-api/internal/infrastructure/config"
	"github.com/ajs/currency-api/internal/transport/http"
	"github.com/ajs/go-common/logger"

	_ "github.com/ajs/currency-api/docs"
)

// @title Currency Exchange API
// @version 2.0.0
// @description A modern currency exchange API built with Go and Gin
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@currencyapi.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	cfg, err := config.Load()
	if err != nil {
		log := logger.New("error")
		log.Fatal("Failed to load config", err)
	}

	log := logger.New(cfg.LogLevel)

	server := http.NewServer(cfg, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Failed to start server", err)
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", err)
	}

	log.Info("Server stopped gracefully")
}
