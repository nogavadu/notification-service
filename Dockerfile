FROM golang:1.24.1-alpine AS builder

COPY . /github.com/nogavadu/notification-service
WORKDIR /github.com/nogavadu/notification-service

RUN go mod download
RUN go build -o ./bin/notification-service cmd/notification-service/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/nogavadu/notification-service/bin/notification-service .
COPY --from=builder /github.com/nogavadu/notification-service/config/config.yaml .

CMD ["./notification-service", "-config=./config.yaml"]