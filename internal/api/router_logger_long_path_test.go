package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// Long query strings should not push the request status and elapsed time off the line.
func TestHTTPLogLongCandidateTargetRemainsReadable(t *testing.T) {
	path := "/api/v1/datalink/source-rules/rule-11111111-1111-4111-8111-111111111111/candidates"
	target := path + "?workspace_id=22222222-2222-4222-8222-222222222222&expected_workspace_revision=33333333-3333-4333-8333-333333333333&revision_id=44444444-4444-4444-8444-444444444444"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, http.NoBody)
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	t.Cleanup(gin.DisableConsoleColor)
	for _, color := range []bool{false, true} {
		if color {
			gin.ForceConsoleColor()
		} else {
			gin.DisableConsoleColor()
		}
		line := ansi.ReplaceAllString(formatHTTPLog(gin.LogFormatterParams{
			TimeStamp: time.Date(2026, time.November, 4, 18, 18, 37, 0, time.UTC),
			Method:    http.MethodGet, Path: target, StatusCode: http.StatusOK,
			Latency: 19 * time.Second, Request: req,
		}), "")
		if !strings.Contains(line, path) || !strings.Contains(line, "200") || !strings.Contains(line, "19s") {
			t.Errorf("color=%t: expected endpoint, status and duration in %q", color, line)
		}
		if n := utf8.RuneCountInString(strings.TrimSuffix(line, "\n")); n > 140 {
			t.Errorf("color=%t: log width %d exceeds 140 characters", color, n)
		}
	}
}

func TestHTTPLoggerKeepsRequestIntactWithLongQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gin.DisableConsoleColor()
	var log bytes.Buffer
	var receivedURI string
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: &log, Formatter: formatHTTPLog}))
	router.GET("/api/v1/datalink/source-rules/:id/candidates", func(c *gin.Context) {
		receivedURI = c.Request.RequestURI
		c.Status(http.StatusOK)
	})
	target := "/api/v1/datalink/source-rules/rule-11111111-1111-4111-8111-111111111111/candidates?workspace_id=" + strings.Repeat("1", 100)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, http.NoBody))
	if response.Code != http.StatusOK || receivedURI != target {
		t.Fatalf("request changed: status=%d, URI=%q", response.Code, receivedURI)
	}
	line := log.String()
	if !strings.Contains(line, "/candidates?") || !strings.Contains(line, "200") || utf8.RuneCountInString(strings.TrimSuffix(line, "\n")) > 140 {
		t.Fatalf("request log should retain endpoint and status with bounded width: %q", line)
	}
}

func TestHTTPLogLongTargetRetainsValidUTF8(t *testing.T) {
	gin.DisableConsoleColor()
	target := "/api/v1/" + strings.Repeat("記憶體", 60)
	line := formatHTTPLog(gin.LogFormatterParams{
		TimeStamp: time.Date(2026, time.November, 4, 18, 18, 37, 0, time.UTC),
		Method:    http.MethodGet, Path: target, StatusCode: http.StatusOK,
		Latency: time.Millisecond,
	})
	if !utf8.ValidString(line) || utf8.RuneCountInString(strings.TrimSuffix(line, "\n")) > 140 {
		t.Fatalf("long UTF-8 target produced an invalid or excessively wide line: %q", line)
	}
}
