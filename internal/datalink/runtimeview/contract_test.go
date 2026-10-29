package runtimeview

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Literal test evidence only; never imported by production or supplied to any reader.
const contractLiteral = `{
 "contract_version":1,"workspace_id":"W","workspace_revision":"W1","device_id":"D",
 "observed_at":"2026-10-27T10:00:00Z","freshness":"fresh",
 "source":{"rule_id":"R","point_id":"P","tag_id":"T","mapping_id":"M","saved_revision":"R2","effective_revision":"R1","state":"saved-only","proof":"unavailable",
 "saved":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"R2"},"mapping_revision":"M2","enabled":true,"address":"40001","protocol_function":"FC03","read_type":"int16","result_type":"float64","polling":{"group_id":"poll","interval_ms":1000,"source":"polling-group"},"transforms":[{"kind":"scale","order":null,"multiplier":0.1,"offset":0,"target_type":null}]},"effective":null},
 "groups":[{"group_id":"G","saved_revision":"B","applied_effective_at":null,"applied_revision":"A","effective_revision":"A","as_of":"2026-10-27T10:00:00Z","state":"effective","proof":"applied-snapshot-and-layout",
 "saved":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"B"},"recording_profile":null,"allowed_lateness_seconds":null,"ends_at":null,"effective_at":null,"enabled":false,"destination":{"connector_id":"C2","connector_revision":"8","alias":"Meter archive","kind":"postgres","schema":"public","table":"meter_new"},"members":[{"device_id":"D","point_id":"P","tag_id":"T","column":"value","source_value_type":null,"circuit_id":null},{"device_id":"E","point_id":"Q","tag_id":"U","column":"other","source_value_type":null,"circuit_id":null}],"layout":"managed_wide_profile_v1","storage_mode":"plain","interval_seconds":30,"reduction":"mean","write_mode":"append","incomplete_policy":"partial"},
 "effective":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"A"},"recording_profile":null,"allowed_lateness_seconds":null,"ends_at":null,"effective_at":"2026-10-27T09:00:00Z","enabled":true,"destination":{"connector_id":"C1","connector_revision":"7","alias":"Meter archive","kind":"postgres","schema":"public","table":"meter_old"},"members":[{"device_id":"D","point_id":"P","tag_id":"T","column":"value","source_value_type":null,"circuit_id":null},{"device_id":"E","point_id":"Q","tag_id":"U","column":"other","source_value_type":null,"circuit_id":null}],"layout":"managed_wide_profile_v1","storage_mode":"plain","interval_seconds":15,"reduction":"mean","write_mode":"append","incomplete_policy":"partial"},
 "delivery":{"oldest_pending_scope":"group-all-revisions","oldest_pending_seconds":null,"current":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"A"},"intake":null,"totals_scope":"group","shared":true,"stages":{"collecting":1,"queued":3,"retrying":0,"blocked":0,"quarantined":0,"unknown":0,"sql_committed":0,"skipped":0},"last_commit":null,"sources":[],"source_state":"unavailable","next_due_at":null,"next_retry_at":null,"attempt_duration_ms":null,"quota":null,"storage":null},"backlog":[{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"A0"},"destination":{"connector_id":"C1","connector_revision":"7","alias":null,"kind":"postgres","schema":"public","table":"old_table"},"pending":3,"reasons":["counter-replay-proof-missing"]}],"counter_replay_quarantines":null,"test_write":null}}],
 "locator":{"step":4,"workspace_id":"W","device_id":"D","rule_id":null,"mapping_id":null,"group_id":"G","group_revision":"A0"}
}`

var contractScope = Scope{"W", "W1", "D"}

func fixture(t *testing.T) RuleEvidenceView {
	t.Helper()
	var v RuleEvidenceView
	if err := json.Unmarshal([]byte(contractLiteral), &v); err != nil {
		t.Fatal(err)
	}
	return v
}
func parseFixture(t *testing.T, v RuleEvidenceView) (*RuleEvidenceView, error) {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return ParseRuleEvidence(raw, contractScope)
}
func mustParse(t *testing.T, v RuleEvidenceView) *RuleEvidenceView {
	t.Helper()
	got, err := parseFixture(t, v)
	if err != nil {
		t.Fatalf("expected valid literal contract, got %v", err)
	}
	return got
}
func str(v string) *string   { return &v }
func num(v float64) *float64 { return &v }

func TestNamedLiteralContractCases(t *testing.T) {
	t.Run("saved draft B differs from applied A", func(t *testing.T) {
		v := mustParse(t, fixture(t))
		g := v.Groups[0]
		if g.Effective.Destination.Table != "meter_old" || *g.Effective.IntervalSeconds != 15 || g.Saved.Destination.Table != "meter_new" || *g.Saved.IntervalSeconds != 30 {
			t.Fatal("draft replaced applied facts")
		}
		if *v.Source.Saved.Polling != (PollingEvidence{"poll", 1000, "polling-group"}) || *v.Source.Saved.Transforms[0].Multiplier != 0.1 || v.Source.State != "saved-only" || v.Source.Effective != nil {
			t.Fatal("poll/transform provenance lost")
		}
	})
	t.Run("queued without receipt not committed", func(t *testing.T) {
		d := mustParse(t, fixture(t)).Groups[0].Delivery.Current
		if d.Stages.Queued != 3 || d.Stages.SQLCommitted != 0 || d.LastCommit != nil || !d.Shared || d.TotalsScope != "group" || d.SourceState != "unavailable" {
			t.Fatal("queue or shared totals became source receipt")
		}
	})
	t.Run("historical backlog retains old destination revision", func(t *testing.T) {
		b := mustParse(t, fixture(t)).Groups[0].Delivery.Backlog[0]
		if b.Binding.EntityRevision != "A0" || b.Destination.ConnectorID != "C1" || b.Destination.ConnectorRevision != "7" || b.Destination.Table != "old_table" || b.Pending != 3 {
			t.Fatal("backlog relabelled")
		}
	})
	t.Run("unknown timing stays null", func(t *testing.T) {
		v := mustParse(t, fixture(t))
		d := v.Groups[0].Delivery.Current
		if d.NextDueAt != nil || d.NextRetryAt != nil || d.AttemptDurationMS != nil {
			t.Fatal("invented timing")
		}
		raw, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"next_due_at", "next_retry_at", "attempt_duration_ms"} {
			if !strings.Contains(string(raw), `"`+key+`":null`) {
				t.Fatal("nullable key omitted", key)
			}
		}
	})
	t.Run("missing effective proof unavailable", func(t *testing.T) {
		v := fixture(t)
		v.Source.State = "effective"
		got := mustParse(t, v)
		if got.Source.State != "unavailable" || got.Source.Effective != nil || *got.Source.Saved.Transforms[0].Multiplier != 0.1 {
			t.Fatal("draft claimed effective")
		}
		v = fixture(t)
		v.Groups[0].Proof = "unavailable"
		got = mustParse(t, v)
		if got.Groups[0].State != "unavailable" || got.Groups[0].Effective != nil || got.Groups[0].Delivery.Current != nil || got.Groups[0].Delivery.Backlog[0].Pending != 3 {
			t.Fatal("unproven active retained or historical evidence lost")
		}
	})
	t.Run("foreign locator mixed scope refused", func(t *testing.T) {
		tests := []struct {
			mutate func(*RuleEvidenceView)
			reason Reason
		}{
			{func(v *RuleEvidenceView) { v.Locator.WorkspaceID = "foreign" }, ReasonForeign},
			{func(v *RuleEvidenceView) { v.Groups[0].Effective.Binding.WorkspaceRevision = "W2" }, ReasonStale},
			{func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Binding.EntityRevision = "B" }, ReasonStale},
			{func(v *RuleEvidenceView) { v.DeviceID = "E" }, ReasonForeign},
		}
		for _, tc := range tests {
			v := fixture(t)
			tc.mutate(&v)
			_, err := parseFixture(t, v)
			if !IsContractError(err, tc.reason) {
				t.Fatalf("want %s got %v", tc.reason, err)
			}
		}
	})
	t.Run("hostile secret fields rejected", func(t *testing.T) {
		for _, injection := range []string{`"host":"secret-host",`, `"password":"secret-password",`, `"dsn":"secret-dsn",`, `"path":"/private/db",`, `"raw_error":"driver-secret",`, `"params":{"nested":{"token":"secret-token"}},`} {
			raw := strings.Replace(contractLiteral, `"alias":"Meter archive",`, injection+`"alias":"Meter archive",`, 1)
			_, err := ParseRuleEvidence([]byte(raw), contractScope)
			if !IsContractError(err, ReasonUnavailable) || strings.Contains(err.Error(), "secret") {
				t.Fatalf("unsafe payload accepted/exposed: %v", err)
			}
		}
	})
	t.Run("unknown enums safe", func(t *testing.T) {
		for _, pair := range [][2]string{{`"freshness":"fresh"`, `"freshness":"future"`}, {`"state":"effective"`, `"state":"future"`}, {`"kind":"scale"`, `"kind":"formula"`}, {`"reduction":"mean"`, `"reduction":"magic"`}, {`"counter-replay-proof-missing"`, `"raw-error-secret"`}} {
			_, err := ParseRuleEvidence([]byte(strings.Replace(contractLiteral, pair[0], pair[1], 1)), contractScope)
			if !IsContractError(err, ReasonUnavailable) {
				t.Fatalf("unknown enum accepted: %v", err)
			}
		}
	})
	t.Run("runtime commit separate from test write readback cleanup", func(t *testing.T) {
		v := fixture(t)
		v.Groups[0].Delivery.TestWrite = &TestWriteEvidence{OperationID: "test1", Provenance: "test-write", WrittenAt: str("2026-10-27T10:00:00Z"), ReadbackAt: str("2026-10-27T10:00:01Z"), CleanupAt: str("2026-10-27T10:00:02Z")}
		got := mustParse(t, v)
		if got.Groups[0].Delivery.Current.Stages.SQLCommitted != 0 || got.Groups[0].Delivery.Current.LastCommit != nil {
			t.Fatal("test receipt promoted")
		}
	})
	t.Run("historical locator is selection only", func(t *testing.T) {
		l := mustParse(t, fixture(t)).Locator
		if l.Step != 4 || *l.GroupID != "G" || *l.GroupRevision != "A0" {
			t.Fatal("historical selection lost")
		}
	})
}

func TestLiteralParityWithFrontend(t *testing.T) {
	raw, err := os.ReadFile("../../../frontend/tests/unit/services/runtime-rule-evidence-parser.test.ts")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(mustParse(t, fixture(t)), func() *RuleEvidenceView { v := fixture(t); return &v }()) {
		t.Fatal("Go literal projection drift")
	}
	if !strings.Contains(string(raw), contractLiteral) {
		t.Fatal("Go/TypeScript literal contract cases drifted")
	}
}

func TestValidatedSerializer(t *testing.T) {
	raw, err := MarshalRuleEvidence(fixture(t), contractScope)
	if err != nil || !strings.Contains(string(raw), `"next_retry_at":null`) {
		t.Fatalf("safe serialization missing: %v", err)
	}
	v := fixture(t)
	v.Groups[0].Saved.Destination.Alias = str("postgres://private-host/password")
	if _, err := MarshalRuleEvidence(v, contractScope); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("unsafe alias serialized", err)
	}
}
func TestUnprovenEffectiveTransformUnavailable(t *testing.T) {
	v := fixture(t)
	s := v.Source
	facts := *s.Saved
	facts.Binding.EntityRevision = "R1"
	facts.Transforms = []TransformStep{{Kind: "unsupported"}}
	s.State = "effective"
	s.Proof = "immutable-source-facts"
	s.Effective = &facts
	got := mustParse(t, v)
	if got.Source.State != "unavailable" || got.Source.Effective != nil {
		t.Fatal("unsupported transform presented as proven effective")
	}
}

func TestEffectiveHandoffLegacyNeverClaimsCurrent(t *testing.T) {
	raw := strings.Replace(contractLiteral, `"effective_revision":"A","as_of":"2026-10-27T10:00:00Z",`, "", 1)
	v, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil {
		t.Fatal(err)
	}
	if v.Groups[0].State != "unavailable" || v.Groups[0].Effective != nil || v.Groups[0].Delivery.Current != nil {
		t.Fatal("legacy requested revision claimed effective without as-of proof")
	}
}
func TestEffectiveHandoffScheduledAndQuarantine(t *testing.T) {
	raw := strings.Replace(contractLiteral, `"applied_revision":"A"`, `"applied_revision":"B"`, 1)
	raw = strings.Replace(raw, `"counter_replay_quarantines":null,"test_write":null`, `"counter_replay_quarantines":[{"workspace_id":"W","group_id":"G","group_revision":"A0","pending_samples":7,"reason":"counter-replay-proof-missing","action":"preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance"}],"test_write":null`, 1)
	got, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil {
		t.Fatal("scheduled B must not replace effective A:", err)
	}
	if got.Groups[0].Effective.Binding.EntityRevision != "A" || got.Groups[0].Delivery.Current.Stages.SQLCommitted != 0 || got.Groups[0].Delivery.Backlog[0].Pending != 3 {
		t.Fatal("accepted input quarantine became output")
	}
}

func TestEffectiveHandoffAsOfAndHistoricalSafety(t *testing.T) {
	for _, field := range []string{"revision", "asof"} {
		v := fixture(t)
		if field == "revision" {
			v.Groups[0].EffectiveRevision = nil
		} else {
			v.Groups[0].AsOf = nil
		}
		g := mustParse(t, v).Groups[0]
		if g.State != "unavailable" || g.Effective != nil || g.Delivery.Current != nil || g.Delivery.Backlog[0].Pending != 3 {
			t.Fatal("missing proof inferred current")
		}
	}
	for _, at := range []string{"2026-10-27T10:00:00.000000001Z", "2026-10-27T09:59:59Z"} {
		v := fixture(t)
		v.Groups[0].AsOf = str(at)
		if _, err := parseFixture(t, v); !IsContractError(err, ReasonStale) {
			t.Fatal("mixed as-of accepted", err)
		}
	}
	v := fixture(t)
	v.Groups[0].Effective.EffectiveAt = str("2026-10-27T10:00:00.000000001Z")
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonStale) {
		t.Fatal("future effective accepted", err)
	}
	v = fixture(t)
	v.Groups[0].EffectiveRevision = str("B")
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonStale) {
		t.Fatal("mixed effective revision accepted", err)
	}
	v = fixture(t)
	v.Groups[0].Effective.EndsAt = v.ObservedAt
	if mustParse(t, v).Groups[0].Delivery.Current.Intake != nil {
		t.Fatal("ended snapshot invented intake")
	}
	v.Groups[0].Delivery.Current.Intake = &IntakeEvidence{GroupRevision: "A", State: "active"}
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("ended active contradiction accepted", err)
	}
	v = fixture(t)
	v.Groups[0].Effective.Enabled = false
	if mustParse(t, v).Groups[0].Delivery.Current.Intake != nil {
		t.Fatal("disabled snapshot invented intake")
	}
	v.Groups[0].Delivery.Current.Intake = &IntakeEvidence{GroupRevision: "A", State: "not_running"}
	if mustParse(t, v).Groups[0].Delivery.Current.Intake.State != "not_running" {
		t.Fatal("historical intake lost")
	}
}

func TestCounterReplayQuarantineSafety(t *testing.T) {
	q := CounterReplayQuarantine{WorkspaceID: "W", GroupID: "G", GroupRevision: "A0", PendingSamples: 7, Reason: ReasonCounterMissing, Action: "preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance"}
	for _, reason := range []Reason{ReasonCounterMissing, ReasonCounterConflict, ReasonCounterUnsupported} {
		v := fixture(t)
		q.Reason = reason
		qs := []CounterReplayQuarantine{q}
		v.Groups[0].Delivery.CounterReplayQuarantines = &qs
		g := mustParse(t, v).Groups[0]
		if (*g.Delivery.CounterReplayQuarantines)[0].PendingSamples != 7 || g.Delivery.Current.Stages.Quarantined != 0 {
			t.Fatal("input count changed output count")
		}
	}
	for _, tc := range []struct {
		name   string
		edit   func(*CounterReplayQuarantine)
		reason Reason
	}{
		{"workspace", func(q *CounterReplayQuarantine) { q.WorkspaceID = "foreign" }, ReasonForeign},
		{"group", func(q *CounterReplayQuarantine) { q.GroupID = "other" }, ReasonForeign},
		{"revision", func(q *CounterReplayQuarantine) { q.GroupRevision = "" }, ReasonUnavailable},
		{"unknown", func(q *CounterReplayQuarantine) { q.Reason = "future" }, ReasonUnavailable},
		{"wrong-reason", func(q *CounterReplayQuarantine) { q.Reason = ReasonFinalizedOverlap }, ReasonUnavailable},
		{"action", func(q *CounterReplayQuarantine) { q.Action = "retry" }, ReasonUnavailable},
		{"negative", func(q *CounterReplayQuarantine) { q.PendingSamples = -1 }, ReasonUnavailable},
		{"unsafe-number", func(q *CounterReplayQuarantine) { q.PendingSamples = 9007199254740992 }, ReasonUnavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := fixture(t)
			entry := q
			tc.edit(&entry)
			qs := []CounterReplayQuarantine{entry}
			v.Groups[0].Delivery.CounterReplayQuarantines = &qs
			if _, err := parseFixture(t, v); !IsContractError(err, tc.reason) {
				t.Fatal(err)
			}
		})
	}
	v := fixture(t)
	if mustParse(t, v).Groups[0].Delivery.CounterReplayQuarantines != nil {
		t.Fatal("missing treated empty")
	}
	empty := []CounterReplayQuarantine{}
	v.Groups[0].Delivery.CounterReplayQuarantines = &empty
	if mustParse(t, v).Groups[0].Delivery.CounterReplayQuarantines == nil {
		t.Fatal("observed empty lost")
	}
}

func TestEffectiveHandoffSavedCRequestedBEffectiveAThenB(t *testing.T) {
	v := fixture(t)
	g := &v.Groups[0]
	g.SavedRevision = str("C")
	g.Saved.Binding.EntityRevision = "C"
	g.AppliedRevision = str("B")
	g.AppliedEffectiveAt = str("2026-10-27T10:00:01Z")
	got := mustParse(t, v).Groups[0]
	if *got.SavedRevision != "C" || *got.AppliedRevision != "B" || *got.EffectiveRevision != "A" || got.Effective.Destination.Table != "meter_old" {
		t.Fatal("three revision identities conflated")
	}
	v.ObservedAt = g.AppliedEffectiveAt
	g.AsOf = v.ObservedAt
	g.EffectiveRevision = str("B")
	g.Effective.Binding.EntityRevision = "B"
	g.Effective.EffectiveAt = g.AppliedEffectiveAt
	g.Delivery.Current.Binding.EntityRevision = "B"
	if *mustParse(t, v).Groups[0].EffectiveRevision != "B" {
		t.Fatal("boundary B not represented")
	}
}
func TestCounterReplayQuarantineMissingAndUnsafeWireFields(t *testing.T) {
	raw := strings.Replace(contractLiteral, `"counter_replay_quarantines":null`, `"counter_replay_quarantines":[{"workspace_id":"W","group_id":"G","group_revision":"A0","pending_samples":7,"reason":"counter-replay-proof-missing","action":"preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance"}]`, 1)
	for _, change := range [][2]string{{`"pending_samples":7,`, ""}, {`"reason":"counter-replay-proof-missing",`, ""}, {`"pending_samples":7`, `"pending_samples":null`}, {`"pending_samples":7`, `"pending_samples":7,"host":"secret-host"`}, {`"pending_samples":7`, `"pending_samples":7,"database":"/private/db"`}} {
		if _, err := ParseRuleEvidence([]byte(strings.Replace(raw, change[0], change[1], 1)), contractScope); !IsContractError(err, ReasonUnavailable) {
			t.Fatal("unsafe or missing evidence accepted", err)
		}
	}
}

func TestEffectiveHandoffEquivalentTimeAndDuplicateQuarantine(t *testing.T) {
	v := fixture(t)
	v.Groups[0].AsOf = str("2026-10-27T10:00:00.000000000Z")
	if mustParse(t, v).Groups[0].State != "effective" {
		t.Fatal("equal timestamp lost")
	}
	q := CounterReplayQuarantine{WorkspaceID: "W", GroupID: "G", GroupRevision: "A0", PendingSamples: 0, Reason: ReasonCounterMissing, Action: "preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance"}
	qs := []CounterReplayQuarantine{q, q}
	v.Groups[0].Delivery.CounterReplayQuarantines = &qs
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("duplicate input totals accepted", err)
	}
}

func TestPersistedCircuitIdentityRoundTrip(t *testing.T) {
	raw := strings.Replace(contractLiteral, `"circuit_id":null`, `"circuit_id":"C1"`, -1)
	got, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil {
		t.Fatal("persisted frozen C1 rejected", err)
	}
	encoded, err := MarshalRuleEvidence(*got, contractScope)
	if err != nil || !strings.Contains(string(encoded), `"circuit_id":"C1"`) {
		t.Fatal("persisted C1 lost", err)
	}
}

func TestPersistedCircuitIdentityMissingUnsafeAndScope(t *testing.T) {
	raw := strings.ReplaceAll(contractLiteral, `,"circuit_id":null`, "")
	got, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil || got.Groups[0].Effective.Members[0].CircuitID != nil {
		t.Fatal("legacy circuit invented", err)
	}
	for _, wire := range []string{`"circuit_id":""`, `"circuit_id":"/private/C1"`, `"circuit_id":"postgres://secret"`, `"circuit_id":1`, `"circuit_id":{"name":"C1"}`, `"circuit_id":"C1","circuit_name":"secret"`} {
		if _, err := ParseRuleEvidence([]byte(strings.Replace(contractLiteral, `"circuit_id":null`, wire, 1)), contractScope); !IsContractError(err, ReasonUnavailable) {
			t.Fatal("unsafe circuit accepted", err)
		}
	}
	v := fixture(t)
	g := &v.Groups[0]
	g.SavedRevision = str("C")
	g.Saved.Binding.EntityRevision = "C"
	g.Saved.Members = g.Saved.Members[1:]
	g.Effective.Members[0].CircuitID = str("C1")
	result := mustParse(t, v).Groups[0]
	if *result.Effective.Members[0].CircuitID != "C1" || result.Saved.Members[0].CircuitID != nil {
		t.Fatal("effective C1 inferred from newer saved")
	}
	g.Effective.Binding.WorkspaceID = "foreign"
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonForeign) {
		t.Fatal("foreign circuit facts accepted", err)
	}
	g.Effective.Binding.WorkspaceID = "W"
	g.Effective.Binding.EntityRevision = "C"
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonStale) {
		t.Fatal("mixed revision circuit facts accepted", err)
	}
}

func TestSafeSummaryFieldsRoundTrip(t *testing.T) {
	raw := strings.Replace(contractLiteral, `"recording_profile":null,"allowed_lateness_seconds":null,"ends_at":null,"effective_at":"2026-10-27T09:00:00Z"`, `"ends_at":null,"effective_at":"2026-10-27T09:00:00Z","recording_profile":{"kind":"instantaneous","interval":15,"unit":"second"},"allowed_lateness_seconds":2`, 1)
	raw = strings.ReplaceAll(raw, `"source_value_type":null`, `"source_value_type":"float64"`)
	got, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil {
		t.Fatal("safe summary rejected", err)
	}
	encoded, err := MarshalRuleEvidence(*got, contractScope)
	if err != nil || !strings.Contains(string(encoded), `"allowed_lateness_seconds":2`) || !strings.Contains(string(encoded), `"source_value_type":"float64"`) {
		t.Fatal("summary lost", err)
	}
}

func TestSafeSummaryProfileConsistency(t *testing.T) {
	for _, tc := range []struct {
		kind, unit, reduction string
		interval, seconds     int64
	}{{"instantaneous", "second", "mean", 15, 16}, {"instantaneous", "minute", "mean", 1, 60}, {"instantaneous", "second", "cumulative_delta", 15, 15}, {"instantaneous", "second", "mean", 301, 301}, {"cumulative", "minute", "mean", 1, 60}, {"cumulative", "second", "cumulative_delta", 1, 1}} {
		v := fixture(t)
		f := v.Groups[0].Effective
		f.RecordingProfile = &RecordingProfile{tc.kind, tc.interval, tc.unit}
		f.IntervalSeconds = &tc.seconds
		f.Reduction = tc.reduction
		if _, err := parseFixture(t, v); !IsContractError(err, ReasonUnavailable) {
			t.Fatal("contradictory profile accepted", tc, err)
		}
	}
}

func TestSafeSummaryBoundsNullableAndUnknown(t *testing.T) {
	for _, tc := range []struct {
		kind, unit, reduction string
		interval, seconds     int64
	}{{"instantaneous", "second", "latest", 1, 1}, {"instantaneous", "second", "mean", 300, 300}, {"cumulative", "minute", "cumulative_delta", 1, 60}, {"cumulative", "minute", "cumulative_delta", 600, 36000}} {
		v := fixture(t)
		f := v.Groups[0].Effective
		f.RecordingProfile = &RecordingProfile{tc.kind, tc.interval, tc.unit}
		f.IntervalSeconds = &tc.seconds
		f.Reduction = tc.reduction
		if mustParse(t, v).Groups[0].Effective.RecordingProfile.Interval != tc.interval {
			t.Fatal("profile changed")
		}
	}
	raw := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(contractLiteral, `"recording_profile":null,`, ""), `"allowed_lateness_seconds":null,`, ""), `"source_value_type":null,`, "")
	got, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil {
		t.Fatal(err)
	}
	f := got.Groups[0].Effective
	if f.RecordingProfile != nil || f.AllowedLatenessSeconds != nil || f.Members[0].SourceValueType != nil {
		t.Fatal("missing summary fabricated")
	}
	for _, change := range [][2]string{{`"source_value_type":null`, `"source_value_type":"future"`}, {`"allowed_lateness_seconds":null`, `"allowed_lateness_seconds":-1`}, {`"allowed_lateness_seconds":null`, `"allowed_lateness_seconds":1.5`}, {`"recording_profile":null`, `"recording_profile":{"kind":"future","interval":15,"unit":"second"}`}, {`"recording_profile":null`, `"recording_profile":{"kind":"instantaneous","interval":0,"unit":"second"}`}, {`"recording_profile":null`, `"recording_profile":{"kind":"instantaneous","interval":15,"unit":"hour"}`}, {`"recording_profile":null`, `"recording_profile":{"kind":"instantaneous","interval":15,"unit":"second","physical_unit":"secret"}`}} {
		if _, err := ParseRuleEvidence([]byte(strings.Replace(contractLiteral, change[0], change[1], 1)), contractScope); !IsContractError(err, ReasonUnavailable) {
			t.Fatal("unsafe summary accepted", err)
		}
	}
	v := fixture(t)
	f = v.Groups[0].Saved
	f.RecordingProfile = &RecordingProfile{"instantaneous", 15, "second"}
	f.IntervalSeconds = nil
	f.Reduction = "unavailable"
	if mustParse(t, v).Groups[0].Saved.IntervalSeconds != nil {
		t.Fatal("profile fabricated normalized evidence")
	}
}
