FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./gateway/go.mod ./gateway/go.sum ./

COPY ./proto /proto

RUN go mod download

COPY ./gateway .

RUN go build -ldflags "-s -w" -o gateway ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gateway .

RUN chmod +x /app/gateway

CMD ["./gateway"]