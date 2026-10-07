package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const logTestInstance = "0123456789abcdef0123456789abcdef"

func TestRuntimeLogQueryStrictBounds(t *testing.T) {
	for _, query := range []string{"limit=0", "limit=501", "limit=-1", "limit=bad", "limit=", "limit=1&limit=2", "level=trace", "source=invalid", "q=%FF", "q=" + strings.Repeat("x", 257), "q=x&q=y", "regex=.*", "q=x;limit=2", "before=", "before=bad", "before=" + logTestInstance + ":01", "before=" + logTestInstance + ":-1", "before=" + logTestInstance + ":18446744073709551616", "unknown=value", "q=%", "after=" + logTestInstance + ":1"} {
		t.Run(query[:min(45, len(query))], func(t *testing.T) {
			r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/?"+query, http.NoBody)
			if _, _, err := parseRuntimeLogQuery(r, false); err == nil {
				t.Fatal("invalid query accepted")
			}
		})
	}
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/?level=warn&source=http&q="+url.QueryEscape(strings.Repeat("界", 256))+"&limit=500", http.NoBody)
	q, _, err := parseRuntimeLogQuery(r, false)
	if err != nil || q.Level != "warn" || q.Source != "http" || q.Limit != 500 {
		t.Fatalf("valid query=%+v error=%v", q, err)
	}
}

func TestRuntimeLogStreamCursorConflict(t *testing.T) {
	tests := []struct {
		query  string
		header []string
		want   string
		valid  bool
	}{
		{"", nil, "", true},
		{"after=" + logTestInstance + ":123", nil, logTestInstance + ":123", true},
		{"", []string{logTestInstance + ":123"}, logTestInstance + ":123", true},
		{"after=" + logTestInstance + ":123", []string{logTestInstance + ":123"}, logTestInstance + ":123", true},
		{"after=" + logTestInstance + ":123", []string{logTestInstance + ":124"}, "", false},
		{"", []string{logTestInstance + ":1", logTestInstance + ":1"}, "", false},
		{"after=bad", nil, "", false}, {"after=", nil, "", false}, {"before=" + logTestInstance + ":1", nil, "", false}, {"limit=500", nil, "", false},
	}
	for _, tt := range tests {
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/?"+tt.query, http.NoBody)
		if tt.header != nil {
			r.Header["Last-Event-Id"] = tt.header
		}
		_, after, err := parseRuntimeLogQuery(r, true)
		if (err == nil) != tt.valid || after != tt.want {
			t.Errorf("query=%q header=%v after=%q err=%v", tt.query, tt.header, after, err)
		}
	}
}
