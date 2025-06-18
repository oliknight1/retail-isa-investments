FROM golang:1.23.10-alpine AS builder

WORKDIR /app

COPY repository/funds.json ./repository/funds.json

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o fund-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/fund-service .
COPY --from=builder /app/repository/funds.json ./repository/funds.json

EXPOSE 8080

CMD ["./fund-service"]

