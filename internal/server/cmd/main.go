package cmd

import (
	"context"
	"sitemapGenerator/WordOfWisdom/internal/repository"
	"sitemapGenerator/WordOfWisdom/internal/server"
)

func main() {
	tcpSrvr := server.NewTCPServer("localhost", "8080", repository.NewInMemoryDB())
	tcpSrvr.Start(context.Background())
}
