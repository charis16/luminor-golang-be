# Stage 1: Build Go App
FROM golang:1.23-bullseye AS builder

# Enable Go modules and optimizations
ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy the full source code
COPY src/ ./src
WORKDIR /app/src

# Build statically-linked binary
RUN go build -o main .

# Stage 2: Minimal runtime container
FROM alpine:latest

# Create non-root user for security
RUN adduser -D appuser
USER appuser

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/src/main .

# Expose the default port (can still be overridden by env)
EXPOSE 8080

# Use ENTRYPOINT for better flexibility (e.g. app arguments override)
ENTRYPOINT ["./main"]