package main

import (
	"context"
	"encoding/json"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/diagnostics"
)

// The callback consumes only metadata copied with the actual poll. It does not
// look up current settings or expose sample, raw error or connection fields.
func newCollectorEventObserver(logs *diagnostics.Broker) func(collector.AcquisitionObservation) {
	if logs == nil {
		return nil
	}
	return func(observation collector.AcquisitionObservation) {
		ctx := context.Background()
		binding := observation.Binding
		scope := runtimeview.Scope{WorkspaceID: binding.WorkspaceID, WorkspaceRevision: binding.WorkspaceRevision, DeviceID: binding.DeviceID}
		// Validate the original rule identity without emitting it as a Step 1
		// locator or substituting the acquisition fingerprint for its revision.
		ruleLocator := runtimeview.EntityLocator{Step: 2, WorkspaceID: binding.WorkspaceID, DeviceID: &binding.DeviceID, RuleID: &binding.RuleID}
		if binding.DeviceID != observation.DeviceID || binding.PointID != observation.PointID ||
			binding.ConfigFingerprint == "" || binding.ConfigFingerprint != observation.ConfigFingerprint ||
			observation.ReceivedAt.IsZero() || runtimeview.ValidateEntityLocator(ruleLocator, scope) != nil ||
			!validCollectorEventQuality(observation.Quality, observation.QualityReason) {
			logs.EmitDomainUnavailable(ctx)
			return
		}
		code, reason, action := diagnostics.DomainEventCode("domain.collector.failed"), runtimeview.ReasonUnavailable, runtimeview.ActionRefresh
		if observation.Quality == schema.QualityGood {
			code, reason, action = "domain.collector.recovered", runtimeview.ReasonNone, runtimeview.ActionNone
		}
		ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: code,
			Binding: runtimeview.Binding{WorkspaceID: binding.WorkspaceID, WorkspaceRevision: binding.WorkspaceRevision, EntityRevision: binding.RuleRevision},
			Locator: runtimeview.EntityLocator{Step: 1, WorkspaceID: binding.WorkspaceID, DeviceID: &binding.DeviceID}, PointID: &binding.PointID,
			Stage: runtimeview.StageAcquisition, Reason: reason, Action: action,
			Evidence: diagnostics.DomainEvidenceRef{Kind: "acquisition-result", ReferenceID: observation.AcquisitionID}}
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

func validCollectorEventQuality(quality schema.QualityFlag, reason string) bool {
	switch quality {
	case schema.QualityGood:
		return reason == ""
	case schema.QualityBad, schema.QualityMissing, schema.QualityStale, schema.QualityInvalid, schema.QualityUncertain:
		return reason != "" // the existing normalizer supplies a reason for every non-good result
	default:
		return false
	}
}
