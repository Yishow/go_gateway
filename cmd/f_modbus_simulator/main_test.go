package main

import (
	"net"
	"testing"

	"go-gateway/internal/virtual/memory"
)

func TestStartSimulatorBindsLoopback(t *testing.T) {
	server, err := startSimulator(memory.NewMemoryBank(128), 0)
	if err != nil {
		t.Fatalf("start simulator: %v", err)
	}
	defer server.Stop()
	if got := server.BindAddress(); got != "127.0.0.1" {
		t.Fatalf("bind address = %q, want 127.0.0.1", got)
	}
	conn, err := (&net.Dialer{}).DialContext(t.Context(), "tcp", server.Address())
	if err != nil {
		t.Fatalf("dial loopback simulator: %v", err)
	}
	_ = conn.Close()
}
