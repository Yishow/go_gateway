package datalink

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestPostgresHandshakeConnectHandlesCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := (postgresHandshakeConnector{dsn: "host=127.0.0.1 port=1 user=fixture dbname=fixture sslmode=disable"}).Connect(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled connection, got %v", err)
	}
}

// A server that accepts TCP but never answers the PostgreSQL handshake must
// not defeat the connection manager's five-second connection timeout.
func TestPostgresConnectTimesOutDuringHandshake(t *testing.T) {
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	accepted := make(chan net.Conn, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr == nil {
			accepted <- conn
		}
	}()
	address := listener.Addr().(*net.TCPAddr)
	manager := NewDBManager(DBConfig{
		Type: DBTypePostgres,
		DSN:  fmt.Sprintf("host=127.0.0.1 port=%d user=fixture dbname=fixture sslmode=disable", address.Port),
	})
	done := make(chan error, 1)
	go func() { done <- manager.Connect() }()
	var peer net.Conn
	select {
	case peer = <-accepted:
		t.Cleanup(func() { _ = peer.Close() })
	case <-time.After(time.Second):
		t.Fatal("fixture did not receive the connection")
	}
	select {
	case connectErr := <-done:
		if connectErr == nil {
			t.Fatal("an unanswered handshake must not produce a connected database")
		}
	case <-time.After(6 * time.Second):
		_ = peer.Close()
		<-done
		t.Fatal("connection manager did not stop a stalled handshake within its five-second timeout")
	}
}

// Ending the initial ping context must not close a successfully connected
// pooled socket; subsequent metadata requests still need to use it.
func TestPostgresConnectKeepsSuccessfulConnectionUsable(t *testing.T) {
	for _, test := range []struct {
		name   string
		suffix string
	}{
		{name: "default"},
		{name: "explicit dial timeout", suffix: " connect_timeout=2"},
	} {
		t.Run(test.name, func(t *testing.T) {
			testPostgresConnectionRemainsUsable(t, test.suffix)
		})
	}
}

func testPostgresConnectionRemainsUsable(t *testing.T, dsnSuffix string) {
	t.Helper()
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	serverDone := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverDone <- acceptErr
			return
		}
		defer conn.Close()
		serverDone <- servePostgresPingFixture(conn)
	}()
	address := listener.Addr().(*net.TCPAddr)
	manager := NewDBManager(DBConfig{
		Type: DBTypePostgres,
		DSN:  fmt.Sprintf("host=127.0.0.1 port=%d user=fixture dbname=fixture sslmode=disable%s", address.Port, dsnSuffix),
	})
	if err := manager.Connect(); err != nil {
		t.Fatalf("connect to fixture: %v", err)
	}
	t.Cleanup(func() { _ = manager.Close() })
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	if err := manager.DB().PingContext(ctx); err != nil {
		t.Fatalf("pooled connection failed after initial context ended: %v", err)
	}
	if err := manager.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverDone; err != nil {
		t.Fatalf("postgres fixture: %v", err)
	}
}

func servePostgresPingFixture(conn net.Conn) error {
	var length [4]byte
	if _, err := io.ReadFull(conn, length[:]); err != nil {
		return err
	}
	if _, err := io.CopyN(io.Discard, conn, int64(binary.BigEndian.Uint32(length[:]))-4); err != nil {
		return err
	}
	if _, err := conn.Write([]byte{'R', 0, 0, 0, 8, 0, 0, 0, 0, 'Z', 0, 0, 0, 5, 'I'}); err != nil {
		return err
	}
	for {
		var header [5]byte
		if _, err := io.ReadFull(conn, header[:]); err != nil {
			return err
		}
		if _, err := io.CopyN(io.Discard, conn, int64(binary.BigEndian.Uint32(header[1:]))-4); err != nil {
			return err
		}
		switch header[0] {
		case 'Q':
			if _, err := conn.Write([]byte{'I', 0, 0, 0, 4, 'Z', 0, 0, 0, 5, 'I'}); err != nil {
				return err
			}
		case 'X':
			return nil
		default:
			return fmt.Errorf("unexpected postgres fixture message: %q", header[0])
		}
	}
}
