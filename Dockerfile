# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and build essentials if CGO is needed for sqlite
RUN apk add --no-cache git build-base

# Copy dependency specifications
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with CGO enabled (for sqlite support) or CGO_ENABLED=0 for pure Go
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Stage 2: Minimal runtime image
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/.env.example .env

# Expose server port
EXPOSE 8080

# Run server
CMD ["./server"]
