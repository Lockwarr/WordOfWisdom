FROM golang:1.17.8 AS builder

WORKDIR /apps

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./client/cmd

CMD ["main"]