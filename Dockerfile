FROM golang:1.19 AS builder

WORKDIR /app

COPY . .

RUN go build -o /app/cinema ./cmd

FROM alpine:latest

WORKDIR /app

RUN apk add gcompat

COPY --from=builder /app/cinema .

ENTRYPOINT ["./cinema"]
