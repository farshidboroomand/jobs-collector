# ---------- Stage 1: Build ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (required for go mod in some cases)
RUN apk add --no-cache git

# Cache go modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o jobs-collector ./cmd/main.go


# ---------- Stage 2: Runtime ----------
FROM alpine:3.19

WORKDIR /app

# Install certificates (needed for HTTPS calls later)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/jobs-collector .

# Expose API port
EXPOSE 8080

# Run application
CMD ["./jobs-collector"]
