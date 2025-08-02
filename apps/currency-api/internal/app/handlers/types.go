package handlers

import "github.com/ajs/currency-api/internal/domain/entities"

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid input provided"`
}

type HealthResponse struct {
	Status      string          `json:"status" example:"healthy"`
	Service     string          `json:"service" example:"currency-exchange-api"`
	Version     string          `json:"version" example:"2.0.0"`
	Timestamp   int64           `json:"timestamp"`
	Environment EnvironmentInfo `json:"environment"`
	Framework   string          `json:"framework" example:"gin-gonic"`
	NxPlugin    string          `json:"nx_plugin" example:"@naxodev/gonx"`
	GoVersion   string          `json:"go_version" example:"1.24"`
	Features    []string        `json:"features"`
	Endpoints   EndpointsInfo   `json:"endpoints"`
}

type EnvironmentInfo struct {
	Mode    string `json:"mode" example:"development"`
	GinMode string `json:"gin_mode" example:"debug"`
	Port    string `json:"port" example:"8080"`
}

type EndpointsInfo struct {
	Health   string `json:"health" example:"/health"`
	Rates    string `json:"rates" example:"/rates?currencies=USD,EUR,GBP"`
	Exchange string `json:"exchange" example:"/exchange?from=WBTC&to=USDT&amount=1.0"`
}

type RatesResponse struct {
	SourceInfo string                  `json:"source_info" example:"🔑 API key provided: Using live rates"`
	Rates      []entities.ExchangeRate `json:"rates"`
}

type RatesErrorResponse struct {
	Error   string `json:"error" example:"currencies parameter is required"`
	Example string `json:"example,omitempty" example:"GET /rates?currencies=USD,EUR,GBP"`
}
