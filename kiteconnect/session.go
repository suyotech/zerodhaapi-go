package kiteconnect

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Session struct {
	UserID        string   `json:"user_id"`
	UserName      string   `json:"user_name"`
	UserShortName string   `json:"user_shortname"`
	Email         string   `json:"email"`
	Broker        string   `json:"broker"`
	Exchanges     []string `json:"exchanges"`
	Products      []string `json:"products"`
	OrderTypes    []string `json:"order_types"`
	APIKey        string   `json:"api_key"`
	AccessToken   string   `json:"access_token"`
	PublicToken   string   `json:"public_token"`
	LoginTime     string   `json:"login_time"`
}

func (c *Client) LoginURL() string {
	u, _ := url.Parse(c.loginURL)
	q := u.Query()
	q.Set("v", KiteVersion)
	q.Set("api_key", c.apiKey)
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) GenerateSession(ctx context.Context, requestToken string) (*Session, error) {
	requestToken, err := normalizeRequestToken(requestToken)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("api_key", c.apiKey)
	form.Set("request_token", requestToken)
	form.Set("checksum", checksum(c.apiKey, requestToken, c.apiSecret))

	var session Session
	if err := c.do(ctx, http.MethodPost, "/session/token", nil, form, &session); err != nil {
		return nil, err
	}
	c.accessToken = session.AccessToken
	return &session, nil
}

func (c *Client) InvalidateSession(ctx context.Context) error {
	q := url.Values{}
	q.Set("api_key", c.apiKey)
	q.Set("access_token", c.accessToken)
	return c.do(ctx, http.MethodDelete, "/session/token", q, nil, nil)
}

func checksum(apiKey, requestToken, apiSecret string) string {
	sum := sha256.Sum256([]byte(apiKey + requestToken + apiSecret))
	return hex.EncodeToString(sum[:])
}

func normalizeRequestToken(input string) (string, error) {
	token := strings.TrimSpace(input)
	if token == "" {
		return "", fmt.Errorf("kiteconnect: request token is empty")
	}

	if u, err := url.Parse(token); err == nil {
		if requestToken := strings.TrimSpace(u.Query().Get("request_token")); requestToken != "" {
			return requestToken, nil
		}
	}

	if values, err := url.ParseQuery(token); err == nil {
		if requestToken := strings.TrimSpace(values.Get("request_token")); requestToken != "" {
			return requestToken, nil
		}
	}

	return token, nil
}
