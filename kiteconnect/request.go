package kiteconnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) do(ctx context.Context, method, path string, query url.Values, form url.Values, out any) error {
	var body io.Reader
	if len(form) > 0 {
		body = strings.NewReader(form.Encode())
	}
	req, err := c.newRequest(ctx, method, path, query, body)
	if err != nil {
		return err
	}
	if len(form) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return decodeAPIError(resp.StatusCode, payload)
	}
	if out == nil {
		return nil
	}

	var wrapper struct {
		Status    string          `json:"status"`
		Data      json.RawMessage `json:"data"`
		Message   string          `json:"message"`
		ErrorType string          `json:"error_type"`
	}
	if err := json.Unmarshal(payload, &wrapper); err != nil {
		return fmt.Errorf("kiteconnect: decode response: %w", err)
	}
	if wrapper.Status != "success" {
		return &APIError{StatusCode: resp.StatusCode, Status: wrapper.Status, Message: wrapper.Message, ErrorType: wrapper.ErrorType}
	}
	if len(wrapper.Data) == 0 {
		return nil
	}
	if err := json.Unmarshal(wrapper.Data, out); err != nil {
		return fmt.Errorf("kiteconnect: decode data: %w", err)
	}
	return nil
}

func (c *Client) raw(ctx context.Context, method, path string, query url.Values) ([]byte, http.Header, error) {
	req, err := c.newRequest(ctx, method, path, query, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp.Header, decodeAPIError(resp.StatusCode, payload)
	}
	return payload, resp.Header, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Kite-Version", KiteVersion)
	if c.apiKey != "" && c.accessToken != "" {
		req.Header.Set("Authorization", "token "+c.apiKey+":"+c.accessToken)
	}
	return req, nil
}

func decodeAPIError(statusCode int, payload []byte) error {
	apiErr := &APIError{StatusCode: statusCode}
	if len(payload) > 0 && json.Unmarshal(payload, apiErr) == nil && (apiErr.Message != "" || apiErr.ErrorType != "") {
		return apiErr
	}
	apiErr.Message = string(bytes.TrimSpace(payload))
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(statusCode)
	}
	return apiErr
}
