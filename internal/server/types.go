package server

import "github.com/jshaw/virtualcamera/internal/messages"

type Server struct {
    ControlChannel chan *ControlMessage
    DataChannel chan *MessageEnvelope
}

type MessageEnvelope struct {
    Msg *messages.Message
    EOL bool

}

type ControlMessage struct {
    Command string
}
