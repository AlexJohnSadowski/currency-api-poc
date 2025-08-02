package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ajs/currency-api/internal/domain/repositories"
	"github.com/ajs/currency-api/internal/infrastructure/config"
	"github.com/ajs/go-common/logger"
	"github.com/sony/gobreaker"
)

type RatesRepositoryImpl struct {
	config         *config.Config
	httpClient     *http.Client
	logger         logger.Logger
	circuitBreaker *gobreaker.CircuitBreaker
}

type OpenExchangeResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func NewRatesRepositoryImpl(cfg *config.Config, log logger.Logger) repositories.RatesRepository {
	settings := gobreaker.Settings{
		Name:        "openexchange-api",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Info("🔌 Circuit breaker state changed",
				"service", name,
				"from", from.String(),
				"to", to.String(),
			)
		},
	}

	return &RatesRepositoryImpl{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:         log,
		circuitBreaker: gobreaker.NewCircuitBreaker(settings),
	}
}

func (r *RatesRepositoryImpl) GetRates(ctx context.Context, currencies []string) (map[string]float64, string, error) {
	if r.config.OpenExchangeAPIKey == "" {
		info := "🤖 No API key: Using mock rates"
		r.logger.Info(info)
		return r.getMockRates(currencies), info, nil
	}

	result, err := r.circuitBreaker.Execute(func() (interface{}, error) {
		return r.fetchRatesFromAPI(ctx, currencies)
	})

	if err != nil {
		if err == gobreaker.ErrOpenState {
			r.logger.Error("⚡ Circuit breaker is OPEN - external API unavailable", err)
			return nil, "", fmt.Errorf("external rates API is currently unavailable (service protection active)")
		}

		if err == gobreaker.ErrTooManyRequests {
			r.logger.Error("🚦 Circuit breaker limiting requests", err)
			return nil, "", fmt.Errorf("external rates API is being rate limited (too many requests)")
		}

		r.logger.Error("External API failed", err,
			"circuit_state", r.circuitBreaker.State().String(),
		)
		return nil, "", fmt.Errorf("failed to fetch live exchange rates: %w", err)
	}

	rates := result.(map[string]float64)
	info := "🔑 API key provided: Using live rates"
	r.logger.Info("✅ Successfully fetched live rates",
		"currencies", len(currencies),
		"circuit_state", r.circuitBreaker.State().String(),
	)
	return rates, info, nil
}

func (r *RatesRepositoryImpl) fetchRatesFromAPI(ctx context.Context, currencies []string) (map[string]float64, error) {
	currenciesParam := strings.Join(currencies, ",")
	url := fmt.Sprintf("%s/latest.json?app_id=%s&symbols=%s",
		r.config.OpenExchangeBaseURL,
		r.config.OpenExchangeAPIKey,
		currenciesParam,
	)

	r.logger.Debug("🌐 Fetching rates from external API", "currencies", currenciesParam)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var openExchangeResp OpenExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&openExchangeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := make(map[string]float64)

	for _, currency := range currencies {
		if currency == "USD" {
			result["USD"] = 1.0
			break
		}
	}

	for _, currency := range currencies {
		if currency != "USD" {
			if rate, exists := openExchangeResp.Rates[currency]; exists {
				result[currency] = rate
			} else {
				return nil, fmt.Errorf("currency '%s' is not supported by the exchange rates provider", currency)
			}
		}
	}

	return result, nil
}

func (r *RatesRepositoryImpl) getMockRates(currencies []string) map[string]float64 {
	mockRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.73,
		"JPY": 110.0,
		"CAD": 1.25,
		"AUD": 1.35,
		"CHF": 0.92,
		"CNY": 7.2,
		"SEK": 10.5,
		"NOK": 11.2,
	}

	result := make(map[string]float64)
	for _, currency := range currencies {
		if rate, exists := mockRates[currency]; exists {
			result[currency] = rate
		}
		// Skip unknown currencies - they'll be caught by the query handler
	}

	return result
}
