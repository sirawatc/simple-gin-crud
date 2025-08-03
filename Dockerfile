FROM golang:1.24.5-alpine3.22

RUN adduser -D app

USER app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux go build -a -o main ./cmd/main

CMD ["./main"]