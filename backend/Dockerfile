# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Install swag for Swagger documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Ensure all dependencies are downloaded and go.sum is in sync
RUN go mod download && go mod verify

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Build the recurring transaction processor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o process-recurring ./cmd/process-recurring

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /root/

# Create bin directory
RUN mkdir -p bin

# Copy the binaries from builder
COPY --from=builder /app/main ./bin/api
COPY --from=builder /app/process-recurring ./bin/process-recurring

# Expose port
EXPOSE 8080

# Run the API by default
CMD ["./bin/api"]

