# # Stage 1: Download dependencies (smaller image for downloads)
# FROM golang:1.16 AS builder
# WORKDIR /app

# # Create the /app/go directory (ensure it exists before copying)
# RUN mkdir -p /app/go

# # Copy Go modules files
# COPY go.mod go.sum ./

# # Download and install Go dependencies
# RUN go mod download

# # Stage 2: Build and copy application (smaller final image)
# FROM alpine AS final
# WORKDIR /app

# # Copy dependencies only from builder stage
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# # Copy the rest of the application source code
# COPY . .
# # Update Alpine package list
# RUN apk update
# # Install Golang in the final stage
# RUN apk add --no-cache golang
# # Build the Go application
# RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix 'static' -o logger .

# # Expose the port on which the Go application will run
# EXPOSE 9000

# # Command to run the Go application
# CMD ["./app"]
# Use the official Golang image as the base image for builder stage
FROM golang:1.16 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application, excluding problematic packages
RUN CGO_ENABLED=0 GOOS=linux go build -o logger .

# Use a minimal base image for the final stage
FROM alpine:latest

# Install necessary certificates bundle
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /

# Copy the built Go binary from the builder stage
COPY --from=builder /app/hubpay-service .

# Expose the port on which the Go application will run
EXPOSE 80

# Command to run the Go application
CMD ["./hubpay-service"]
