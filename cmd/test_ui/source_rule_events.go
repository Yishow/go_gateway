package main

import (
	"context"
	"encoding/json"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/diagnostics"
)

func newSourceRuleEventObserver(logs *diagnostics.Broker) datalinkruntime.SourceRuleObserver {
	if logs == nil {
		return nil
	}
	return func(ctx context.Context, observation sourcerule.RuntimeSourceObservation) {
		if !observation.Verified || !observation.Attempted || (observation.Result != sourcerule.SourceObservationApplied && observation.Result != sourcerule.SourceObservationRejected) {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		code, stage, reason, action := diagnostics.DomainEventCode("domain.source.rejected"), runtimeview.Stage("blocked"), runtimeview.Reason("unavailable"), runtimeview.Action("refresh")
		if observation.Result == sourcerule.SourceObservationApplied {
			if !observation.Applied {
				logs.EmitDomainUnavailable(ctx)
				return
			}
			code, stage, reason, action = "domain.source.applied", "accepted", "none", "none"
		}
		ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: code,
			Binding: runtimeview.Binding{WorkspaceID: observation.Scope.WorkspaceID, WorkspaceRevision: observation.Scope.WorkspaceRevision, EntityRevision: observation.RuleRevision},
			Locator: runtimeview.EntityLocator{Step: 2, WorkspaceID: observation.Scope.WorkspaceID, DeviceID: &observation.Scope.DeviceID, RuleID: &observation.RuleID},
			Stage:   stage, Reason: reason, Action: action, Evidence: diagnostics.DomainEvidenceRef{Kind: "source-runtime-transition", ReferenceID: observation.ReferenceID}}
		data, err := json.Marshal(ref)
		if err != nil {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		if _, err := diagnostics.ParseDomainEventRef(data, observation.Scope); err != nil {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		logs.EmitDomain(ctx, data, observation.Scope)
	}
}
