# Multi-stage Dockerfile for PayBridge AI Go API Service
# Syntax: docker build -t paybridge-api -f services/api/Dockerfile .

# Stage 1: Build binary
FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache ca-certificates git tzdata

# Cache Go modules
COPY services/api/go.mod ./
RUN go mod download || true

# Copy source code
COPY services/api/ ./

# Build static, stripped Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/migrate ./cmd/migrate

# Stage 2: Minimal runtime container
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -g 10001 -S paybridge && \
    adduser -u 10001 -S paybridge -G paybridge

COPY --from=builder /bin/api /usr/local/bin/api
COPY --from=builder /bin/migrate /usr/local/bin/migrate

USER paybridge:paybridge

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/usr/local/bin/api"]
