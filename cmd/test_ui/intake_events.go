package main

import (
	"context"
	"encoding/json"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/diagnostics"
)

// ScopeRevision is the already-created delivery ReadSnapshot seam. The adapter
// never opens a workspace, probes a destination or derives native intake state.
func newIntakeQuotaEventObserver(logs *diagnostics.Broker, scopeRevision func(context.Context, groupdelivery.GroupKey) (string, bool)) grouppipeline.IntakeQuotaObservationConfig {
	if logs == nil || !logs.DomainWarningsEnabled() {
		return grouppipeline.IntakeQuotaObservationConfig{}
	}
	return grouppipeline.IntakeQuotaObservationConfig{
		RecordLimit: diagnostics.AdmissionRecords, OnDrop: logs.RecordDomainCaptureLoss,
		IntakeEnabled: logs.DomainInfoEnabled(), QuotaEnabled: true, ScopeRevision: scopeRevision,
		Observe: func(o grouppipeline.IntakeQuotaObservation) {
			ref, ok := intakeQuotaEventRef(o, false)
			if !ok {
				logs.EmitDomainUnavailable(context.Background())
				return
			}
			emitDeliveryEvent(logs, ref)
		},
		TakeQuotaRecoveryHint: func(ctx context.Context) (grouppipeline.QuotaRecoveryHint, bool) {
			hint, ok := logs.TakeQuotaRecoveryHint(ctx)
			if !ok {
				return grouppipeline.QuotaRecoveryHint{}, false
			}
			ref := hint.Ref
			if ref.Code != "domain.quota.refused" || ref.Locator.GroupID == nil || ref.ConnectorID == nil || ref.ConnectorRevision == nil || ref.Evidence.Quota == nil {
				logs.EmitDomainUnavailable(ctx)
				return grouppipeline.QuotaRecoveryHint{}, false
			}
			native := grouppipeline.QuotaRecoveryHint{Key: groupdelivery.GroupKey{WorkspaceID: ref.Binding.WorkspaceID, GroupID: *ref.Locator.GroupID, GroupRevision: ref.Binding.EntityRevision}, WorkspaceRevision: ref.Binding.WorkspaceRevision, ConnectorID: *ref.ConnectorID, ConnectorRevision: *ref.ConnectorRevision, ReferenceID: ref.Evidence.ReferenceID, Scope: groupdelivery.QuotaScope(ref.Evidence.Quota.Scope)}
			native.Observe = func(o grouppipeline.IntakeQuotaObservation) {
				// A native observation cannot swap the selected historical episode. The
				// broker additionally CAS-checks its opaque token before retiring anything.
				if o.Key != native.Key || o.WorkspaceRevision != native.WorkspaceRevision || o.ConnectorID != native.ConnectorID || o.ConnectorRevision != native.ConnectorRevision || o.ObservationID != native.ReferenceID || o.Quota.Scope != native.Scope || o.Kind != grouppipeline.IntakeQuotaRecovered {
					logs.EmitDomainUnavailable(ctx)
					return
				}
				event, ok := intakeQuotaEventRef(o, true)
				if !ok {
					logs.EmitDomainUnavailable(ctx)
					return
				}
				scope := runtimeview.Scope{WorkspaceID: o.Key.WorkspaceID, WorkspaceRevision: o.WorkspaceRevision}
				raw, err := json.Marshal(event)
				if err != nil {
					logs.EmitDomainUnavailable(ctx)
					return
				}
				if _, err = diagnostics.ParseDomainEventRef(raw, scope); err != nil {
					logs.EmitDomainUnavailable(ctx)
					return
				}
				logs.EmitQuotaRecovery(ctx, hint, raw, scope)
			}
			return native, true
		},
	}
}

func intakeQuotaEventRef(o grouppipeline.IntakeQuotaObservation, recovery bool) (diagnostics.DomainEventRef, bool) {
	if !o.Available {
		return diagnostics.DomainEventRef{}, false
	}
	ref := deliveryEventRef(o.Key, o.WorkspaceRevision, o.ConnectorID, o.ConnectorRevision)
	ref.Reason, ref.Action = runtimeview.ReasonNone, runtimeview.ActionNone
	ref.Evidence = diagnostics.DomainEvidenceRef{ReferenceID: o.ObservationID}
	switch o.Kind {
	case grouppipeline.IntakeGroupApplied, grouppipeline.IntakeGroupPaused, grouppipeline.IntakeGroupRetiring:
		expected := "active"
		ref.Code, ref.Stage = "domain.group.applied", runtimeview.StageCollecting
		if o.Kind == grouppipeline.IntakeGroupPaused {
			expected = "not_running"
			ref.Code, ref.Stage = "domain.group.paused", runtimeview.StageBlocked
		}
		if o.Kind == grouppipeline.IntakeGroupRetiring {
			expected = "retiring"
			ref.Code = "domain.group.retiring"
		}
		if o.IntakeState != expected || o.Quota != (groupdelivery.QuotaStatus{}) {
			return diagnostics.DomainEventRef{}, false
		}
		ref.Evidence.Kind = "intake-status"
		ref.Evidence.Intake = &runtimeview.IntakeEvidence{GroupRevision: o.Key.GroupRevision, State: o.IntakeState}
	case grouppipeline.IntakeQuotaRefused, grouppipeline.IntakeQuotaRecovered:
		q := o.Quota
		if o.IntakeState != "" || !q.Configured || q.MaxBytes <= 0 || (q.Scope != groupdelivery.QuotaScopeGlobal && q.Scope != groupdelivery.QuotaScopeGroup) || (q.Scope == groupdelivery.QuotaScopeGlobal && q.GroupID != "") || (q.Scope == groupdelivery.QuotaScopeGroup && q.GroupID != o.Key.GroupID) {
			return diagnostics.DomainEventRef{}, false
		}
		ref.Code, ref.Stage = "domain.quota.refused", runtimeview.StageBlocked
		ref.Reason, ref.Action = runtimeview.ReasonUnavailable, runtimeview.ActionRefresh
		ref.Evidence.Kind = "quota-refusal"
		if o.Kind == grouppipeline.IntakeQuotaRecovered {
			if !recovery || (q.State != groupdelivery.QuotaOK && q.State != groupdelivery.QuotaWarning) || q.IntakeRefused {
				return diagnostics.DomainEventRef{}, false
			}
			ref.Code, ref.Stage = "domain.quota.recovered", runtimeview.StageAccepted
			ref.Reason, ref.Action = runtimeview.ReasonNone, runtimeview.ActionNone
			ref.Evidence.Kind = "quota-status"
		} else if q.State != groupdelivery.QuotaHard || !q.IntakeRefused {
			return diagnostics.DomainEventRef{}, false
		}
		// Do not copy LossRiskNotice or infer ratios, thresholds or per-device usage.
		ref.Evidence.Quota = &runtimeview.QuotaEvidence{Scope: string(q.Scope), State: q.State, UsedBytes: q.UsedBytes, MaxBytes: q.MaxBytes, IntakeRefused: q.IntakeRefused}
	default:
		return diagnostics.DomainEventRef{}, false
	}
	return ref, true
}
