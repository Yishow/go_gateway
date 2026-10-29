package runtimeview

import (
	"math"
	"strings"
	"testing"
)

func TestMissingUnknownAndUnsafeInputs(t *testing.T) {
	cases := []struct{ name, from, to string }{
		{"missing-required", "\"contract_version\":1,", ""}, {"null-primitive", "\"enabled\":true", "\"enabled\":null"},
		{"missing-boolean", "\"enabled\":false,", ""}, {"missing-array", "\"groups\":[", "\"unknown_groups\":["},
		{"unknown-root", "\"contract_version\":1", "\"contract_version\":1,\"password\":\"secret\""},
		{"wrong-case", "\"workspace_id\":\"W\"", "\"Workspace_id\":\"W\""},
		{"duplicate-key", "\"contract_version\":1", "\"contract_version\":1,\"contract_version\":1"},
		{"negative-count", "\"queued\":3", "\"queued\":-1"}, {"fractional-count", "\"queued\":3", "\"queued\":0.5"},
		{"unsafe-integer", "\"queued\":3", "\"queued\":9007199254740992"}, {"zero-poll", "\"interval_ms\":1000", "\"interval_ms\":0"},
		{"unknown-version", "\"contract_version\":1", "\"contract_version\":2"}, {"null-array", "\"sources\":[]", "\"sources\":null"},
		{"wrong-scalar", "\"enabled\":true", "\"enabled\":\"true\""}, {"wrong-object", "\"locator\":{", "\"locator\":[{"},
		{"bad-timestamp", "2026-10-27T10:00:00Z", "2026-02-30T10:00:00Z"}, {"non-utc", "2026-10-27T10:00:00Z", "2026-10-27T10:00:00+01:00"},
		{"alias-endpoint", "Meter archive", "postgres://host/private"}, {"private-table", "meter_new", "/private/db.sqlite"},
		{"sql-table", "meter_new", "SELECT * FROM secret"}, {"unknown-datatype", "\"int16\"", "\"number\""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseRuleEvidence([]byte(strings.Replace(contractLiteral, tc.from, tc.to, 1)), contractScope)
			if !IsContractError(err, ReasonUnavailable) {
				t.Fatalf("want unavailable got %v", err)
			}
		})
	}
	for _, raw := range []string{"", "null", "[]", "true", contractLiteral + " {}"} {
		if _, err := ParseRuleEvidence([]byte(raw), contractScope); !IsContractError(err, ReasonUnavailable) {
			t.Fatalf("accepted invalid JSON %v", err)
		}
	}
	raw := strings.Replace(contractLiteral, `"next_retry_at":null,`, "", 1)
	v, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil || v.Groups[0].Delivery.Current.NextRetryAt != nil {
		t.Fatal("missing nullable timing fabricated", err)
	}
	if _, err := ParseRuleEvidence([]byte(contractLiteral), Scope{"", "W1", "D"}); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("invalid scope accepted")
	}
	if _, err := ParseRuleEvidence([]byte(contractLiteral), Scope{"W", "W2", "D"}); !IsContractError(err, ReasonStale) {
		t.Fatal("stale observation accepted")
	}
}

func TestProvenSourceAndTransforms(t *testing.T) {
	v := fixture(t)
	s := v.Source
	f := *s.Saved
	f.Binding.EntityRevision = "R1"
	s.Effective = &f
	s.State = "effective"
	s.Proof = "immutable-source-facts"
	f.Transforms = []TransformStep{{Kind: "byte-order", Order: str("CDAB")}, {Kind: "scale", Multiplier: num(-2), Offset: num(3)}, {Kind: "cast", TargetType: str("float64")}}
	got := mustParse(t, v)
	if got.Source.State != "effective" || *got.Source.Effective.Transforms[1].Multiplier != -2 || got.Source.Saved.Binding.EntityRevision != "R2" {
		t.Fatal("source revisions collapsed")
	}
	for _, mutate := range []func(*SourceFacts){func(f *SourceFacts) { f.MappingRevision = nil }, func(f *SourceFacts) { f.Address = nil }, func(f *SourceFacts) { f.ReadType = nil }, func(f *SourceFacts) { f.ResultType = nil }, func(f *SourceFacts) { f.Polling = nil }} {
		copy := f
		mutate(&copy)
		s.Effective = &copy
		if mustParse(t, v).Source.State != "unavailable" {
			t.Fatal("missing effective fact promoted")
		}
	}
	s.Effective = &f
	for _, step := range []TransformStep{{Kind: "scale"}, {Kind: "cast", TargetType: str("float64"), Offset: num(1)}, {Kind: "byte-order"}, {Kind: "unsupported", Offset: num(1)}} {
		f.Transforms = []TransformStep{step}
		if _, err := parseFixture(t, v); !IsContractError(err, ReasonUnavailable) {
			t.Fatal("invalid transform accepted", err)
		}
	}
	f.Transforms = []TransformStep{}
	f.Binding.WorkspaceID = "foreign"
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonForeign) {
		t.Fatal("foreign source accepted")
	}
}

func committedFixture(t *testing.T) RuleEvidenceView {
	v := fixture(t)
	c := v.Groups[0].Delivery.Current
	c.Stages.SQLCommitted = 1
	c.LastCommit = &Receipt{"A", "7", "record1", "effect1", "digest1", "2026-10-27T10:00:00Z", "runtime"}
	c.SourceState = "available"
	unknown := false
	c.Sources = []SourceEvidence{{DeviceID: "D", PointID: "P", TagID: "T", Latest: &SourceResult{State: "unknown", WindowStart: str("2026-10-27T10:00:00Z"), WindowEnd: str("2026-10-27T10:00:15Z"), DeltaKnown: &unknown}, LastCommitted: &SourceResult{State: "partial", WindowStart: str("2026-10-27T09:59:45Z"), WindowEnd: str("2026-10-27T10:00:00Z"), Coverage: str("partial"), Receipt: c.LastCommit}}}
	return v
}
func TestRuntimeReceiptAndSourceBoundary(t *testing.T) {
	v := committedFixture(t)
	got := mustParse(t, v)
	s := got.Groups[0].Delivery.Current.Sources[0]
	if s.Latest.State != "unknown" || *s.Latest.DeltaKnown || s.Latest.Receipt != nil || s.LastCommitted.State != "partial" || s.LastCommitted.Receipt.RecordID != "record1" {
		t.Fatal("source evidence conflated")
	}
	cases := []struct {
		name   string
		mutate func(*RuleEvidenceView)
		reason Reason
	}{
		{"receipt-group", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.LastCommit.GroupRevision = "B" }, ReasonStale},
		{"receipt-connector", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.LastCommit.ConnectorRevision = "8" }, ReasonStale},
		{"test-as-runtime", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.LastCommit.Provenance = "test-write" }, ReasonUnavailable},
		{"missing-receipt", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.LastCommit = nil }, ReasonUnavailable},
		{"receipt-with-zero", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Stages.SQLCommitted = 0 }, ReasonUnavailable},
		{"shared-hidden", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Shared = false }, ReasonUnavailable},
		{"foreign-source", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Sources[0].DeviceID = "foreign" }, ReasonForeign},
		{"duplicate-source", func(v *RuleEvidenceView) {
			c := v.Groups[0].Delivery.Current
			c.Sources = append(c.Sources, c.Sources[0])
		}, ReasonUnavailable},
		{"unavailable-sources", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.SourceState = "unavailable" }, ReasonUnavailable},
		{"half-window", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Sources[0].Latest.WindowEnd = nil }, ReasonUnavailable},
		{"reverse-window", func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.WindowEnd = str("2026-10-27T09:00:00Z")
		}, ReasonUnavailable},
		{"unknown-as-zero", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Sources[0].Latest.State = "available" }, ReasonUnavailable},
		{"partial-as-complete", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Sources[0].LastCommitted.State = "available" }, ReasonUnavailable},
		{"unmatched-source-receipt", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Current.Sources[0].LastCommitted.Receipt = nil }, ReasonUnavailable},
		{"unavailable-committed", func(v *RuleEvidenceView) {
			r := v.Groups[0].Delivery.Current.Sources[0].LastCommitted
			r.Coverage = nil
			r.State = "unavailable"
		}, ReasonUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := committedFixture(t)
			tc.mutate(&v)
			_, err := parseFixture(t, v)
			if !IsContractError(err, tc.reason) {
				t.Fatalf("want %s got %v", tc.reason, err)
			}
		})
	}
}

func TestLocatorStepsAndActions(t *testing.T) {
	locators := []EntityLocator{{Step: 1, WorkspaceID: "W", DeviceID: str("D")}, {Step: 2, WorkspaceID: "W", DeviceID: str("D"), RuleID: str("R")}, {Step: 3, WorkspaceID: "W", DeviceID: str("D"), MappingID: str("M")}, {Step: 4, WorkspaceID: "W", GroupID: str("G"), GroupRevision: str("A0")}}
	for _, l := range locators {
		if err := ValidateEntityLocator(l, contractScope); err != nil {
			t.Fatal(err)
		}
		v := fixture(t)
		v.Locator = &l
		mustParse(t, v)
		wrong := l
		if l.Step < 4 {
			wrong.DeviceID = nil
		} else {
			wrong.GroupID = nil
		}
		if err := ValidateEntityLocator(wrong, contractScope); !IsContractError(err, ReasonUnavailable) {
			t.Fatal("missing locator identity accepted")
		}
		wrong = l
		wrong.MappingID = str("other")
		if l.Step == 3 {
			wrong.RuleID = str("other")
		}
		if err := ValidateEntityLocator(wrong, contractScope); !IsContractError(err, ReasonUnavailable) {
			t.Fatal("unrelated locator identity accepted")
		}
	}
	for _, pair := range []struct {
		reason Reason
		action Action
	}{{ReasonNone, ActionNone}, {ReasonForeign, ActionOpenWorkspaceList}, {ReasonDeleted, ActionOpenWorkspaceList}, {ReasonUnsupported, ActionInspectSettings}, {ReasonMissingProof, ActionRefresh}, {ReasonCounterMissing, ActionInspectImmutableEvidence}, {ReasonCounterConflict, ActionInspectImmutableEvidence}, {ReasonCounterUnsupported, ActionInspectImmutableEvidence}, {ReasonOwnershipUnavailable, ActionInspectImmutableEvidence}, {ReasonOwnershipConflict, ActionInspectImmutableEvidence}, {ReasonHandoffConflict, ActionInspectImmutableEvidence}, {ReasonFinalizedOverlap, ActionInspectImmutableEvidence}, {ReasonStale, ActionRefresh}, {ReasonUnavailable, ActionRefresh}} {
		if ParseReason(string(pair.reason)) != pair.reason || ActionForReason(pair.reason) != pair.action {
			t.Fatal("reason/action contract drift", pair)
		}
	}
	for _, stage := range []Stage{StageUnavailable, StageAcquisition, StageAccepted, StageCollecting, StageQueued, StageRetrying, StageBlocked, StageQuarantined, StageUnknown, StageSkipped, StageSQLCommitted, StageReadbackVerified, StageCleanupCompleted, StageTransportEvent} {
		if ParseStage(string(stage)) != stage {
			t.Fatal("stage drift")
		}
	}
	if ParseStage("secret-driver-text") != StageUnavailable || ParseReason("secret-driver-text") != ReasonUnavailable || ActionForReason(Reason("future")) != ActionRefresh || IsContractError(nil, ReasonNone) {
		t.Fatal("unsafe default")
	}
}

func TestFreshnessStatesAndOptionalEvidence(t *testing.T) {
	v := fixture(t)
	v.Freshness = "stale"
	mustParse(t, v)
	v.ObservedAt = nil
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("stale without observation")
	}
	v.Freshness = "unavailable"
	v.Source = nil
	v.Groups = []GroupView{}
	v.Locator = nil
	mustParse(t, v)
	v = fixture(t)
	g := &v.Groups[0]
	g.State = "saved-only"
	g.Effective = nil
	g.Proof = "unavailable"
	g.Delivery = nil
	mustParse(t, v)
	g.State = "unavailable"
	g.Saved = nil
	g.SavedRevision = nil
	mustParse(t, v)
	v = fixture(t)
	c := v.Groups[0].Delivery.Current
	c.NextDueAt = str("2026-10-27T10:00:15Z")
	c.NextRetryAt = str("2026-10-27T10:00:16Z")
	c.AttemptDurationMS = num(12.5)
	c.Quota = &QuotaEvidence{Scope: "global", State: "hard_limit", UsedBytes: 95, MaxBytes: 100, IntakeRefused: true}
	c.Storage = &StorageEvidence{Generation: 2, HistoryMayBeLost: true, Reconstructions: 1, LastReconstructedAt: str("2026-10-27T09:00:00Z"), RepairState: "blocked", Events: []StorageEvent{{Kind: "repair_blocked", At: str("2026-10-27T10:00:00Z")}}}
	got := mustParse(t, v)
	if *got.Groups[0].Delivery.Current.AttemptDurationMS != 12.5 || got.Groups[0].Delivery.Current.Quota.State != "hard_limit" {
		t.Fatal("actual evidence dropped")
	}
}

func TestIncompleteProofCannotHideReceiptRevisionMismatch(t *testing.T) {
	v := committedFixture(t)
	v.Groups[0].Proof = "unavailable"
	v.Groups[0].Delivery.Current.LastCommit.GroupRevision = "B"
	if _, err := parseFixture(t, v); !IsContractError(err, ReasonStale) {
		t.Fatal("missing proof concealed mixed receipt revision", err)
	}
}

func TestReviewReceiptConflictDraftRemovalAndNanoseconds(t *testing.T) {
	t.Run("repeated receipt digest conflict", func(t *testing.T) {
		v := committedFixture(t)
		r := *v.Groups[0].Delivery.Current.LastCommit
		r.PayloadDigest = "different"
		v.Groups[0].Delivery.Current.Sources[0].LastCommitted.Receipt = &r
		if _, err := parseFixture(t, v); !IsContractError(err, ReasonStale) {
			t.Fatal("conflicting receipt accepted", err)
		}
	})
	t.Run("different older source receipt allowed", func(t *testing.T) {
		v := committedFixture(t)
		r := *v.Groups[0].Delivery.Current.LastCommit
		r.RecordID = "older-record"
		r.EffectKey = "older-effect"
		r.PayloadDigest = "older-digest"
		r.CommittedAt = "2026-10-27T09:59:45Z"
		v.Groups[0].Delivery.Current.Sources[0].LastCommitted.Receipt = &r
		mustParse(t, v)
	})
	t.Run("draft removal retains applied device", func(t *testing.T) {
		v := fixture(t)
		v.Groups[0].Saved.Members = v.Groups[0].Saved.Members[1:]
		got := mustParse(t, v)
		if len(got.Groups[0].Effective.Members) != 2 || len(got.Groups[0].Saved.Members) != 1 {
			t.Fatal("membership merged")
		}
	})
	t.Run("nanosecond window", func(t *testing.T) {
		v := committedFixture(t)
		r := v.Groups[0].Delivery.Current.Sources[0].Latest
		r.WindowStart = str("2026-10-27T10:00:00.000000001Z")
		r.WindowEnd = str("2026-10-27T10:00:00.000000002Z")
		mustParse(t, v)
		for _, end := range []string{"2026-10-27T10:00:00.000000001Z", "2026-10-27T10:00:00Z", "2026-10-27T10:00:00.000000002+00:00"} {
			r.WindowEnd = str(end)
			if _, err := parseFixture(t, v); !IsContractError(err, ReasonUnavailable) {
				t.Fatal("invalid exact window accepted", err)
			}
		}
	})
	t.Run("cleanup does not imply readback", func(t *testing.T) {
		v := fixture(t)
		v.Groups[0].Delivery.TestWrite = &TestWriteEvidence{OperationID: "test2", Provenance: "test-write", WrittenAt: str("2026-10-27T10:00:00Z"), CleanupAt: str("2026-10-27T10:00:02Z")}
		got := mustParse(t, v)
		if got.Groups[0].Delivery.TestWrite.ReadbackAt != nil {
			t.Fatal("cleanup implied readback")
		}
	})
}

func TestCurrentIntakeReasonIsRevisionBound(t *testing.T) {
	raw := strings.Replace(contractLiteral, `"intake":null`, `"intake":{"group_revision":"A","state":"blocked","reason":"counter-replay-proof-missing"}`, 1)
	got, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil {
		t.Fatal("current intake reason lost", err)
	}
	encoded, err := MarshalRuleEvidence(*got, contractScope)
	if err != nil || !strings.Contains(string(encoded), `"intake":{"group_revision":"A","state":"blocked","reason":"counter-replay-proof-missing"}`) {
		t.Fatal("current reason not projected", err)
	}
	raw = strings.Replace(raw, `"state":"blocked","reason":"counter-replay-proof-missing"`, `"state":"active","reason":null`, 1)
	if _, err := ParseRuleEvidence([]byte(raw), contractScope); err != nil {
		t.Fatal("healthy reason absence invented failure", err)
	}
	raw = strings.Replace(raw, `"state":"active","reason":null`, `"state":"active","reason":"secret-error"`, 1)
	if _, err := ParseRuleEvidence([]byte(raw), contractScope); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("raw current reason accepted", err)
	}
}

func TestStateScopeAndDefaultBranches(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*RuleEvidenceView)
		reason Reason
	}{
		{"source-saved-missing", func(v *RuleEvidenceView) { v.Source.Saved = nil }, ReasonUnavailable},
		{"source-unavailable-with-effective", func(v *RuleEvidenceView) {
			v.Source.State = "unavailable"
			v.Source.Effective = v.Source.Saved
			v.Source.EffectiveRevision = v.Source.SavedRevision
		}, ReasonUnavailable},
		{"source-saved-with-effective", func(v *RuleEvidenceView) {
			v.Source.Effective = v.Source.Saved
			v.Source.EffectiveRevision = v.Source.SavedRevision
		}, ReasonUnavailable},
		{"duplicate-group", func(v *RuleEvidenceView) { v.Groups = append(v.Groups, v.Groups[0]) }, ReasonUnavailable},
		{"duplicate-member", func(v *RuleEvidenceView) { f := v.Groups[0].Saved; f.Members = append(f.Members, f.Members[0]) }, ReasonUnavailable},
		{"foreign-all-members", func(v *RuleEvidenceView) {
			v.Groups[0].Saved.Members = v.Groups[0].Saved.Members[1:]
			v.Groups[0].Effective.Members = v.Groups[0].Effective.Members[1:]
		}, ReasonForeign},
		{"saved-group-with-effective", func(v *RuleEvidenceView) { v.Groups[0].State = "saved-only" }, ReasonUnavailable},
		{"unavailable-group-with-effective", func(v *RuleEvidenceView) { v.Groups[0].State = "unavailable" }, ReasonUnavailable},
		{"backlog-current-revision", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Backlog[0].Binding.EntityRevision = "A" }, ReasonStale},
		{"backlog-mixed-workspace", func(v *RuleEvidenceView) { v.Groups[0].Delivery.Backlog[0].Binding.WorkspaceID = "foreign" }, ReasonForeign},
		{"locator-stale-group", func(v *RuleEvidenceView) { v.Locator.GroupID = str("other") }, ReasonStale},
		{"locator-stale-rule", func(v *RuleEvidenceView) {
			v.Locator = &EntityLocator{Step: 2, WorkspaceID: "W", DeviceID: str("D"), RuleID: str("other")}
		}, ReasonStale},
		{"locator-stale-mapping", func(v *RuleEvidenceView) {
			v.Locator = &EntityLocator{Step: 3, WorkspaceID: "W", DeviceID: str("D"), MappingID: str("other")}
		}, ReasonStale},
		{"locator-foreign-device", func(v *RuleEvidenceView) { v.Locator.DeviceID = str("other") }, ReasonForeign},
		{"locator-invalid-step", func(v *RuleEvidenceView) { v.Locator.Step = 0 }, ReasonUnavailable},
		{"trailing-newline-label", func(v *RuleEvidenceView) { v.Groups[0].Saved.Destination.Alias = str("Meter archive\n") }, ReasonUnavailable},
		{"trailing-newline-id", func(v *RuleEvidenceView) { v.Groups[0].Saved.Destination.Table = "meter_new\n" }, ReasonUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := fixture(t)
			tc.mutate(&v)
			_, err := parseFixture(t, v)
			if !IsContractError(err, tc.reason) {
				t.Fatalf("want %s got %v", tc.reason, err)
			}
		})
	}
	for _, mutate := range []func(*GroupView){func(g *GroupView) { g.Effective.EffectiveAt = nil }, func(g *GroupView) { g.Effective.Layout = "unavailable" }, func(g *GroupView) { g.Effective.IntervalSeconds = nil }, func(g *GroupView) { g.Effective = nil }} {
		v := fixture(t)
		mutate(&v.Groups[0])
		if mustParse(t, v).Groups[0].State != "unavailable" {
			t.Fatal("partial group proof claimed effective")
		}
	}
	v := fixture(t)
	v.Source.Saved.Transforms = []TransformStep{{Kind: "unsupported"}}
	mustParse(t, v)
	v = fixture(t)
	v.Groups[0].Saved.Members = v.Groups[0].Saved.Members[:1]
	v.Groups[0].Effective.Members = v.Groups[0].Effective.Members[:1]
	v.Groups[0].Delivery.Current.Shared = false
	mustParse(t, v)
	v = fixture(t)
	v.Groups[0].Delivery.Current = nil
	mustParse(t, v)
	v = fixture(t)
	v.Groups[0].Delivery.Current.AttemptDurationMS = num(math.NaN())
	if _, err := MarshalRuleEvidence(v, contractScope); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("non-finite serialization not safely refused", err)
	}
}

func TestAuthorityScopeLabels(t *testing.T) {
	raw := strings.Replace(contractLiteral, `"intake":null`, `"intake":{"group_revision":"A","state":"active","reason":null}`, 1)
	raw = strings.Replace(raw, `"quota":null`, `"quota":{"scope":"global","state":"hard_limit","used_bytes":95,"max_bytes":100,"intake_refused":true}`, 1)
	raw = strings.Replace(raw, `"oldest_pending_seconds":null`, `"oldest_pending_seconds":120`, 1)
	v, err := ParseRuleEvidence([]byte(raw), contractScope)
	if err != nil {
		t.Fatal("scoped authority evidence rejected", err)
	}
	if v.Groups[0].Delivery.OldestPendingScope != "group-all-revisions" || *v.Groups[0].Delivery.OldestPendingSeconds != 120 {
		t.Fatal("pending age lost aggregate scope")
	}
	groupQuota := strings.Replace(raw, `"scope":"global"`, `"scope":"group"`, 1)
	if _, err := ParseRuleEvidence([]byte(groupQuota), contractScope); err != nil {
		t.Fatal("group quota rejected", err)
	}
	if _, err := ParseRuleEvidence([]byte(strings.Replace(raw, `"scope":"global"`, `"scope":"future"`, 1)), contractScope); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("unknown quota scope accepted", err)
	}
	encoded, err := MarshalRuleEvidence(*v, contractScope)
	if err != nil || !strings.Contains(string(encoded), `"scope":"global"`) {
		t.Fatal("global quota scope lost", err)
	}
	mixed := strings.Replace(raw, `"group_revision":"A","state":"active"`, `"group_revision":"B","state":"active"`, 1)
	if _, err := ParseRuleEvidence([]byte(mixed), contractScope); !IsContractError(err, ReasonStale) {
		t.Fatal("foreign intake revision accepted", err)
	}
	aggregate := strings.Replace(raw, `"group_revision":"A","state":"active"`, `"state":"active"`, 1)
	if _, err := ParseRuleEvidence([]byte(aggregate), contractScope); !IsContractError(err, ReasonUnavailable) {
		t.Fatal("aggregate-only intake became current", err)
	}
}
