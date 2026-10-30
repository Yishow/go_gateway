package diagnostics

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

func domainQueryRef(t *testing.T, code, workspace, device, rule, group, revision, ref string) DomainEventRef {
	t.Helper()
	var value DomainEventRef
	for _, candidate := range projectionCases {
		if candidate.code == code {
			value = projectionRef(candidate)
			break
		}
	}
	if value.Code == "" {
		t.Fatal("unknown fixture code")
	}
	value.Binding.WorkspaceID, value.Binding.EntityRevision = workspace, revision
	value.Locator.WorkspaceID = workspace
	if value.Locator.DeviceID != nil {
		value.Locator.DeviceID = domainText(device)
	}
	if value.Locator.RuleID != nil {
		value.Locator.RuleID = domainText(rule)
	}
	if value.Locator.GroupID != nil {
		value.Locator.GroupID, value.Locator.GroupRevision = domainText(group), domainText(revision)
	}
	value.Evidence.ReferenceID = ref
	return value
}

func emitDomainQuery(t *testing.T, b *Broker, value DomainEventRef) {
	t.Helper()
	scope := runtimeview.Scope{WorkspaceID: value.Binding.WorkspaceID, WorkspaceRevision: value.Binding.WorkspaceRevision}
	if value.Locator.DeviceID != nil {
		scope.DeviceID = *value.Locator.DeviceID
	}
	raw := projectionJSON(t, value)
	if _, err := ParseDomainEventRef(raw, scope); err != nil {
		t.Fatal(err)
	}
	if !b.EmitDomain(t.Context(), raw, scope) {
		t.Fatal("valid fixture was not accepted")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func queryReferences(records []Event) []string {
	out := make([]string, len(records))
	for i, event := range records {
		out[i], _ = event.Fields["reference_id"].(string)
	}
	return out
}

func TestDomainQueryFiltersMatchSnapshotReplayAndLive(t *testing.T) {
	b := testBroker(t)
	for _, v := range []DomainEventRef{
		domainQueryRef(t, "domain.source.rejected", "W", "D", "S", "", "R1", "a"),
		domainQueryRef(t, "domain.source.rejected", "W", "D", "S", "", "R2", "b"),
		domainQueryRef(t, "domain.source.rejected", "W", "D2", "S2", "", "R1", "other-device"),
		domainQueryRef(t, "domain.source.rejected", "W2", "D", "S", "", "R1", "other-workspace"),
		domainQueryRef(t, "domain.delivery.blocked", "W", "", "", "G", "R1", "group"),
	} {
		emitDomainQuery(t, b, v)
	}
	cases := []struct {
		name  string
		query Query
		refs  []string
	}{
		{"workspace", Query{WorkspaceID: "W2"}, []string{"other-workspace"}},
		{"device", Query{DeviceID: "D2"}, []string{"other-device"}},
		{"rule", Query{RuleID: "S2"}, []string{"other-device"}},
		{"group", Query{GroupID: "G"}, []string{"group"}},
		{"and", Query{WorkspaceID: "W", DeviceID: "D", RuleID: "S", Stage: "blocked", Actionable: true}, []string{"a", "b"}},
		{"missing-group-device", Query{GroupID: "G", DeviceID: "D"}, []string{}},
		{"missing-group-rule", Query{GroupID: "G", RuleID: "S"}, []string{}},
		{"different-stage", Query{Stage: "sql-committed"}, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			snapshot, err := b.Snapshot(tc.query)
			if err != nil || !reflect.DeepEqual(queryReferences(snapshot.Records), tc.refs) {
				t.Fatalf("snapshot=%v want=%v error=%v", queryReferences(snapshot.Records), tc.refs, err)
			}
			sub, err := b.Subscribe(tc.query, "")
			if err != nil {
				t.Fatal(err)
			}
			defer sub.Close()
			var replay []Event
			for range tc.refs {
				ctx, cancel := context.WithTimeout(t.Context(), time.Second)
				event, nextErr := sub.Next(ctx)
				cancel()
				if nextErr != nil {
					t.Fatal(nextErr)
				}
				replay = append(replay, event)
				sub.Ack(event.Cursor())
			}
			if !reflect.DeepEqual(queryReferences(replay), tc.refs) {
				t.Fatalf("replay mismatch: %v", queryReferences(replay))
			}
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Millisecond)
			defer cancel()
			if _, err := sub.Next(ctx); !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("unexpected extra match: %v", err)
			}
			if sub.Metadata().Progress != snapshot.Latest {
				t.Fatal("filtered progress did not advance past nonmatches")
			}
		})
	}
	q := Query{WorkspaceID: "W", DeviceID: "D", RuleID: "S", Stage: "blocked", Actionable: true}
	snapshot, _ := b.Snapshot(q)
	sub, err := b.Subscribe(q, snapshot.Latest)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	emitDomainQuery(t, b, domainQueryRef(t, "domain.source.rejected", "W", "D2", "S2", "", "R9", "late-foreign"))
	emitDomainQuery(t, b, domainQueryRef(t, "domain.source.rejected", "W", "D", "S", "", "R9", "live"))
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	live, err := sub.Next(ctx)
	if err != nil || live.Fields["reference_id"] != "live" {
		t.Fatalf("live mismatch: %#v %v", live, err)
	}
}

func TestDomainQueryRejectsInvalidTypedFilters(t *testing.T) {
	setters := map[string]func(*Query, string){
		"workspace": func(q *Query, s string) { q.WorkspaceID = s }, "device": func(q *Query, s string) { q.DeviceID = s },
		"rule": func(q *Query, s string) { q.RuleID = s }, "group": func(q *Query, s string) { q.GroupID = s },
	}
	for name, set := range setters {
		t.Run(name, func(t *testing.T) {
			for _, value := range []string{" W", "W ", "../private", "/private", "C:\\private", "https://private.example", "private.example", "node.LOCAL", "127.0.0.1", "[::1]", "fe80::1", "token=secret", "node\n", "設備", strings.Repeat("x", 129)} {
				q := Query{}
				set(&q, value)
				if _, err := ValidateQuery(q); !errors.Is(err, ErrInvalidQuery) {
					t.Errorf("invalid ID accepted: %q", value)
				}
			}
			for _, value := range []string{"", "W", "device-1", "_id", "-id", "R.1", "A_1-2.3", strings.Repeat("x", 128)} {
				q := Query{}
				set(&q, value)
				if _, err := ValidateQuery(q); err != nil {
					t.Errorf("canonical ID rejected: %q %v", value, err)
				}
			}
		})
	}
	for _, stage := range []runtimeview.Stage{"Unavailable", "future", " blocked", "password=secret"} {
		if _, err := ValidateQuery(Query{Stage: stage}); !errors.Is(err, ErrInvalidQuery) {
			t.Errorf("invalid stage accepted: %q", stage)
		}
	}
	for _, stage := range []runtimeview.Stage{"", "unavailable", "acquisition", "accepted", "collecting", "queued", "retrying", "blocked", "quarantined", "unknown", "skipped", "sql-committed", "readback-verified", "cleanup-completed", "transport-event"} {
		if _, err := ValidateQuery(Query{Stage: stage}); err != nil {
			t.Errorf("canonical stage rejected: %q", stage)
		}
	}
}

func TestDomainQueryActionableUsesFixedTypedActions(t *testing.T) {
	b := testBroker(t)
	emitDomainQuery(t, b, domainQueryRef(t, "domain.source.rejected", "W", "D", "S", "", "R1", "fault"))
	emitDomainQuery(t, b, domainQueryRef(t, "domain.source.applied", "W", "D", "S", "", "R2", "applied"))
	b.EmitDomainUnavailable(t.Context())
	b.Emit(Input{Code: "startup.failed", Fields: map[string]any{"action": "refresh", "stage": "blocked"}})
	b.Emit(Input{Code: "domain.unknown", Fields: map[string]any{"action": "refresh"}})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := b.Snapshot(Query{Actionable: true})
	if err != nil || len(snapshot.Records) != 2 || snapshot.Records[0].Fields["reference_id"] != "fault" || snapshot.Records[1].Code != "domain.unavailable" {
		t.Fatalf("actionable historical events=%#v error=%v", snapshot.Records, err)
	}
	all, _ := b.Snapshot(Query{})
	unrestricted, _ := b.Snapshot(Query{Actionable: false})
	if !reflect.DeepEqual(all.Records, unrestricted.Records) {
		t.Fatal("false means unrestricted, not non-actionable")
	}
	unknown, _ := b.Snapshot(Query{Stage: "unavailable", Actionable: true})
	if len(unknown.Records) != 1 || unknown.Records[0].Code != "domain.unavailable" {
		t.Fatal("fixed unavailable must be selectable")
	}
	unknown, _ = b.Snapshot(Query{WorkspaceID: "W", Stage: "unavailable", Actionable: true})
	if len(unknown.Records) != 0 {
		t.Fatal("unattributed event acquired a workspace")
	}
	for _, code := range []string{"domain.future", "startup.failed", "raw.suppressed"} {
		if matches(Query{Actionable: true}, Event{Code: code, Level: "error", Source: "runtime", Message: "refresh repair blocked", Fields: map[string]any{"action": "refresh", "stage": "blocked"}}) {
			t.Fatalf("raw code/text guessed actionable: %s", code)
		}
	}
	for _, action := range []any{"retry", "delete", "", true, map[string]any{"action": "refresh"}} {
		event := ProjectDomainEvent(sourceDomainJSON(t), domainScope())
		event.Fields["action"] = action
		if matches(Query{Actionable: true}, event) {
			t.Fatalf("invalid typed action accepted: %v", action)
		}
	}
}

func TestDomainQueryPreservesStringCursorAndRevisionHistory(t *testing.T) {
	b := testBroker(t)
	const precise uint64 = 9007199254740993
	b.mu.Lock()
	b.next, b.latest = precise-1, precise-1
	b.mu.Unlock()
	emitDomainQuery(t, b, domainQueryRef(t, "domain.source.rejected", "W", "A", "S", "", "R1", "fault-A"))
	a, err := b.Snapshot(Query{DeviceID: "A", Actionable: true})
	if err != nil || len(a.Records) != 1 || a.Records[0].Sequence != "9007199254740993" {
		t.Fatalf("lost exact cursor: %#v %v", a, err)
	}
	emitDomainQuery(t, b, domainQueryRef(t, "domain.source.rejected", "W", "B", "S2", "", "R2", "fault-B"))
	bView, _ := b.Snapshot(Query{DeviceID: "B"})
	if len(bView.Records) != 1 || bView.Records[0].Sequence != strconv.FormatUint(precise+1, 10) {
		t.Fatal("scope B borrowed A rows or rounded sequence")
	}
	before, err := b.Snapshot(Query{DeviceID: "B", Before: a.Latest})
	if err != nil || len(before.Records) != 0 {
		t.Fatal("before boundary did not remain exact")
	}
	sub, err := b.Subscribe(Query{DeviceID: "B"}, a.Latest)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	event, err := sub.Next(ctx)
	if err != nil || event.Sequence != "9007199254740994" {
		t.Fatalf("replay cursor changed: %s %v", event.Sequence, err)
	}
	if sub.Metadata().Progress != a.Latest {
		t.Fatal("unsent filtered replay must block progress")
	}
	sub.Ack(event.Cursor())
	if sub.Metadata().Progress != event.Cursor() {
		t.Fatal("ack did not preserve precise progress")
	}
}

func TestDomainQueryLegacyUnicodeSearchAndBoundsRemainCompatible(t *testing.T) {
	b := testBroker(t)
	emitDomainQuery(t, b, domainQueryRef(t, "domain.source.rejected", "W", "D", "STRASSE", "", "R1", "fault"))
	b.Emit(Input{Code: "runtime.degraded"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	for _, q := range []Query{{Q: "straße", WorkspaceID: "W"}, {Q: "SOURCE", DeviceID: "D", Level: "warn", Source: "runtime"}} {
		out, err := b.Snapshot(q)
		if err != nil || len(out.Records) != 1 {
			t.Fatalf("legacy folded search changed: %#v %v", out, err)
		}
	}
	for _, q := range []Query{{Q: ".*"}, {Q: "[A-Z]"}, {Level: "error"}, {Source: "http"}} {
		out, err := b.Snapshot(q)
		if err != nil || len(out.Records) != 0 {
			t.Fatal("legacy literal/severity/source filter changed")
		}
	}
	for _, q := range []Query{{Limit: 501}, {Limit: -1}, {Q: strings.Repeat("界", 257)}, {Q: string([]byte{0xff})}} {
		if _, err := b.Snapshot(q); !errors.Is(err, ErrInvalidQuery) {
			t.Fatal("legacy bound weakened")
		}
	}
	if RingRecords != 2000 || ReplayRecords != 500 || LiveRecords != 256 || MaxSubscribers != 8 {
		t.Fatal("transport caps changed")
	}
}
