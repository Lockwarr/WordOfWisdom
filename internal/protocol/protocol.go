package protocol

import "encoding/json"

const (
	ChallengeRequest = iota
	ChallengeResponse
	QuoteRequest
	QuoteResponse
	Stop
)

// Message - represents a message to be used for communication between tcp server and its' connected clients
type Message struct {
	// Accepted types of messages are ChallengeRequest, ChallengeResponse, QuoteRequest, QuoteResponse, Stop
	Type int `json:"type"`
	// Data could be a challenge in the from of json encoded haschash.Stamp or a quote
	Data string `json:"data"`
}

// ToJsonString - encodes protocol.Message to json string
func (m *Message) ToJsonString() string {
	msgBytes, _ := json.Marshal(m)
	return string(msgBytes)
}

// GetType - returns message type
func (m *Message) GetType(t int, d string) int {
	return m.Type
}

// GetData - returns message data
func (m *Message) GetData() string {
	return m.Data
}

// ParseMessage - decodes json message to protocol.Message
func ParseMessage(msgBytes []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
