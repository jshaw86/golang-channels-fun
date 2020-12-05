package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jshaw/virtualcamera/internal/messages"
)
func handleReq(resp *http.Request) ([]byte, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        return nil, err
	}

    return body, nil
}

func (s *Server) HandleControlChannel(w http.ResponseWriter, r *http.Request) {
    log.Printf("[Handle Channel] \t Waiting...")
    select {
    case <-r.Context().Done():
        log.Printf("[Handle Channel] \t Client went away")
        return
    case command := <-s.ControlChannel:
        if command == "send-stats" {
            // Log the request protocol
            log.Printf("[Handle Channel] \t Cleared wait... %s", r.Proto)
            // Send a message back to the client
            w.Write([]byte("ok"))
        }
    }
}

func (s *Server) HandleSendStats(w http.ResponseWriter, r *http.Request) {
    body, err := handleReq(r)
    if err != nil {
        // TODO
        log.Printf("[Handle Stats] \t Handle request body error %s", err)
        return
    }

    msgs, unmarshalErr := messages.Deserialize(body)

    if unmarshalErr != nil {
        log.Printf("[Handle Stats] \t Unmarshal send-stats error... %s", err)
        return
    }

    log.Printf("[Handle Stats] \t Messages... %v", msgs)

    for _, msg := range msgs {
        if msg == nil {
            continue
        }
        log.Printf("[Handle Stats] \t Sending message... %v", msg)
        s.DataChannel <- &MessageEnvelope{msg, false}
        log.Printf("[Handle Stats] \t Sent message... %v", msg)
    }

    s.DataChannel <- &MessageEnvelope{nil, true}

    // Send a message back to the client
    w.Write([]byte("ok"))
}

func (s *Server) HandleLogs(w http.ResponseWriter, r *http.Request) {
    log.Printf("[Handle Logs] \t Publish to control channel... ")
    select {
    case s.ControlChannel <- "send-stats":
        break
    case <-time.After(5 * time.Second):
        w.WriteHeader(http.StatusRequestTimeout)
        return
    }


    var msgs []*messages.Message

    L:
    for {
        log.Printf("[Handle Logs] \t Waiting for data....")
        select {
        case msgEnvelope := <-s.DataChannel:
            if msgEnvelope.EOL == true {
                break L
            }

            log.Printf("[Handle Logs] \t Read from data channel... %v", msgEnvelope)
            msgs = append(msgs, msgEnvelope.Msg)

        }
    }

    if msgs == nil {
        msgs = make([]*messages.Message,0)
    }

    responseBytes, err := messages.Serialize(msgs)
    if err != nil {
        log.Printf("[Handle Logs] \t Serialize error... %s", err)
    }

    w.Write(responseBytes)

}
