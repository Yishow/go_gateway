package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

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

	want := "[HTTP] 00:10:57 PUT    /api/v1/datalink/studio-v2/workspace/source-rules/rule-A 500 2.219125ms\n"
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

	want := "[HTTP] 00:20:17 GET    /api/v1/datalink/runtime/stream?device_id=dev-1234567890&point_ids=point-1,point-2 200 3.698792ms\n"
	if got != want {
		t.Fatalf("unexpected query log format:\nwant: %q\ngot:  %q", want, got)
	}
}
