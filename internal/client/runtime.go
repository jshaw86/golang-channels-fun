package client

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"log"

	"github.com/jshaw/virtualcamera/internal/messages"
)

var thingsToSee = []string{"dog", "cat", "bike", "person"}

func GenerateMessages(ctx context.Context, c *Client,  msgs chan *messages.Message){
    rand.Seed(time.Now().Unix())
    L:
    for  {
        select {
        case <-ctx.Done():
           break L
        case <-time.After(10 * time.Second):
            log.Printf("Generating Message....")
            msgs <- messages.New(fmt.Sprintf("saw a %s", thingsToSee[rand.Intn(len(thingsToSee) - 0) + 0]))

        }

        log.Printf("Messages... %v", msgs)

    }

    log.Printf("Stopping message generation....")
}

func connectLoop(ctx context.Context, c *Client, connectResp chan []byte) {
    L:
    for {
        select {
        case <-ctx.Done():
            break L
        default:
            log.Printf("Connecting...")
            resp, err := c.Connect(ctx)
            if err != nil {
                log.Printf("err %s", err)
                break
            }
            connectResp <- resp
        }
    }
}

func RequestLoop(ctx context.Context, c *Client, msgs chan *messages.Message){
    var messageSet []*messages.Message

    connectResp := make(chan []byte)
    defer close(connectResp)

    go connectLoop(ctx, c, connectResp)

    L:
    for {
        select {
            case <-ctx.Done():
                break L
            case msg := <-msgs:
                messageSet = append(messageSet, msg)
            case <-connectResp:
                log.Printf("Sending stats...")
                c.SendStats(ctx, messageSet)
                log.Printf("Sent stats... %v", msgs)
        }
    }

    log.Printf("Stopping request loop...")
}
