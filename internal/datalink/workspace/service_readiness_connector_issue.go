package workspace

import (
	"strings"

	"go-gateway/internal/datalink/schema"
)

func databaseConnectorReadinessIssue(connector *schema.DatabaseConnector) (ReadinessIssue, bool) {
	if connector == nil || !connector.Enabled {
		return ReadinessIssue{}, false
	}

	code := ""
	switch connector.Status {
	case schema.DatabaseConnectorStatusUnreachable:
		code = "database-connector-unreachable"
	case schema.DatabaseConnectorStatusAuthFailed:
		code = "database-connector-auth-failed"
	case schema.DatabaseConnectorStatusError:
		code = "database-connector-error"
	default:
		return ReadinessIssue{}, false
	}

	message := firstReadinessMessage(connector.LastCheckError, "database connector requires attention")
	return ReadinessIssue{
		Code:     code,
		Severity: ReadinessSeverityWarning,
		Step:     ReadinessStep4,
		Scope:    strings.TrimSpace(connector.ID),
		Message:  message,
	}, true
}

func firstReadinessMessage(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
