package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sitemapGenerator/WordOfWisdom/internal/hashcash"
	"sitemapGenerator/WordOfWisdom/protocol"
)

const maxIterations = 10000000

// Run - connect to given address and send request
func Run(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Println("connected to", address)

	quote, err := requestQuote(ctx, conn)
	if err != nil {
		return err
	}
	fmt.Println("quote result:", quote)
	return nil
}

func requestQuote(ctx context.Context, conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	// Request challenge
	message := protocol.Message{Type: protocol.ChallengeRequest, Data: "empty"}
	err := sendMsg(message, conn)
	if err != nil {
		return "", fmt.Errorf("err send message: %w", err)
	}

	// We need to solve the returned challenge
	resp, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("err read connection: %w", err)
	}
	quoteRequest, err := handleChallengeResponse(resp)
	if err != nil {
		return "", fmt.Errorf("err handle challenge response: %w", err)
	}

	// Request quote with solved challenge
	err = sendMsg(*quoteRequest, conn)
	if err != nil {
		return "", fmt.Errorf("err send message: %w", err)
	}

	// Read quote response
	respQuote, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("err read quote response: %w", err)
	}
	quoteResponseMessage := protocol.Message{}
	err = json.Unmarshal([]byte(respQuote), &quoteResponseMessage)
	if err != nil {
		return "", fmt.Errorf("err unmarshal quote response: %w", err)
	}
	return quoteResponseMessage.Data, nil
}

func handleChallengeResponse(resp string) (*protocol.Message, error) {
	stamp := hashcash.Stamp{}

	challengeResponseMessage := protocol.Message{}
	err := json.Unmarshal([]byte(resp), &challengeResponseMessage)
	if err != nil {
		return nil, fmt.Errorf("err unmarshal message: %w", err)
	}
	err = json.Unmarshal([]byte(challengeResponseMessage.Data), &stamp)
	if err != nil {
		return nil, fmt.Errorf("err unmarshal message data: %w", err)
	}
	solvedStamp, err := stamp.ComputeHashcash(maxIterations)
	if err != nil {
		return nil, fmt.Errorf("err compute hashcash: %w", err)
	}

	solvedStampMarshalled, err := json.Marshal(solvedStamp)
	if err != nil {
		return nil, fmt.Errorf("err marshal stamp: %w", err)
	}
	quoteRequest := protocol.Message{Type: protocol.QuoteRequest, Data: string(solvedStampMarshalled)}
	return &quoteRequest, nil
}

// sendMsg - send protocol message to connection
func sendMsg(msg protocol.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToJsonString())
	_, err := conn.Write([]byte(msgStr))
	return err
}
