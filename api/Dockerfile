# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o ./docs --parseInternal --parseDependency

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o reddit-content-analyzer cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /root/

# Copy the binary and config from builder
COPY --from=builder /app/reddit-content-analyzer .
COPY --from=builder /app/config.yaml .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./reddit-content-analyzer"]

