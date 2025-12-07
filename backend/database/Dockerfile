FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o database-app cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup


COPY --from=builder /app/database-app .

COPY --from=builder /app/config/config.yaml ./config/config.yaml

COPY --from=builder /app/internal/migrations ./internal/migrations

RUN touch app.log && chown appuser:appgroup app.log

USER appuser

EXPOSE 2012

CMD ["./database-app", "--config=./config/config.yaml"]