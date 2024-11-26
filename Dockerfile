# Stage 1: Build the Go application
FROM golang:1.20 AS builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application binary for Linux (cross-compile)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags=-static" -o main ./cmd

# Stage 2: Create a lightweight runtime container
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

COPY wait-for-it.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/wait-for-it.sh



EXPOSE 8080

# Use the wait-for-it script to ensure the database is ready before starting the app
CMD ["wait-for-it", "db:3306", "--", "./main"]
