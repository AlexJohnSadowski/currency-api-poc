package queries

import (
	"context"
	"fmt"
	"strings"

	"github.com/ajs/currency-api/internal/domain/entities"
	"github.com/shopspring/decimal"
)

type ExchangeQuery struct {
	From   string
	To     string
	Amount string
}

type ExchangeQueryHandler struct{}

func NewExchangeQueryHandler() *ExchangeQueryHandler {
	return &ExchangeQueryHandler{}
}

func (h *ExchangeQueryHandler) Handle(ctx context.Context, query ExchangeQuery) (*entities.ExchangeResult, error) {
	from := strings.ToUpper(strings.TrimSpace(query.From))
	to := strings.ToUpper(strings.TrimSpace(query.To))

	if from == "" || to == "" || query.Amount == "" {
		return nil, fmt.Errorf("from, to, and amount parameters are required")
	}

	amount, err := decimal.NewFromString(query.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("amount must be positive")
	}

	fromCurrency, err := entities.GetCurrency(from)
	if err != nil {
		return nil, fmt.Errorf("unsupported currency %s", from)
	}

	toCurrency, err := entities.GetCurrency(to)
	if err != nil {
		return nil, fmt.Errorf("unsupported currency %s", to)
	}

	usdAmount := amount.Mul(fromCurrency.RateToUSD)
	resultAmount := usdAmount.Div(toCurrency.RateToUSD)

	finalAmount := toCurrency.RoundToDecimalPlaces(resultAmount)

	return &entities.ExchangeResult{
		From:   from,
		To:     to,
		Amount: finalAmount,
	}, nil
}
