package kiteconnect

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginURL(t *testing.T) {
	client := NewClient("abc", "secret")
	got := client.LoginURL()
	want := "https://kite.zerodha.com/connect/login?api_key=abc&v=3"
	if got != want {
		t.Fatalf("LoginURL() = %q, want %q", got, want)
	}
}

func TestGenerateSessionSetsAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/token" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Kite-Version"); got != "3" {
			t.Fatalf("missing kite version header: %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("api_key") != "key" || r.Form.Get("request_token") != "request" || r.Form.Get("checksum") == "" {
			t.Fatalf("unexpected form: %#v", r.Form)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"data": map[string]any{
				"user_id":      "AB1234",
				"access_token": "access",
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", "secret", WithBaseURL(server.URL))
	session, err := client.GenerateSession(context.Background(), "request")
	if err != nil {
		t.Fatal(err)
	}
	if session.AccessToken != "access" || client.AccessToken() != "access" {
		t.Fatalf("access token was not set: %#v", session)
	}
}

func TestGenerateSessionAcceptsRedirectURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if got := r.Form.Get("request_token"); got != "request" {
			t.Fatalf("request token was not extracted: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"data":   map[string]any{"access_token": "access"},
		})
	}))
	defer server.Close()

	client := NewClient("key", "secret", WithBaseURL(server.URL))
	_, err := client.GenerateSession(context.Background(), "https://example.com/callback?request_token=request&action=login&status=success")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateSessionRejectsEmptyRequestToken(t *testing.T) {
	client := NewClient("key", "secret")
	if _, err := client.GenerateSession(context.Background(), " \n\t "); err == nil {
		t.Fatal("expected empty request token error")
	}
}

func TestPlaceOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orders/regular" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "token key:access" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("tradingsymbol") != "INFY" || r.Form.Get("transaction_type") != "BUY" {
			t.Fatalf("unexpected order form: %#v", r.Form)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"data":   map[string]any{"order_id": "123"},
		})
	}))
	defer server.Close()

	client := NewClient("key", "secret", WithBaseURL(server.URL))
	client.SetAccessToken("access")
	order, err := client.PlaceOrder(context.Background(), PlaceOrderParams{
		Exchange:        ExchangeNSE,
		TradingSymbol:   "INFY",
		TransactionType: TransactionBuy,
		Quantity:        1,
		Product:         ProductMIS,
		OrderType:       OrderTypeMarket,
		Validity:        ValidityDay,
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.OrderID != "123" {
		t.Fatalf("unexpected order response: %#v", order)
	}
}
