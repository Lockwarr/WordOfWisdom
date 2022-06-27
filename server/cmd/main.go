package main

import (
	"context"

	"github.com/Lockwarr/WordOfWisdom/internal/repository"
	"github.com/Lockwarr/WordOfWisdom/internal/server"
)

// TODO: start using config file
const address = ":8080"
const host = "localhost"

func main() {
	tcpSrvr := server.NewTCPServer(host, address, repository.NewInMemoryDB())
	tcpSrvr.Start(context.Background())
}
