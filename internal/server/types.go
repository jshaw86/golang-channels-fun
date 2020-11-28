package server

import "github.com/jshaw/virtualcamera/internal/messages"

type Server struct {
    ControlChannel chan string
    DataChannel chan *MessageEnvelope
}

type MessageEnvelope struct {
    Msg *messages.Message
    EOL bool

}
