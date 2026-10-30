package diagnostics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

type projectionCase struct {
	code, stage, kind, reason, action, level, message string
	step                                              int
}

var projectionCases = []projectionCase{
	{"domain.source.applied", "accepted", "source-runtime-transition", "none", "none", "info", "The source rule was applied to acquisition runtime.", 2},
	{"domain.source.rejected", "blocked", "source-runtime-transition", "unavailable", "refresh", "warn", "A source-rule runtime update was rejected.", 2},
	{"domain.collector.failed", "acquisition", "acquisition-result", "unavailable", "refresh", "warn", "Acquisition failed for this source.", 1},
	{"domain.collector.recovered", "acquisition", "acquisition-result", "none", "none", "info", "Acquisition recovered for this source.", 1},
	{"domain.mapping.failed", "blocked", "mapping-result", "unavailable", "refresh", "warn", "Mapping failed for this source.", 3},
	{"domain.mapping.recovered", "acquisition", "mapping-result", "none", "none", "info", "Mapping recovered for this source.", 3},
	{"domain.group.applied", "collecting", "intake-status", "none", "none", "info", "Confirmed reconciliation observation: intake was active for this original revision.", 4},
	{"domain.group.paused", "blocked", "intake-status", "none", "none", "info", "Recording intake is paused; historical delivery may remain pending.", 4},
	{"domain.group.retiring", "collecting", "intake-status", "none", "none", "info", "The recording group revision is retiring; delivery may remain pending.", 4},
	{"domain.delivery.retrying", "retrying", "durable-transition", "unavailable", "refresh", "warn", "Delivery is retrying for this group revision.", 4},
	{"domain.delivery.blocked", "blocked", "durable-transition", "unavailable", "refresh", "warn", "Delivery is blocked for this group revision.", 4},
	{"domain.delivery.quarantined", "quarantined", "durable-transition", "unavailable", "refresh", "warn", "A delivery output was quarantined.", 4},
	{"domain.delivery.unknown", "unknown", "durable-transition", "unavailable", "refresh", "warn", "The delivery outcome is unknown; commitment is not confirmed.", 4},
	{"domain.delivery.first-confirmed", "sql-committed", "runtime-receipt", "none", "none", "info", "First newly observed SQL commitment in this capture scope; this receipt proves only its own effect.", 4},
	{"domain.delivery.recovered", "sql-committed", "runtime-receipt", "none", "none", "info", "Delivery recovered with a matching SQL commitment.", 4},
	{"domain.delivery.summary", "sql-committed", "runtime-receipt", "none", "none", "info", "Newly observed SQL commitment count in this capture scope; the attached receipt proves only its own effect.", 4},
	{"domain.counter-replay.quarantined", "quarantined", "counter-replay-inputs", "counter-replay-proof-missing", "inspect-immutable-evidence", "warn", "Accepted input samples are quarantined pending original revision evidence.", 4},
	{"domain.quota.refused", "blocked", "quota-refusal", "unavailable", "refresh", "warn", "Durable intake was refused by the configured storage quota.", 4},
	{"domain.repair.completed", "accepted", "operation-result", "none", "none", "info", "The storage repair operation completed; data recovery is not implied.", 4},
	{"domain.repair.blocked", "blocked", "operation-result", "unsupported", "inspect-settings", "warn", "The storage repair operation was blocked.", 4},
	{"domain.cutover.applied", "accepted", "operation-result", "none", "none", "info", "The recording cutover operation was applied.", 4},
	{"domain.cutover.rejected", "blocked", "operation-result", "recording-cutover-finalized-overlap", "inspect-immutable-evidence", "warn", "The recording cutover operation was rejected.", 4},
}

func projectionRef(tc projectionCase) DomainEventRef {
	v := domainFixture()
	v.Code, v.Stage, v.Reason, v.Action = DomainEventCode(tc.code), runtimeview.Stage(tc.stage), runtimeview.Reason(tc.reason), runtimeview.Action(tc.action)
	v.Evidence.Kind = DomainEvidenceKind(tc.kind)
	if tc.step < 4 {
		v.Locator = runtimeview.EntityLocator{Step: tc.step, WorkspaceID: "W", DeviceID: domainText("D")}
		v.ConnectorID, v.ConnectorRevision = nil, nil
	}
	switch tc.step {
	case 1:
		v.PointID = domainText("P")
	case 2:
		v.Locator.RuleID = domainText("S")
	case 3:
		v.Locator.MappingID, v.PointID, v.TagID = domainText("M"), domainText("P"), domainText("T")
	}
	switch tc.kind {
	case "runtime-receipt":
		v.Evidence.Receipt = domainCommitFixture().Evidence.Receipt
	case "intake-status":
		state := "active"
		if tc.code == "domain.group.paused" {
			state = "not_running"
		}
		if tc.code == "domain.group.retiring" {
			state = "retiring"
		}
		v.Evidence.Intake = &runtimeview.IntakeEvidence{GroupRevision: "R1", State: state}
	case "quota-refusal":
		v.Evidence.Quota = &runtimeview.QuotaEvidence{Scope: "global", State: "hard_limit", UsedBytes: 95, MaxBytes: 100, IntakeRefused: true}
	case "counter-replay-inputs":
		v.Evidence.PendingSamples = domainText("18446744073709551615")
	}
	return v
}

func projectionJSON(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func requireProjectionSuppressed(t *testing.T, e Event) {
	t.Helper()
	if e.Code != "raw.suppressed" || e.Message != "Unclassified diagnostic content was suppressed." || e.Level != "info" || e.Source != "application" || len(e.Fields) != 0 {
		t.Fatalf("want fixed, unattributed suppression: %#v", e)
	}
	if e.InstanceID != "" || e.Sequence != "" || !e.Timestamp.IsZero() {
		t.Fatal("projector fabricated admission identity")
	}
}

func TestDomainProjectionFixedCodesAndScalarFields(t *testing.T) {
	for _, tc := range projectionCases {
		t.Run(tc.code, func(t *testing.T) {
			e := ProjectDomainEvent(projectionJSON(t, projectionRef(tc)), domainScope())
			if e.Code != tc.code || e.Message != tc.message || e.Level != tc.level || e.Source != "runtime" || e.Truncated {
				t.Fatalf("fixed template mismatch: %#v", e)
			}
			if e.InstanceID != "" || e.Sequence != "" || !e.Timestamp.IsZero() {
				t.Fatal("projection must leave admission envelope unassigned")
			}
			want := map[string]any{"event_version": 1, "workspace_id": "W", "workspace_revision": "W1", "entity_revision": "R1", "locator_step": tc.step, "stage": tc.stage, "reason": tc.reason, "action": tc.action, "evidence_kind": tc.kind, "reference_id": "E"}
			switch tc.step {
			case 1:
				want["device_id"], want["point_id"] = "D", "P"
			case 2:
				want["device_id"], want["rule_id"] = "D", "S"
			case 3:
				want["device_id"], want["mapping_id"], want["point_id"], want["tag_id"] = "D", "M", "P", "T"
			case 4:
				want["group_id"], want["group_revision"], want["connector_id"], want["connector_revision"] = "G", "R1", "C", "C1"
			}
			switch tc.kind {
			case "runtime-receipt":
				want["record_id"], want["effect_key"], want["payload_digest"], want["committed_at"], want["receipt_provenance"] = "record-E", "E", "digest-E", "2026-10-27T09:00:00Z", "runtime"
			case "counter-replay-inputs":
				want["pending_samples"] = "18446744073709551615"
			case "quota-refusal":
				want["quota_scope"], want["quota_state"], want["quota_used_bytes"], want["quota_max_bytes"], want["intake_refused"] = "global", "hard_limit", int64(95), int64(100), true
			case "intake-status":
				states := map[string]string{"domain.group.applied": "active", "domain.group.paused": "not_running", "domain.group.retiring": "retiring"}
				want["intake_state"] = states[tc.code]
			}
			if !reflect.DeepEqual(e.Fields, want) {
				t.Fatalf("scalar allowlist mismatch\ngot: %#v\nwant: %#v", e.Fields, want)
			}
		})
	}
}

func TestDomainProjectionHostileAndSecretBearingFields(t *testing.T) {
	base := projectionJSON(t, domainFixture())
	for _, extra := range []string{
		`"raw_error":"password_fixture","credential":{"password":"nested_password_fixture"}`,
		`"message":"driver_exception_fixture"`, `"params":{"token":"token_fixture"}`,
		`"connector":{"dsn":"dsn_fixture","username":"username_fixture"}`,
		`"payload":{"value":123}`, `"endpoint":"https://private-host:5432"`,
		`"path":"/home/private/database"`, `"sql":"SELECT secret_fixture"`,
		`"host":"private-host","port":5432`, `"sample_value":42`,
		`"instance_id":"forged","sequence":"99","timestamp":"2026-10-27T10:00:00Z"`,
	} {
		raw := []byte(strings.TrimSuffix(string(base), "}") + "," + extra + "}")
		e := ProjectDomainEvent(raw, domainScope())
		requireProjectionSuppressed(t, e)
		encoded := string(projectionJSON(t, e))
		for _, secret := range []string{"password_fixture", "token_fixture", "driver_exception_fixture", "dsn_fixture", "username_fixture", "secret_fixture", "private-host", "/home/private", "forged"} {
			if strings.Contains(encoded, secret) {
				t.Fatalf("projection retained %q", secret)
			}
		}
	}
	for _, hostile := range []string{`<img src=x onerror=alert(1)>`, "postgres://user:password@host/db", "../private/path", "SELECT * FROM private", "host:5432", "line\nraw-error", strings.Repeat("x", 129)} {
		v := domainFixture()
		v.Evidence.ReferenceID = hostile
		requireProjectionSuppressed(t, ProjectDomainEvent(projectionJSON(t, v), domainScope()))
	}
	for _, unknown := range []DomainEventCode{"future.domain.code", "password_fixture", "http.access", "runtime.group_failed", "domain.unavailable"} {
		v := domainFixture()
		v.Code = unknown
		requireProjectionSuppressed(t, ProjectDomainEvent(projectionJSON(t, v), domainScope()))
	}
}

func TestDomainProjectionValidatesEvidenceAndScope(t *testing.T) {
	for _, mutate := range []func(*DomainEventRef){
		func(v *DomainEventRef) { v.Evidence.Kind = "buffer-ack" },
		func(v *DomainEventRef) { v.Evidence.Receipt = nil },
		func(v *DomainEventRef) { v.Evidence.Receipt.GroupRevision = "R0" },
		func(v *DomainEventRef) { v.Evidence.Receipt.CommittedAt = "2026-02-30T09:00:00Z" },
		func(v *DomainEventRef) { v.Action = "retry" },
		func(v *DomainEventRef) { v.Reason = "raw_driver_error" },
		func(v *DomainEventRef) { v.Stage = "readback-verified" },
		func(v *DomainEventRef) { v.Stage = "cleanup-completed" },
	} {
		v := domainCommitFixture()
		mutate(&v)
		requireProjectionSuppressed(t, ProjectDomainEvent(projectionJSON(t, v), domainScope()))
	}
	for _, scope := range []runtimeview.Scope{{}, {WorkspaceID: "foreign", WorkspaceRevision: "W1"}, {WorkspaceID: "W", WorkspaceRevision: "W0"}} {
		requireProjectionSuppressed(t, ProjectDomainEvent(projectionJSON(t, domainCommitFixture()), scope))
	}
	e := ProjectDomainEvent(projectionJSON(t, domainCommitFixture()), runtimeview.Scope{WorkspaceID: "W", WorkspaceRevision: "W1"})
	if e.Fields["group_revision"] != "R1" || e.Fields["effect_key"] != "E" {
		t.Fatalf("historical G/R1 evidence lost: %#v", e)
	}
	if _, attributed := e.Fields["device_id"]; attributed {
		t.Fatal("group receipt falsely attributed to device")
	}
}

func TestDomainProjectionSizeAndIndependentCopies(t *testing.T) {
	base := projectionJSON(t, domainFixture())
	exact := append(append([]byte{}, base...), []byte(strings.Repeat(" ", MaxEventBytes-len(base)))...)
	if e := ProjectDomainEvent(exact, domainScope()); e.Code != "domain.delivery.retrying" {
		t.Fatal("valid exact-size candidate refused")
	}
	tooLarge := append(exact, ' ')
	e := ProjectDomainEvent(tooLarge, domainScope())
	requireProjectionSuppressed(t, e)
	if !e.Truncated {
		t.Fatal("oversized candidate must report truncation without retaining text")
	}
	for _, tc := range projectionCases {
		e := ProjectDomainEvent(projectionJSON(t, projectionRef(tc)), domainScope())
		e.InstanceID, e.Sequence = strings.Repeat("f", 32), "18446744073709551615"
		e.Timestamp = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
		if len(projectionJSON(t, e)) > 8192 {
			t.Fatal("projected record exceeded existing event cap")
		}
	}
	first := ProjectDomainEvent(base, domainScope())
	first.Fields["workspace_id"] = "modified"
	for i := range base {
		base[i] = ' '
	}
	second := ProjectDomainEvent(projectionJSON(t, domainFixture()), domainScope())
	if second.Fields["workspace_id"] != "W" || first.Fields["reference_id"] != "E" {
		t.Fatal("projection output aliased input or another result")
	}
}

func TestDomainProjectionDoesNotActivateLegacyAdmission(t *testing.T) {
	for _, tc := range projectionCases {
		e := ProjectDomainEvent(projectionJSON(t, projectionRef(tc)), domainScope())
		legacy := project(Input{Code: e.Code, Fields: e.Fields}, nil)
		requireProjectionSuppressed(t, legacy)
	}
}

func TestDomainProjectionTemplateGrowthFailsClosed(t *testing.T) {
	const code = DomainEventCode("domain.delivery.retrying")
	prior := domainProjectionTemplates[code]
	defer func() { domainProjectionTemplates[code] = prior }()
	delete(domainProjectionTemplates, code)
	e := ProjectDomainEvent(projectionJSON(t, domainFixture()), domainScope())
	requireProjectionSuppressed(t, e)
	if e.Truncated {
		t.Fatal("missing fixed template is not truncation")
	}
	grown := prior
	grown.message = strings.Repeat("safe fixed text ", 8192)
	domainProjectionTemplates[code] = grown
	e = ProjectDomainEvent(projectionJSON(t, domainFixture()), domainScope())
	requireProjectionSuppressed(t, e)
	if !e.Truncated {
		t.Fatal("oversized template must be suppressed, never partially labeled delivery")
	}
	if domainProjectionFits(Event{Fields: map[string]any{"invalid-internal-scalar": make(chan int)}}) {
		t.Fatal("encoding failure must fail closed")
	}
}

func TestDomainProjectionMaximumIdentifiersAndMalformedInput(t *testing.T) {
	v := domainCommitFixture()
	maxID := strings.Repeat("x", 128)
	v.Binding = runtimeview.Binding{WorkspaceID: maxID, WorkspaceRevision: maxID, EntityRevision: maxID}
	v.Locator.WorkspaceID, v.Locator.GroupID, v.Locator.GroupRevision = maxID, &maxID, &maxID
	v.ConnectorID, v.ConnectorRevision, v.Evidence.ReferenceID = &maxID, &maxID, maxID
	v.Evidence.Receipt = &runtimeview.Receipt{GroupRevision: maxID, ConnectorRevision: maxID, RecordID: maxID, EffectKey: maxID, PayloadDigest: maxID, CommittedAt: "9999-12-31T23:59:59.999999999Z", Provenance: "runtime"}
	e := ProjectDomainEvent(projectionJSON(t, v), runtimeview.Scope{WorkspaceID: maxID, WorkspaceRevision: maxID})
	if e.Code != "domain.delivery.first-confirmed" || e.Truncated || e.Fields["record_id"] != maxID {
		t.Fatalf("maximum safe identity should remain intact: %#v", e)
	}
	e.InstanceID, e.Sequence = strings.Repeat("f", 32), "18446744073709551615"
	e.Timestamp = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	if len(projectionJSON(t, e)) > 8192 {
		t.Fatal("max identifiers plus admission envelope overflowed")
	}
	base := string(projectionJSON(t, domainFixture()))
	for _, raw := range []string{"", "null", "[]", "{", base + "{}", strings.Replace(base, `"step":4`, `"step":4,"step":4`, 1), strings.Replace(base, `"reference_id":"E"`, `"reference_id":{"password":"password_fixture"}`, 1)} {
		requireProjectionSuppressed(t, ProjectDomainEvent([]byte(raw), domainScope()))
	}
}

func TestDomainProjectionHostileIdentityCannotBecomeAScalar(t *testing.T) {
	v := domainFixture()
	v.Evidence.ReferenceID = "<script>hostile_fixture</script>"
	e := ProjectDomainEvent(projectionJSON(t, v), domainScope())
	requireProjectionSuppressed(t, e)
	if strings.Contains(string(projectionJSON(t, e)), "hostile_fixture") {
		t.Fatal("hostile identity retained")
	}
}

func TestDomainProjectionSummaryMetadata(t *testing.T) {
	data := projectionJSON(t, domainFixture())
	meta := DomainEventSummary{FirstSeenAt: "2026-10-27T10:00:00.000000001Z", LastSeenAt: "2026-10-27T10:00:00.999999999Z", Occurrences: "18446744073709551615"}
	e := ProjectDomainSummary(data, domainScope(), meta)
	if e.Fields["first_seen_at"] != "2026-10-27T10:00:00.000000001Z" || e.Fields["last_seen_at"] != "2026-10-27T10:00:00.999999999Z" || e.Fields["occurrences"] != "18446744073709551615" || len(e.Fields) != 17 {
		t.Fatalf("summary metadata was lost or widened: %#v", e)
	}
	for _, count := range []string{"0", "01", "-1", "1.5", "18446744073709551616", "password_fixture"} {
		bad := meta
		bad.Occurrences = count
		requireProjectionSuppressed(t, ProjectDomainSummary(data, domainScope(), bad))
	}
	for _, stamp := range []string{"", "2026-02-30T10:00:00Z", "2026-10-27T10:00:00+00:00", "2026-10-27T10:00:00.1234567890Z", "2026-10-27T10:00:00,1Z"} {
		bad := meta
		bad.FirstSeenAt = stamp
		requireProjectionSuppressed(t, ProjectDomainSummary(data, domainScope(), bad))
		bad = meta
		bad.LastSeenAt = stamp
		requireProjectionSuppressed(t, ProjectDomainSummary(data, domainScope(), bad))
	}
	reversed := meta
	reversed.FirstSeenAt, reversed.LastSeenAt = meta.LastSeenAt, meta.FirstSeenAt
	requireProjectionSuppressed(t, ProjectDomainSummary(data, domainScope(), reversed))
	meta.LastSeenAt = meta.FirstSeenAt
	if e := ProjectDomainSummary(data, domainScope(), meta); e.Code != "domain.delivery.retrying" {
		t.Fatal("first occurrence may have equal timestamps")
	}
	requireProjectionSuppressed(t, ProjectDomainSummary([]byte(`{"raw_error":"password_fixture"}`), domainScope(), meta))
}

func TestDomainProjectionSummaryKeepsTheExistingEncodedCap(t *testing.T) {
	const code = DomainEventCode("domain.delivery.retrying")
	prior := domainProjectionTemplates[code]
	defer func() { domainProjectionTemplates[code] = prior }()
	data := projectionJSON(t, domainFixture())
	base := ProjectDomainEvent(data, domainScope())
	base.InstanceID, base.Sequence = strings.Repeat("f", 32), "18446744073709551615"
	base.Timestamp = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	grown := prior
	grown.message = strings.Repeat("m", 8192-len(projectionJSON(t, base))+len(prior.message)-50)
	domainProjectionTemplates[code] = grown
	if e := ProjectDomainEvent(data, domainScope()); e.Code != "domain.delivery.retrying" {
		t.Fatal("base event should fit before summary metadata")
	}
	e := ProjectDomainSummary(data, domainScope(), DomainEventSummary{FirstSeenAt: "2026-10-27T10:00:00.000000001Z", LastSeenAt: "2026-10-27T10:00:00.999999999Z", Occurrences: "18446744073709551615"})
	requireProjectionSuppressed(t, e)
	if !e.Truncated {
		t.Fatal("summary metadata overflow must be separately bounded")
	}
}
