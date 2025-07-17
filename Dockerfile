FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app .
COPY .env .

ENV PORT=8080
ENV DB_URL=postgres://user:password@db:5432/infotecs?sslmode=disable

EXPOSE ${PORT}

CMD ["sh", "-c", "goose -dir ./sql/schema postgres ${DB_URL} up && ./main"]