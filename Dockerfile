# Step 1: Modules caching
FROM golang:1.23.4 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.23.4 as builder
RUN apt-get update && \
    apt-get install -y librdkafka-dev && \
    rm -rf /var/lib/apt/lists/*

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=1 GOOS=linux go build -o /go/bin/main ./cmd/app

# Step 3: Final for production
FROM debian:bullseye-slim as production
# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y librdkafka1 ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -r -u 1001 -g root appuser

# Create app directory and set permissions
RUN mkdir /app && chown appuser:root /app

# Copy the binary from builder
COPY --from=builder /go/bin/main /app/

# Use the non-root user
USER appuser

# Set the working directory
WORKDIR /app

# Command to run the application
CMD ["./main"]