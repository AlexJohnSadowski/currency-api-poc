package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	app := createApp()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		ginMode := getEnv("GIN_MODE", "debug")
		hotReload := getEnv("ENABLE_HOT_RELOAD", "false")
		envMode := detectEnvironmentMode()
		
		log.Printf("🚀 Currency Exchange API")
		log.Printf("🔧 Powered by @naxodev/gonx + Gin")
		log.Printf("🧠 Go 1.24")
		log.Printf("🌐 Starting on port %s", port)
		log.Printf("🔧 Environment: %s", envMode)
		log.Printf("🎯 Gin Mode: %s", ginMode)
		log.Printf("🔥 Hot Reload: %s", hotReload)
		log.Printf("")
		log.Printf("📍 Endpoints:")
		log.Printf("   🏥 Health:   http://localhost:%s/health", port)
		log.Printf("   💱 Rates:    http://localhost:%s/rates?currencies=USD,EUR,GBP", port)
		log.Printf("   🪙 Exchange: http://localhost:%s/exchange?from=WBTC&to=USDT&amount=1.0", port)
		log.Printf("")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	<-ctx.Done()
	stop()

	log.Println("🛑 Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("❌ Forced shutdown: %v", err)
	}

	log.Println("✅ Server stopped")
}

func createApp() *gin.Engine {
	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		statusEmoji := "✅"
		if param.StatusCode >= 400 {
			statusEmoji = "❌"
		}
		return fmt.Sprintf("%s [%s] %s %s %d %s\n",
			statusEmoji,
			param.TimeStamp.Format("15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	r.Use(gin.Recovery())

	r.Use(corsMiddleware())

	setupRoutes(r)

	return r
}

func setupRoutes(r *gin.Engine) {
	healthHandler := func(c *gin.Context) {
		envMode := detectEnvironmentMode()
		ginMode := getEnv("GIN_MODE", "debug")
		hotReload := getEnv("ENABLE_HOT_RELOAD", "false")
		
		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"service":     "currency-exchange-api",
			"version":     "1.0.0",
			"timestamp":   time.Now().Unix(),
			"environment": map[string]interface{}{
				"mode":        envMode,
				"gin_mode":    ginMode,
				"hot_reload":  hotReload == "true",
				"port":        getEnv("PORT", "8080"),
			},
			"framework":   "gin-gonic",
			"nx_plugin":   "@naxodev/gonx",
			"go_version":  "1.24",
			"ready_for": []string{
				"CQRS Implementation",
				"Redis Caching",
				"OpenExchange API",
				"Crypto Conversion",
			},
		})
	}
	
	r.GET("/health", healthHandler)
	r.HEAD("/health", healthHandler)

	api := r.Group("/")
	{
		api.GET("/rates", handleRates)
		api.GET("/exchange", handleExchange)
	}
}

func handleRates(c *gin.Context) {
	currencies := c.Query("currencies")

	if currencies == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "currencies parameter is required",
			"example": "GET /rates?currencies=USD,EUR,GBP",
		})
		return
	}
}

func handleExchange(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amount := c.Query("amount")

	if from == "" || to == "" || amount == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "from, to, and amount parameters are required",
			"example": "GET /exchange?from=WBTC&to=USDT&amount=1.0",
		})
		return
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, HEAD")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func detectEnvironmentMode() string {
	ginMode := getEnv("GIN_MODE", "debug")
	hotReload := getEnv("ENABLE_HOT_RELOAD", "false")
	
	if hotReload == "true" {
		return "development"
	}
	
	switch ginMode {
	case "release":
		return "production"
	case "test":
		return "testing"
	default:
		return "development"
	}
}