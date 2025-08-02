package handlers

import (
	"net/http"
	"strings"

	"github.com/ajs/currency-api/internal/app/queries"
	"github.com/ajs/go-common/logger"
	"github.com/gin-gonic/gin"
)

type RatesHandler struct {
	queryHandler *queries.GetRatesQueryHandler
	logger       logger.Logger
}

func NewRatesHandler(queryHandler *queries.GetRatesQueryHandler, logger logger.Logger) *RatesHandler {
	return &RatesHandler{
		queryHandler: queryHandler,
		logger:       logger,
	}
}

// @Summary		Get exchange rates
// @Description	Get exchange rates for a list of currencies (minimum 2 required)
// @Tags			Rates
// @Accept			json
// @Produce		json
// @Param			currencies	query		string	true	"Comma-separated list of currency codes (e.g., USD,EUR,GBP)"
// @Success		200			{object}	RatesResponse
// @Failure		400			{object}	RatesErrorResponse
// @Router			/api/v1/rates [get]
func (h *RatesHandler) GetRates(c *gin.Context) {
	currenciesParam := c.Query("currencies")

	if currenciesParam == "" {
		c.JSON(http.StatusBadRequest, RatesErrorResponse{
			Error:   "currencies parameter is required",
			Example: "GET /rates?currencies=USD,EUR,GBP",
		})
		return
	}

	currencies := strings.Split(currenciesParam, ",")

	query := queries.GetRatesQuery{
		Currencies: currencies,
	}

	rates, info, err := h.queryHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Failed to get rates", err)
		c.JSON(http.StatusBadRequest, RatesErrorResponse{
			Error: "Failed to retrieve exchange rates. Ensure currency codes are valid.",
		})
		return
	}

	response := RatesResponse{
		SourceInfo: info,
		Rates:      rates,
	}

	c.JSON(http.StatusOK, response)
}
