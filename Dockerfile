# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o sharer main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/sharer .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Create directory for SQLite database
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV DB_PATH=/app/data/sharer.db

# Run the application
CMD ["./sharer"]