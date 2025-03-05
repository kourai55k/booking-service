FROM golang:1.24 AS builder

WORKDIR /root

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/booking-service/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /root/main .

EXPOSE 8080

CMD ["./main"]