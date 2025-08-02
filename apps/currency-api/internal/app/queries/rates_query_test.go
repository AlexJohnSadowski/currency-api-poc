package queries

import (
	"context"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestRatesRepository struct {
	rates map[string]float64
	info  string
	err   error
}

func NewTestRatesRepository() *TestRatesRepository {
	return &TestRatesRepository{
		rates: make(map[string]float64),
		info:  "test repository",
	}
}

func (r *TestRatesRepository) SetRates(rates map[string]float64) {
	r.rates = rates
}

func (r *TestRatesRepository) SetError(err error) {
	r.err = err
}

func (r *TestRatesRepository) SetInfo(info string) {
	r.info = info
}

func (r *TestRatesRepository) GetRates(ctx context.Context, currencies []string) (map[string]float64, string, error) {
	if r.err != nil {
		return nil, "", r.err
	}

	result := make(map[string]float64)
	for _, currency := range currencies {
		if rate, exists := r.rates[currency]; exists {
			result[currency] = rate
		}
	}

	return result, r.info, nil
}

func TestGetRatesQueryHandler_Handle_WithDecimal(t *testing.T) {
	tests := []struct {
		name          string
		query         GetRatesQuery
		repoRates     map[string]float64
		repoInfo      string
		repoError     error
		expectedRates []struct {
			from string
			to   string
			rate string
		}
		expectedInfo  string
		expectedError string
	}{
		{
			name: "successful USD EUR GBP rates",
			query: GetRatesQuery{
				Currencies: []string{"USD", "EUR", "GBP"},
			},
			repoRates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
				"GBP": 0.73,
			},
			repoInfo: "🔑 API key provided: Using live rates",
			expectedRates: []struct {
				from string
				to   string
				rate string
			}{
				{"USD", "EUR", "0.85"},
				{"USD", "GBP", "0.73"},
				{"EUR", "USD", "1.1764705882352941"},
				{"EUR", "GBP", "0.8588235294117647"},
				{"GBP", "USD", "1.3698630136986301"},
				{"GBP", "EUR", "1.1643835616438356"},
			},
			expectedInfo: "🔑 API key provided: Using live rates",
		},
		{
			name: "successful two currency pair",
			query: GetRatesQuery{
				Currencies: []string{"USD", "EUR"},
			},
			repoRates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
			},
			repoInfo: "🤖 No API key: Using mock rates",
			expectedRates: []struct {
				from string
				to   string
				rate string
			}{
				{"USD", "EUR", "0.85"},
				{"EUR", "USD", "1.1764705882352941"},
			},
			expectedInfo: "🤖 No API key: Using mock rates",
		},
		{
			name: "case insensitive currency handling",
			query: GetRatesQuery{
				Currencies: []string{"usd", "eur"},
			},
			repoRates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
			},
			repoInfo: "test rates",
			expectedRates: []struct {
				from string
				to   string
				rate string
			}{
				{"USD", "EUR", "0.85"},
				{"EUR", "USD", "1.1764705882352941"},
			},
			expectedInfo: "test rates",
		},
		// Error cases
		{
			name: "insufficient currencies - one currency",
			query: GetRatesQuery{
				Currencies: []string{"USD"},
			},
			expectedError: "at least two currencies are required",
		},
		{
			name: "insufficient currencies - empty list",
			query: GetRatesQuery{
				Currencies: []string{},
			},
			expectedError: "at least two currencies are required",
		},
		{
			name: "repository error",
			query: GetRatesQuery{
				Currencies: []string{"USD", "EUR"},
			},
			repoError:     fmt.Errorf("external API unavailable"),
			expectedError: "failed to get rates",
		},
		{
			name: "unsupported currency",
			query: GetRatesQuery{
				Currencies: []string{"USD", "INVALID"},
			},
			repoRates: map[string]float64{
				"USD": 1.0,
				// INVALID currency not provided
			},
			expectedError: "currency 'INVALID' is not supported or not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTestRatesRepository()
			if tt.repoRates != nil {
				repo.SetRates(tt.repoRates)
			}
			if tt.repoInfo != "" {
				repo.SetInfo(tt.repoInfo)
			}
			if tt.repoError != nil {
				repo.SetError(tt.repoError)
			}

			handler := NewGetRatesQueryHandler(repo)
			ctx := context.Background()

			rates, info, err := handler.Handle(ctx, tt.query)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedInfo, info)
			assert.Len(t, rates, len(tt.expectedRates))

			rateMap := make(map[string]decimal.Decimal)
			for _, rate := range rates {
				key := fmt.Sprintf("%s-%s", rate.From, rate.To)
				rateMap[key] = rate.Rate
			}

			for _, expectedRate := range tt.expectedRates {
				key := fmt.Sprintf("%s-%s", expectedRate.from, expectedRate.to)
				actualRate, exists := rateMap[key]

				assert.True(t, exists, "missing exchange rate from %s to %s", expectedRate.from, expectedRate.to)
				if exists {
					expectedDecimal, err := decimal.NewFromString(expectedRate.rate)
					require.NoError(t, err)

					assert.True(t, expectedDecimal.Equal(actualRate),
						"rate from %s to %s: expected %s, got %s",
						expectedRate.from, expectedRate.to, expectedDecimal.String(), actualRate.String())
				}
			}
		})
	}
}

func TestGetRatesQueryHandler_CalculateRate_WithDecimal(t *testing.T) {
	handler := &GetRatesQueryHandler{}

	tests := []struct {
		name          string
		rates         map[string]float64
		from          string
		to            string
		expectedRate  string
		expectedError string
	}{
		{
			name: "USD to EUR",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
			},
			from:         "USD",
			to:           "EUR",
			expectedRate: "0.85",
		},
		{
			name: "EUR to USD",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
			},
			from:         "EUR",
			to:           "USD",
			expectedRate: "1.1764705882352941",
		},
		{
			name: "EUR to GBP (cross-rate calculation)",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
				"GBP": 0.73,
			},
			from:         "EUR",
			to:           "GBP",
			expectedRate: "0.8588235294117647",
		},
		{
			name: "missing from currency",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
			},
			from:          "GBP",
			to:            "USD",
			expectedError: "rate not available for currency GBP",
		},
		{
			name: "missing to currency",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
			},
			from:          "USD",
			to:            "GBP",
			expectedError: "rate not available for currency GBP",
		},
		{
			name: "zero from rate",
			rates: map[string]float64{
				"USD": 0.0,
				"EUR": 0.85,
			},
			from:          "USD",
			to:            "EUR",
			expectedError: "invalid rate",
		},
		{
			name: "zero to rate",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.0,
			},
			from:          "USD",
			to:            "EUR",
			expectedError: "invalid rate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, err := handler.calculateRate(tt.rates, tt.from, tt.to)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)

			expectedDecimal, err := decimal.NewFromString(tt.expectedRate)
			require.NoError(t, err)

			assert.True(t, expectedDecimal.Equal(rate),
				"expected rate %s, got %s", expectedDecimal.String(), rate.String())
		})
	}
}
