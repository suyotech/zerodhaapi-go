package kiteconnect

import (
	"net/http"
	"strings"
)

const (
	DefaultBaseURL  = "https://api.kite.trade"
	DefaultLoginURL = "https://kite.zerodha.com/connect/login"
	KiteVersion     = "3"
)

type Client struct {
	apiKey      string
	apiSecret   string
	accessToken string
	baseURL     string
	loginURL    string
	httpClient  *http.Client
}

type Option func(*Client)

func NewClient(apiKey, apiSecret string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		baseURL:    DefaultBaseURL,
		loginURL:   DefaultLoginURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.baseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

func (c *Client) AccessToken() string {
	return c.accessToken
}
