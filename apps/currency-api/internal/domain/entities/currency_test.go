package entities

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrency_RoundToDecimalPlaces_WithDecimal(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		amount   string
		expected string
	}{
		{
			name: "USDT with 6 decimal places",
			currency: Currency{
				Code:          "USDT",
				DecimalPlaces: 6,
				RateToUSD:     decimal.NewFromFloat(0.999),
			},
			amount:   "57094.314314159",
			expected: "57094.314314",
		},
		{
			name: "WBTC with 8 decimal places",
			currency: Currency{
				Code:          "WBTC",
				DecimalPlaces: 8,
				RateToUSD:     decimal.NewFromFloat(57037.22),
			},
			amount:   "1.123456789",
			expected: "1.12345679",
		},
		{
			name: "BEER with 18 decimal places",
			currency: Currency{
				Code:          "BEER",
				DecimalPlaces: 18,
				RateToUSD:     decimal.NewFromFloat(0.00002461),
			},
			amount:   "40593.254769230769230769999",
			expected: "40593.254769230769230770",
		},
		{
			name: "exact precision maintained",
			currency: Currency{
				Code:          "USDT",
				DecimalPlaces: 6,
				RateToUSD:     decimal.NewFromFloat(0.999),
			},
			amount:   "100.0",
			expected: "100.000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount := decimal.RequireFromString(tt.amount)
			expected := decimal.RequireFromString(tt.expected)
			result := tt.currency.RoundToDecimalPlaces(amount)
			assert.True(t, expected.Equal(result),
				"RoundToDecimalPlaces() = %s, want %s", result.String(), expected.String())
		})
	}
}

func TestCurrency_IsValid_WithDecimal(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		expected bool
	}{
		{
			name: "valid currency",
			currency: Currency{
				Code:          "USDT",
				DecimalPlaces: 6,
				RateToUSD:     decimal.NewFromFloat(0.999),
			},
			expected: true,
		},
		{
			name: "empty code",
			currency: Currency{
				Code:          "",
				DecimalPlaces: 6,
				RateToUSD:     decimal.NewFromFloat(0.999),
			},
			expected: false,
		},
		{
			name: "zero rate",
			currency: Currency{
				Code:          "USDT",
				DecimalPlaces: 6,
				RateToUSD:     decimal.Zero,
			},
			expected: false,
		},
		{
			name: "negative rate",
			currency: Currency{
				Code:          "USDT",
				DecimalPlaces: 6,
				RateToUSD:     decimal.NewFromInt(-1),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.currency.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCryptoCurrencies_DecimalPrecision(t *testing.T) {
	expectedCurrencies := []string{"BEER", "FLOKI", "GATE", "USDT", "WBTC"}
	assert.Len(t, CryptoCurrencies, len(expectedCurrencies))

	for _, code := range expectedCurrencies {
		t.Run("currency_"+code, func(t *testing.T) {
			currency, exists := CryptoCurrencies[code]
			require.True(t, exists, "currency %s not found", code)
			assert.Equal(t, code, currency.Code)
			assert.True(t, currency.IsValid(), "currency %s should be valid", code)
			assert.True(t, currency.RateToUSD.GreaterThan(decimal.Zero),
				"currency %s should have positive rate", code)
		})
	}
}

func TestCryptoCurrencies_ExactRates(t *testing.T) {
	tests := []struct {
		code          string
		decimalPlaces int32
		rateToUSD     string
	}{
		{"BEER", 18, "0.00002461"},
		{"FLOKI", 18, "0.0001428"},
		{"GATE", 18, "6.87"},
		{"USDT", 6, "0.999"},
		{"WBTC", 8, "57037.22"},
	}

	for _, tt := range tests {
		t.Run(tt.code+"_exact_values", func(t *testing.T) {
			currency, exists := CryptoCurrencies[tt.code]
			require.True(t, exists, "currency %s not found", tt.code)

			assert.Equal(t, tt.decimalPlaces, currency.DecimalPlaces)

			expectedRate := decimal.RequireFromString(tt.rateToUSD)
			assert.True(t, expectedRate.Equal(currency.RateToUSD),
				"currency %s: expected rate %s, got %s",
				tt.code, expectedRate.String(), currency.RateToUSD.String())
		})
	}
}
