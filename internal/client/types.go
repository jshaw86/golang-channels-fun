package client

import "net/http"

type Client struct {
    transport *http.Client
    host string
    port string
    scheme string
}


type ConnectResponse struct{
    Resp []byte
    Err error
}
