FROM golang:1.17.8

WORKDIR /apps

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o serverExecutable ./server/cmd/main.go
RUN chmod +x serverExecutable

EXPOSE 8080

CMD "./serverExecutable"