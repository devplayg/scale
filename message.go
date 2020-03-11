package scale

import "encoding/json"

type Message struct {
	Data  []byte
	Value string
}

func NewMessage(data []byte, value string) *Message {
	return &Message{
		Data:  data,
		Value: value,
	}
}

func (m *Message) Marshal() []byte {
	j, _ := json.Marshal(m)
	return j
}
