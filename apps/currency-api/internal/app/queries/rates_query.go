package queries

import (
	"context"
	"fmt"
	"strings"

	"github.com/ajs/currency-api/internal/domain/entities"
	"github.com/ajs/currency-api/internal/domain/repositories"
	"github.com/shopspring/decimal"
)

type GetRatesQuery struct {
	Currencies []string
}

type GetRatesQueryHandler struct {
	ratesRepo repositories.RatesRepository
}

func NewGetRatesQueryHandler(ratesRepo repositories.RatesRepository) *GetRatesQueryHandler {
	return &GetRatesQueryHandler{ratesRepo: ratesRepo}
}

func (h *GetRatesQueryHandler) Handle(ctx context.Context, query GetRatesQuery) ([]entities.ExchangeRate, string, error) {
	if len(query.Currencies) < 2 {
		return nil, "", fmt.Errorf("at least two currencies are required")
	}

	currencies := make([]string, len(query.Currencies))
	for i, currency := range query.Currencies {
		currencies[i] = strings.ToUpper(strings.TrimSpace(currency))
	}

	rates, info, err := h.ratesRepo.GetRates(ctx, currencies)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get rates: %w", err)
	}

	for _, currency := range currencies {
		if _, exists := rates[currency]; !exists {
			return nil, "", fmt.Errorf("currency '%s' is not supported or not available", currency)
		}
	}

	result := make([]entities.ExchangeRate, 0, len(currencies)*(len(currencies)-1))

	for _, from := range currencies {
		for _, to := range currencies {
			if from != to {
				rate, err := h.calculateRate(rates, from, to)
				if err != nil {
					return nil, "", fmt.Errorf("failed to calculate rate from %s to %s: %w", from, to, err)
				}

				result = append(result, entities.ExchangeRate{
					From: from,
					To:   to,
					Rate: rate,
				})
			}
		}
	}

	return result, info, nil
}

func (h *GetRatesQueryHandler) calculateRate(rates map[string]float64, from, to string) (decimal.Decimal, error) {
	fromRate, fromExists := rates[from]
	toRate, toExists := rates[to]

	if !fromExists {
		return decimal.Zero, fmt.Errorf("rate not available for currency %s", from)
	}

	if !toExists {
		return decimal.Zero, fmt.Errorf("rate not available for currency %s", to)
	}

	if fromRate == 0 || toRate == 0 {
		return decimal.Zero, fmt.Errorf("invalid rate: %s=%.6f, %s=%.6f", from, fromRate, to, toRate)
	}

	fromDecimal := decimal.NewFromFloat(fromRate)
	toDecimal := decimal.NewFromFloat(toRate)

	rate := toDecimal.Div(fromDecimal)

	return rate, nil
}
