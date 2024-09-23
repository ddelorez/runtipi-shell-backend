# Build Stage
FROM golang:1.23.0-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o backend cmd/main.go

# Final Stage
FROM alpine:latest

# Install bash
RUN apk add --no-cache bash

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/backend .

# Expose port 8080
EXPOSE 8080

# Set environment variables (optional)
ENV ALLOWED_ORIGINS=http://localhost
ENV SHELL_COMMAND=bash

# Run the application
CMD ["./backend"]
