package api

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func localLogRequest(t *testing.T) *http.Request {
	t.Helper()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://127.0.0.1:8080/api/v1/system/logs", http.NoBody)
	r.RemoteAddr = "127.0.0.1:45678"
	return r.WithContext(context.WithValue(r.Context(), http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}))
}

func TestLocalLogGuardConnectionBoundary(t *testing.T) {
	tests := []struct {
		name string
		edit func(*http.Request)
		want bool
	}{
		{"IPv4", func(*http.Request) {}, true},
		{"origin", func(r *http.Request) { r.Header.Set("Origin", "http://127.0.0.1:8080") }, true},
		{"localhost", func(r *http.Request) { r.Host = "localhost:8080" }, true},
		{"IPv6", func(r *http.Request) { r.Host = "[::1]:8080"; r.RemoteAddr = "[::1]:5000" }, true},
		{"fetch same-origin", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "same-origin") }, true},
		{"fetch none", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "none") }, true},
		{"TLS", func(r *http.Request) {
			r.TLS = &tls.ConnectionState{}
			r.Header.Set("Origin", "https://127.0.0.1:8080")
		}, true},
		{"LAN peer", func(r *http.Request) { r.RemoteAddr = "192.168.1.2:5000" }, false},
		{"named peer", func(r *http.Request) { r.RemoteAddr = "localhost:5000" }, false},
		{"missing peer port", func(r *http.Request) { r.RemoteAddr = "127.0.0.1" }, false},
		{"invalid peer port", func(r *http.Request) { r.RemoteAddr = "127.0.0.1:abc" }, false},
		{"excessive peer port", func(r *http.Request) { r.RemoteAddr = "127.0.0.1:65536" }, false},
		{"rebinding", func(r *http.Request) { r.Host = "evil.example:8080" }, false},
		{"alias loopback", func(r *http.Request) { r.Host = "127.0.0.2:8080" }, false},
		{"wrong port", func(r *http.Request) { r.Host = "127.0.0.1:80" }, false},
		{"missing port", func(r *http.Request) { r.Host = "127.0.0.1" }, false},
		{"padded port", func(r *http.Request) { r.Host = "127.0.0.1:08080" }, false},
		{"bracket IPv4", func(r *http.Request) { r.Host = "[127.0.0.1]:8080" }, false},
		{"userinfo", func(r *http.Request) { r.Host = "u@localhost:8080" }, false},
		{"foreign origin", func(r *http.Request) { r.Header.Set("Origin", "http://evil.example") }, false},
		{"local alias origin", func(r *http.Request) { r.Header.Set("Origin", "http://localhost:8080") }, false},
		{"wrong scheme", func(r *http.Request) { r.Header.Set("Origin", "https://127.0.0.1:8080") }, false},
		{"null origin", func(r *http.Request) { r.Header.Set("Origin", "null") }, false},
		{"empty origin", func(r *http.Request) { r.Header["Origin"] = []string{""} }, false},
		{"duplicate origin", func(r *http.Request) { r.Header["Origin"] = []string{"http://127.0.0.1:8080", "http://127.0.0.1:8080"} }, false},
		{"same-site no origin", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "same-site") }, false},
		{"cross-site", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "cross-site") }, false},
		{"fetch list", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "same-origin, none") }, false},
		{"fetch whitespace", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", " same-origin") }, false},
		{"fetch empty", func(r *http.Request) { r.Header["Sec-Fetch-Site"] = []string{""} }, false},
		{"fetch duplicate", func(r *http.Request) { r.Header["Sec-Fetch-Site"] = []string{"same-origin", "same-origin"} }, false},
		{"fetch case duplicate", func(r *http.Request) {
			r.Header["Sec-Fetch-Site"] = []string{"same-origin"}
			r.Header["sec-fetch-site"] = []string{"none"}
		}, false},
		{"missing local address", func(r *http.Request) { *r = *r.WithContext(t.Context()) }, false},
		{"LAN destination", func(r *http.Request) {
			*r = *r.WithContext(context.WithValue(r.Context(), http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 8080}))
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := localLogRequest(t)
			tt.edit(r)
			if got := localLogAllowed(r); got != tt.want {
				t.Fatalf("allowed=%v want=%v", got, tt.want)
			}
		})
	}
}

func TestLocalLogGuardRejectsProxyHeaders(t *testing.T) {
	for _, name := range []string{"Forwarded", "X-Forwarded", "X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Proto", "X-Forwarded-Unknown", "X-Real-IP", "x-forwarded-for"} {
		t.Run(name, func(t *testing.T) {
			r := localLogRequest(t)
			r.Header[name] = []string{""}
			if localLogAllowed(r) {
				t.Fatal("proxy header accepted")
			}
		})
	}
}
