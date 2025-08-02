package queries

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchangeQueryHandler_Handle_WithDecimal(t *testing.T) {
	handler := NewExchangeQueryHandler()
	ctx := context.Background()

	tests := []struct {
		name           string
		query          ExchangeQuery
		expectedAmount string
		expectedError  string
	}{
		{
			name: "successful WBTC to USDT exchange",
			query: ExchangeQuery{
				From:   "WBTC",
				To:     "USDT",
				Amount: "1.0",
			},
			expectedAmount: "57094.314314",
		},
		{
			name: "successful USDT to BEER exchange",
			query: ExchangeQuery{
				From:   "USDT",
				To:     "BEER",
				Amount: "1.0",
			},
			expectedAmount: "40593.2547744819179195",
		},
		{
			name: "successful BEER to FLOKI exchange",
			query: ExchangeQuery{
				From:   "BEER",
				To:     "FLOKI",
				Amount: "1000.0",
			},
			expectedAmount: "172.3389355742296919",
		},
		{
			name: "GATE to WBTC exchange",
			query: ExchangeQuery{
				From:   "GATE",
				To:     "WBTC",
				Amount: "100.0",
			},
			expectedAmount: "0.01204477",
		},
		{
			name: "very precise small amount",
			query: ExchangeQuery{
				From:   "BEER",
				To:     "WBTC",
				Amount: "100000.0",
			},
			expectedAmount: "0.00004315",
		},
		{
			name: "same currency exchange",
			query: ExchangeQuery{
				From:   "USDT",
				To:     "USDT",
				Amount: "100.0",
			},
			expectedAmount: "100.000000",
		},
		{
			name: "missing from parameter",
			query: ExchangeQuery{
				From:   "",
				To:     "USDT",
				Amount: "1.0",
			},
			expectedError: "from, to, and amount parameters are required",
		},
		{
			name: "invalid amount format",
			query: ExchangeQuery{
				From:   "WBTC",
				To:     "USDT",
				Amount: "not-a-number",
			},
			expectedError: "invalid amount",
		},
		{
			name: "negative amount",
			query: ExchangeQuery{
				From:   "WBTC",
				To:     "USDT",
				Amount: "-1.0",
			},
			expectedError: "amount must be positive",
		},
		{
			name: "unsupported currency",
			query: ExchangeQuery{
				From:   "MATIC",
				To:     "USDT",
				Amount: "1.0",
			},
			expectedError: "unsupported currency MATIC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handler.Handle(ctx, tt.query)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.query.From, result.From)
			assert.Equal(t, tt.query.To, result.To)

			expectedAmount, err := decimal.NewFromString(tt.expectedAmount)
			require.NoError(t, err)

			assert.True(t, expectedAmount.Equal(result.Amount),
				"Exchange %s->%s: expected %s, got %s",
				result.From, result.To,
				expectedAmount.String(), result.Amount.String())
		})
	}
}

func TestExchangeQueryHandler_AllCryptoPairs_WithDecimal(t *testing.T) {
	handler := NewExchangeQueryHandler()
	ctx := context.Background()
	cryptos := []string{"BEER", "FLOKI", "GATE", "USDT", "WBTC"}

	for _, from := range cryptos {
		for _, to := range cryptos {
			t.Run(from+"_to_"+to, func(t *testing.T) {
				testAmount := "10.0"
				query := ExchangeQuery{
					From:   from,
					To:     to,
					Amount: testAmount,
				}

				result, err := handler.Handle(ctx, query)
				require.NoError(t, err, "Exchange from %s to %s should succeed", from, to)
				require.NotNil(t, result)
				assert.Equal(t, from, result.From)
				assert.Equal(t, to, result.To)

				assert.True(t, result.Amount.GreaterThanOrEqual(decimal.Zero),
					"Amount should be positive or zero for %s->%s: got %s", from, to, result.Amount.String())

				if from == to {
					expectedAmount, err := decimal.NewFromString("10.0")
					require.NoError(t, err)
					assert.True(t, expectedAmount.Equal(result.Amount),
						"Same currency exchange should return same amount: expected %s, got %s",
						expectedAmount.String(), result.Amount.String())
				}
			})
		}
	}
}
