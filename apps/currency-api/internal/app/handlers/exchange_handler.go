package handlers

import (
	"net/http"

	"github.com/ajs/currency-api/internal/app/queries"
	"github.com/ajs/go-common/logger"
	"github.com/gin-gonic/gin"
)

type ExchangeHandler struct {
	queryHandler *queries.ExchangeQueryHandler
	logger       logger.Logger
}

func NewExchangeHandler(queryHandler *queries.ExchangeQueryHandler, logger logger.Logger) *ExchangeHandler {
	return &ExchangeHandler{
		queryHandler: queryHandler,
		logger:       logger,
	}
}

// @Summary Exchange cryptocurrencies
// @Description Convert one cryptocurrency to another using predefined exchange rates
// @Tags Exchange
// @Accept json
// @Produce json
// @Param from query string true "Source cryptocurrency code" Enums(BEER,FLOKI,GATE,USDT,WBTC)
// @Param to query string true "Target cryptocurrency code" Enums(BEER,FLOKI,GATE,USDT,WBTC)
// @Param amount query number true "Amount to exchange" minimum(0.000001)
// @Success 200 {object} entities.ExchangeResult
// @Failure 400 {object} HTTPError
// @Router /api/v1/exchange [get]
func (h *ExchangeHandler) Exchange(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amount := c.Query("amount")

	query := queries.ExchangeQuery{
		From:   from,
		To:     to,
		Amount: amount,
	}

	result, err := h.queryHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Failed to process exchange", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, result)
}
