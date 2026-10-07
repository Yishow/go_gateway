package diagnostics

import (
	"encoding/json"
	"math"
	"slices"
	"strings"
	"unicode/utf8"
)

type template struct{ level, source, message string }

var templates = map[string]template{
	"startup.begin":            {"info", "startup", "Gateway startup has begun."},
	"startup.ready":            {"info", "startup", "The embedded web service is ready."},
	"startup.failed":           {"error", "startup", "Gateway startup failed; review the stable diagnostic code."},
	"startup.config_invalid":   {"error", "startup", "Startup configuration is invalid. Check the configured values."},
	"startup.path_ambiguous":   {"error", "startup", "The database location is ambiguous. Configure an explicit absolute database path."},
	"startup.path_unusable":    {"error", "startup", "The selected data location is unavailable. Check access permissions."},
	"startup.owner_busy":       {"error", "startup", "This database already has an active gateway owner."},
	"startup.log_unavailable":  {"error", "startup", "Managed diagnostic storage is unavailable. Check its location and permissions."},
	"startup.tray_failed":      {"error", "startup", "The desktop control icon could not be initialized."},
	"startup.bind_failed":      {"error", "startup", "The configured listener could not be bound. Check its port."},
	"startup.assets_missing":   {"error", "startup", "Embedded application assets are missing. Use a complete product build."},
	"startup.database_linked":  {"error", "startup", "The selected database has multiple hardlinks, so its WAL and SHM locations are ambiguous. Use one canonical database location."},
	"startup.db_open_failed":   {"error", "startup", "The selected database could not be opened."},
	"startup.db_ping_failed":   {"error", "startup", "The selected database did not pass its connection check."},
	"startup.migration_failed": {"error", "startup", "Database initialization failed. Preserve the data and review diagnostics."},
	"startup.service_failed":   {"error", "startup", "Application services could not be initialized."},
	"startup.pipeline_failed":  {"error", "startup", "The data pipeline could not be initialized."},
	"runtime.started":          {"info", "runtime", "Acquisition runtime has started."},
	"runtime.degraded":         {"warn", "runtime", "Acquisition runtime is degraded. Review runtime status."},
	"runtime.share_degraded":   {"warn", "runtime", "Local sharing is degraded. Review output status."},
	"runtime.group_failed":     {"error", "runtime", "A recording group operation failed. Review recording status."},
	"runtime.stopped":          {"info", "runtime", "Acquisition runtime has stopped."},
	"shutdown.begin":           {"info", "shutdown", "Gateway shutdown has begun."},
	"shutdown.finalizing":      {"info", "shutdown", "Finalizing local diagnostics. Delivery completion is not implied."},
	"shutdown.complete":        {"info", "shutdown", "Gateway shutdown has completed."},
	"shutdown.timeout":         {"warn", "shutdown", "Shutdown is still in progress; completion has not been confirmed."},
	"shutdown.failed":          {"error", "shutdown", "A shutdown phase failed; completion has not been confirmed."},
	"http.access":              {"info", "http", "An HTTP request completed."},
	"http.request":             {"info", "http", "An HTTP request completed."},
	"http.recovered":           {"error", "http", "An HTTP application fault was recovered; raw details were suppressed."},
	"http.recovery":            {"error", "http", "An HTTP application fault was recovered; raw details were suppressed."},
	"http.server_error":        {"error", "http", "The HTTP server reported an error; raw details were suppressed."},
	"raw.suppressed":           {"info", "application", "Unclassified diagnostic content was suppressed."},
}

func project(in Input, routes map[string]struct{}) Event {
	code := in.Code
	if len(code) > 256 {
		code = "raw.suppressed"
	}
	tmpl, ok := templates[code]
	if !ok {
		code = "raw.suppressed"
		tmpl = templates[code]
	}
	e := Event{Code: strings.Clone(code), Level: tmpl.level, Source: tmpl.source, Message: tmpl.message, Truncated: len(in.Code) > 256}
	if code == "raw.suppressed" {
		if source, ok := in.Fields["source"].(string); ok && slices.Contains([]string{"standard", "slog", "http", "runtime"}, source) {
			e.Source = strings.Clone(source)
		}
		return e
	}
	e.Fields = make(map[string]any, 6)
	if strings.HasPrefix(code, "http.") {
		if method, ok := in.Fields["method"].(string); ok && slices.Contains([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"}, method) {
			e.Fields["method"] = strings.Clone(method)
		}
		if route, ok := in.Fields["route"].(string); ok {
			if len(route) <= 1024 {
				if _, allowed := routes[route]; allowed {
					e.Fields["route"] = strings.Clone(route)
				}
			} else {
				e.Truncated = true
			}
		}
		if n, ok := number(in.Fields["status"]); ok && n >= 100 && n <= 599 && math.Trunc(n) == n {
			e.Fields["status"] = int(n)
		}
		if n, ok := number(in.Fields["duration_ms"]); ok && n >= 0 && n <= 86400000 {
			e.Fields["duration_ms"] = n
		}
		if id, ok := in.Fields["request_id"].(string); ok && opaqueID(id) {
			e.Fields["request_id"] = strings.Clone(id)
		}
	}
	if strings.HasPrefix(code, "shutdown.") {
		if phase, ok := in.Fields["phase"].(string); ok && slices.Contains([]string{"startup", "http", "runtime", "pipeline", "share", "connections", "database", "diagnostics", "tray", "owner"}, phase) {
			e.Fields["phase"] = strings.Clone(phase)
		}
	}
	return e
}
func number(v any) (float64, bool) {
	var n float64
	switch x := v.(type) {
	case int:
		n = float64(x)
	case int64:
		n = float64(x)
	case uint64:
		n = float64(x)
	case float64:
		n = x
	default:
		return 0, false
	}
	return n, !math.IsInf(n, 0) && !math.IsNaN(n)
}
func opaqueID(s string) bool {
	if len(s) != 32 && len(s) != 36 {
		return false
	}
	for i, c := range s {
		if len(s) == 36 && (i == 8 || i == 13 || i == 18 || i == 23) {
			if c != '-' {
				return false
			}
			continue
		}
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}
func levelRank(level string) int {
	switch level {
	case "debug":
		return 0
	case "info":
		return 1
	case "warn":
		return 2
	case "error":
		return 3
	default:
		return -1
	}
}

// boundEvent is called only on projected safe content, never on raw input.
func boundEvent(e Event) (event Event, size int) {
	e.Sequence = "18446744073709551615"
	encoded, err := json.Marshal(e)
	if err != nil {
		return boundEvent(Event{Code: "raw.suppressed", Message: templates["raw.suppressed"].message, Truncated: true})
	}
	for len(encoded) > MaxEventBytes {
		e.Truncated = true
		if e.Message != "" {
			e.Message = truncateUTF8(e.Message, max(0, len(e.Message)-(len(encoded)-MaxEventBytes)-16))
		} else {
			e.Fields = nil
		}
		encoded, err = json.Marshal(e)
		if err != nil {
			return boundEvent(Event{Code: "raw.suppressed", Message: templates["raw.suppressed"].message, Truncated: true})
		}
	}
	return e, len(encoded)
}
func truncateUTF8(s string, n int) string {
	if len(s) <= n {
		return strings.Clone(s)
	}
	for n > 0 && !utf8.RuneStart(s[n]) {
		n--
	}
	return strings.Clone(s[:n])
}
