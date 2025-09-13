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

# Build a statically linked, stripped binary
# -s -w removes symbol/debug info to shrink binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /app/bank-service ./cmd/main.go

# ---------- Stage 2: Runtime ----------
FROM alpine:3.18

LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      app="eBankSEP-BE" \
      description="Bank backend service"

# Install CA certs (needed for TLS connections) and tzdata
RUN apk add --no-cache ca-certificates tzdata \
 && update-ca-certificates || true

# Create a non-root user for security
RUN addgroup -S app && adduser -S -G app -u 10001 app

WORKDIR /app

# copy binary from builder
COPY --from=builder /app/bank-service /app/bank-service

# copy seed scripts into image
COPY --from=builder /src/scripts /app/scripts

# make sure binary and scripts are owned by non-root user
RUN chown -R app:app /app

# run as non-root
USER app

# sensible defaults; Azure will override PORT when binding
ENV APP_ENV=acquirer \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/bank-service"]
