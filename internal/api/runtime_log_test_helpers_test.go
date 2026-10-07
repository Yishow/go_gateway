package api

import (
	"bufio"
	"context"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/diagnostics"
)

func newRuntimeLogBroker(t *testing.T) *diagnostics.Broker {
	t.Helper()
	b, err := diagnostics.New(diagnostics.Options{Level: "debug"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := b.Close(ctx); err != nil {
			t.Error(err)
		}
	})
	return b
}

func readRuntimeFrame(t *testing.T, r *bufio.Reader) (name, id string, data []byte) {
	t.Helper()
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			t.Fatalf("read stream: %v", err)
		}
		line = strings.TrimSuffix(line, "\n")
		if line == "" {
			return name, id, data
		}
		if value, ok := strings.CutPrefix(line, "event: "); ok {
			name = value
		}
		if value, ok := strings.CutPrefix(line, "id: "); ok {
			id = value
		}
		if value, ok := strings.CutPrefix(line, "data: "); ok {
			data = []byte(value)
		}
	}
}
