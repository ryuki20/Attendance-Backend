# Build stage
FROM golang:1.23-alpine AS builder

# Install air for hot reload
RUN go install github.com/air-verse/air@v1.52.3

# Set working directory
WORKDIR /app

# Development stage
FROM golang:1.23-alpine

# Install necessary packages
RUN apk add --no-cache git curl

# Install air (固定バージョン)
RUN go install github.com/air-verse/air@v1.52.3

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Create tmp directory
RUN mkdir -p /app/tmp

# Run air for hot reload
CMD ["air", "-c", ".air.toml"]
