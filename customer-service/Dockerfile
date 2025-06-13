FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o customer-service .

EXPOSE 8080

CMD ["./customer-service"]
