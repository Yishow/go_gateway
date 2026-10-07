package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadinessRequiresOwnedEmbeddedResponse(t *testing.T) {
	expected := []byte("owned-embedded-index")
	for _, correct := range []bool{true, false} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if correct {
				_, _ = w.Write(expected)
			} else {
				_, _ = w.Write([]byte("another application"))
			}
		}))
		err := probeOwnHTTP(t.Context(), server.Listener.Addr().String(), expected)
		server.Close()
		if (err == nil) != correct {
			t.Fatalf("correct=%v err=%v", correct, err)
		}
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := probeOwnHTTP(ctx, "127.0.0.1:1", expected); err == nil {
		t.Fatal("canceled startup became ready")
	}
}
