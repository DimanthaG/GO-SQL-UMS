# Stage 1: Build React frontend
FROM --platform=linux/arm64 node:16 AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go backend
FROM --platform=linux/arm64 golang:1.20 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main ./cmd

# Stage 3: Final container
FROM --platform=linux/arm64 alpine:latest
WORKDIR /root/
# Copy the Go binary
COPY --from=backend-builder /app/main .
# Copy frontend static files
COPY --from=frontend-builder /app/build ./frontend/build
RUN apk --no-cache add ca-certificates bash
# Add wait-for-it.sh script
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

# Expose the Go backend port
EXPOSE 8080

# Run the backend and wait for the database
CMD ["sh", "-c", "wait-for-it.sh db:3306 -- ./main"]
