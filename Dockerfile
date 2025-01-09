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
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/app

# Step 3: Final
FROM debian:bookworm-slim
RUN apt-get update && \
    apt-get install -y librdkafka1 librdkafka++1 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/app /app
COPY --from=builder /app/config /config
CMD ["/app"]