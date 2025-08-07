FROM golang:1.24.6-alpine3.22 AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

EXPOSE 8080

CMD ["/build/main"]