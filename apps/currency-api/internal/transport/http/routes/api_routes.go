package routes

import (
	"github.com/ajs/currency-api/internal/app/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	r *gin.Engine,
	healthHandler *handlers.HealthHandler,
	ratesHandler *handlers.RatesHandler,
	exchangeHandler *handlers.ExchangeHandler,
) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	r.GET("/health", healthHandler.Health)
	r.HEAD("/health", healthHandler.Health)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/rates", ratesHandler.GetRates)
		v1.GET("/exchange", exchangeHandler.Exchange)
	}
}
