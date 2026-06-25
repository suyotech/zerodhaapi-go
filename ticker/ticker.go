package ticker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	URL = "wss://ws.kite.trade"
)

type Mode string

const (
	ModeLTP   Mode = "ltp"
	ModeQuote Mode = "quote"
	ModeFull  Mode = "full"
)

type Client struct {
	apiKey        string
	accessToken   string
	endpoint      string
	dialer        *websocket.Dialer
	conn          *websocket.Conn
	mu            sync.RWMutex
	writeMu       sync.Mutex
	subscriptions map[uint32]struct{}
	modes         map[Mode]map[uint32]struct{}
	onTick        func([]Tick)
	onError       func(error)
	onConnect     func()
	onClose       func(error)
}

type Message struct {
	Action string `json:"a"`
	Value  any    `json:"v"`
}

func NewClient(apiKey, accessToken string) *Client {
	return &Client{
		apiKey:        apiKey,
		accessToken:   accessToken,
		endpoint:      URL,
		dialer:        websocket.DefaultDialer,
		subscriptions: make(map[uint32]struct{}),
		modes:         make(map[Mode]map[uint32]struct{}),
	}
}

func (c *Client) Endpoint() string {
	u, _ := url.Parse(URL)
	if c.endpoint != "" {
		u, _ = url.Parse(c.endpoint)
	}
	q := u.Query()
	q.Set("api_key", c.apiKey)
	q.Set("access_token", c.accessToken)
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) OnTick(fn func([]Tick)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onTick = fn
}

func (c *Client) OnError(fn func(error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onError = fn
}

func (c *Client) OnConnect(fn func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onConnect = fn
}

func (c *Client) OnClose(fn func(error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onClose = fn
}

func (c *Client) Connect(ctx context.Context) error {
	for {
		if err := c.connectOnce(ctx); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			c.emitError(err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}

func (c *Client) Close() error {
	c.mu.Lock()
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	if conn == nil {
		return nil
	}
	return conn.Close()
}

func (c *Client) Subscribe(mode Mode, tokens ...uint32) error {
	if mode == "" {
		mode = ModeFull
	}

	c.mu.Lock()
	if c.modes[mode] == nil {
		c.modes[mode] = make(map[uint32]struct{})
	}
	for _, token := range tokens {
		c.subscriptions[token] = struct{}{}
		for _, modeTokens := range c.modes {
			delete(modeTokens, token)
		}
		c.modes[mode][token] = struct{}{}
	}
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		return nil
	}
	if err := c.write(Subscribe(tokens...)); err != nil {
		return err
	}
	return c.write(SetMode(mode, tokens...))
}

func (c *Client) Unsubscribe(tokens ...uint32) error {
	c.mu.Lock()
	for _, token := range tokens {
		delete(c.subscriptions, token)
		for _, modeTokens := range c.modes {
			delete(modeTokens, token)
		}
	}
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		return nil
	}
	return c.write(Unsubscribe(tokens...))
}

func (c *Client) connectOnce(ctx context.Context) error {
	conn, _, err := c.dialer.DialContext(ctx, c.Endpoint(), nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	defer func() {
		_ = conn.Close()
		c.mu.Lock()
		if c.conn == conn {
			c.conn = nil
		}
		c.mu.Unlock()
	}()

	c.emitConnect()
	if err := c.resubscribe(); err != nil {
		c.emitClose(err)
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()
	go c.keepAlive(ctx, conn)

	err = c.readLoop(conn)
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	c.emitClose(err)
	return err
}

func (c *Client) readLoop(conn *websocket.Conn) error {
	for {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if messageType != websocket.BinaryMessage {
			continue
		}

		ticks, err := Parse(payload)
		if err != nil {
			c.emitError(err)
			continue
		}
		if len(ticks) == 0 {
			continue
		}
		c.emitTick(ticks)
	}
}

func (c *Client) keepAlive(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.writeMu.Lock()
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second))
			c.writeMu.Unlock()
			if err != nil {
				_ = conn.Close()
				return
			}
		}
	}
}

func (c *Client) resubscribe() error {
	c.mu.RLock()
	tokens := make([]uint32, 0, len(c.subscriptions))
	for token := range c.subscriptions {
		tokens = append(tokens, token)
	}

	modes := make(map[Mode][]uint32, len(c.modes))
	for mode, modeTokens := range c.modes {
		for token := range modeTokens {
			modes[mode] = append(modes[mode], token)
		}
	}
	c.mu.RUnlock()

	if len(tokens) > 0 {
		if err := c.write(Subscribe(tokens...)); err != nil {
			return err
		}
	}
	for mode, modeTokens := range modes {
		if len(modeTokens) == 0 {
			continue
		}
		if err := c.write(SetMode(mode, modeTokens...)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) write(payload []byte) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return fmt.Errorf("ticker: websocket is not connected")
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return conn.WriteMessage(websocket.TextMessage, payload)
}

func (c *Client) emitTick(ticks []Tick) {
	c.mu.RLock()
	fn := c.onTick
	c.mu.RUnlock()
	if fn != nil {
		fn(ticks)
	}
}

func (c *Client) emitError(err error) {
	c.mu.RLock()
	fn := c.onError
	c.mu.RUnlock()
	if fn != nil {
		fn(err)
	}
}

func (c *Client) emitConnect() {
	c.mu.RLock()
	fn := c.onConnect
	c.mu.RUnlock()
	if fn != nil {
		fn()
	}
}

func (c *Client) emitClose(err error) {
	c.mu.RLock()
	fn := c.onClose
	c.mu.RUnlock()
	if fn != nil {
		fn(err)
	}
}

func Subscribe(tokens ...uint32) []byte {
	body, _ := json.Marshal(Message{Action: "subscribe", Value: tokens})
	return body
}

func Unsubscribe(tokens ...uint32) []byte {
	body, _ := json.Marshal(Message{Action: "unsubscribe", Value: tokens})
	return body
}

func SetMode(mode Mode, tokens ...uint32) []byte {
	body, _ := json.Marshal(Message{Action: "mode", Value: []any{mode, tokens}})
	return body
}
