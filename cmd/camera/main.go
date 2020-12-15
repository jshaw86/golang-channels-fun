package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jshaw/virtualcamera/internal/client"
	"github.com/jshaw/virtualcamera/internal/messages"
)

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {
    // setup clean shutdown
    log.Printf("Starting camera...")
    ctx, cancel := context.WithCancel(context.Background())
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

    defer close(stop)

    c, err := client.New(getEnv("SCHEME", "http"), getEnv("HOST", "localhost"), getEnv("PORT", "8080"))

    if (err != nil){
        log.Fatal(err)
    }

    messageChannel := make(chan *messages.Message)
    go client.GenerateMessages(ctx, c, messageChannel)
    go client.RequestLoop(ctx, c, messageChannel)

    <-stop
    log.Printf("Shutting down...")
    cancel()


}
