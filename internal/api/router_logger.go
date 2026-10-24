package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

// customLoggerMiddleware 自定義日誌中間件，過濾頻繁的 debug 與狀態輪詢請求
func customLoggerMiddleware() gin.HandlerFunc {
	// 配置 Logger，跳過頻繁輪詢的 API 靜態路徑
	skipPaths := []string{
		strings.TrimRight(config.Get().API.BasePath, "/") + "/system/logs",
		strings.TrimRight(config.Get().API.BasePath, "/") + "/system/logs/stream",
		"/api/v1/debug/packets",
		"/api/v1/debug/logs",
		"/api/v1/test/status",
	}

	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: skipPaths,
		Skip:      isRoutinePollingLogSkipped,
		Formatter: formatHTTPLog,
	})
}

// isRoutinePollingLogSkipped 過濾高頻且狀態正常的定期輪詢心跳日誌，保留異常（>=400）與非 GET 請求
func isRoutinePollingLogSkipped(c *gin.Context) bool {
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		return false
	}
	if c.Writer.Status() >= http.StatusBadRequest {
		return false
	}
	path := c.Request.URL.Path
	if strings.HasSuffix(path, "/delivery") && strings.Contains(path, "/workspace/write-groups/") {
		return true
	}
	if path == "/api/v1/datalink/runtime/stream" || path == "/api/v1/datalink/modbus-share/status" {
		return true
	}
	return false
}

// formatDurationNoDecimals 將耗時四捨五入為整數單位，避免小數點造成日誌不齊
func formatDurationNoDecimals(d time.Duration) string {
	if d <= 0 {
		return "0ms"
	}
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		us := (d + time.Microsecond/2) / time.Microsecond
		if us >= 1000 {
			return "1ms"
		}
		return fmt.Sprintf("%dµs", us)
	}
	if d < 2*time.Second {
		ms := (d + time.Millisecond/2) / time.Millisecond
		return fmt.Sprintf("%dms", ms)
	}
	if d < time.Minute {
		sec := (d + time.Second/2) / time.Second
		return fmt.Sprintf("%ds", sec)
	}
	m := d / time.Minute
	s := (d%time.Minute + time.Second/2) / time.Second
	if s >= 60 {
		m++
		s = 0
	}
	return fmt.Sprintf("%dm%ds", m, s)
}

func formatHTTPLog(param gin.LogFormatterParams) string {
	errorMessage := strings.TrimSpace(param.ErrorMessage)
	if errorMessage != "" {
		errorMessage = " " + errorMessage
	}

	latencyStr := formatDurationNoDecimals(param.Latency)

	if !param.IsOutputColor() {
		return fmt.Sprintf("[HTTP] %s  %-7s  %-60s  %3d  %7s%s\n",
			param.TimeStamp.Format("15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			latencyStr,
			errorMessage,
		)
	}

	// 終端彩色輸出模式
	tagColor := "\033[90m" // 灰階弱化前綴與時間戳，突出主要路徑與狀態
	resetColor := "\033[0m"
	methodColor := param.MethodColor()
	statusColor := param.StatusCodeColor()

	// 耗時色彩梯度：極快(<50ms)灰階、正常預設、慢(>=500ms)黃色、超時/阻塞(>=2s)紅色
	latencyColor := ""
	switch {
	case param.Latency < 50*time.Millisecond:
		latencyColor = "\033[90m"
	case param.Latency >= 2*time.Second:
		latencyColor = "\033[31m"
	case param.Latency >= 500*time.Millisecond:
		latencyColor = "\033[33m"
	}

	errColor := ""
	if errorMessage != "" {
		errColor = "\033[31m"
	}

	return fmt.Sprintf("%s[HTTP] %s%s  %s%-7s%s  %-60s  %s%3d%s  %s%7s%s%s%s%s\n",
		tagColor,
		param.TimeStamp.Format("15:04:05"),
		resetColor,
		methodColor,
		param.Method,
		resetColor,
		param.Path,
		statusColor,
		param.StatusCode,
		resetColor,
		latencyColor,
		latencyStr,
		resetColor,
		errColor,
		errorMessage,
		resetColor,
	)
}
