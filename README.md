# Currency Exchange API

A modern currency exchange POC built with Go, Gin, and clean architecture principles. Features real-time exchange rates from OpenExchangeRates.org and cryptocurrency conversion capabilities.

## 🚀 Features

- **Clean Architecture**: Domain-driven design with CQRS pattern
- **Real-time Exchange Rates**: Integration with OpenExchangeRates.org API
- **API Versioning**: RESTful API with proper versioning (`/api/v1/`)
- **Auto-generated Documentation**: Swagger/OpenAPI 3.0 documentation
- **Circuit Breaker Pattern**: Resilient external API integration
- **Docker Support**: Containerized with hot-reload for development
- **API Gateway**: Traefik reverse proxy with rate limiting and load balancing
- **Health Checks**: Kubernetes-ready health endpoints
- **Graceful Fallbacks**: Mock data for development when API key not provided

## 🏗️ Architecture

This application follows **Clean Architecture** principles with **Domain-Driven Design (DDD)** and **CQRS** patterns:

```
apps/currency-api/
├── cmd/server/              # Application entry point
│   └── main.go             # Swagger annotations & startup logic
├── internal/
│   ├── app/                # Application layer
│   │   ├── handlers/       # HTTP handlers (Presentation layer)
│   │   │   ├── health_handler.go
│   │   │   ├── rates_handler.go
│   │   │   ├── exchange_handler.go
│   │   │   └── types.go    # Response/Request DTOs
│   │   └── queries/        # CQRS Query handlers (Application logic)
│   │       ├── rates_query.go
│   │       └── exchange_query.go
│   ├── domain/             # Domain layer (Business logic)
│   │   ├── entities/       # Domain entities & value objects
│   │   │   └── currency.go
│   │   └── repositories/   # Repository interfaces
│   │       └── rates_repository.go
│   ├── infrastructure/     # Infrastructure layer
│   │   ├── config/         # Configuration management
│   │   │   └── config.go
│   │   └── repositories/   # Repository implementations
│   │       └── rates_repository_impl.go
│   └── transport/http/     # HTTP transport layer
│       ├── middleware/     # HTTP middleware (future)
│       └── routes/         # Route definitions
│           └── api_routes.go
├── docs/                   # Auto-generated Swagger documentation
├── test/                   # Test suites (to be added)
└── Dockerfile             # Container definition
```


## 🛠️ Tech Stack

- **Language**: Go 1.24
- **Framework**: Gin-Gonic
- **Architecture**: Clean Architecture + CQRS + DDD
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker + Docker Compose
- **API Gateway**: Traefik with host-based routing
- **Logging**: Structured logging with slog
- **Build Tool**: Nx with @naxodev/gonx
- **External APIs**: OpenExchangeRates.org

## 🌐 API Gateway & Routing

This application uses **Traefik** as an API gateway with **host-based routing** for modern microservices architecture:

### **Host-Based Routing Strategy**
- **Primary API**: `http://api.localhost` → Currency API service
- **Traefik Dashboard**: `http://traefik.localhost` → Management interface
- **Direct Access**: `http://localhost:8080` → Container direct access (development only)

### **Gateway Features**
- **Load Balancing**: Automatic request distribution
- **Health Checks**: Upstream service monitoring
- **Rate Limiting**: 100 requests/second with burst capability
- **CORS Support**: Cross-origin resource sharing headers
- **Auto-Discovery**: Services automatically registered via Docker labels

### **Service Discovery**
Traefik automatically discovers services through Docker labels:
```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.currency-api.rule=Host(`api.localhost`)"
  - "traefik.http.services.currency-api.loadbalancer.server.port=8080"
```

## 🚀 Quick Start

### 1. Clone and Setup
```bash
git clone <repository-url>
cd currency-exchange-api
npm install
```

### 2. Environment Configuration
```bash
# Copy environment template (NOTE: Add your own OpenExchange AppId!)
cp .env.example .env.development

# Edit configuration
code .env.development
```

### 3. Host Configuration (Optional)
For easier local development, add these entries to your `/etc/hosts` file, cause sometimes it can get funky inside of WSL (if you use one - and you probably should):
```bash
# Add to /etc/hosts
127.0.0.1 api.localhost
127.0.0.1 traefik.localhost
```

### 4. Start Development Server
```bash
# Start with hot reload
npm run dev

# The API will be available at:
# - Via Traefik Gateway: http://api.localhost (recommended)
# - Direct Container: http://localhost:8080 (development only)
# - Swagger UI: http://api.localhost/swagger/index.html
# - Traefik Dashboard: http://traefik.localhost/dashboard/
```

## 🔧 Configuration

### Environment Variables

Create `.env.development` for local development:

```env
# Server Configuration
PORT=8080
GIN_MODE=debug
ENV=development
# External APIs
OPEN_EXCHANGE_API_KEY=your_api_key_here

```

### Production Configuration (`.env.production`)
```env
PORT=8080
GIN_MODE=release
ENV=production
OPEN_EXCHANGE_API_KEY=your_production_api_key
```

### Development Without API Key
If `OPEN_EXCHANGE_API_KEY` is not provided, the API automatically uses mock data for development purposes.

## 📚 API Documentation

### Base URLs
- **Production/Recommended**: `http://api.localhost` (via Traefik Gateway)
- **Development Direct**: `http://localhost:8080` (container direct access - development only)
- **API Version**: All endpoints are prefixed with `/api/v1`

> **📝 Note**: Direct container access (`localhost:8080`) is exposed only for development convenience. In production, all traffic should go through the API gateway for proper load balancing, rate limiting, and monitoring.

### Interactive Documentation
- **Swagger UI**: http://api.localhost/swagger/index.html
- **OpenAPI JSON**: http://api.localhost/swagger/doc.json

## 🔗 API Endpoints

### Health Check
```bash
# Via Traefik Gateway (recommended)
curl -X GET "http://api.localhost/health" \
  -H "accept: application/json"

# Direct container access (development only)
curl -X GET "http://localhost:8080/health" \
  -H "accept: application/json"
```

**Response:**
```json
{
  "status": "healthy",
  "service": "currency-exchange-api",
  "version": "2.0.0",
  "timestamp": 1691234567,
  "environment": {
    "mode": "development",
    "gin_mode": "debug",
    "port": "8080"
  },
  "framework": "gin-gonic",
  "nx_plugin": "@naxodev/gonx",
  "go_version": "1.24",
  "features": [
    "CQRS Pattern",
    "Domain-Driven Design",
    "Repository Pattern",
    "API Versioning",
    "OpenAPI Documentation"
  ],
  "endpoints": {
    "health": "/health",
    "rates": "/api/v1/rates?currencies=USD,EUR,GBP",
    "exchange": "/api/v1/exchange?from=WBTC&to=USDT&amount=1.0"
  }
}
```

### Exchange Rates

#### Get Currency Exchange Rates
```bash
# Via Traefik Gateway (recommended)
curl -X GET "http://api.localhost/api/v1/rates?currencies=USD,EUR,GBP" \
  -H "accept: application/json"

# Multiple currencies
curl -X GET "http://api.localhost/api/v1/rates?currencies=USD,EUR,GBP,JPY,CAD" \
  -H "accept: application/json"

# Minimum 2 currencies
curl -X GET "http://api.localhost/api/v1/rates?currencies=EUR,GBP" \
  -H "accept: application/json"

# Direct container access (development only)
curl -X GET "http://localhost:8080/api/v1/rates?currencies=USD,EUR,GBP" \
  -H "accept: application/json"
```

**Successful Response:**
```json
{
  "source_info": "🔑 API key provided: Using live rates",
  "rates": [
    {"from": "USD", "to": "EUR", "rate": 0.86255},
    {"from": "USD", "to": "GBP", "rate": 0.752955},
    {"from": "EUR", "to": "USD", "rate": 1.1593530809808126},
    {"from": "EUR", "to": "GBP", "rate": 0.8729406990899078},
    {"from": "GBP", "to": "USD", "rate": 1.3281006169027365},
    {"from": "GBP", "to": "EUR", "rate": 1.1455531871094553}
  ]
}
```

#### Error Cases
```bash
# Missing currencies parameter
curl -X GET "http://api.localhost/api/v1/rates" \
  -H "accept: application/json"

# Empty currencies parameter  
curl -X GET "http://api.localhost/api/v1/rates?currencies=" \
  -H "accept: application/json"

# Only one currency (minimum 2 required)
curl -X GET "http://api.localhost/api/v1/rates?currencies=USD" \
  -H "accept: application/json"

# Invalid currency code
curl -X GET "http://api.localhost/api/v1/rates?currencies=EUR,XYZ" \
  -H "accept: application/json"
```

**Error Response:**
```json
{
  "error": "Failed to retrieve exchange rates. Ensure currency codes are valid."
}
```

### Cryptocurrency Exchange

#### Convert Cryptocurrencies
```bash
# Via Traefik Gateway (recommended)
# WBTC to USDT
curl -X GET "http://api.localhost/api/v1/exchange?from=WBTC&to=USDT&amount=1.0" \
  -H "accept: application/json"

# USDT to BEER (large number result)
curl -X GET "http://api.localhost/api/v1/exchange?from=USDT&to=BEER&amount=1.0" \
  -H "accept: application/json"

# GATE to FLOKI
curl -X GET "http://api.localhost/api/v1/exchange?from=GATE&to=FLOKI&amount=0.5" \
  -H "accept: application/json"

# Small amount exchange
curl -X GET "http://api.localhost/api/v1/exchange?from=WBTC&to=GATE&amount=0.001" \
  -H "accept: application/json"

# Direct container access (development only)
curl -X GET "http://localhost:8080/api/v1/exchange?from=WBTC&to=USDT&amount=1.0" \
  -H "accept: application/json"
```

**Successful Response:**
```json
{
  "from": "WBTC",
  "to": "USDT",
  "amount": 57094.314314
}
```

#### Supported Cryptocurrencies (Mock Values)
| Symbol | Name | Decimal Places | Rate (to USD) |
|--------|------|----------------|---------------|
| BEER | BEER Token | 18 | $0.00002461 |
| FLOKI | FLOKI | 18 | $0.0001428 |
| GATE | Gate Token | 18 | $6.87 |
| USDT | Tether | 6 | $0.999 |
| WBTC | Wrapped Bitcoin | 8 | $57,037.22 |

#### Error Cases
```bash
# Missing 'from' parameter
curl -X GET "http://api.localhost/api/v1/exchange?to=USDT&amount=1.0" \
  -H "accept: application/json"

# Missing 'to' parameter
curl -X GET "http://api.localhost/api/v1/exchange?from=WBTC&amount=1.0" \
  -H "accept: application/json"

# Missing 'amount' parameter
curl -X GET "http://api.localhost/api/v1/exchange?from=WBTC&to=USDT" \
  -H "accept: application/json"

# Invalid cryptocurrency
curl -X GET "http://api.localhost/api/v1/exchange?from=INVALID&to=USDT&amount=1.0" \
  -H "accept: application/json"

# Invalid amount
curl -X GET "http://api.localhost/api/v1/exchange?from=WBTC&to=USDT&amount=invalid" \
  -H "accept: application/json"

# Negative amount
curl -X GET "http://api.localhost/api/v1/exchange?from=WBTC&to=USDT&amount=-1.0" \
  -H "accept: application/json"
```

**Error Response:**
```json
{
  "error": "from, to, and amount parameters are required"
}
```

## 🧪 Testing Your API Gateway

### Quick Test Suite
```bash
echo "🧪 Testing Traefik API Gateway..."
echo "1. Health check:"
curl -s http://api.localhost/health | jq .status

echo -e "\n2. Exchange rates:"
curl -s http://api.localhost/api/v1/rates?currencies=USD,EUR | jq .rates[0]

echo -e "\n3. Crypto exchange:"
curl -s http://api.localhost/api/v1/exchange?from=WBTC&to=USDT&amount=1.0 | jq .

echo -e "\n4. CORS headers:"
curl -s -I http://api.localhost/api/v1/rates?currencies=USD,EUR | grep -i "access-control"

echo -e "\n✅ All tests completed!"
```

### CORS Testing
```bash
# Test CORS preflight request
curl -i -X OPTIONS http://api.localhost/api/v1/rates \
  -H 'Origin: http://localhost:3000' \
  -H 'Access-Control-Request-Method: GET' \
  -H 'Access-Control-Request-Headers: Content-Type'
```

### Rate Limiting Tests
```bash
# Test rate limiting (100 requests/second limit)
for i in {1..10}; do 
  curl -s -o /dev/null -w "%{http_code} " http://api.localhost/api/v1/rates?currencies=USD,EUR
done; echo
```

### Traefik Dashboard & Monitoring
```bash
# Traefik dashboard
curl http://traefik.localhost/dashboard/

# Traefik API endpoints
curl http://localhost:8090/api/rawdata

# Check Traefik health
curl http://localhost:8090/ping

# View active routes
curl http://localhost:8090/api/http/routers | jq .
```

## 🏗️ Development Commands

### Using Nx (Recommended)
```bash
# EXAMPLE (run from the root using nx): Generate Swagger documentation
nx docs:swagger currency-api
```

### Using npm scripts
```bash
# Development with docker watch
npm run dev

# Production mode
npm run prod

# Stop services
npm run stop

# Health check (via gateway)
npm run health
```

## 📊 Monitoring & Observability

### Health Monitoring
- **Endpoint**: `/health` (accessible via gateway and direct)
- **Docker**: Built-in healthcheck every 30s
- **Kubernetes**: Ready for readiness and liveness probes

### Logging
- **Format**: Structured JSON logging via Go's slog
- **Levels**: DEBUG, INFO, WARN, ERROR
- **Context**: Request tracing, error details, performance metrics


## 🔌 Circuit Breaker Testing

The API includes a circuit breaker that protects against external API failures. You can test it without restarting the application:

### Trigger Circuit Breaker Failures
```bash
# Test via gateway (recommended)
# Failure 1
curl "http://api.localhost/api/v1/rates?currencies=USD,INVALID_CURRENCY"

# Failure 2
curl "http://api.localhost/api/v1/rates?currencies=USD,INVALID_CURRENCY"

# Failure 3 (Circuit opens!)
curl "http://api.localhost/api/v1/rates?currencies=USD,INVALID_CURRENCY"

# Fast fail (Circuit is OPEN)
curl "http://api.localhost/api/v1/rates?currencies=USD,EUR"
```

### Test Circuit Breaker Recovery
```bash
# Wait for half-open state
sleep 30

# Test recovery with valid currencies
curl "http://api.localhost/api/v1/rates?currencies=USD,EUR"
```

### Monitor Circuit Breaker Activity
```bash
# Watch circuit breaker logs in real-time
npm run logs | grep -E "(🔌|⚡|circuit|Circuit)"
```

### Expected Behavior
- **Failures 1-3**: API errors with external service failures
- **Failure 4+**: Fast circuit breaker errors (no API calls made)
- **After 30 seconds**: Half-open state - tests recovery automatically
- **Recovery**: If valid API call succeeds, circuit closes

### Circuit Breaker States
- **CLOSED**: Normal operation, requests pass through
- **OPEN**: Circuit breaker active, requests fail fast
- **HALF-OPEN**: Testing recovery, limited requests allowed

## 🐳 Docker Services

The application runs as part of a multi-service Docker Compose setup:

- **currency-api**: Main API service (port 8080)
- **redis**: Caching layer (port 6379)  
- **traefik**: API gateway with load balancing and rate limiting (ports 80, 8090)

### Service URLs
- **API via Gateway**: http://api.localhost (recommended)
- **API Direct**: http://localhost:8080 (development only)
- **Traefik Dashboard**: http://traefik.localhost/dashboard/
- **Redis**: localhost:6379

## 🛡️ Security Features

- **Rate Limiting**: Request rate limiting per IP via Traefik (100/second)
- **CORS**: Configurable cross-origin resource sharing
- **Security Headers**: X-Content-Type-Options, X-Frame-Options, etc.
- **Input Validation**: Comprehensive request parameter validation
- **Timeout Protection**: Request timeout middleware
- **Error Handling**: Sanitized error responses (no sensitive data exposure)
- **Gateway Protection**: All traffic routed through secure gateway

## 📝 License

GNU GENERAL PUBLIC LICENSE
Version 3, 29 June 2007

Copyright (C) 2025 AlexJohnSadowski

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

---

**Built with ❤️ using Go, Gin, Clean Architecture, Traefik Gateway, and modern development practices.**