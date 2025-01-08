# Step 1: Modules caching
FROM golang:1.23.4 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.23.4 as builder
# Install necessary build dependencies
RUN apt-get update && \
    apt-get install -y librdkafka-dev && \
    rm -rf /var/lib/apt/lists/*

# Copy the module files
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

# Build with specific flags to ensure compatibility
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o main ./cmd/app

# Step 3: Final for production
FROM ubuntu:22.04 as production
# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y librdkafka1 ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Copy binary and config files
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

RUN chmod +x main

# Command to run the application
CMD ["/app/main"]