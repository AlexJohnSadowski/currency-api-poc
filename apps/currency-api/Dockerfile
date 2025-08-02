FROM golang:1.24-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum* ./
COPY apps/currency-api/go.* ./apps/currency-api/
COPY libs/go-common/go.* ./libs/go-common/

WORKDIR /app/apps/currency-api
RUN go mod download

WORKDIR /app

COPY . .

CMD ["echo", "Ready to be run with Docker Compose Watch."]
