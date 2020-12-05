package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jshaw/virtualcamera/internal/messages"
)

func handleResp(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
        return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        return nil, err
	}

    return body, nil
}

func (c *Client) generateUrl(uri string) string {
    return fmt.Sprintf("%s://%s:%s/%s", c.scheme, c.host, c.port, uri)

}

func (c *Client) SendStats(ctx context.Context, msgs []*messages.Message) ([]byte, error) {
    body, err := messages.Serialize(msgs)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", c.generateUrl("send-stats"), bytes.NewReader(body))
    if err != nil {
       return nil, err
	}

    req = req.WithContext(ctx)
	resp, err := c.transport.Do(req)

    return handleResp(resp, err)

}

func (c *Client) Connect(ctx context.Context) ([]byte, error) {
    req, err := http.NewRequest("GET", c.generateUrl("control-channel"), nil)
    if err != nil {
        return nil, err
	}
    req = req.WithContext(ctx)
	httpResp, httpErr := c.transport.Do(req)
    resp, err := handleResp(httpResp, httpErr)

    return resp, err

}

func New(scheme string, host string, port string) (*Client, error) {
	client := &http.Client{
	    Timeout: 60 * time.Second,
	}

    return &Client{
        transport: client,
        host: host,
        port: port,
        scheme: scheme,

    }, nil
}
