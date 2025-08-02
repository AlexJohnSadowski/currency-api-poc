package handlers

import (
	"net/http"
	"time"

	"github.com/ajs/currency-api/internal/infrastructure/config"
	"github.com/ajs/go-common/logger"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	config *config.Config
	logger logger.Logger
}

func NewHealthHandler(cfg *config.Config, log logger.Logger) *HealthHandler {
	return &HealthHandler{
		config: cfg,
		logger: log,
	}
}

// @Summary Health check
// @Description Get the current health status of the API
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"service":   "currency-exchange-api",
		"version":   "2.0.0",
		"timestamp": time.Now().Unix(),
		"environment": map[string]interface{}{
			"mode":     h.config.Environment,
			"gin_mode": h.config.GinMode,
			"port":     h.config.Port,
		},
		"framework":  "gin-gonic",
		"nx_plugin":  "@naxodev/gonx",
		"go_version": "1.24",
		"features": []string{
			"CQRS Pattern",
			"Domain-Driven Design",
			"Repository Pattern",
			"Dependency Injection",
			"Graceful Shutdown",
			"Structured Logging",
			"OpenExchange API Ready",
			"Redis Ready",
			"Kafka Ready",
		},
		"endpoints": map[string]string{
			"health":   "/health",
			"rates":    "/rates?currencies=USD,EUR,GBP",
			"exchange": "/exchange?from=WBTC&to=USDT&amount=1.0",
		},
	}

	c.JSON(http.StatusOK, response)
}
