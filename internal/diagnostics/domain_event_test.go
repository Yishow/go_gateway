package diagnostics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"go-gateway/internal/datalink/runtimeview"
)

func domainScope() runtimeview.Scope {
	return runtimeview.Scope{WorkspaceID: "W", WorkspaceRevision: "W1", DeviceID: "D"}
}

func domainText(s string) *string { return &s }

func domainFixture() DomainEventRef {
	return DomainEventRef{
		ContractVersion: 1, Code: "domain.delivery.retrying",
		Binding:     runtimeview.Binding{WorkspaceID: "W", WorkspaceRevision: "W1", EntityRevision: "R1"},
		Locator:     runtimeview.EntityLocator{Step: 4, WorkspaceID: "W", GroupID: domainText("G"), GroupRevision: domainText("R1")},
		ConnectorID: domainText("C"), ConnectorRevision: domainText("C1"),
		Stage: "retrying", Reason: "unavailable", Action: "refresh",
		Evidence: DomainEvidenceRef{Kind: "durable-transition", ReferenceID: "E"},
	}
}

func domainCommitFixture() DomainEventRef {
	v := domainFixture()
	v.Code, v.Stage, v.Reason, v.Action = "domain.delivery.first-confirmed", "sql-committed", "none", "none"
	v.Evidence.Kind = "runtime-receipt"
	v.Evidence.Receipt = &runtimeview.Receipt{GroupRevision: "R1", ConnectorRevision: "C1", RecordID: "record-E", EffectKey: "E", PayloadDigest: "digest-E", CommittedAt: "2026-10-27T09:00:00Z", Provenance: "runtime"}
	return v
}

func parseDomainFixture(t *testing.T, v DomainEventRef, scope runtimeview.Scope) (DomainEventRef, error) {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return ParseDomainEventRef(raw, scope)
}

func requireDomainUnavailable(t *testing.T, got DomainEventRef, err error) {
	t.Helper()
	if err == nil || err.Error() != "unavailable" {
		t.Fatalf("want safe unavailable error, got %v", err)
	}
	want := DomainEventRef{ContractVersion: 1, Code: "domain.unavailable", Stage: "unavailable", Reason: "unavailable", Action: "refresh", Evidence: DomainEvidenceRef{Kind: "unavailable"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("rejected input retained attribution/evidence: %#v", got)
	}
}

func TestDomainEventRequiredContracts(t *testing.T) {
	t.Run("missing scope unavailable", func(t *testing.T) {
		got, err := parseDomainFixture(t, domainFixture(), runtimeview.Scope{})
		requireDomainUnavailable(t, got, err)
	})
	t.Run("wrong revision not accepted", func(t *testing.T) {
		v := domainCommitFixture()
		v.Evidence.Receipt.GroupRevision = "R0" // Spec: effect E's only receipt is R0, not R1.
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	})
	t.Run("unknown code safe", func(t *testing.T) {
		v := domainFixture()
		v.Code = "password_fixture"
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	})
	t.Run("proof of queued cannot label committed", func(t *testing.T) {
		v := domainCommitFixture()
		v.Evidence = DomainEvidenceRef{Kind: "buffer-ack", ReferenceID: "E"}
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	})
	t.Run("no mutation locator", func(t *testing.T) {
		v := domainFixture()
		v.Action = "retry"
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	})
	t.Run("event reference carries original group/revision for historical failure, not latest draft", func(t *testing.T) {
		got, err := parseDomainFixture(t, domainFixture(), domainScope())
		if err != nil || got.Locator.GroupID == nil || *got.Locator.GroupID != "G" || got.Locator.GroupRevision == nil || *got.Locator.GroupRevision != "R1" || got.Binding.EntityRevision != "R1" {
			t.Fatalf("historical G/R1 must survive a newer draft: %#v, %v", got, err)
		}
	})
}

func TestDomainEventNoInventedDeliveryProof(t *testing.T) {
	for _, kind := range []DomainEvidenceKind{"buffer-ack", "http-200", "connection", "process", "durable-transition", "unavailable"} {
		t.Run(string(kind), func(t *testing.T) {
			v := domainCommitFixture()
			v.Evidence.Kind = kind
			got, err := parseDomainFixture(t, v, domainScope())
			requireDomainUnavailable(t, got, err)
		})
	}
	for _, mutate := range []struct {
		name   string
		change func(*DomainEventRef)
	}{
		{"absent receipt", func(v *DomainEventRef) { v.Evidence.Receipt = nil }},
		{"empty record", func(v *DomainEventRef) { v.Evidence.Receipt.RecordID = "" }},
		{"empty digest", func(v *DomainEventRef) { v.Evidence.Receipt.PayloadDigest = "" }},
		{"wrong effect", func(v *DomainEventRef) { v.Evidence.Receipt.EffectKey = "other-E" }},
		{"wrong connector revision", func(v *DomainEventRef) { v.Evidence.Receipt.ConnectorRevision = "C0" }},
		{"test-write is not runtime", func(v *DomainEventRef) { v.Evidence.Receipt.Provenance = "test-write" }},
		{"bad timestamp", func(v *DomainEventRef) { v.Evidence.Receipt.CommittedAt = "not-time" }},
		{"readback not inferred", func(v *DomainEventRef) { v.Stage = "readback-verified" }},
		{"cleanup not inferred", func(v *DomainEventRef) { v.Stage = "cleanup-completed" }},
	} {
		t.Run(mutate.name, func(t *testing.T) {
			v := domainCommitFixture()
			mutate.change(&v)
			got, err := parseDomainFixture(t, v, domainScope())
			requireDomainUnavailable(t, got, err)
		})
	}
	for _, code := range []DomainEventCode{"domain.delivery.first-confirmed", "domain.delivery.recovered", "domain.delivery.summary"} {
		v := domainCommitFixture()
		v.Code = code
		got, err := parseDomainFixture(t, v, domainScope())
		if err != nil || got.Stage != "sql-committed" || got.Evidence.Receipt == nil || got.Evidence.Receipt.EffectKey != "E" || got.Evidence.Receipt.GroupRevision != "R1" {
			t.Fatalf("G/R1 recovery must reference E: %#v, %v", got, err)
		}
	}
}

func TestDomainEventIdentityAndClosedShape(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*DomainEventRef)
	}{
		{"foreign workspace", func(v *DomainEventRef) { v.Binding.WorkspaceID = "foreign" }},
		{"stale workspace", func(v *DomainEventRef) { v.Binding.WorkspaceRevision = "W0" }},
		{"foreign device", func(v *DomainEventRef) { v.Locator.DeviceID = domainText("foreign") }},
		{"foreign locator", func(v *DomainEventRef) { v.Locator.WorkspaceID = "foreign" }},
		{"missing revision", func(v *DomainEventRef) { v.Binding.EntityRevision = "" }},
		{"historical hint mismatch", func(v *DomainEventRef) { v.Locator.GroupRevision = domainText("R2") }},
		{"missing group", func(v *DomainEventRef) { v.Locator.GroupID = nil }},
		{"missing connector", func(v *DomainEventRef) { v.ConnectorID = nil }},
		{"partial source identity", func(v *DomainEventRef) { v.PointID = domainText("P") }},
		{"not a per-tag success", func(v *DomainEventRef) { v.PointID, v.TagID = domainText("P"), domainText("T") }},
		{"unknown reason", func(v *DomainEventRef) { v.Reason = "raw_driver_error" }},
		{"unknown stage", func(v *DomainEventRef) { v.Stage = "future-stage" }},
		{"wrong shared action", func(v *DomainEventRef) { v.Action = "inspect-settings" }},
		{"missing evidence reference", func(v *DomainEventRef) { v.Evidence.ReferenceID = "" }},
		{"endpoint is not identity", func(v *DomainEventRef) { v.ConnectorID = domainText("https://private-host/") }},
		{"oversized identity", func(v *DomainEventRef) { v.Evidence.ReferenceID = strings.Repeat("a", 129) }},
		{"mixed evidence", func(v *DomainEventRef) { v.Evidence.Receipt = domainCommitFixture().Evidence.Receipt }},
		{"future version", func(v *DomainEventRef) { v.ContractVersion = 2 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := domainFixture()
			tc.change(&v)
			got, err := parseDomainFixture(t, v, domainScope())
			requireDomainUnavailable(t, got, err)
		})
	}
	raw, _ := json.Marshal(domainFixture())
	for _, candidate := range []string{
		strings.Replace(string(raw), `"contract_version":1`, `"contract_version":null`, 1),
		strings.Replace(string(raw), `"contract_version":1`, `"contract_version":1,"contract_version":1`, 1),
		strings.Replace(string(raw), `"contract_version":1`, `"Contract_Version":1`, 1),
		strings.Replace(string(raw), `"contract_version":1,`, "", 1),
		strings.Replace(string(raw), `"reference_id":"E"`, `"reference_id":"E","committed":true`, 1),
		strings.Replace(string(raw), `"step":4`, `"step":4,"mutation":"retry"`, 1),
		strings.Replace(string(raw), `"step":4`, `"step":4,"payload":{"raw_error":"password_fixture"}`, 1),
		string(raw) + `{}`, `[]`, `null`, `{"code":`, strings.Repeat(" ", MaxEventBytes+1),
	} {
		got, err := ParseDomainEventRef([]byte(candidate), domainScope())
		requireDomainUnavailable(t, got, err)
	}
}

func TestDomainEventSourceTransitions(t *testing.T) {
	for _, tc := range []struct {
		code, stage, kind string
		step              int
		success           bool
	}{
		{"domain.source.applied", "accepted", "source-runtime-transition", 2, true},
		{"domain.source.rejected", "blocked", "source-runtime-transition", 2, false},
		{"domain.collector.failed", "acquisition", "acquisition-result", 1, false},
		{"domain.collector.recovered", "acquisition", "acquisition-result", 1, true},
		{"domain.mapping.failed", "blocked", "mapping-result", 3, false},
		{"domain.mapping.recovered", "acquisition", "mapping-result", 3, true},
	} {
		t.Run(tc.code, func(t *testing.T) {
			v := domainFixture()
			v.Code, v.Stage, v.Evidence.Kind = DomainEventCode(tc.code), runtimeview.Stage(tc.stage), DomainEvidenceKind(tc.kind)
			v.Locator = runtimeview.EntityLocator{Step: tc.step, WorkspaceID: "W", DeviceID: domainText("D")}
			v.ConnectorID, v.ConnectorRevision = nil, nil
			if tc.step == 2 {
				v.Locator.RuleID = domainText("S")
			}
			if tc.step == 1 || tc.step == 3 {
				v.PointID = domainText("P")
			}
			if tc.step == 3 {
				v.Locator.MappingID, v.TagID = domainText("M"), domainText("T")
			}
			if tc.success {
				v.Reason, v.Action = "none", "none"
			}
			got, err := parseDomainFixture(t, v, domainScope())
			if err != nil || string(got.Code) != tc.code || got.Binding.EntityRevision != "R1" {
				t.Fatalf("confirmed source outcome lost: %#v, %v", got, err)
			}
			v.Evidence.Kind = "draft-saved"
			got, err = parseDomainFixture(t, v, domainScope())
			requireDomainUnavailable(t, got, err)
		})
	}
}

func TestDomainEventIntakeRevisionAndQuotaScope(t *testing.T) {
	for _, tc := range []struct{ code, stage, state string }{
		{"domain.group.applied", "collecting", "active"},
		{"domain.group.paused", "blocked", "not_running"},
		{"domain.group.retiring", "collecting", "retiring"},
	} {
		v := domainFixture()
		v.Code, v.Stage, v.Reason, v.Action = DomainEventCode(tc.code), runtimeview.Stage(tc.stage), "none", "none"
		v.Evidence.Kind = "intake-status"
		v.Evidence.Intake = &runtimeview.IntakeEvidence{GroupRevision: "R1", State: tc.state}
		got, err := parseDomainFixture(t, v, domainScope())
		if err != nil || got.Evidence.Intake.State != tc.state {
			t.Fatalf("%s: %#v, %v", tc.code, got, err)
		}
		v.Evidence.Intake.GroupRevision = "R0"
		got, err = parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
		v.Evidence.Intake = nil
		got, err = parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	}
	for _, quotaScope := range []string{"global", "group"} {
		v := domainFixture()
		v.Code, v.Stage, v.Evidence.Kind = "domain.quota.refused", "blocked", "quota-refusal"
		v.Evidence.Quota = &runtimeview.QuotaEvidence{Scope: quotaScope, State: "hard_limit", UsedBytes: 100, MaxBytes: 100, IntakeRefused: true}
		got, err := parseDomainFixture(t, v, domainScope())
		if err != nil || got.Evidence.Quota.Scope != quotaScope {
			t.Fatalf("quota scope relabeled: %#v, %v", got, err)
		}
		v.Evidence.Quota.Scope = "current-revision"
		got, err = parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	}
}

func TestDomainEventCounterReplayCountsInputsOnly(t *testing.T) {
	for _, reason := range []runtimeview.Reason{"counter-replay-proof-missing", "counter-replay-proof-conflict", "counter-replay-proof-unsupported"} {
		v := domainFixture()
		v.Code, v.Stage, v.Reason, v.Action = "domain.counter-replay.quarantined", "quarantined", reason, "inspect-immutable-evidence"
		v.Evidence.Kind, v.Evidence.PendingSamples = "counter-replay-inputs", domainText("18446744073709551615")
		got, err := parseDomainFixture(t, v, domainScope())
		if err != nil || got.Evidence.PendingSamples == nil || *got.Evidence.PendingSamples != "18446744073709551615" || got.Action != "inspect-immutable-evidence" {
			t.Fatalf("pending input count/action lost: %#v, %v", got, err)
		}
		for _, invalid := range []string{"-1", "01", "1.5", "18446744073709551616", "", "attention.length"} {
			v.Evidence.PendingSamples = domainText(invalid)
			got, err := parseDomainFixture(t, v, domainScope())
			requireDomainUnavailable(t, got, err)
		}
		v.Evidence.PendingSamples = nil
		got, err = parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	}
}

func TestDomainProducerMatrix(t *testing.T) {
	want := []string{"source", "collector", "mapping", "intake", "delivery", "commit", "counter-replay", "recording-refusal", "quota", "repair", "cutover"}
	matrix := DomainProducerMatrix()
	if len(matrix) != len(want) {
		t.Fatalf("matrix categories = %d, want 11", len(matrix))
	}
	seen := map[DomainEventCode]bool{}
	for i, p := range matrix {
		if p.Category != want[i] || p.Authority == "" || p.EmitBoundary == "" || p.Identity == "" || p.Evidence == "" || p.UnknownOrFailure == "" || len(p.Events) == 0 {
			t.Fatalf("incomplete producer plan: %#v", p)
		}
		for _, e := range p.Events {
			if seen[e.Code] || !strings.HasPrefix(string(e.Code), "domain.") {
				t.Fatalf("duplicate/invalid code: %s", e.Code)
			}
			seen[e.Code] = true
		}
	}
	for _, code := range []DomainEventCode{"domain.source.applied", "domain.source.rejected", "domain.collector.failed", "domain.collector.recovered", "domain.mapping.failed", "domain.mapping.recovered", "domain.group.applied", "domain.group.paused", "domain.group.retiring", "domain.delivery.retrying", "domain.delivery.blocked", "domain.delivery.quarantined", "domain.delivery.unknown", "domain.delivery.recovered", "domain.delivery.first-confirmed", "domain.delivery.summary", "domain.counter-replay.quarantined", "domain.recording.refused", "domain.quota.refused", "domain.repair.completed", "domain.repair.blocked", "domain.cutover.applied", "domain.cutover.rejected"} {
		if !seen[code] {
			t.Errorf("producer missing %s", code)
		}
	}
	if len(seen) != 24 {
		t.Fatalf("event code count = %d, want 24", len(seen))
	}
	// Returned documentation cannot mutate the parser's code policy.
	matrix[0].Events[0].Code = "domain.injected"
	if DomainProducerMatrix()[0].Events[0].Code != "domain.source.applied" {
		t.Fatal("shared mutable producer policy")
	}
}

func TestDomainEventRemainingTransitionsAndReasons(t *testing.T) {
	for _, tc := range []struct{ code, stage, kind, reason, action string }{
		{"domain.delivery.blocked", "blocked", "durable-transition", "unavailable", "refresh"},
		{"domain.delivery.quarantined", "quarantined", "durable-transition", "unavailable", "refresh"},
		{"domain.delivery.unknown", "unknown", "durable-transition", "unavailable", "refresh"},
		{"domain.repair.completed", "accepted", "operation-result", "none", "none"},
		{"domain.repair.blocked", "blocked", "operation-result", "unsupported", "inspect-settings"},
		{"domain.cutover.applied", "accepted", "operation-result", "none", "none"},
		{"domain.cutover.rejected", "blocked", "operation-result", "recording-ownership-unavailable", "inspect-immutable-evidence"},
		{"domain.cutover.rejected", "blocked", "operation-result", "recording-ownership-conflict", "inspect-immutable-evidence"},
		{"domain.cutover.rejected", "blocked", "operation-result", "recording-handoff-proof-conflict", "inspect-immutable-evidence"},
		{"domain.cutover.rejected", "blocked", "operation-result", "recording-cutover-finalized-overlap", "inspect-immutable-evidence"},
	} {
		v := domainFixture()
		v.Code, v.Stage, v.Evidence.Kind = DomainEventCode(tc.code), runtimeview.Stage(tc.stage), DomainEvidenceKind(tc.kind)
		v.Reason, v.Action = runtimeview.Reason(tc.reason), runtimeview.Action(tc.action)
		got, err := parseDomainFixture(t, v, domainScope())
		if err != nil || string(got.Code) != tc.code || string(got.Action) != tc.action {
			t.Fatalf("transition lost: %#v, %v", got, err)
		}
		if tc.reason == "none" {
			v.Reason, v.Action = "unavailable", "refresh"
		} else {
			v.Reason, v.Action = "none", "none"
		}
		got, err = parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	}
	t.Run("input quarantine never becomes finalized output", func(t *testing.T) {
		v := domainFixture()
		v.Code, v.Stage, v.Reason, v.Action = "domain.delivery.quarantined", "quarantined", "counter-replay-proof-missing", "inspect-immutable-evidence"
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	})
	t.Run("receipt time obeys shared wire contract", func(t *testing.T) {
		v := domainCommitFixture()
		v.Evidence.Receipt.CommittedAt = "2026-10-27T09:00:00,1Z"
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	})
}

func TestDomainEventIntakeQuotaFailureBranches(t *testing.T) {
	for _, mutate := range []func(*runtimeview.IntakeEvidence){
		func(i *runtimeview.IntakeEvidence) { i.State = "retiring" },
		func(i *runtimeview.IntakeEvidence) { r := runtimeview.ReasonUnavailable; i.Reason = &r },
	} {
		v := domainFixture()
		v.Code, v.Stage, v.Reason, v.Action, v.Evidence.Kind = "domain.group.applied", "collecting", "none", "none", "intake-status"
		v.Evidence.Intake = &runtimeview.IntakeEvidence{GroupRevision: "R1", State: "active"}
		mutate(v.Evidence.Intake)
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	}
	for _, mutate := range []func(*runtimeview.QuotaEvidence){
		func(q *runtimeview.QuotaEvidence) { q.State = "ok" },
		func(q *runtimeview.QuotaEvidence) { q.IntakeRefused = false },
		func(q *runtimeview.QuotaEvidence) { q.UsedBytes = -1 },
		func(q *runtimeview.QuotaEvidence) { q.MaxBytes = -1 },
		func(q *runtimeview.QuotaEvidence) { q.UsedBytes = 9007199254740992 },
		func(q *runtimeview.QuotaEvidence) { q.MaxBytes = 9007199254740992 },
	} {
		v := domainFixture()
		v.Code, v.Stage, v.Evidence.Kind = "domain.quota.refused", "blocked", "quota-refusal"
		v.Evidence.Quota = &runtimeview.QuotaEvidence{Scope: "global", State: "hard_limit", UsedBytes: 95, MaxBytes: 100, IntakeRefused: true}
		mutate(v.Evidence.Quota)
		got, err := parseDomainFixture(t, v, domainScope())
		requireDomainUnavailable(t, got, err)
	}
}

func TestDomainEventProducerReasonCompatibility(t *testing.T) {
	for _, tc := range []struct {
		code, kind, stage string
		step              int
		ownerReason       bool
	}{
		{"domain.collector.failed", "acquisition-result", "acquisition", 1, false},
		{"domain.source.rejected", "source-runtime-transition", "blocked", 2, false},
		{"domain.mapping.failed", "mapping-result", "blocked", 3, false},
		{"domain.delivery.blocked", "durable-transition", "blocked", 4, true}, // Task 2.4 preserves owner reasons here.
		{"domain.quota.refused", "quota-refusal", "blocked", 4, false},
		{"domain.repair.blocked", "operation-result", "blocked", 4, false},
	} {
		for _, reason := range []runtimeview.Reason{"recording-ownership-unavailable", "recording-ownership-conflict", "recording-handoff-proof-conflict", "recording-cutover-finalized-overlap"} {
			t.Run(tc.code+"/"+string(reason), func(t *testing.T) {
				v := domainFixture()
				v.Code, v.Stage, v.Evidence.Kind = DomainEventCode(tc.code), runtimeview.Stage(tc.stage), DomainEvidenceKind(tc.kind)
				v.Reason, v.Action = reason, "inspect-immutable-evidence"
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
				if tc.kind == "quota-refusal" {
					v.Evidence.Quota = &runtimeview.QuotaEvidence{Scope: "global", State: "hard_limit", UsedBytes: 95, MaxBytes: 100, IntakeRefused: true}
				}
				got, err := parseDomainFixture(t, v, domainScope())
				if tc.ownerReason {
					if err != nil || got.Reason != reason || got.Action != "inspect-immutable-evidence" {
						t.Fatalf("owner reason/action must survive group delivery: %#v, %v", got, err)
					}
				} else {
					requireDomainUnavailable(t, got, err)
				}
			})
		}
	}
}

func TestDomainEventGroupOnlyScope(t *testing.T) {
	scope := runtimeview.Scope{WorkspaceID: "W", WorkspaceRevision: "W1"}
	got, err := parseDomainFixture(t, domainCommitFixture(), scope)
	if err != nil || got.Code != "domain.delivery.first-confirmed" || got.Locator.DeviceID != nil || got.Binding.EntityRevision != "R1" {
		t.Fatalf("group receipt must not require or fabricate device attribution: %#v, %v", got, err)
	}
	for _, invalidScope := range []runtimeview.Scope{{WorkspaceRevision: "W1"}, {WorkspaceID: "W"}, {WorkspaceID: "W", WorkspaceRevision: "W1", DeviceID: "https://private-host"}} {
		got, err = parseDomainFixture(t, domainCommitFixture(), invalidScope)
		requireDomainUnavailable(t, got, err)
	}
	for _, step := range []int{1, 2, 3, 4} {
		v := domainFixture()
		v.Locator = runtimeview.EntityLocator{Step: step, WorkspaceID: "W", DeviceID: domainText("D")}
		switch step {
		case 1:
			v.Code, v.Stage, v.Evidence.Kind, v.PointID = "domain.collector.failed", "acquisition", "acquisition-result", domainText("P")
		case 2:
			v.Code, v.Stage, v.Evidence.Kind, v.Locator.RuleID = "domain.source.rejected", "blocked", "source-runtime-transition", domainText("S")
		case 3:
			v.Code, v.Stage, v.Evidence.Kind, v.Locator.MappingID, v.PointID, v.TagID = "domain.mapping.failed", "blocked", "mapping-result", domainText("M"), domainText("P"), domainText("T")
		case 4:
			v.Locator.GroupID, v.Locator.GroupRevision = domainText("G"), domainText("R1")
		}
		if step < 4 {
			v.ConnectorID, v.ConnectorRevision = nil, nil
		}
		got, err = parseDomainFixture(t, v, scope)
		requireDomainUnavailable(t, got, err)
		got, err = parseDomainFixture(t, v, domainScope())
		if err != nil || got.Locator.DeviceID == nil || *got.Locator.DeviceID != "D" {
			t.Fatalf("matching device must remain accepted: %#v, %v", got, err)
		}
	}
}
