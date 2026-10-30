package main

import (
	"context"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"
)

// Native owners prove the original frozen identity and newly settled ledger
// result. This adapter adds no schema lookup, operation ID or repair authority.
func newOperationEventObservers(logs *diagnostics.Broker, scopeRevision func(context.Context, groupdelivery.GroupKey) (string, bool)) (grouppipeline.StorageObservationConfig, workspace.CutoverObservationConfig) {
	if logs == nil || !logs.DomainInfoEnabled() && !logs.DomainWarningsEnabled() {
		return grouppipeline.StorageObservationConfig{}, workspace.CutoverObservationConfig{}
	}
	storage := grouppipeline.StorageObservationConfig{ScopeRevision: scopeRevision,
		InfoEnabled: logs.DomainInfoEnabled(), WarnEnabled: logs.DomainWarningsEnabled(),
		RecordLimit: diagnostics.AdmissionRecords, OnDrop: logs.RecordDomainCaptureLoss,
	}
	storage.Observe = func(o grouppipeline.StorageOperationObservation) {
		unavailable := func() { logs.EmitDomainUnavailable(context.Background()) }
		if !o.Available {
			unavailable()
			return
		}
		ref := deliveryEventRef(o.Key, o.WorkspaceRevision, o.ConnectorID, o.ConnectorRevision)
		ref.Evidence = diagnostics.DomainEvidenceRef{Kind: "operation-result", ReferenceID: o.OperationID}
		switch o.Outcome {
		case grouppipeline.StorageObservationCompleted:
			if o.Reason != string(runtimeview.ReasonNone) {
				unavailable()
				return
			}
			ref.Code, ref.Stage, ref.Reason = "domain.repair.completed", runtimeview.StageAccepted, runtimeview.ReasonNone
		case grouppipeline.StorageObservationBlocked:
			if o.Reason != string(runtimeview.ReasonUnavailable) && o.Reason != string(runtimeview.ReasonMissingProof) {
				unavailable()
				return
			}
			ref.Code, ref.Stage, ref.Reason = "domain.repair.blocked", runtimeview.StageBlocked, runtimeview.Reason(o.Reason)
		default:
			unavailable()
			return
		}
		ref.Action = runtimeview.ActionForReason(ref.Reason)
		emitDeliveryEvent(logs, ref)
	}
	cutover := workspace.CutoverObservationConfig{InfoEnabled: storage.InfoEnabled, WarnEnabled: storage.WarnEnabled,
		RecordLimit: storage.RecordLimit, OnDrop: storage.OnDrop,
		ScopeRevision: func(ctx context.Context, w, g, r, _, _ string) (string, bool) {
			if scopeRevision == nil {
				return "", false
			}
			return scopeRevision(ctx, groupdelivery.GroupKey{WorkspaceID: w, GroupID: g, GroupRevision: r})
		},
	}
	cutover.Observe = func(ctx context.Context, o workspace.CutoverObservation) {
		if !o.Available {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		key := groupdelivery.GroupKey{WorkspaceID: o.WorkspaceID, GroupID: o.GroupID, GroupRevision: o.GroupRevision}
		ref := deliveryEventRef(key, o.WorkspaceRevision, o.ConnectorID, o.ConnectorRevision)
		ref.Evidence = diagnostics.DomainEvidenceRef{Kind: "operation-result", ReferenceID: o.OperationID}
		switch o.Outcome {
		case "applied":
			ref.Code, ref.Stage, ref.Reason = "domain.cutover.applied", runtimeview.StageAccepted, runtimeview.ReasonNone
		case "rejected":
			ref.Code, ref.Stage, ref.Reason = "domain.cutover.rejected", runtimeview.StageBlocked, runtimeview.ReasonUnavailable
		default:
			logs.EmitDomainUnavailable(ctx)
			return
		}
		ref.Action = runtimeview.ActionForReason(ref.Reason)
		emitDeliveryEvent(logs, ref) //nolint:contextcheck // A settled cutover remains observable after request cancellation.
	}
	return storage, cutover
}
