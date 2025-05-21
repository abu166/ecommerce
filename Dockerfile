# Use the official Go image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to install dependencies
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application (will be overridden by docker-compose command)
RUN go build -o main ./cmd/apigateway/main.go

# Expose ports (will be overridden by docker-compose)
EXPOSE 8080

# Default command (will be overridden by docker-compose)
CMD ["./main"]