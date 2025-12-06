FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./dice/go.mod ./dice/go.sum ./

COPY ./proto /proto

RUN go mod download

COPY ./dice .

RUN go build -ldflags "-s -w" -o dice ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dice .

RUN chmod +x /app/dice

CMD ["./dice"]