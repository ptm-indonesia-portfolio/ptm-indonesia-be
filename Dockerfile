FROM golang:1.26-alpine AS base

WORKDIR /app

RUN apk add --no-cache build-base ca-certificates git tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base AS development

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]

FROM base AS builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.22 AS production

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/api /app/api
COPY --from=builder /app/lang /app/lang
COPY --from=builder /app/db/migrations /app/db/migrations
COPY --from=builder /app/env.example /app/env.example

RUN mkdir -p /app/logs

EXPOSE 9100

ENTRYPOINT ["/app/api"]
