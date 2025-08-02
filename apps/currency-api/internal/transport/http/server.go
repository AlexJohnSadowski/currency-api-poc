package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ajs/currency-api/internal/app/handlers"
	"github.com/ajs/currency-api/internal/app/queries"
	"github.com/ajs/currency-api/internal/infrastructure/config"
	"github.com/ajs/currency-api/internal/infrastructure/repositories"
	"github.com/ajs/currency-api/internal/transport/http/routes"
	"github.com/ajs/go-common/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	logger logger.Logger
	server *http.Server
}

func NewServer(cfg *config.Config, log logger.Logger) *Server {
	return &Server{
		config: cfg,
		logger: log,
	}
}

func (s *Server) Start() error {
	gin.SetMode(s.config.GinMode)

	r := gin.New()
	r.Use(gin.Recovery())

	ratesRepo := repositories.NewRatesRepositoryImpl(s.config, s.logger)

	ratesQueryHandler := queries.NewGetRatesQueryHandler(ratesRepo)
	exchangeQueryHandler := queries.NewExchangeQueryHandler()

	healthHandler := handlers.NewHealthHandler(s.config, s.logger)
	ratesHandler := handlers.NewRatesHandler(ratesQueryHandler, s.logger)
	exchangeHandler := handlers.NewExchangeHandler(exchangeQueryHandler, s.logger)

	routes.SetupRoutes(r, healthHandler, ratesHandler, exchangeHandler)

	s.server = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info(fmt.Sprintf("🚀 Starting server on port %s", s.config.Port))
	s.logger.Info(fmt.Sprintf("🔧 Environment: %s", s.config.Environment))
	s.logger.Info(fmt.Sprintf("⚙️ Gin Mode: %s", s.config.GinMode))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("🛑 Shutting down server...")
	return s.server.Shutdown(ctx)
}
