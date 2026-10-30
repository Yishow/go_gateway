package diagnostics

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"go-gateway/internal/datalink/runtimeview"

	"golang.org/x/text/cases"
)

// Stable errors never include input, file paths or raw application details.
var (
	ErrInvalidQuery  = errors.New("invalid diagnostic query")
	ErrInvalidCursor = errors.New("invalid diagnostic cursor")
	ErrCapacity      = errors.New("diagnostic stream capacity reached")
	ErrClosed        = errors.New("diagnostic stream closed")
	ErrOverflow      = errors.New("diagnostic stream overflow")
)

// Query bounds a current-process history query; Level is minimum severity.
type Query struct {
	Level, Source, Q, Before               string
	Limit                                  int
	WorkspaceID, DeviceID, RuleID, GroupID string
	Stage                                  runtimeview.Stage
	Actionable                             bool
}

// domainQueryNetwork matches the existing frontend safeRuntimeIdentity network
// exclusion. Other URL, IPv6 and private-path syntax fails domainID itself.
var domainQueryNetwork = regexp.MustCompile(`(?i)(?:[a-z0-9-]+\.)+[a-z]{2,}|(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)

// Gap identifies an explicitly unavailable interval or process reset.
type Gap struct {
	Reason string `json:"reason"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
}

// Metadata reports retained boundaries, delivered progress and separate sink loss.
type Metadata struct {
	InstanceID         string `json:"instance_id"`
	Oldest             string `json:"oldest_cursor"`
	Latest             string `json:"latest_cursor"`
	Progress           string `json:"progress_cursor"`
	CaptureDropped     uint64 `json:"capture_dropped"`
	Evicted            uint64 `json:"evicted"`
	SubscriberOverflow uint64 `json:"subscriber_overflow"`
	FileDropped        uint64 `json:"file_dropped"`
	FileRecoveredGaps  uint64 `json:"file_recovered_gaps"`
	DiskErrors         uint64 `json:"disk_errors"`
	SinkHealth         string `json:"sink_health"`
	Gaps               []Gap  `json:"gaps,omitempty"`
}

// Snapshot contains chronological safe records at an atomic append boundary.
type Snapshot struct {
	Metadata
	Records []Event `json:"records"`
}

func validateQuery(q Query) (Query, error) {
	if q.Level != "" && levelRank(q.Level) < 0 {
		return q, ErrInvalidQuery
	}
	if q.Limit == 0 {
		q.Limit = 200
	}
	if q.Limit < 1 || q.Limit > 500 || !utf8.ValidString(q.Q) || utf8.RuneCountInString(q.Q) > 256 {
		return q, ErrInvalidQuery
	}
	switch q.Source {
	case "", sourceStartup, sourceRuntime, sourceShutdown, sourceHTTP, "standard", "slog", "application":
	default:
		return q, ErrInvalidQuery
	}
	for _, id := range []string{q.WorkspaceID, q.DeviceID, q.RuleID, q.GroupID} {
		if id != "" && (!domainID.MatchString(id) || domainQueryNetwork.MatchString(id)) {
			return q, ErrInvalidQuery
		}
	}
	if q.Stage != "" && runtimeview.ParseStage(string(q.Stage)) != q.Stage {
		return q, ErrInvalidQuery
	}
	q.WorkspaceID, q.DeviceID = strings.Clone(q.WorkspaceID), strings.Clone(q.DeviceID)
	q.RuleID, q.GroupID = strings.Clone(q.RuleID), strings.Clone(q.GroupID)
	q.Stage = runtimeview.Stage(strings.Clone(string(q.Stage)))
	q.Q = strings.Clone(cases.Fold().String(q.Q))
	q.Level = strings.Clone(q.Level)
	q.Source = strings.Clone(q.Source)
	q.Before = strings.Clone(q.Before)
	return q, nil
}

// ValidateQuery validates external limits without exposing normalized text whose
// Unicode case folding could expand past the original character limit.
func ValidateQuery(q Query) (Query, error) { _, err := validateQuery(q); return q, err }
func matches(q Query, e Event) bool {
	if q.Level != "" && levelRank(e.Level) < levelRank(q.Level) || q.Source != "" && q.Source != e.Source {
		return false
	}
	if !matchesDomainQuery(q, e) {
		return false
	}
	if q.Q == "" {
		return true
	}
	fold := cases.Fold()
	if strings.Contains(fold.String(e.Message), q.Q) || strings.Contains(fold.String(e.Code), q.Q) {
		return true
	}
	for _, v := range e.Fields {
		if s, ok := v.(string); ok && strings.Contains(fold.String(s), q.Q) {
			return true
		}
	}
	return false
}

// matchesDomainQuery consumes already projected scalar metadata. It does not
// reconstruct a runtime authority, resolve names, or infer present fault state.
func matchesDomainQuery(q Query, e Event) bool {
	if q.WorkspaceID == "" && q.DeviceID == "" && q.RuleID == "" && q.GroupID == "" && q.Stage == "" && !q.Actionable {
		return true
	}
	if e.Code == "domain.unavailable" && e.Source == sourceRuntime {
		// The fixed observation is actionable but has no attributable identity.
		return q.WorkspaceID == "" && q.DeviceID == "" && q.RuleID == "" && q.GroupID == "" &&
			(q.Stage == "" || q.Stage == runtimeview.StageUnavailable)
	}
	if _, fixed := domainProjectionTemplates[DomainEventCode(e.Code)]; !fixed || e.Source != sourceRuntime || e.Fields["event_version"] != 1 {
		return false
	}
	stage, stageOK := e.Fields["stage"].(string)
	reason, reasonOK := e.Fields["reason"].(string)
	action, actionOK := e.Fields["action"].(string)
	if !stageOK || stage == "" || runtimeview.ParseStage(stage) != runtimeview.Stage(stage) ||
		!reasonOK || runtimeview.ParseReason(reason) != runtimeview.Reason(reason) ||
		!actionOK || runtimeview.ActionForReason(runtimeview.Reason(reason)) != runtimeview.Action(action) {
		return false
	}
	if q.Stage != "" && string(q.Stage) != stage || q.Actionable && action == string(runtimeview.ActionNone) {
		return false
	}
	for _, identity := range []struct{ key, value string }{
		{"workspace_id", q.WorkspaceID}, {"device_id", q.DeviceID}, {"rule_id", q.RuleID}, {"group_id", q.GroupID},
	} {
		if identity.value != "" && e.Fields[identity.key] != identity.value {
			return false
		}
	}
	return true
}

func parseCursor(cursor string) (instanceID string, number uint64, err error) {
	instance, sequence, ok := strings.Cut(cursor, ":")
	if !ok || !opaqueID(instance) || sequence == "" || len(sequence) > 20 || len(sequence) > 1 && sequence[0] == '0' {
		return "", 0, ErrInvalidCursor
	}
	for _, c := range sequence {
		if c < '0' || c > '9' {
			return "", 0, ErrInvalidCursor
		}
	}
	n, err := strconv.ParseUint(sequence, 10, 64)
	if err != nil {
		return "", 0, ErrInvalidCursor
	}
	return instance, n, nil
}

// ValidateCursor checks opaque-instance and decimal-sequence wire syntax.
func ValidateCursor(cursor string) error { _, _, err := parseCursor(cursor); return err }
