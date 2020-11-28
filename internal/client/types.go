package client

import "net/http"

type Client struct {
    transport *http.Client
}


type ConnectResponse struct{
    Resp []byte
    Err error
}
