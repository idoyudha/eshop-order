# Step 1: Modules caching
FROM golang:1.23.4 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.23.4 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

# Build with CGO enabled for Kafka
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/app

# Step 3: Final for production
FROM debian:bullseye-slim as production

# Install required packages in a single layer
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    libsasl2-2 \
    libsasl2-modules \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -r -u 1001 -g root appuser

# Copy the binary from builder
COPY --from=builder /app/main /app/main

# Set ownership
RUN chown appuser:root /app/main

# Use the non-root user
USER appuser

# Set the working directory
WORKDIR /app

# Command to run the application
ENTRYPOINT ["/app/main"]