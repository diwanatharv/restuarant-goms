# Build Stage
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final Stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=builder /app/main .

# Expose the port your application listens on
EXPOSE 8000

# Set the entry point to the built binary
ENTRYPOINT ["./main"]