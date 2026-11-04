package api

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func TestHTTPLogStatusColumnRemainsAligned(t *testing.T) {
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	status := regexp.MustCompile(` {2}(200|503) {2}`)
	paths := []string{
		"/api/v1/datalink/studio-v2/workspace/devices",
		"/api/v1/datalink/studio-v2/workspace/database-catalog?workspace_id=" + strings.Repeat("x", 100),
		"/api/v1/設備/狀態",
	}
	t.Cleanup(gin.DisableConsoleColor)
	for _, color := range []bool{false, true} {
		if color {
			gin.ForceConsoleColor()
		} else {
			gin.DisableConsoleColor()
		}
		columns := make([]int, 0, len(paths))
		for i, path := range paths {
			code := http.StatusOK
			if i == 1 {
				code = http.StatusServiceUnavailable
			}
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, http.NoBody)
			line := ansi.ReplaceAllString(formatHTTPLog(gin.LogFormatterParams{
				TimeStamp: time.Date(2026, time.November, 4, 18, 36, 50, 0, time.UTC),
				Method:    http.MethodGet, Path: path, StatusCode: code, Latency: 19 * time.Second, Request: req,
			}), "")
			match := status.FindStringIndex(line)
			if match == nil {
				t.Fatalf("color=%t: missing status in %q", color, line)
			}
			columns = append(columns, utf8.RuneCountInString(line[:match[0]]))
			if !strings.Contains(line, req.URL.Path) || !strings.Contains(line, "19s") {
				t.Fatalf("color=%t: endpoint or duration missing from %q", color, line)
			}
		}
		for _, column := range columns[1:] {
			if columns[0] != column {
				t.Errorf("color=%t: status columns differ: %v", color, columns)
			}
		}
	}
}
