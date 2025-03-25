FROM --platform=linux/amd64 golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

FROM --platform=linux/amd64 alpine:3.18

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/api .
COPY --from=builder /app/.env .

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["./api"]