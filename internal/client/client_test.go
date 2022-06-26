package client_test

import (
	"context"
	"sitemapGenerator/WordOfWisdom/internal/client"
	"sitemapGenerator/WordOfWisdom/internal/repository"
	"sitemapGenerator/WordOfWisdom/internal/server"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func (s *clientTestSuite) TestRequestQuote() {
// 	// Arrange
// 	conn, err := net.Dial("tcp", ":8080")
// 	s.NoError(err)
// 	defer conn.Close()

// 	// Act
// 	quote, err := client.RequestQuote(context.Background(), conn)

// 	// Assert
// 	s.Equal(true, strings.Contains(quote, "Quote"))
// 	s.Equal(nil, err)
// }

func TestClientRun(t *testing.T) {
	//Arrange
	repo := repository.NewInMemoryDB()
	tcpServer := server.NewTCPServer("localhost", "8080", repo)
	go tcpServer.Start(context.Background())
	time.Sleep(time.Second)

	// Act
	err := client.Run(context.Background(), ":8080")

	// Assert
	assert.Equal(t, nil, err)
}

func TestClientRun_WithoutRunningServer_ThenFail(t *testing.T) {
	//Arrange

	// Act
	err := client.Run(context.Background(), ":5555")

	// Assert
	assert.Equal(t, "dial tcp :5555: connect: connection refused", err.Error())
}
