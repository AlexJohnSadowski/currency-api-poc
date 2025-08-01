FROM golang:1.24.5-alpine

# Install dependencies
RUN apk add --no-cache git ca-certificates wget

# Set working directory
WORKDIR /app

# Always install Air and make sure it's in PATH
RUN echo "Installing Air for potential hot reloading..." && \
    go install github.com/air-verse/air@latest && \
    cp /go/bin/air /usr/local/bin/air && \
    chmod +x /usr/local/bin/air

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build production binary (always available as fallback)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Simple entrypoint that checks .env variable
CMD sh -c 'if [ "$ENABLE_HOT_RELOAD" = "true" ]; then echo "🔥 Starting with hot reloading..." && air; else echo "🚀 Starting production mode..." && ./main; fi'