FROM golang:1.25-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

FROM debian:bullseye-slim

RUN groupadd -r appuser && useradd -r -g appuser appuser
RUN install -d -o appuser -g appuser /data
ENV FILES_BASE_DIR=/data

WORKDIR /app

COPY --chown=appuser:appuser --from=builder /build/server .

COPY --chown=appuser:appuser --from=builder /build/migrations ./migrations

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER appuser

EXPOSE 8080

CMD ["./server"]