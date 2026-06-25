package ticker

import (
	"encoding/json"
	"net/url"
)

const (
	ModeLTP   = "ltp"
	ModeQuote = "quote"
	ModeFull  = "full"
	URL       = "wss://ws.kite.trade"
)

type Client struct {
	apiKey      string
	accessToken string
}

type Message struct {
	Action string `json:"a"`
	Value  any    `json:"v"`
}

func NewClient(apiKey, accessToken string) *Client {
	return &Client{apiKey: apiKey, accessToken: accessToken}
}

func (c *Client) Endpoint() string {
	u, _ := url.Parse(URL)
	q := u.Query()
	q.Set("api_key", c.apiKey)
	q.Set("access_token", c.accessToken)
	u.RawQuery = q.Encode()
	return u.String()
}

func Subscribe(tokens ...uint32) []byte {
	body, _ := json.Marshal(Message{Action: "subscribe", Value: tokens})
	return body
}

func Unsubscribe(tokens ...uint32) []byte {
	body, _ := json.Marshal(Message{Action: "unsubscribe", Value: tokens})
	return body
}

func SetMode(mode string, tokens ...uint32) []byte {
	body, _ := json.Marshal(Message{Action: "mode", Value: []any{mode, tokens}})
	return body
}
