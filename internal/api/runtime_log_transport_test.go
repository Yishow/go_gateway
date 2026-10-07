package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

type logTestConn struct {
	net.Conn
	notify func()
}

func (c *logTestConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
}
func (c *logTestConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 40000}
}
func (c *logTestConn) Write(data []byte) (int, error) { c.notify(); return c.Conn.Write(data) }

type logTestListener struct {
	conn      net.Conn
	accepted  bool
	done      chan struct{}
	closeOnce func()
}

func (l *logTestListener) Accept() (net.Conn, error) {
	if !l.accepted {
		l.accepted = true
		return l.conn, nil
	}
	<-l.done
	return nil, net.ErrClosed
}
func (l *logTestListener) Close() error   { l.closeOnce(); return nil }
func (l *logTestListener) Addr() net.Addr { return l.conn.LocalAddr() }

func TestRuntimeLogStalledTransportDeadlineDoesNotLockBroker(t *testing.T) {
	b := newRuntimeLogBroker(t)
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	serverPipe, clientPipe := net.Pipe()
	writing := make(chan struct{})
	conn := &logTestConn{Conn: serverPipe, notify: sync.OnceFunc(func() { close(writing) })}
	listener := &logTestListener{conn: conn, done: make(chan struct{})}
	listener.closeOnce = sync.OnceFunc(func() { close(listener.done) })
	handlerDone := make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { defer close(handlerDone); r.ServeHTTP(w, req) }), ReadHeaderTimeout: time.Second}
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Serve(listener) }()
	t.Cleanup(func() {
		_ = clientPipe.Close()
		_ = server.Close()
		if err := <-serverDone; !errors.Is(err, http.ErrServerClosed) {
			t.Error(err)
		}
	})
	if _, err := clientPipe.Write([]byte("GET /api/v1/system/logs/stream HTTP/1.1\r\nHost: 127.0.0.1:8080\r\n\r\n")); err != nil {
		t.Fatal(err)
	}
	select {
	case <-writing:
	case <-time.After(time.Second):
		t.Fatal("no stream write")
	}
	start := time.Now()
	if !b.Emit(diagnostics.Input{Code: "runtime.started"}) {
		t.Fatal("producer rejected")
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	if err := b.Flush(ctx); err != nil {
		t.Fatalf("writer retained broker lock: %v", err)
	}
	snapshot, err := b.Snapshot(diagnostics.Query{})
	if err != nil || len(snapshot.Records) != 1 {
		t.Fatalf("broker stalled %+v %v", snapshot, err)
	}
	select {
	case <-handlerDone:
	case <-time.After(7 * time.Second):
		t.Fatal("stream exceeded 5-second write deadline")
	}
	if elapsed := time.Since(start); elapsed < 4*time.Second || elapsed > 6*time.Second {
		t.Fatalf("write duration %s", elapsed)
	}
	subscriptions := make([]*diagnostics.Subscription, 0, 8)
	defer func() {
		for _, sub := range subscriptions {
			sub.Close()
		}
	}()
	for range 8 {
		sub, err := b.Subscribe(diagnostics.Query{}, snapshot.Latest)
		if err != nil {
			t.Fatalf("stalled stream leaked slot: %v", err)
		}
		subscriptions = append(subscriptions, sub)
	}
}

func TestRuntimeLogIPv6Socket(t *testing.T) {
	cfg := net.ListenConfig{}
	listener, err := cfg.Listen(t.Context(), "tcp6", "[::1]:0")
	if err != nil {
		t.Skipf("IPv6 unavailable: %v", err)
	}
	b := newRuntimeLogBroker(t)
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	server := &http.Server{Handler: r, ReadHeaderTimeout: time.Second}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	t.Cleanup(func() {
		_ = server.Close()
		if err := <-done; !errors.Is(err, http.ErrServerClosed) {
			t.Error(err)
		}
	})
	client := &http.Client{Transport: &http.Transport{Proxy: nil}, Timeout: time.Second}
	t.Cleanup(client.CloseIdleConnections)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+listener.Addr().String()+"/api/v1/system/logs", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "http://"+listener.Addr().String())
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("IPv6 denied: %d", resp.StatusCode)
	}
}
