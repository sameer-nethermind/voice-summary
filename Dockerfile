# ---------------------------------
# Stage 1: Build Go binary
# ---------------------------------
FROM golang:1.22.0 AS builder

WORKDIR /app

# Copy go.mod and go.sum for dependency download
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of your source code
COPY .env .env
COPY . .

# Build the Go application (static binary), pointing to ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# ---------------------------------
# Stage 2: Minimal runtime image
# ---------------------------------
FROM alpine:3.17

WORKDIR /app

# Copy just the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file as well
COPY --from=builder /app/.env .

# Expose the port your app listens on
EXPOSE 8080

# Start the Go server
CMD ["./main"]
