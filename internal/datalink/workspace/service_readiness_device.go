package workspace

import (
	"strings"

	"go-gateway/internal/datalink/schema"
)

func deviceReadinessIssue(deviceID string, readiness *schema.DeviceReadiness) (ReadinessIssue, bool) {
	if readiness == nil || readiness.ActivationAllowed {
		return ReadinessIssue{}, false
	}

	code := "device-activation-blocked"
	switch {
	case readiness.ConnectStatus == schema.ReadinessStageStatusSuccess && readiness.ProbeStatus != schema.ReadinessStageStatusSuccess:
		code = "device-probe-required"
	case readiness.ConnectStatus != schema.ReadinessStageStatusSuccess:
		code = "device-connect-required"
	}

	messageInputs := append([]string{}, readiness.BlockingReasons...)
	messageInputs = append(messageInputs, readiness.AvailabilityReason, "device is not ready for activation")
	return ReadinessIssue{
		Code:     code,
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep1,
		Scope:    strings.TrimSpace(deviceID),
		Message:  safeReadinessMessage(firstReadinessMessage(messageInputs...)),
	}, true
}

// safeReadinessMessage preserves an actionable category while removing
// endpoint, credential, and transport details from public readiness payloads.
func safeReadinessMessage(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return "device is not ready for activation"
	}
	lower := strings.ToLower(message)
	cut := len(message)
	for _, marker := range []string{"host=", "dsn=", "password=", "token=", "secret=", "dial tcp "} {
		if index := strings.Index(lower, marker); index >= 0 && index < cut {
			cut = index
		}
	}
	if cut == len(message) {
		return message
	}
	prefix := strings.TrimRight(strings.TrimSpace(message[:cut]), ":,;-")
	if prefix == "" {
		return "device diagnostics failed; connection details redacted"
	}
	return prefix + ": connection details redacted"
}
