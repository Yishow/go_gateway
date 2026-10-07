package handlers

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestWebSocketHandlersCloseWhenServerContextEnds(t *testing.T) {
	for _, route := range []string{"/ws", "/monitor"} {
		t.Run(route, func(t *testing.T) {
			h := NewWebSocketHandler()
			returned := make(chan struct{})
			r := gin.New()
			r.GET("/ws", func(c *gin.Context) { defer close(returned); h.HandleWebSocket(c) })
			r.GET("/monitor", func(c *gin.Context) { defer close(returned); h.HandleMonitorStream(c) })
			base, cancel := context.WithCancel(context.Background())
			defer cancel()
			server := httptest.NewUnstartedServer(r)
			server.Config.BaseContext = func(net.Listener) context.Context { return base }
			server.Start()
			defer server.Close()
			conn, resp, err := websocket.DefaultDialer.DialContext(t.Context(), "ws"+strings.TrimPrefix(server.URL, "http")+route, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusSwitchingProtocols {
				t.Fatalf("status %d", resp.StatusCode)
			}
			cancel()
			select {
			case <-returned:
			case <-time.After(2 * time.Second):
				t.Fatal("hijacked websocket outlived server shutdown context")
			}
		})
	}
}
