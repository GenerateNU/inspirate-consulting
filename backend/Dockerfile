# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the executable
# -o specifies output location
# CGO_ENABLED=0 creates a static binary (better for alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/inspirate-consulting ./cmd

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from builder stage
COPY --from=builder /app/bin/inspirate-consulting .

EXPOSE 8080

# Lets docker-compose.prod.yml's `condition: service_healthy` resolve. The
# deploy workflow polls the same endpoint through Caddy before going live.
HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- http://localhost:8080/api/v1/health || exit 1

# Run the executable
CMD ["./inspirate-consulting"]