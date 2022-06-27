FROM golang:1.17.8

WORKDIR /apps

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o clientExecutable ./client/cmd/main.go
RUN chmod +x clientExecutable

CMD "./clientExecutable"