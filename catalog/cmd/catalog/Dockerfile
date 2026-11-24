FROM golang:1.26.0-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/catalog ./cmd/catalog

FROM scratch

COPY --from=builder /out/catalog /app/catalog
COPY --from=builder /src/internal/store/postgres/migrations /app/internal/store/postgres/migrations

USER 10001:10001
WORKDIR /app

EXPOSE 8080

CMD ["/app/catalog"]
