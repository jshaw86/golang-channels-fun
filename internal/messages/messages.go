package messages

import (
	"encoding/json"
	"time"
)

func New(description string) *Message {
    return &Message{
        Description: description,
        Timestamp: time.Now(),
    }
}

func Deserialize(body []byte) ([]*Message, error) {
    messages :=[]*Message{}

    err := json.Unmarshal(body, &messages)

    return messages, err
}

func Serialize(m []*Message) ([]byte, error) {
    b, err := json.Marshal(m)
    return b, err
}
