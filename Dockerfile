# Build stage
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the Go binary explicitly for Linux
RUN GOOS=linux GOARCH=amd64 go build -o /app/server cmd/server/main.go

# Final image
FROM alpine:latest

# Install necessary system dependencies
RUN apk --no-cache add ca-certificates libc6-compat

# Set the working directory
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/server /server

# Expose the port the server runs on
EXPOSE 8080

# Command to run the server
CMD ["/server"]
