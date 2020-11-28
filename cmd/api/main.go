package main

import (
	"log"
	"net/http"

	"github.com/jshaw/virtualcamera/internal/server"
)

func main() {
	srv := &http.Server{Addr: ":8080"}

    s := server.Server{
        ControlChannel: make(chan string),
        DataChannel: make(chan *server.MessageEnvelope),

    }
    http.HandleFunc("/logs", s.HandleLogs)
    http.HandleFunc("/send-stats", s.HandleSendStats)
    http.HandleFunc("/control-channel", s.HandleControlChannel)

	log.Printf("Listening https://0.0.0.0:8080")
	log.Fatal(srv.ListenAndServe())
}


