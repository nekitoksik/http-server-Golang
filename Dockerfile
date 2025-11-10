FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init -g cmd/app/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

COPY --from=builder /app/docs ./docs

COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./server"]