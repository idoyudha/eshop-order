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
# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Step 3: Final for production
FROM alpine:3.19
# Add CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# Create a non-root user
RUN adduser -D -g '' appuser

# Copy the binary from builder
COPY --from=builder /app/main /app/main

# Use the non-root user
USER appuser

# Set the working directory
WORKDIR /app

# Command to run the application
ENTRYPOINT ["/app/main"]