package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// RuntimeRecordingGroupReader reads the canonical saved snapshot, never effective facts.
type RuntimeRecordingGroupReader interface {
	List(context.Context) (*workspace.WriteGroupListResult, error)
}

// RuntimeAppliedConfigurationReader verifies an immutable snapshot at one server observation.
type RuntimeAppliedConfigurationReader interface {
	ResolveAppliedConfigurationAt(context.Context, string, time.Time) (*workspace.WriteGroupAppliedSnapshot, error)
}

// WithRuleEvidenceServices installs read-only authorities, without starting a pipeline.
func (h *RuntimeHandler) WithRuleEvidenceServices(groups RuntimeRecordingGroupReader, applied RuntimeAppliedConfigurationReader) *RuntimeHandler {
	h.ruleGroups, h.ruleApplied = groups, applied
	return h
}

// WithRuleEvidenceClock sets the server observation clock, not a client-selectable time.
func (h *RuntimeHandler) WithRuleEvidenceClock(now func() time.Time) *RuntimeHandler {
	h.ruleEvidenceNow = now
	return h
}

func runtimeEvidenceError(reason runtimeview.Reason) error {
	return &runtimeview.ContractError{Reason: reason, Action: runtimeview.ActionForReason(reason)}
}

func runtimeEvidenceRequestScope(c *gin.Context) (*runtimeview.Scope, bool) {
	query, err := url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		renderRuntimeHandlerError(c, http.StatusBadRequest, "unavailable", false)
		return nil, false
	}
	keys := []string{"workspace_id", "workspace_revision", "device_id"}
	present := false
	for _, key := range keys {
		_, exists := query[key]
		present = present || exists
	}
	if !present {
		return nil, true
	}
	for _, key := range keys {
		values := query[key]
		if len(values) != 1 || values[0] == "" || runtimeSafeID(values[0]) != values[0] {
			renderRuntimeHandlerError(c, http.StatusBadRequest, "unavailable", false)
			return nil, false
		}
	}
	return &runtimeview.Scope{WorkspaceID: query.Get(keys[0]), WorkspaceRevision: query.Get(keys[1]), DeviceID: query.Get(keys[2])}, true
}

func validateRuntimeEvidenceScope(scope runtimeview.Scope, record *workspace.Record) error {
	if record == nil || record.ID != scope.WorkspaceID || !slices.Contains(record.OrderedDeviceIDs, scope.DeviceID) {
		return runtimeEvidenceError(runtimeview.ReasonForeign)
	}
	if record.DatabaseSetupRevision != scope.WorkspaceRevision {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	return nil
}

func renderRuntimeEvidenceError(c *gin.Context, err error) {
	status, reason := http.StatusServiceUnavailable, runtimeview.ReasonUnavailable
	if runtimeview.IsContractError(err, runtimeview.ReasonForeign) {
		status, reason = http.StatusForbidden, runtimeview.ReasonForeign
	} else if runtimeview.IsContractError(err, runtimeview.ReasonStale) {
		status, reason = http.StatusConflict, runtimeview.ReasonStale
	}
	renderRuntimeHandlerError(c, status, string(reason), reason != runtimeview.ReasonForeign)
}

// Scoped reads bypass the legacy setup/readiness chain, which can bootstrap a
// workspace. Empty compatibility fields are not synthetic saved/effective data.
func (h *RuntimeHandler) scopedRuntimeRuleEvidence(c *gin.Context, scope runtimeview.Scope, selection *runtimeSourceSelection) {
	if h.workspaceSvc == nil || h.deviceSvc == nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	record, err := h.workspaceSvc.ReadSnapshot(c.Request.Context())
	if err != nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	if err := validateRuntimeEvidenceScope(scope, record); err != nil {
		renderRuntimeEvidenceError(c, err)
		return
	}
	if _, err := h.deviceSvc.GetByID(c.Request.Context(), scope.DeviceID); err != nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
		return
	}
	evidence, err := h.readRuntimeRuleEvidence(c.Request.Context(), record, scope, selection)
	if err != nil {
		renderRuntimeEvidenceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: runtimeWorkspaceContextResponse{
		WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, Devices: []runtimeWorkspaceContextDeviceEntry{}, RuleEvidence: evidence,
	}})
}

func (h *RuntimeHandler) readRuntimeRuleEvidence(ctx context.Context, record *workspace.Record, scope runtimeview.Scope, selection *runtimeSourceSelection) (*runtimeview.RuleEvidenceView, error) {
	if h.ruleGroups == nil || h.ruleApplied == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if err := validateRuntimeEvidenceScope(scope, record); err != nil {
		return nil, err
	}
	asOf := time.Now().UTC()
	if h.ruleEvidenceNow != nil {
		asOf = h.ruleEvidenceNow().UTC()
	}
	if asOf.IsZero() {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	var sourceBefore *runtimeSourceRecords
	if selection != nil {
		var err error
		sourceBefore, err = h.readRuntimeSourceRecords(ctx, scope, *selection)
		if err != nil {
			return nil, err
		}
	}
	before, err := h.ruleGroups.List(ctx)
	if err != nil || before == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if before.WorkspaceID != scope.WorkspaceID || before.WorkspaceRevision != scope.WorkspaceRevision {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	observed := asOf.Format(time.RFC3339Nano)
	view := runtimeview.RuleEvidenceView{ContractVersion: 1, WorkspaceID: scope.WorkspaceID,
		WorkspaceRevision: scope.WorkspaceRevision, DeviceID: scope.DeviceID,
		ObservedAt: &observed, Freshness: "fresh", Groups: []runtimeview.GroupView{}}
	if sourceBefore != nil {
		view.Source, err = projectRuntimeSource(sourceBefore, scope)
		if err != nil {
			return nil, err
		}
	}
	snapshots := make(map[string]*workspace.WriteGroupAppliedSnapshot, len(before.Groups))
	seen := make(map[string]bool, len(before.Groups))
	for _, group := range before.Groups {
		if group == nil || group.WorkspaceID != scope.WorkspaceID || seen[group.ID] {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		seen[group.ID] = true
		snapshot, err := h.runtimeRuleSnapshotAt(ctx, group, scope, asOf)
		if err != nil {
			return nil, err
		}
		snapshots[group.ID] = snapshot
		if !runtimeGroupHasDevice(group, scope.DeviceID) && (snapshot == nil || !runtimeGroupHasDevice(snapshot.Group, scope.DeviceID)) {
			continue
		}
		projected, err := projectRuntimeRecordingGroup(group, snapshot, scope, observed)
		if err != nil {
			return nil, err
		}
		projected.Delivery, err = h.readRuntimeRuleDelivery(ctx, group, snapshot, scope, asOf)
		if err != nil {
			return nil, err
		}
		view.Groups = append(view.Groups, projected)
	}
	// Check immutable evidence again at the SAME time, including dynamic revision ends.
	// A stable saved revision alone cannot prove that the applied read stayed stable.
	for _, group := range before.Groups {
		after, err := h.runtimeRuleSnapshotAt(ctx, group, scope, asOf)
		if err != nil || !reflect.DeepEqual(snapshots[group.ID], after) {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
	}
	after, err := h.ruleGroups.List(ctx)
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if !reflect.DeepEqual(before, after) {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	if selection != nil {
		sourceAfter, err := h.readRuntimeSourceRecords(ctx, scope, *selection)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(sourceBefore, sourceAfter) {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
	}
	afterRecord, err := h.workspaceSvc.ReadSnapshot(ctx)
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if !reflect.DeepEqual(record, afterRecord) {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	encoded, err := runtimeview.MarshalRuleEvidence(view, scope)
	if err != nil {
		return nil, err
	}
	return runtimeview.ParseRuleEvidence(encoded, scope)
}

func (h *RuntimeHandler) runtimeRuleSnapshotAt(ctx context.Context, group *workspace.WriteGroup, scope runtimeview.Scope, asOf time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	if group.AppliedRevision == "" {
		return nil, nil
	}
	snapshot, err := h.ruleApplied.ResolveAppliedConfigurationAt(ctx, group.ID, asOf)
	if errors.Is(err, grouppipeline.ErrAppliedConfigurationStale) || errors.Is(err, workspace.ErrWriteGroupNotFound) {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	if err != nil || snapshot == nil {
		return nil, nil // Missing proof is unavailable, never a saved-as-effective fallback.
	}
	if snapshot.Group == nil || snapshot.RuntimeLayout == nil || snapshot.WorkspaceID != scope.WorkspaceID ||
		snapshot.GroupID != group.ID || snapshot.Group.ID != group.ID || snapshot.Group.WorkspaceID != scope.WorkspaceID ||
		snapshot.GroupRevision != snapshot.AppliedRevision || snapshot.Group.Revision != snapshot.AppliedRevision ||
		snapshot.Group.AppliedRevision != snapshot.AppliedRevision || snapshot.EffectiveAt.IsZero() || snapshot.EffectiveAt.After(asOf) {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	return snapshot, nil
}

func runtimeGroupHasDevice(group *workspace.WriteGroup, deviceID string) bool {
	if group == nil {
		return false
	}
	for _, member := range group.Members {
		if member.DeviceID == deviceID {
			return true
		}
	}
	return false
}
