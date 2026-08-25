package workspace

import (
	"strings"

	"go-gateway/internal/datalink/schema"
)

func databaseConnectorReadinessIssue(connector *schema.DatabaseConnector) (ReadinessIssue, bool) {
	if connector == nil || !connector.Enabled {
		return ReadinessIssue{}, false
	}

	var code string
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

	return ReadinessIssue{
		Code:     code,
		Severity: ReadinessSeverityWarning,
		Step:     ReadinessStep4,
		Scope:    strings.TrimSpace(connector.ID),
		Message:  databaseConnectorReadinessMessage(code),
	}, true
}

func databaseConnectorReadinessMessage(code string) string {
	switch code {
	case "database-connector-unreachable":
		return "database connector is currently unreachable"
	case "database-connector-auth-failed":
		return "database connector authentication failed"
	case "database-connector-error":
		return "database connector is currently in an error state"
	default:
		return "database connector requires attention"
	}
}

func firstReadinessMessage(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
