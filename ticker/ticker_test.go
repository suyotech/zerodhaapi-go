package ticker

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestModesUseKiteWireValues(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want string
	}{
		{name: "ltp", mode: ModeLTP, want: "ltp"},
		{name: "quote", mode: ModeQuote, want: "quote"},
		{name: "full", mode: ModeFull, want: "full"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got struct {
				Action string `json:"a"`
				Value  []any  `json:"v"`
			}
			if err := json.Unmarshal(SetMode(tt.mode, 408065), &got); err != nil {
				t.Fatal(err)
			}
			if got.Action != "mode" || got.Value[0] != tt.want {
				t.Fatalf("unexpected mode message: %#v", got)
			}
		})
	}
}

func TestClientConnectSubscribesAndReceivesTicks(t *testing.T) {
	upgrader := websocket.Upgrader{}
	received := make(chan string, 2)
	done := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		for i := 0; i < 2; i++ {
			_, payload, err := conn.ReadMessage()
			if err != nil {
				t.Error(err)
				return
			}
			received <- string(payload)
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, ltpPayload(408065, 152345)); err != nil {
			t.Error(err)
			return
		}
		<-done
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewClient("key", "access")
	client.endpoint = "ws" + strings.TrimPrefix(server.URL, "http")

	ticksCh := make(chan []Tick, 1)
	client.OnTick(func(ticks []Tick) {
		ticksCh <- ticks
		cancel()
		close(done)
	})

	if err := client.Subscribe(ModeFull, 408065); err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Connect(ctx)
	}()

	for i := 0; i < 2; i++ {
		select {
		case message := <-received:
			if message == "" {
				t.Fatal("empty websocket command")
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for websocket command")
		}
	}

	select {
	case ticks := <-ticksCh:
		if len(ticks) != 1 || ticks[0].InstrumentToken != 408065 || ticks[0].LastPrice != 1523.45 {
			t.Fatalf("unexpected ticks: %#v", ticks)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ticks")
	}

	select {
	case err := <-errCh:
		if err != context.Canceled {
			t.Fatalf("unexpected connect error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("connect did not stop after context cancellation")
	}
}

func TestClientConnectStopsOnContextCancelWhileIdle(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	client := NewClient("key", "access")
	client.endpoint = "ws" + strings.TrimPrefix(server.URL, "http")

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Connect(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != context.Canceled {
			t.Fatalf("unexpected connect error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("connect did not stop after context cancellation")
	}
}

func ltpPayload(token uint32, price uint32) []byte {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint16(payload[0:2], 1)
	binary.BigEndian.PutUint16(payload[2:4], 8)
	binary.BigEndian.PutUint32(payload[4:8], token)
	binary.BigEndian.PutUint32(payload[8:12], price)
	return payload
}
