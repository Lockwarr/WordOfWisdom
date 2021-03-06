package main

import (
	"context"

	"github.com/Lockwarr/WordOfWisdom/internal/repository"
	"github.com/Lockwarr/WordOfWisdom/server"
)

// TODO: start using config file
const port = "8080"
const host = "0.0.0.0"

func main() {
	tcpSrvr := server.NewTCPServer(host, port, repository.NewInMemoryDB())
	tcpSrvr.Start(context.Background())
}
