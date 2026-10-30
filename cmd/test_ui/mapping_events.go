package main

import (
	"context"
	"encoding/json"

	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/diagnostics"
)

// Only the actual executed, frozen mapping observation can reach this adapter.
// It carries no value, raw transform, error or configuration to serialize.
func newMappingEventObserver(logs *diagnostics.Broker) datalinkruntime.MappingObserver {
	if logs == nil {
		return nil
	}
	return func(observation datalinkruntime.MappingObservation) {
		ctx := context.Background()
		code, stage, reason, action := diagnostics.DomainEventCode("domain.mapping.failed"), runtimeview.StageBlocked, runtimeview.ReasonUnavailable, runtimeview.ActionRefresh
		switch observation.Result {
		case datalinkruntime.MappingObservationFailed:
		case datalinkruntime.MappingObservationSucceeded:
			code, stage, reason, action = "domain.mapping.recovered", runtimeview.StageAcquisition, runtimeview.ReasonNone, runtimeview.ActionNone
		default:
			logs.EmitDomainUnavailable(ctx)
			return
		}
		if observation.ReceivedAt.IsZero() {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		scope := runtimeview.Scope{WorkspaceID: observation.WorkspaceID, WorkspaceRevision: observation.WorkspaceRevision, DeviceID: observation.DeviceID}
		ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: code,
			Binding: runtimeview.Binding{WorkspaceID: observation.WorkspaceID, WorkspaceRevision: observation.WorkspaceRevision, EntityRevision: observation.MappingRevision},
			Locator: runtimeview.EntityLocator{Step: 3, WorkspaceID: observation.WorkspaceID, DeviceID: &observation.DeviceID, MappingID: &observation.MappingID},
			PointID: &observation.PointID, TagID: &observation.TagID, Stage: stage, Reason: reason, Action: action,
			Evidence: diagnostics.DomainEvidenceRef{Kind: "mapping-result", ReferenceID: observation.AcquisitionID}}
		data, err := json.Marshal(ref)
		if err != nil {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		if _, err := diagnostics.ParseDomainEventRef(data, scope); err != nil {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		logs.EmitDomain(ctx, data, scope)
	}
}
