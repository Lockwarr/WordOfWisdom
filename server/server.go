package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Lockwarr/WordOfWisdom/internal/hashcash"
	"github.com/Lockwarr/WordOfWisdom/internal/repository"
	"github.com/Lockwarr/WordOfWisdom/protocol"
)

// Quotes - const array of quotes to respond on client's request
var Quotes = []string{
	"Quote 1",
	"Quote 2",
	"Quote 3",
	"Quote 4",
	"Quote 5",
}

type Server interface {
	Start(context.Context)
	ProcessRequest(context.Context, string, string) (*protocol.Message, error)
}

type tcpServer struct {
	port string
	host string
	repo repository.Repository
}

func NewTCPServer(host, port string, repo repository.Repository) Server {
	return &tcpServer{
		port: port,
		host: host,
		repo: repo,
	}
}

// implement Start function
func (s *tcpServer) Start(ctx context.Context) {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", s.host+":"+s.port)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + s.host + ":" + s.port)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go s.handleConnection(ctx, conn)

	}
}

func (s *tcpServer) handleConnection(ctx context.Context, conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err read connection:", err)
			return
		}
		msg, err := s.ProcessRequest(ctx, req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}
}

// ProcessRequest handles incoming requests.
func (s *tcpServer) ProcessRequest(ctx context.Context, message, clientDetails string) (*protocol.Message, error) {
	parsedMessage, err := protocol.ParseMessage([]byte(message))
	if err != nil {
		log.Println("Error parsing:", err.Error())
	}

	switch parsedMessage.Type {
	case protocol.ChallengeRequest:
		log.Println("Challenge request received")
		indicator := rand.Intn(200000)
		rand := strconv.Itoa(indicator)
		stamp := hashcash.Stamp{
			Version:    1,
			ZerosCount: 5,
			Date:       time.Now().Unix(),
			Resource:   parsedMessage.Data,
			Rand:       rand,
			Counter:    0,
		}

		err := s.repo.AddIndicator(ctx, int64(indicator))
		if err != nil {
			return nil, fmt.Errorf("Error adding indicator: %w", err)
		}

		marshaledStamp, err := json.Marshal(stamp)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling stamp: %w", err)
		}
		respMsg := protocol.Message{
			Type: protocol.ChallengeResponse,
			Data: string(marshaledStamp),
		}

		return &respMsg, nil
	case protocol.QuoteRequest:
		fmt.Printf("client %s requests quote %s\n", clientDetails, parsedMessage.Data)
		// parse client's solution
		var stamp hashcash.Stamp
		err := json.Unmarshal([]byte(parsedMessage.Data), &stamp)
		if err != nil {
			return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
		}
		fmt.Println(stamp)
		// validate hashcash params
		if !stamp.ValidStamp(ctx, stamp, s.repo) {
			return nil, fmt.Errorf("invalid hashcash %w", err)
		}

		randValue, err := strconv.Atoi(stamp.Rand)
		if err != nil {
			return nil, fmt.Errorf("err decode rand: %w", err)
		}

		// if rand exists in inmemory db, it means, that hashcash is valid and really challenged by this server in past
		_, err = s.repo.GetIndicator(ctx, int64(randValue))
		if err != nil {
			return nil, fmt.Errorf("err get rand from cache: %w", err)
		}

		if !stamp.IsHashSolved() {
			return nil, fmt.Errorf("challenge is not solved")
		}

		//get random quote
		fmt.Printf("client %s succesfully computed hashcash %s\n", clientDetails, parsedMessage.Data)

		msg := protocol.Message{
			Type: protocol.QuoteResponse,
			Data: Quotes[rand.Intn(5)],
		}

		// delete rand from cache to prevent duplicated request with same hashcash value
		s.repo.RemoveIndicator(ctx, int64(randValue))

		// respond to client
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown request received")
	}
}

// sendMsg - send protocol message to connection
func sendMsg(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToJsonString())
	_, err := conn.Write([]byte(msgStr))
	return err
}
