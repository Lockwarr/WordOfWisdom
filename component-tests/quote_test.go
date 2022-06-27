//go:build component
// +build component

package component_tests

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Lockwarr/WordOfWisdom/internal/hashcash"
	"github.com/Lockwarr/WordOfWisdom/internal/protocol"
	"github.com/Lockwarr/WordOfWisdom/internal/repository"
	"github.com/Lockwarr/WordOfWisdom/server"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

const maxIterations = 10000000

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	f := &quoteFeature{}

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: f.InitializeTestSuite,
		ScenarioInitializer:  f.InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts) // godog v0.11.0 (latest)
}

type quoteFeature struct {
	tcpConn            net.Conn
	challenge          string
	quote              string
	solvedQuoteRequest protocol.Message
}

func (f *quoteFeature) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	repo := repository.NewInMemoryDB()
	tcpSrvr := server.NewTCPServer("localhost", "8080", repo)
	go tcpSrvr.Start(context.Background())
	// give time to start the server
	time.Sleep(1 * time.Second)
}

func (f *quoteFeature) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		f.quote = ""
		f.challenge = ""
		f.solvedQuoteRequest = protocol.Message{}
	})
	ctx.Step(`^a tcp connection exists with server running$`, f.aTcpConnectionWithServerRunning)
	ctx.Step(`^I send a "([^"]*)" to the server$`, f.iSendRequestToTheServer)
	ctx.Step(`^I receive "([^"]*)"$`, f.iReceiveResponse)
	ctx.Step(`^I solve the challenge$`, f.iSolveTheChallenge)
	ctx.Step(`^the received Quote is "([^"]*)"$`, f.isValidQuote)

}

// Given
func (f *quoteFeature) aTcpConnectionWithServerRunning() error {
	var err error
	f.tcpConn, err = net.Dial("tcp", ":8080")
	return err
}

// When
func (f *quoteFeature) iSendRequestToTheServer(messageType string) error {
	var message protocol.Message
	switch messageType {
	case "ChallengeRequest":
		message = protocol.Message{Type: protocol.ChallengeRequest, Data: "empty"}
	case "QuoteRequest":
		message = f.solvedQuoteRequest
	}
	msgStr := fmt.Sprintf("%s\n", message.ToJsonString())
	_, err := f.tcpConn.Write([]byte(msgStr))
	return err
}

// Then
func (f *quoteFeature) iReceiveResponse(requestType string) error {
	var err error
	reader := bufio.NewReader(f.tcpConn)
	switch requestType {
	case "ChallengeResponse":
		f.challenge, err = reader.ReadString('\n')
		return err
	case "QuoteResponse":
		f.quote, err = reader.ReadString('\n')
		return err
	}
	return nil
}

func (f *quoteFeature) iSolveTheChallenge() error {
	stamp := hashcash.Stamp{}
	challengeResponseMessage := protocol.Message{}
	err := json.Unmarshal([]byte(f.challenge), &challengeResponseMessage)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(challengeResponseMessage.Data), &stamp)
	if err != nil {
		return err
	}
	solvedStamp, err := stamp.ComputeHashcash(maxIterations)
	if err != nil {
		return err
	}

	solvedStampMarshalled, err := json.Marshal(solvedStamp)
	if err != nil {
		return err
	}
	f.solvedQuoteRequest = protocol.Message{Type: protocol.QuoteRequest, Data: string(solvedStampMarshalled)}
	return nil
}

func (f *quoteFeature) isValidQuote(isValid string) error {
	// Read quote response

	switch isValid {
	case "valid":
		quoteResponseMessage := protocol.Message{}
		err := json.Unmarshal([]byte(f.quote), &quoteResponseMessage)
		if err != nil {
			return err
		}
		if !strings.Contains(quoteResponseMessage.Data, "Quote") {
			fmt.Println(quoteResponseMessage)
			return fmt.Errorf("Invalid response")
		}
	case "invalid":
		if f.quote != "" {
			return fmt.Errorf("expected invalid response")
		}
	}

	return nil
}
