package protocol

import "encoding/json"

const (
	ChallengeRequest = iota
	ChallengeResponse
	QuoteRequest
	QuoteResponse
	Stop
)

// define struct for the message struct
type Message struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

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
