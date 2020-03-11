package scale

import (
	"encoding/json"
	"strings"
)

type Message struct {
	Body  string `json:"body"`
	Value string `json:"value"`
}

func NewMessage(body, value []byte) *Message {
	return &Message{
		Body:  strings.TrimSpace(string(body)),
		Value: strings.TrimSpace(string(value)),
	}
}

func (m *Message) Marshal() []byte {
	j, _ := json.Marshal(m)
	return j
}
