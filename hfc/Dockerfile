# Stage 1: Build the Go app
FROM golang:1.23.1-alpine AS builder

# Create a directory for the app
WORKDIR /app

# Copy the Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o hfc .

# Stage 2: Create a minimal final image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/hfc /usr/local/bin/hfc

# Command to run the application
CMD ["hfc"]