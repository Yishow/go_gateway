package api

import (
	"net/http"
	"net/url"
	"strconv"

	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/diagnostics"
)

// Parsing rejects duplicates, unknown operations and unbounded inputs before any
// ring read. The original Unicode query is normalized exactly once by the broker.
func parseRuntimeLogQuery(r *http.Request, stream bool) (diagnostics.Query, string, error) {
	var q diagnostics.Query
	invalid := func() (diagnostics.Query, string, error) { return diagnostics.Query{}, "", diagnostics.ErrInvalidQuery }
	if len(r.URL.RawQuery) > 4096 {
		return invalid()
	}
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return invalid()
	}
	var after string
	for key, entries := range values {
		if len(entries) != 1 {
			return invalid()
		}
		value := entries[0]
		switch key {
		case "workspace_id", "device_id", "rule_id", "group_id", "stage", "actionable":
			// Independent exact filters, not an ownership tuple or authority.
			// Presence is distinct from omission; never trim or silently drop it.
			if value == "" {
				return invalid()
			}
			switch key {
			case "workspace_id":
				q.WorkspaceID = value
			case "device_id":
				q.DeviceID = value
			case "rule_id":
				q.RuleID = value
			case "group_id":
				q.GroupID = value
			case "stage":
				q.Stage = runtimeview.Stage(value)
			case "actionable":
				if value != "true" && value != "false" {
					return invalid()
				}
				q.Actionable = value == "true"
			}
		case "level":
			q.Level = value
		case "source":
			q.Source = value
		case "q":
			q.Q = value
		case "limit":
			if stream {
				return invalid()
			}
			q.Limit, err = strconv.Atoi(value)
			if err != nil || q.Limit < 1 || q.Limit > 500 {
				return invalid()
			}
		case "before":
			if stream || value == "" || diagnostics.ValidateCursor(value) != nil {
				return invalid()
			}
			q.Before = value
		case "after":
			if !stream || value == "" || diagnostics.ValidateCursor(value) != nil {
				return invalid()
			}
			after = value
		default:
			return invalid()
		}
	}
	lastIDs, present := exactHeaderValues(r.Header, "Last-Event-ID")
	if present {
		if !stream || len(lastIDs) != 1 || lastIDs[0] == "" || diagnostics.ValidateCursor(lastIDs[0]) != nil {
			return invalid()
		}
		if after != "" && after != lastIDs[0] {
			return invalid()
		}
		after = lastIDs[0]
	}
	if _, err := diagnostics.ValidateQuery(q); err != nil {
		return invalid()
	}
	return q, after, nil
}
