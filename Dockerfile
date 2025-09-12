# Stage 1: Build Go Backend
FROM golang:1.23 AS ebanksep-be-builder
LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      project="eBankSEP-BE" \
      description="Secure bank backend service with TLS, PostgreSQL, and environment-based config"

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bank-service ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest AS ebanksep-be-server
LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      app="eBankSEP-BE" \
      role="Bank Backend Server" \
      version="1.0"

WORKDIR /app

# Copy binary
COPY --from=ebanksep-be-builder /app/bank-service .

# Copy TLS certificates
COPY cert/cert.pem cert/key.pem ./cert/

# Copy environment files
COPY .env.acquirer .env.acquirer
COPY .env.issuer .env.issuer

# Set default environment (can be overridden)
ENV APP_ENV=acquirer

# Expose HTTPS port
EXPOSE 8443

# Run backend
ENTRYPOINT ["./bank-service"]