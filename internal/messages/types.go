package messages

import "time"

type Message struct {
    Description string `json:"description"`
    Timestamp time.Time `json:"timestamp"`
}

