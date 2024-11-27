# Stage 1: Build React frontend
FROM node:16 AS frontend-builder
WORKDIR /app
# Copy package.json and package-lock.json for dependency installation
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
# Copy the rest of the frontend source files
COPY frontend/ .
# Build the React frontend
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.20 AS backend-builder
WORKDIR /app
# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the backend source files
COPY . .
# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Stage 3: Final container
FROM alpine:3.18
WORKDIR /root/
# Install necessary tools
RUN apk --no-cache add ca-certificates bash
# Copy the Go binary from the backend build stage
COPY --from=backend-builder /app/main .
# Copy the frontend static files
COPY --from=frontend-builder /app/build ./frontend/build
# Copy the wait-for-it.sh script to wait for the database
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

# Serve frontend files with Go backend (configure in backend code)
EXPOSE 8080

# Start the application and wait for the database to be ready
CMD ["sh", "-c", "wait-for-it.sh db:3306 -- ./main"]
