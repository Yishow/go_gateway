package diagnostics

import (
	"encoding/json"
	"strconv"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

// Fixed operator text only. Candidate messages, labels, errors and parameters are
// never used as templates. Future localized renderers can key off the same code.
var domainProjectionTemplates = map[DomainEventCode]struct{ level, message string }{
	"domain.source.applied":             {levelInfo, "The source rule was applied to acquisition runtime."},
	"domain.source.rejected":            {levelWarn, "A source-rule runtime update was rejected."},
	"domain.collector.failed":           {levelWarn, "Acquisition failed for this source."},
	"domain.collector.recovered":        {levelInfo, "Acquisition recovered for this source."},
	"domain.mapping.failed":             {levelWarn, "Mapping failed for this source."},
	"domain.mapping.recovered":          {levelInfo, "Mapping recovered for this source."},
	"domain.group.applied":              {levelInfo, "Confirmed reconciliation observation: intake was active for this original revision."},
	"domain.group.paused":               {levelInfo, "Recording intake is paused; historical delivery may remain pending."},
	"domain.group.retiring":             {levelInfo, "The recording group revision is retiring; delivery may remain pending."},
	"domain.delivery.retrying":          {levelWarn, "Delivery is retrying for this group revision."},
	"domain.delivery.blocked":           {levelWarn, "Delivery is blocked for this group revision."},
	"domain.delivery.quarantined":       {levelWarn, "A delivery output was quarantined."},
	"domain.delivery.unknown":           {levelWarn, "The delivery outcome is unknown; commitment is not confirmed."},
	"domain.delivery.first-confirmed":   {levelInfo, "First newly observed SQL commitment in this capture scope; this receipt proves only its own effect."},
	"domain.delivery.recovered":         {levelInfo, "Delivery recovered with a matching SQL commitment."},
	"domain.delivery.summary":           {levelInfo, "Newly observed SQL commitment count in this capture scope; the attached receipt proves only its own effect."},
	"domain.counter-replay.quarantined": {levelWarn, "Accepted input samples are quarantined pending original revision evidence."},
	"domain.recording.refused":          {levelWarn, "Recording was refused by original ownership evidence; pending inputs are not delivery outputs."},
	"domain.quota.recovered":            {levelInfo, "Storage capacity is below the hard threshold; this is not resumed intake or SQL commitment."},
	"domain.quota.refused":              {levelWarn, "Durable intake was refused by the configured storage quota."},
	"domain.repair.completed":           {levelInfo, "The storage repair operation completed; data recovery is not implied."},
	"domain.repair.blocked":             {levelWarn, "The storage repair operation was blocked."},
	"domain.cutover.applied":            {levelInfo, "The recording cutover operation was applied."},
	"domain.cutover.rejected":           {levelWarn, "The recording cutover operation was rejected."},
}

// ProjectDomainEvent is an unmounted, pure projection; it performs no admission or I/O.
// It consumes the shared domain/runtime validators, then builds an exhaustive scalar
// allowlist, never serializing a candidate DTO or copying arbitrary fields. Invalid
// projection is suppressed; the producer's unavailable result remains a separate
// contract outcome, not permission to pass a partially validated event through here.
// InstanceID/Sequence/Timestamp remain unassigned for the future existing admission
// path. This does not verify producer origin, authorize a locator or loosen local access.
func ProjectDomainEvent(data []byte, scope runtimeview.Scope) Event {
	v, err := ParseDomainEventRef(data, scope)
	if err != nil {
		return suppressedDomainProjection(len(data) > MaxEventBytes)
	}
	tmpl, ok := domainProjectionTemplates[v.Code]
	if !ok {
		return suppressedDomainProjection(false)
	}
	f := map[string]any{
		"event_version": v.ContractVersion, "workspace_id": v.Binding.WorkspaceID,
		"workspace_revision": v.Binding.WorkspaceRevision, "entity_revision": v.Binding.EntityRevision,
		"locator_step": v.Locator.Step, "stage": string(v.Stage), "reason": string(v.Reason),
		"action": string(v.Action), "evidence_kind": string(v.Evidence.Kind), "reference_id": v.Evidence.ReferenceID,
	}
	// Only applicable identities survive the parser; no display names or endpoints.
	for _, identity := range []struct {
		key   string
		value *string
	}{
		{"device_id", v.Locator.DeviceID}, {"rule_id", v.Locator.RuleID}, {"mapping_id", v.Locator.MappingID},
		{"group_id", v.Locator.GroupID}, {"group_revision", v.Locator.GroupRevision},
		{"point_id", v.PointID}, {"tag_id", v.TagID}, {"connector_id", v.ConnectorID}, {"connector_revision", v.ConnectorRevision},
	} {
		if identity.value != nil {
			f[identity.key] = *identity.value
		}
	}
	// The parser's closed evidence union and per-code matrix gate each branch.
	if r := v.Evidence.Receipt; r != nil {
		f["record_id"], f["effect_key"], f["payload_digest"] = r.RecordID, r.EffectKey, r.PayloadDigest
		f["committed_at"], f["receipt_provenance"] = r.CommittedAt, r.Provenance
	}
	if i := v.Evidence.Intake; i != nil {
		f["intake_state"] = i.State
	}
	if q := v.Evidence.Quota; q != nil {
		f["quota_scope"], f["quota_state"] = q.Scope, q.State
		f["quota_used_bytes"], f["quota_max_bytes"], f["intake_refused"] = q.UsedBytes, q.MaxBytes, q.IntakeRefused
	}
	if n := v.Evidence.PendingSamples; n != nil {
		f["pending_samples"] = *n
	}
	e := Event{Code: string(v.Code), Message: tmpl.message, Level: tmpl.level, Source: sourceRuntime, Fields: f}
	if !domainProjectionFits(e) {
		return suppressedDomainProjection(true)
	}
	return e
}

func suppressedDomainProjection(truncated bool) Event {
	e := project(Input{Code: codeRawSuppressed}, nil)
	e.Truncated = truncated
	return e
}

// Reserve the existing broker's complete envelope when checking the unchanged
// 8 KiB cap. These maximum-length sizing values never enter the returned event.
// Future admission still owns its own boundEvent/drop accounting and actual clock.
func domainProjectionFits(e Event) bool {
	e.InstanceID = "ffffffffffffffffffffffffffffffff"
	e.Sequence = "18446744073709551615"
	e.Timestamp = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	raw, err := json.Marshal(e)
	return err == nil && len(raw) <= MaxEventBytes
}

// DomainEventSummary is bounded observation metadata, never delivery evidence.
type DomainEventSummary struct{ FirstSeenAt, LastSeenAt, Occurrences string }

// ProjectDomainSummary projects only the three fixed summary scalars onto a validated reference.
func ProjectDomainSummary(data []byte, scope runtimeview.Scope, summary DomainEventSummary) Event {
	e := ProjectDomainEvent(data, scope)
	if e.Code == codeRawSuppressed {
		return e
	}
	first, firstErr := time.Parse(time.RFC3339Nano, summary.FirstSeenAt)
	last, lastErr := time.Parse(time.RFC3339Nano, summary.LastSeenAt)
	n, countErr := strconv.ParseUint(summary.Occurrences, 10, 64)
	if firstErr != nil || lastErr != nil || !domainUTCTime.MatchString(summary.FirstSeenAt) || !domainUTCTime.MatchString(summary.LastSeenAt) || last.Before(first) || countErr != nil || n == 0 || strconv.FormatUint(n, 10) != summary.Occurrences {
		return suppressedDomainProjection(false)
	}
	e.Fields["first_seen_at"], e.Fields["last_seen_at"], e.Fields["occurrences"] = summary.FirstSeenAt, summary.LastSeenAt, summary.Occurrences
	if !domainProjectionFits(e) {
		return suppressedDomainProjection(true)
	}
	return e
}
