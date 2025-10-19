# Build stage - use stable Go version with Alpine for smaller size
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata make

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' \
  -a -installsuffix cgo \
  -tags netgo \
  -o caribex-backend ./cmd/api-server/main.go

# Runtime stage - use distroless for security and minimal size
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/caribex-backend /app/caribex-backend

# Copy migrations
COPY --from=builder /app/migrations /app/migrations

# Copy timezone data for proper time handling
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Use non-root user for security (already set in base image)
USER nonroot:nonroot

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV TZ=UTC
ENV GIN_MODE=release

# Health check for container orchestration
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/app/caribex-backend", "--health"] || exit 1

# Run the application
ENTRYPOINT ["/app/caribex-backend"]
