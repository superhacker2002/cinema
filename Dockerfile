FROM golang

WORKDIR /app

COPY . .

RUN go build -o /app/cinema ./cmd

ENTRYPOINT ["./cinema"]
