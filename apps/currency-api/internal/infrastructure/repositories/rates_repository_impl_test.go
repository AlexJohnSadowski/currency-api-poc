package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ajs/currency-api/internal/infrastructure/config"
	"github.com/ajs/go-common/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRatesRepositoryImpl_GetRates_MockData(t *testing.T) {
	// Test with no API key (should use mock data)
	cfg := &config.Config{
		OpenExchangeAPIKey:  "",
		OpenExchangeBaseURL: "https://openexchangerates.org/api",
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx := context.Background()
	currencies := []string{"USD", "EUR", "GBP"}

	rates, info, err := repo.GetRates(ctx, currencies)

	require.NoError(t, err)
	assert.Equal(t, "🤖 No API key: Using mock rates", info)

	for _, currency := range currencies {
		assert.Contains(t, rates, currency, "missing rate for currency %s", currency)
	}

	expectedMockRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.73,
	}

	for currency, expectedRate := range expectedMockRates {
		if assert.Contains(t, rates, currency) {
			assert.InDelta(t, expectedRate, rates[currency], 1e-6,
				"currency %s: expected rate %f, got %f", currency, expectedRate, rates[currency])
		}
	}
}

func TestRatesRepositoryImpl_GetRates_MockData_UnknownCurrency(t *testing.T) {
	cfg := &config.Config{
		OpenExchangeAPIKey:  "",
		OpenExchangeBaseURL: "https://openexchangerates.org/api",
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx := context.Background()
	currencies := []string{"USD", "UNKNOWN"}

	rates, info, err := repo.GetRates(ctx, currencies)

	require.NoError(t, err)
	assert.Equal(t, "🤖 No API key: Using mock rates", info)

	// Should have USD but not UNKNOWN
	assert.Contains(t, rates, "USD", "expected USD rate in mock data")
	assert.NotContains(t, rates, "UNKNOWN", "did not expect UNKNOWN currency in mock data")
}

func TestRatesRepositoryImpl_GetRates_WithAPIKey_Success(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-api-key", r.URL.Query().Get("app_id"), "expected correct API key")

		symbols := r.URL.Query().Get("symbols")
		assert.Equal(t, "USD,EUR", symbols, "expected correct symbols parameter")

		response := OpenExchangeResponse{
			Rates: map[string]float64{
				"EUR": 0.85,
				// USD is not included in OpenExchange response as it's the base
			},
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer testServer.Close()

	cfg := &config.Config{
		OpenExchangeAPIKey:  "test-api-key",
		OpenExchangeBaseURL: testServer.URL,
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx := context.Background()
	currencies := []string{"USD", "EUR"}

	rates, info, err := repo.GetRates(ctx, currencies)

	require.NoError(t, err)
	assert.Equal(t, "🔑 API key provided: Using live rates", info)

	expectedRates := map[string]float64{
		"USD": 1.0,  // USD should always be 1.0
		"EUR": 0.85, // From the mock API response
	}

	for currency, expectedRate := range expectedRates {
		if assert.Contains(t, rates, currency, "missing rate for currency %s", currency) {
			assert.InDelta(t, expectedRate, rates[currency], 1e-6,
				"currency %s: expected rate %f, got %f", currency, expectedRate, rates[currency])
		}
	}
}

func TestRatesRepositoryImpl_GetRates_WithAPIKey_UnsupportedCurrency(t *testing.T) {
 	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := OpenExchangeResponse{
			Rates: map[string]float64{
				"EUR": 0.85,
				// INVALID currency not included
			},
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer testServer.Close()

	cfg := &config.Config{
		OpenExchangeAPIKey:  "test-api-key",
		OpenExchangeBaseURL: testServer.URL,
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx := context.Background()
	currencies := []string{"USD", "EUR", "INVALID"}

	_, _, err := repo.GetRates(ctx, currencies)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "currency 'INVALID' is not supported by the exchange rates provider")
}

func TestRatesRepositoryImpl_GetRates_WithAPIKey_APIError(t *testing.T) {
	// Create a test server that returns an error
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	cfg := &config.Config{
		OpenExchangeAPIKey:  "test-api-key",
		OpenExchangeBaseURL: testServer.URL,
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx := context.Background()
	currencies := []string{"USD", "EUR"}

	_, _, err := repo.GetRates(ctx, currencies)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch live exchange rates")
}

func TestRatesRepositoryImpl_GetRates_WithAPIKey_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json"))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	cfg := &config.Config{
		OpenExchangeAPIKey:  "test-api-key",
		OpenExchangeBaseURL: testServer.URL,
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx := context.Background()
	currencies := []string{"USD", "EUR"}

	_, _, err := repo.GetRates(ctx, currencies)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode response")
}

func TestRatesRepositoryImpl_GetRates_ContextCancellation(t *testing.T) {
	// Create a test server with a delay
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		response := OpenExchangeResponse{
			Rates: map[string]float64{"EUR": 0.85},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer testServer.Close()

	cfg := &config.Config{
		OpenExchangeAPIKey:  "test-api-key",
		OpenExchangeBaseURL: testServer.URL,
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	currencies := []string{"USD", "EUR"}

	_, _, err := repo.GetRates(ctx, currencies)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make request")
}

func TestRatesRepositoryImpl_CircuitBreaker(t *testing.T) {
	// This test verifies circuit breaker behavior
	failureCount := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failureCount++
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Simulated failure"))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	cfg := &config.Config{
		OpenExchangeAPIKey:  "test-api-key",
		OpenExchangeBaseURL: testServer.URL,
	}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log)

	ctx := context.Background()
	currencies := []string{"USD", "EUR"}

	var circuitBreakerTriggered bool
	for i := 0; i < 5; i++ {
		_, _, err := repo.GetRates(ctx, currencies)
		require.Error(t, err, "expected error on attempt %d", i+1)

		// After 3 failures, subsequent calls should be circuit breaker errors
		if i >= 3 && (assert.Contains(t, err.Error(), "external rates API is currently unavailable") ||
			assert.Contains(t, err.Error(), "external rates API is being rate limited")) {
			circuitBreakerTriggered = true
			break
		}
	}

	assert.True(t, circuitBreakerTriggered, "circuit breaker should have been triggered")
	assert.LessOrEqual(t, failureCount, 4, "circuit breaker should have limited HTTP requests")
}

func TestRatesRepositoryImpl_GetMockRates(t *testing.T) {
	cfg := &config.Config{}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log).(*RatesRepositoryImpl)

	tests := []struct {
		name             string
		currencies       []string
		expectedLength   int
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:           "known currencies",
			currencies:     []string{"USD", "EUR", "GBP"},
			expectedLength: 3,
			shouldContain:  []string{"USD", "EUR", "GBP"},
		},
		{
			name:             "mixed known and unknown",
			currencies:       []string{"USD", "UNKNOWN", "EUR"},
			expectedLength:   2, // Only USD and EUR should be returned
			shouldContain:    []string{"USD", "EUR"},
			shouldNotContain: []string{"UNKNOWN"},
		},
		{
			name:           "all unknown currencies",
			currencies:     []string{"UNKNOWN1", "UNKNOWN2"},
			expectedLength: 0,
		},
		{
			name:           "empty currencies list",
			currencies:     []string{},
			expectedLength: 0,
		},
		{
			name:           "all supported mock currencies",
			currencies:     []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NOK"},
			expectedLength: 10,
			shouldContain:  []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NOK"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rates := repo.getMockRates(tt.currencies)

			assert.Len(t, rates, tt.expectedLength)

			// Verify expected currencies are present
			for _, currency := range tt.shouldContain {
				assert.Contains(t, rates, currency, "expected currency %s to be present", currency)
				assert.Positive(t, rates[currency], "rate for %s should be positive", currency)
			}

			// Verify unexpected currencies are not present
			for _, currency := range tt.shouldNotContain {
				assert.NotContains(t, rates, currency, "currency %s should not be present", currency)
			}

			// Verify that all returned currencies were requested
			for currency := range rates {
				assert.Contains(t, tt.currencies, currency, "unexpected currency %s in results", currency)
			}
		})
	}
}

func TestRatesRepositoryImpl_GetMockRates_SpecificValues(t *testing.T) {
	cfg := &config.Config{}
	log := logger.New("error")
	repo := NewRatesRepositoryImpl(cfg, log).(*RatesRepositoryImpl)

	// Test specific mock rate values
	currencies := []string{"USD", "EUR", "GBP", "JPY"}
	rates := repo.getMockRates(currencies)

	expectedRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.73,
		"JPY": 110.0,
	}

	for currency, expectedRate := range expectedRates {
		if assert.Contains(t, rates, currency, "missing rate for %s", currency) {
			assert.Equal(t, expectedRate, rates[currency], "incorrect rate for %s", currency)
		}
	}
}
