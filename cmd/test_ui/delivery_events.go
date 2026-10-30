package main

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"
)

const runtimeReceiptProvenance = "runtime"

// The shared read proves only observed workspace scope. Frozen output/owner
// identity and actual transition proof belong to the native producer, not here.
func newDeliveryEventObservers(logs *diagnostics.Broker, workspaces *workspace.Service) (groupdelivery.DeliveryObservationConfig, grouppipeline.OwnerObservationConfig) {
	if logs == nil {
		return groupdelivery.DeliveryObservationConfig{}, grouppipeline.OwnerObservationConfig{}
	}
	scopeRevision := func(ctx context.Context, key groupdelivery.GroupKey) (string, bool) {
		if workspaces == nil || !key.Valid() {
			return "", false
		}
		record, err := workspaces.ReadSnapshot(ctx)
		if err != nil || record == nil || record.ID != key.WorkspaceID || record.DatabaseSetupRevision == "" {
			return "", false
		}
		return record.DatabaseSetupRevision, true
	}
	delivery := groupdelivery.DeliveryObservationConfig{ScopeRevision: scopeRevision, Observe: newDeliveryEventObserver(logs)}
	// Owner observations are warn-only. Intentional filtering must not create
	// descriptors, proof reads, UUIDs or batch-overflow capture losses.
	if !logs.DomainWarningsEnabled() {
		return delivery, grouppipeline.OwnerObservationConfig{}
	}
	return delivery, grouppipeline.OwnerObservationConfig{ScopeRevision: scopeRevision, Observe: newOwnerEventObserver(logs), ObserveCaptured: newOwnerEventCapture(logs), RecordLimit: diagnostics.AdmissionRecords, OnDrop: logs.RecordDomainCaptureLoss}
}

func deliveryEventRef(key groupdelivery.GroupKey, revision, connectorID, connectorRevision string) diagnostics.DomainEventRef {
	return diagnostics.DomainEventRef{ContractVersion: 1,
		Binding:     runtimeview.Binding{WorkspaceID: key.WorkspaceID, WorkspaceRevision: revision, EntityRevision: key.GroupRevision},
		Locator:     runtimeview.EntityLocator{Step: 4, WorkspaceID: key.WorkspaceID, GroupID: &key.GroupID, GroupRevision: &key.GroupRevision},
		ConnectorID: &connectorID, ConnectorRevision: &connectorRevision}
}

// emitDeliveryEvent acknowledges validated local coalescer acceptance only, not
// broker admission, downstream retention, persistence or the native outcome.
func emitDeliveryEvent(logs *diagnostics.Broker, ref diagnostics.DomainEventRef) bool {
	ctx := context.Background()
	scope := runtimeview.Scope{WorkspaceID: ref.Binding.WorkspaceID, WorkspaceRevision: ref.Binding.WorkspaceRevision}
	raw, err := json.Marshal(ref)
	if err != nil {
		logs.EmitDomainUnavailable(ctx)
		return false
	}
	if _, err := diagnostics.ParseDomainEventRef(raw, scope); err != nil {
		logs.EmitDomainUnavailable(ctx)
		return false
	}
	return logs.EmitDomain(ctx, raw, scope)
}

func newDeliveryEventObserver(logs *diagnostics.Broker) func(groupdelivery.DeliveryObservation) {
	if logs == nil {
		return nil
	}
	return func(o groupdelivery.DeliveryObservation) {
		unavailable := func() { logs.EmitDomainUnavailable(context.Background()) }
		if !o.Available {
			unavailable()
			return
		}
		ref := deliveryEventRef(o.Key, o.WorkspaceRevision, o.ConnectorID, o.ConnectorRevision)
		ref.Reason, ref.Action = o.Reason, runtimeview.ActionForReason(o.Reason)
		ref.Evidence = diagnostics.DomainEvidenceRef{Kind: "durable-transition", ReferenceID: o.EffectKey}
		switch o.State {
		case groupdelivery.StateRetrying:
			ref.Code, ref.Stage = "domain.delivery.retrying", runtimeview.StageRetrying
		case groupdelivery.StateUnknown:
			ref.Code, ref.Stage = "domain.delivery.unknown", runtimeview.StageUnknown
		case groupdelivery.StateBlocked:
			ref.Code, ref.Stage = "domain.delivery.blocked", runtimeview.StageBlocked
		case groupdelivery.StateQuarantined:
			ref.Code, ref.Stage = "domain.delivery.quarantined", runtimeview.StageQuarantined
		case groupdelivery.StateCommitted:
			if !o.ReceiptMatched || o.CommittedAt.IsZero() || o.Reason != runtimeview.ReasonNone {
				unavailable()
				return
			}
			ref.Code, ref.Stage = "domain.delivery.first-confirmed", runtimeview.StageSQLCommitted
			ref.Evidence.Kind = "runtime-receipt"
			ref.Evidence.Receipt = &runtimeview.Receipt{GroupRevision: o.Key.GroupRevision, ConnectorRevision: o.ConnectorRevision, RecordID: o.RecordID, EffectKey: o.EffectKey, PayloadDigest: o.PayloadDigest, CommittedAt: o.CommittedAt.UTC().Format(time.RFC3339Nano), Provenance: runtimeReceiptProvenance}
		default:
			unavailable()
			return
		}
		if o.State != groupdelivery.StateCommitted && (o.ReceiptMatched || !o.CommittedAt.IsZero()) {
			unavailable()
			return
		}
		emitDeliveryEvent(logs, ref)
	}
}

func newOwnerEventObserver(logs *diagnostics.Broker) func(grouppipeline.OwnerObservation) {
	capture := newOwnerEventCapture(logs)
	if capture == nil {
		return nil
	}
	return func(o grouppipeline.OwnerObservation) { capture(o) }
}

func newOwnerEventCapture(logs *diagnostics.Broker) func(grouppipeline.OwnerObservation) bool {
	if logs == nil {
		return nil
	}
	return func(o grouppipeline.OwnerObservation) bool {
		if !o.Available {
			logs.EmitDomainUnavailable(context.Background())
			return false
		}
		ref := deliveryEventRef(o.Key, o.WorkspaceRevision, o.ConnectorID, o.ConnectorRevision)
		ref.Reason = runtimeview.ParseReason(o.Reason)
		ref.Action = runtimeview.ActionForReason(ref.Reason)
		count := strconv.FormatUint(o.PendingSamples, 10)
		ref.Evidence = diagnostics.DomainEvidenceRef{ReferenceID: o.ObservationID, PendingSamples: &count}
		switch o.Kind {
		case grouppipeline.OwnerCounterQuarantine:
			ref.Code, ref.Stage, ref.Evidence.Kind = "domain.counter-replay.quarantined", runtimeview.StageQuarantined, "counter-replay-inputs"
		case grouppipeline.OwnerRecordingRefusal:
			ref.Code, ref.Stage, ref.Evidence.Kind = "domain.recording.refused", runtimeview.StageBlocked, "recording-refusal"
		default:
			logs.EmitDomainUnavailable(context.Background())
			return false
		}
		return emitDeliveryEvent(logs, ref)
	}
}
