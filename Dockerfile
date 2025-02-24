# Use official Golang image as a build stage
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the Go binary
RUN go build -o server cmd/server/main.go

# Use a lightweight Alpine image for the final container
FROM alpine:latest

# Install necessary system dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/server /server

# Copy .env file (optional, if using ENV variables)
COPY config/.env .env

# Expose the port the server runs on
EXPOSE 8080

# Command to run the server
CMD ["/server"]
