package entities

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Currency struct {
	Code          string          `json:"code"`
	DecimalPlaces int32           `json:"decimal_places"`
	RateToUSD     decimal.Decimal `json:"rate_to_usd"`
}

type ExchangeRate struct {
	From string          `json:"from"`
	To   string          `json:"to"`
	Rate decimal.Decimal `json:"rate"`
}

type ExchangeResult struct {
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

var CryptoCurrencies = map[string]Currency{
	"BEER": {
		Code:          "BEER",
		DecimalPlaces: 18,
		RateToUSD:     decimal.NewFromFloat(0.00002461),
	},
	"FLOKI": {
		Code:          "FLOKI",
		DecimalPlaces: 18,
		RateToUSD:     decimal.NewFromFloat(0.0001428),
	},
	"GATE": {
		Code:          "GATE",
		DecimalPlaces: 18,
		RateToUSD:     decimal.NewFromFloat(6.87),
	},
	"USDT": {
		Code:          "USDT",
		DecimalPlaces: 6,
		RateToUSD:     decimal.NewFromFloat(0.999),
	},
	"WBTC": {
		Code:          "WBTC",
		DecimalPlaces: 8,
		RateToUSD:     decimal.NewFromFloat(57037.22),
	},
}

func (c Currency) RoundToDecimalPlaces(amount decimal.Decimal) decimal.Decimal {
	return amount.Round(c.DecimalPlaces)
}

func (c Currency) IsValid() bool {
	return c.Code != "" && c.RateToUSD.GreaterThan(decimal.Zero)
}

func GetCurrency(code string) (Currency, error) {
	currency, exists := CryptoCurrencies[code]
	if !exists {
		return Currency{}, fmt.Errorf("currency %s not supported", code)
	}
	return currency, nil
}
