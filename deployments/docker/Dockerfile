# Build stage
FROM golang:1.26-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o blind-service \
    ./cmd/api

# Final stage
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Copy SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /build/blind-service /blind-service

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/blind-service", "health"]

# Run the binary
ENTRYPOINT ["/blind-service"]