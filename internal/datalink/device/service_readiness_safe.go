package device

import "strings"

func safeReadinessDiagnostic(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return "diagnostics failed"
	}
	lower := strings.ToLower(message)
	markers := []string{"host=", "dsn=", "password=", "token=", "secret=", "dial tcp "}
	cut := len(message)
	marker := ""
	for _, candidate := range markers {
		if index := strings.Index(lower, candidate); index >= 0 && index < cut {
			cut = index
			marker = candidate
		}
	}
	if cut == len(message) {
		return message
	}
	if marker == "dial tcp " {
		if suffixStart := strings.LastIndex(message, ": "); suffixStart > cut {
			suffix := strings.TrimSpace(message[suffixStart+2:])
			if suffix != "" && !containsReadinessSecret(suffix) {
				return "connection failed: " + suffix
			}
		}
		return "connection diagnostics failed; connection details redacted"
	}
	prefix := strings.TrimRight(strings.TrimSpace(message[:cut]), ":,;- ")
	if prefix == "" {
		return "diagnostics failed; connection details redacted"
	}
	return prefix + ": connection details redacted"
}

func containsReadinessSecret(message string) bool {
	lower := strings.ToLower(message)
	for _, marker := range []string{"host=", "dsn=", "password=", "token=", "secret="} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}
