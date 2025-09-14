# ---------- Stage 1: Build ----------
FROM golang:1.23 AS builder

LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      stage="builder"

WORKDIR /src

# copy only module files first for Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# allow building for different architectures if needed
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Build binary (može sa CGO_ENABLED=1, da koristi TLS biblioteke)
RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /app/bank-service ./cmd/main.go

# ---------- Stage 2: Runtime ----------
FROM debian:bookworm-slim

LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      app="eBankSEP-BE" \
      description="Bank backend service"

# Install TLS certificates and tzdata
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata postgresql-client && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r app && useradd -r -g app -u 10001 app

WORKDIR /app

# copy binary from builder
COPY --from=builder /app/bank-service /app/bank-service

# copy seed scripts
COPY --from=builder /src/scripts /app/scripts

# make sure binary and scripts are owned by non-root user
RUN chown -R app:app /app

# run as non-root
USER app

# environment defaults
ENV APP_ENV=acquirer \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/bank-service"]
