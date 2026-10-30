package diagnostics

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"go-gateway/internal/datalink/runtimeview"
)

func quotaRef(scope string, recovered bool) DomainEventRef {
	v := domainFixture()
	v.Code, v.Stage = "domain.quota.refused", runtimeview.StageBlocked
	v.Evidence = DomainEvidenceRef{Kind: "quota-refusal", ReferenceID: "quota-observation", Quota: &runtimeview.QuotaEvidence{Scope: scope, State: "hard_limit", UsedBytes: 100, MaxBytes: 100, IntakeRefused: true}}
	if recovered {
		v.Code, v.Stage, v.Reason, v.Action = "domain.quota.recovered", runtimeview.StageAccepted, runtimeview.ReasonNone, runtimeview.ActionNone
		v.Evidence.Kind = "quota-status"
		v.Evidence.Quota.State, v.Evidence.Quota.UsedBytes, v.Evidence.Quota.IntakeRefused = "ok", 0, false
	}
	return v
}
func intakeRef(state string) DomainEventRef {
	v := domainFixture()
	v.Code, v.Stage, v.Reason, v.Action = DomainEventCode("domain.group."+state), runtimeview.StageCollecting, runtimeview.ReasonNone, runtimeview.ActionNone
	intake := map[string]string{"applied": "active", "paused": "not_running", "retiring": "retiring"}[state]
	if state == "paused" {
		v.Stage = runtimeview.StageBlocked
	}
	v.Evidence = DomainEvidenceRef{Kind: "intake-status", ReferenceID: "intake-observation", Intake: &runtimeview.IntakeEvidence{GroupRevision: v.Binding.EntityRevision, State: intake}}
	return v
}
func TestDomainQuotaRecoveryContractIsCapacityOnly(t *testing.T) {
	for _, scope := range []string{"global", "group"} {
		for _, state := range []string{"ok", "warning"} {
			t.Run(scope+state, func(t *testing.T) {
				v := quotaRef(scope, true)
				v.Evidence.Quota.State = state
				if state == "warning" {
					v.Evidence.Quota.UsedBytes = 85
				}
				raw := projectionJSON(t, v)
				if _, err := ParseDomainEventRef(raw, domainScope()); err != nil {
					t.Fatalf("valid exact capacity recovery suppressed: %v", err)
				}
				e := ProjectDomainEvent(raw, domainScope())
				if e.Code != "domain.quota.recovered" || e.Fields["quota_scope"] != scope || e.Fields["intake_refused"] != false || e.Fields["record_id"] != nil || e.Fields["intake_state"] != nil || !strings.Contains(e.Message, "not resumed intake or SQL") {
					t.Fatalf("wrong capacity-only projection: %#v", e)
				}
			})
		}
	}
}
func TestDomainQuotaRecoveryRejectsInvalidOrMixedProof(t *testing.T) {
	for _, kind := range []string{"hard", "refused", "zero-max", "negative", "over-max", "wrong-kind", "receipt", "intake", "scope", "raw"} {
		t.Run(kind, func(t *testing.T) {
			v := quotaRef("group", true)
			switch kind {
			case "hard":
				v.Evidence.Quota.State = "hard_limit"
			case "refused":
				v.Evidence.Quota.IntakeRefused = true
			case "zero-max":
				v.Evidence.Quota.MaxBytes = 0
			case "negative":
				v.Evidence.Quota.UsedBytes = -1
			case "over-max":
				v.Evidence.Quota.UsedBytes = 100
			case "wrong-kind":
				v.Evidence.Kind = "quota-refusal"
			case "receipt":
				v.Evidence.Receipt = domainCommitFixture().Evidence.Receipt
			case "intake":
				v.Evidence.Intake = intakeRef("applied").Evidence.Intake
			case "scope":
				v.Evidence.Quota.Scope = "device"
			}
			raw := projectionJSON(t, v)
			if kind == "raw" {
				raw = []byte(strings.TrimSuffix(string(raw), "}") + `,"raw_error":"password_fixture"}`)
			}
			requireProjectionSuppressed(t, ProjectDomainEvent(raw, domainScope()))
		})
	}
}
func TestDomainIntakeInfoEligibilityIsImmutableAndNilSafe(t *testing.T) {
	for _, level := range []string{"nil", "empty", "error", "warn", "info"} {
		t.Run(level, func(t *testing.T) {
			var b *Broker
			if level == "empty" {
				b = &Broker{}
			} else if level != "nil" {
				b = pausedDomainBroker(t)
				b.level = levelRank(level)
			}
			enabled, ok := any(b).(interface{ DomainInfoEnabled() bool })
			if !ok {
				t.Fatal("missing immutable info eligibility")
			}
			if enabled.DomainInfoEnabled() != (level == "info") {
				t.Fatal("wrong immutable info eligibility")
			}
		})
	}
}

// Reflection keeps the initial absent-API test executable rather than a compile
// failure. All returned references and emitted events are checked through the
// same production API, and the opaque token is never manufactured by this helper.
func takeQuotaHint(t *testing.T, b *Broker, ctx context.Context) (reflect.Value, bool) {
	t.Helper()
	m := reflect.ValueOf(b).MethodByName("TakeQuotaRecoveryHint")
	if !m.IsValid() {
		t.Fatal("missing bounded retained-quota hint API")
	}
	result := m.Call([]reflect.Value{reflect.ValueOf(ctx)})
	return result[0], result[1].Bool()
}
func quotaHintRef(t *testing.T, h reflect.Value) DomainEventRef {
	t.Helper()
	v, ok := h.FieldByName("Ref").Interface().(DomainEventRef)
	if !ok {
		t.Fatal("hint lacks validated public reference")
	}
	return v
}
func emitQuotaHint(t *testing.T, b *Broker, h reflect.Value, v DomainEventRef, scope runtimeview.Scope) bool {
	t.Helper()
	m := reflect.ValueOf(b).MethodByName("EmitQuotaRecovery")
	if !m.IsValid() {
		t.Fatal("missing exact-token quota recovery API")
	}
	return m.Call([]reflect.Value{reflect.ValueOf(t.Context()), h, reflect.ValueOf(projectionJSON(t, v)), reflect.ValueOf(scope)})[0].Bool()
}

func TestDomainIntakeTextDoesNotClaimLatestConfiguration(t *testing.T) {
	e := ProjectDomainEvent(projectionJSON(t, intakeRef("applied")), domainScope())
	if !strings.Contains(e.Message, "reconciliation observation") || !strings.Contains(e.Message, "original revision") {
		t.Fatalf("current/latest whole configuration implied: %q", e.Message)
	}
}
func TestDomainQuotaRecoveryCannotBypassHint(t *testing.T) {
	b := pausedDomainBroker(t)
	emitDelivery(t, b, quotaRef("group", false))
	if emitDelivery(t, b, quotaRef("group", true)) {
		t.Fatal("recovery accepted without retained opaque hint")
	}
	if len(deliveryEvents(b, "domain.quota.recovered")) != 0 || b.domain.Stats().Entries != 1 {
		t.Fatal("direct recovery changed fault")
	}
}
