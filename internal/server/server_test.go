package server_test

import (
	"context"
	"encoding/json"
	"net"
	"sitemapGenerator/WordOfWisdom/internal/hashcash"
	"sitemapGenerator/WordOfWisdom/internal/repository"
	"sitemapGenerator/WordOfWisdom/internal/server"
	"sitemapGenerator/WordOfWisdom/protocol"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandlingConnection(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()
	tcpSrvr := server.NewTCPServer("localhost", "8080", repo)

	go tcpSrvr.Start(context.Background())
	time.Sleep(time.Second)

	// Start client
	conn, err := net.Dial("tcp", ":8080")
	assert.NoError(t, err)
	defer conn.Close()

	// Request challenge
	message := protocol.Message{Type: protocol.ChallengeRequest, Data: "empty"}
	msgBytes, err := json.Marshal(message)
	assert.NoError(t, err)
	_, err = conn.Write(append(msgBytes, '\n'))
	assert.NoError(t, err)
}

func TestProcessChallengeRequest(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()
	tcpServer := server.NewTCPServer("localhost", "8080", repo)
	message := protocol.Message{Type: protocol.ChallengeRequest, Data: "empty"}
	msgBytes, err := json.Marshal(message)

	// Act
	// Request challenge
	msg, err := tcpServer.ProcessRequest(context.Background(), string(msgBytes), "testClient")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, protocol.ChallengeResponse, msg.Type)
}

func TestProcessQuoteRequest(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()
	tcpServer := server.NewTCPServer("localhost", "8080", repo)
	stamp := hashcash.Stamp{}

	// We need to send a challenge request first so we can have an indicator entry in the repo
	message := protocol.Message{Type: protocol.ChallengeRequest, Data: "empty"}
	msgBytes, err := json.Marshal(message)
	assert.NoError(t, err)
	msg, err := tcpServer.ProcessRequest(context.Background(), string(msgBytes), "testClient")
	assert.NoError(t, err)

	// We need to solve the challenge that server returned as a response to our challenge request
	err = json.Unmarshal([]byte(msg.Data), &stamp)
	assert.NoError(t, err)
	solvedStamp, err := stamp.ComputeHashcash(10000000)
	assert.NoError(t, err)

	solvedStampMarshaled, err := json.Marshal(solvedStamp)
	assert.NoError(t, err)
	message2 := protocol.Message{Type: protocol.QuoteRequest, Data: string(solvedStampMarshaled)}
	msgBytes2, err := json.Marshal(message2)
	assert.NoError(t, err)

	// Act
	// Request quote
	msg2, err := tcpServer.ProcessRequest(context.Background(), string(msgBytes2), "testClient")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, protocol.QuoteResponse, msg2.Type)
}

func TestProcessUnknownRequest(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryDB()
	tcpServer := server.NewTCPServer("localhost", "8080", repo)
	message := protocol.Message{Type: 7, Data: "empty"}
	msgBytes, err := json.Marshal(message)
	assert.NoError(t, err)

	// Act
	// Request challenge
	_, err = tcpServer.ProcessRequest(context.Background(), string(msgBytes), "testClient")

	// Assert
	assert.Error(t, err)
}
