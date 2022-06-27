package client_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Lockwarr/WordOfWisdom/client"
	"github.com/Lockwarr/WordOfWisdom/internal/repository"
	"github.com/Lockwarr/WordOfWisdom/server"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	repo := repository.NewInMemoryDB()
	tcpSrvr := server.NewTCPServer("localhost", "8000", repo)
	go tcpSrvr.Start(context.Background())
	time.Sleep(time.Second)

	code := m.Run()
	tcpSrvr.Stop()
	os.Exit(code)
}

// func (s *clientTestSuite) TestRequestQuote() {
// 	// Arrange
// 	conn, err := net.Dial("tcp", ":8000")
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

	// Act
	err := client.Run(context.Background(), ":8000")

	// Assert
	assert.Equal(t, nil, err)
}
