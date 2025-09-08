# Step 1: Modules caching
FROM golang:1.24-alpine AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.24-alpine AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /bin/avito-pvz ./cmd/main.go

# Step 3: Final
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/avito-pvz /app/avito-pvz
COPY --from=builder /app/config/config.yaml /app/config/config.yaml
COPY --from=builder /app/internal/database/migrations /app/database/migrations

WORKDIR /app
EXPOSE 8080
CMD ["/app/avito-pvz"]
