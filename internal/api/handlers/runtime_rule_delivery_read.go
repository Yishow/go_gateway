package handlers

import (
	"context"
	"reflect"
	"sort"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

// RuntimeRevisionDeliveryReader observes durable delivery and revision-bound
// intake. Neither method authorizes activation, retry or repair.
type RuntimeRevisionDeliveryReader interface {
	DeliveryForRevision(context.Context, string, string) (*grouppipeline.DeliveryView, error)
	Status() []grouppipeline.GroupStatus
}

// WithRuleDeliveryServices adds optional delivery evidence to the shared read.
// Configuration-only callers retain null delivery when no reader is installed.
func (h *RuntimeHandler) WithRuleDeliveryServices(reader RuntimeRevisionDeliveryReader) *RuntimeHandler {
	h.ruleDelivery = reader
	return h
}

func (h *RuntimeHandler) readRuntimeRuleDelivery(ctx context.Context, saved *workspace.WriteGroup, snapshot *workspace.WriteGroupAppliedSnapshot, scope runtimeview.Scope, asOf time.Time) (*runtimeview.DeliveryEvidence, error) {
	if h.ruleDelivery == nil {
		return nil, nil
	}
	revision := ""
	if snapshot != nil {
		revision = snapshot.AppliedRevision
	}
	before := runtimeDeliveryStatuses(h.ruleDelivery.Status(), saved.ID, revision)
	raw, err := h.ruleDelivery.DeliveryForRevision(ctx, saved.ID, revision)
	if err != nil || raw == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	after := runtimeDeliveryStatuses(h.ruleDelivery.Status(), saved.ID, revision)
	if !reflect.DeepEqual(before, after) {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	// Status is a live observation, not derived from the retained configuration.
	var intake *runtimeview.IntakeEvidence
	if snapshot != nil && len(before) == 1 {
		status := before[0]
		compatible := status.State != "active" || (saved.Status != workspace.WriteGroupStatusDisabled && saved.Status != workspace.WriteGroupStatusDeleted && snapshot.Group.Status != workspace.WriteGroupStatusDisabled && snapshot.Group.Status != workspace.WriteGroupStatusDeleted && (snapshot.EndsAt.IsZero() || snapshot.EndsAt.After(asOf)))
		if compatible && runtimeSafeEnum(status.State, "active", "retiring", "blocked", "not_running", "unavailable") != "" {
			intake = &runtimeview.IntakeEvidence{GroupRevision: revision, State: status.State}
			if status.Reason != "" {
				reason := runtimeview.ParseReason(status.Reason)
				if reason != runtimeview.ReasonNone {
					intake.Reason = &reason
				}
			}
		}
	}
	return projectRuntimeRuleDelivery(raw, saved, snapshot, scope, intake)
}

// Detach and compare only the exact observed revision. Unrelated intake and
// volatile delivery counters cannot make a read spuriously stale.
func runtimeDeliveryStatuses(statuses []grouppipeline.GroupStatus, groupID, revision string) []grouppipeline.GroupStatus {
	result := []grouppipeline.GroupStatus{}
	if revision == "" {
		return result
	}
	for _, status := range statuses {
		if status.GroupID == groupID && status.Revision == revision {
			result = append(result, status)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].State != result[j].State {
			return result[i].State < result[j].State
		}
		return result[i].Reason < result[j].Reason
	})
	return result
}
