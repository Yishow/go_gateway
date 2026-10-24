package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestFormatDurationNoDecimals(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{d: 0, want: "0ms"},
		{d: -10 * time.Millisecond, want: "0ms"},
		{d: 500 * time.Nanosecond, want: "500ns"},
		{d: 311541 * time.Nanosecond, want: "312µs"},
		{d: 999600 * time.Nanosecond, want: "1ms"},
		{d: 1540458 * time.Nanosecond, want: "2ms"},
		{d: 11731333 * time.Nanosecond, want: "12ms"},
		{d: 1450 * time.Millisecond, want: "1450ms"},
		{d: 28 * time.Second, want: "28s"},
		{d: 75 * time.Second, want: "1m15s"},
	}

	for _, tc := range tests {
		got := formatDurationNoDecimals(tc.d)
		if got != tc.want {
			t.Errorf("formatDurationNoDecimals(%v) = %q; want %q", tc.d, got, tc.want)
		}
	}
}

func TestFormatHTTPLog_OmitsUserAgentAndNormalizesSpacing(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "PUT", "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", http.NoBody)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	got := formatHTTPLog(gin.LogFormatterParams{
		TimeStamp:  time.Date(2026, time.May, 31, 0, 10, 57, 0, time.UTC),
		Method:     "PUT",
		Path:       req.URL.Path,
		StatusCode: 500,
		Latency:    2219125 * time.Nanosecond,
		Request:    req,
	})

	want := "[HTTP] 00:10:57  PUT      /api/v1/datalink/studio-v2/workspace/source-rules/rule-A      500      2ms\n"
	if got != want {
		t.Fatalf("unexpected log format:\nwant: %q\ngot:  %q", want, got)
	}
	if strings.Contains(got, "Mozilla/5.0") {
		t.Fatalf("expected user agent to be omitted, got %q", got)
	}
}

func TestFormatHTTPLog_PreservesFullQueryString(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/api/v1/datalink/runtime/stream?device_id=dev-1234567890&point_ids=point-1,point-2", http.NoBody)

	got := formatHTTPLog(gin.LogFormatterParams{
		TimeStamp:  time.Date(2026, time.May, 31, 0, 20, 17, 0, time.UTC),
		Method:     "GET",
		Path:       req.URL.Path + "?" + req.URL.RawQuery,
		StatusCode: 200,
		Latency:    3698792 * time.Nanosecond,
		Request:    req,
	})

	want := "[HTTP] 00:20:17  GET      /api/v1/datalink/runtime/stream?device_id=dev-1234567890&point_ids=point-1,point-2  200      4ms\n"
	if got != want {
		t.Fatalf("unexpected query log format:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestFormatHTTPLog_WithColorOutput(t *testing.T) {
	gin.ForceConsoleColor()
	t.Cleanup(func() {
		// Restore default mode
		gin.DisableConsoleColor()
		gin.SetMode(gin.TestMode)
	})

	req := httptest.NewRequestWithContext(t.Context(), "GET", "/api/v1/datalink/studio-v2/workspace", http.NoBody)
	got := formatHTTPLog(gin.LogFormatterParams{
		TimeStamp:  time.Date(2026, time.May, 31, 0, 20, 17, 0, time.UTC),
		Method:     "GET",
		Path:       req.URL.Path,
		StatusCode: 200,
		Latency:    12 * time.Millisecond,
		Request:    req,
	})

	if !strings.Contains(got, "\033[") {
		t.Fatalf("expected ANSI color codes in formatted output, got %q", got)
	}
	if !strings.Contains(got, "/api/v1/datalink/studio-v2/workspace") {
		t.Fatalf("expected path to be preserved, got %q", got)
	}
}

func TestIsRoutinePollingLogSkipped(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
		wantSkip   bool
	}{
		{
			name:       "write-group delivery 200 skipped",
			method:     http.MethodGet,
			path:       "/api/v1/datalink/studio-v2/workspace/write-groups/fde619d6-cd87-4522-9822-79b7cc5d2560/delivery",
			statusCode: http.StatusOK,
			wantSkip:   true,
		},
		{
			name:       "write-group delivery 500 not skipped",
			method:     http.MethodGet,
			path:       "/api/v1/datalink/studio-v2/workspace/write-groups/fde619d6-cd87-4522-9822-79b7cc5d2560/delivery",
			statusCode: http.StatusInternalServerError,
			wantSkip:   false,
		},
		{
			name:       "write-group delivery POST not skipped",
			method:     http.MethodPost,
			path:       "/api/v1/datalink/studio-v2/workspace/write-groups/fde619d6-cd87-4522-9822-79b7cc5d2560/delivery",
			statusCode: http.StatusOK,
			wantSkip:   false,
		},
		{
			name:       "runtime stream 200 skipped",
			method:     http.MethodGet,
			path:       "/api/v1/datalink/runtime/stream",
			statusCode: http.StatusOK,
			wantSkip:   true,
		},
		{
			name:       "modbus share status 200 skipped",
			method:     http.MethodGet,
			path:       "/api/v1/datalink/modbus-share/status",
			statusCode: http.StatusOK,
			wantSkip:   true,
		},
		{
			name:       "standard workspace API not skipped",
			method:     http.MethodGet,
			path:       "/api/v1/datalink/studio-v2/workspace",
			statusCode: http.StatusOK,
			wantSkip:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequestWithContext(t.Context(), tc.method, tc.path, http.NoBody)
			c.Writer.WriteHeader(tc.statusCode)

			got := isRoutinePollingLogSkipped(c)
			if got != tc.wantSkip {
				t.Errorf("isRoutinePollingLogSkipped(%s %s %d) = %v; want %v", tc.method, tc.path, tc.statusCode, got, tc.wantSkip)
			}
		})
	}
}
