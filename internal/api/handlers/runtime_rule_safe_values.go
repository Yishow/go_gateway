package handlers

import (
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var runtimeSafeIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9_.-]{0,127}$`)
var runtimeSafeLabelPattern = regexp.MustCompile(`^[\p{L}\p{M}\p{N}_ .,()\[\]·（）°%+-]{1,128}$`)
var runtimeRegisterPattern = regexp.MustCompile(`^(?:[0-9]{1,6}|0[xX][0-9a-fA-F]{1,8}|(?:[sS]?[wWbB])[0-9a-fA-F]{1,8}|[A-Za-z]{1,4}[0-9]{1,6}(?:\.[0-9]{1,2})?)$`)
var runtimeMCHexPattern = regexp.MustCompile(`(?i)^(?:S?[WB])([0-9A-F]{1,8})$`)
var runtimeNetworkPattern = regexp.MustCompile(`(?i)(?:[a-z0-9-]+\.)+[a-z]{2,}`)
var runtimeIPv4Pattern = regexp.MustCompile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)
var runtimeSQLPattern = regexp.MustCompile(`(?i)\b(?:select\s+.+\s+from|insert\s+into|update\s+.+\s+set|delete\s+from|drop\s+table|alter\s+table|create\s+table|grant\s+.+\s+on)\b`)
var runtimeCredentialPattern = regexp.MustCompile(`(?i)\b(?:password|passwd|token|secret|dsn|username|uid|pwd)\s*[:=]`)
var runtimeBracketedIPPattern = regexp.MustCompile(`\[([^\[\]]+)\]`)
var runtimeDSNPairPattern = regexp.MustCompile(`(?i)(?:^|[;\s])((?:data source|initial catalog|network address|integrated security|user id)|server|host|address|addr|database|dbname|driver|trusted_connection|port)\s*=\s*[^;\s]+`)
var runtimeFileURIPattern = regexp.MustCompile(`(?i)\bfile:`)
var runtimePathPattern = regexp.MustCompile(`(?:^|[\s=:'"()\[\]{},;\x60])(?:/|(?:~|\.\.?)[\\/]|[A-Za-z]:[\\/]|\\\\[^\\/\s]+[\\/][^\\/\s]+)`)

func runtimeSafeID(value string) string {
	if !runtimeSafeIDPattern.MatchString(value) || runtimeIPv4Pattern.MatchString(value) || runtimeNetworkPattern.MatchString(value) {
		return ""
	}
	return value
}

func runtimeSafeLabel(value string) string {
	if !runtimeSafeLabelPattern.MatchString(value) || !runtimeSafeSensorText(value) {
		return ""
	}
	return value
}

func runtimeSafeAddress(value string) string {
	if !runtimeRegisterPattern.MatchString(value) {
		return ""
	}
	if parts := runtimeMCHexPattern.FindStringSubmatch(value); parts != nil {
		// Match hsllogic MC parser's signed 32-bit hex address bound.
		if _, err := strconv.ParseUint(parts[1], 16, 31); err != nil {
			return ""
		}
	}
	return value
}

// Normal scalar text is retained; this rejects identifiable unsafe syntax, not words
// such as SELECT or driver that may themselves be legitimate sensor values.
func runtimeSafeSensorText(value string) bool {
	if !utf8.ValidString(value) || len(value) > 256 || runtimeUnsafeNetworkText(value) || runtimeNetworkPattern.MatchString(value) || runtimeIPv4Pattern.MatchString(value) || runtimeSQLPattern.MatchString(value) || runtimeCredentialPattern.MatchString(value) || runtimePathPattern.MatchString(value) || strings.Contains(value, "://") || runtimeFileURIPattern.MatchString(value) {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

// Unstructured backend diagnostics are never operator text, even if they look harmless.
func runtimeSafeDiagnostic(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return "runtime diagnostic unavailable; inspect settings"
}

func runtimeSafeDiagnosticPointer(value *string) *string {
	if value == nil {
		return nil
	}
	safe := runtimeSafeDiagnostic(*value)
	return &safe
}

func runtimeSafeEnum(value string, allowed ...string) string {
	for _, candidate := range allowed {
		if value == candidate {
			return value
		}
	}
	return ""
}

func runtimeSafeAvailabilityReason(value *string) *string {
	if value != nil && *value == "invalid Step 1 configuration" {
		copy := *value
		return &copy
	}
	return runtimeSafeDiagnosticPointer(value)
}

// Recognize numeric bracketed IPv6 and paired DSN keys, not arbitrary colons or
// the words Server/Database in ordinary readings. No network lookup is performed.
func runtimeUnsafeNetworkText(value string) bool {
	for _, match := range runtimeBracketedIPPattern.FindAllStringSubmatch(value, -1) {
		if address, err := netip.ParseAddr(match[1]); err == nil && address.Is6() {
			return true
		}
	}
	keys := map[string]bool{}
	connection := false
	for _, match := range runtimeDSNPairPattern.FindAllStringSubmatch(value, -1) {
		key := strings.ToLower(match[1])
		keys[key] = true
		switch key {
		case "server", "host", "data source", "network address", "address", "addr", "driver":
			connection = true
		}
	}
	return connection && len(keys) >= 2
}
