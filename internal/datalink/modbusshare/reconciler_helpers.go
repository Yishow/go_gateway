package modbusshare

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"time"
)

func stagedMappings(req ReconcileRequest, capacity int, updatedAt time.Time) map[string]TagMirrorMapping {
	staged := make(map[string]TagMirrorMapping, len(req.DesiredMappings))
	for _, dm := range req.DesiredMappings {
		span := DataTypeSpan(dm.DataType)
		staged[dm.TagID] = TagMirrorMapping{
			WorkspaceID: req.WorkspaceID, SourceRuleID: dm.SourceRuleID, SourceRuleRevision: dm.SourceRuleRevision,
			TagID: dm.TagID, MappingID: dm.MappingID, Register: dm.ZeroBasedRegister, ShareStartRegister: ZeroBasedToHuman(dm.ZeroBasedRegister),
			ZeroBasedRegister: dm.ZeroBasedRegister, SpanRegisters: span, StrideRegisters: dm.StrideRegisters,
			CapacityRegisters: capacity, DataType: dm.DataType, TagKey: dm.TagKey, DisplayName: dm.DisplayName,
			UpdatedAt: updatedAt,
		}
	}
	return staged
}

func mappingOwnershipValid(ctx context.Context, mapping TagMirrorMapping, workspaceID string, validator func(context.Context, DesiredMapping) error) bool {
	if validator == nil {
		return false
	}
	if mapping.WorkspaceID != workspaceID || strings.TrimSpace(mapping.SourceRuleID) == "" || strings.TrimSpace(mapping.SourceRuleRevision) == "" {
		return false
	}
	return validator(ctx, DesiredMapping{WorkspaceID: mapping.WorkspaceID, SourceRuleID: mapping.SourceRuleID, SourceRuleRevision: mapping.SourceRuleRevision, TagID: mapping.TagID, MappingID: mapping.MappingID, DataType: mapping.DataType, ShareStartRegister: mapping.ShareStartRegister, ZeroBasedRegister: mapping.ZeroBasedRegister, SpanRegisters: mapping.SpanRegisters, StrideRegisters: mapping.StrideRegisters, CapacityRegisters: mapping.CapacityRegisters, TagKey: mapping.TagKey, DisplayName: mapping.DisplayName}) == nil
}

func markHydrationFailed(svc *Service, hydration HydrationState) {
	hydration.State = HydrationStateFailed
	hydration.Readiness = false
	svc.SetHydrationState(hydration)
}

func (r *Reconciler) replayOutcome(key, workspaceRevision, settingsRevision string) (ReconcileOutcome, bool) {
	r.idempotencyMu.Lock()
	defer r.idempotencyMu.Unlock()
	record, ok := r.idempotency[key]
	if !ok || record.postWorkspaceRevision != workspaceRevision || record.postSettingsRevision != settingsRevision {
		return ReconcileOutcome{}, false
	}
	return record.outcome, true
}

func (r *Reconciler) rememberOutcome(key string, outcome ReconcileOutcome, workspaceRevision, settingsRevision string) {
	r.idempotencyMu.Lock()
	r.idempotency[key] = idempotencyRecord{outcome: outcome, postWorkspaceRevision: workspaceRevision, postSettingsRevision: settingsRevision}
	r.idempotencyMu.Unlock()
}

// ReconcilePersisted restores durable desired mappings without changing the
// workspace revision.
func (r *Reconciler) ReconcilePersisted(ctx context.Context, workspaceID, workspaceRevision, settingsRevision string) (ReconcileOutcome, error) {
	if strings.TrimSpace(workspaceID) == "" || strings.TrimSpace(workspaceRevision) == "" || strings.TrimSpace(settingsRevision) == "" {
		return ReconcileOutcome{Outcome: HydrationStateFailed}, NewError(ErrCodeRevisionRequired, "workspace and settings revisions are required for restore", false)
	}
	durable, ok := r.revStore.(DurableDesiredMappingStore)
	if !ok {
		return ReconcileOutcome{Outcome: HydrationStateFailed}, &Error{Code: ErrCodeProjectionRequired, Message: "durable desired mapping store is required for restore", Retryable: false, Action: "complete durable Share hydration"}
	}
	mappings, err := durable.GetDesiredMappings(ctx, workspaceID)
	if err != nil {
		return ReconcileOutcome{Outcome: HydrationStateFailed}, err
	}
	hydration, err := r.svc.CheckHydration(ctx)
	if err != nil {
		return ReconcileOutcome{Outcome: HydrationStateFailed}, err
	}
	return r.Reconcile(ctx, ReconcileRequest{WorkspaceID: workspaceID, ExpectedWorkspaceRevision: workspaceRevision, ExpectedSettingsRevision: settingsRevision, ReadinessToken: hydration.ReadinessToken, DesiredMappings: mappings, Restore: true})
}

func reconcileKey(req ReconcileRequest) (string, error) {
	type keyedMapping struct {
		mapping DesiredMapping
		key     string
	}
	keyed := make([]keyedMapping, 0, len(req.DesiredMappings))
	for _, mapping := range req.DesiredMappings {
		encoded, err := json.Marshal(mapping)
		if err != nil {
			return "", err
		}
		keyed = append(keyed, keyedMapping{mapping: mapping, key: string(encoded)})
	}
	sort.Slice(keyed, func(i, j int) bool { return keyed[i].key < keyed[j].key })
	mappings := make([]DesiredMapping, 0, len(keyed))
	for _, item := range keyed {
		mappings = append(mappings, item.mapping)
	}
	payload, err := json.Marshal(struct {
		WorkspaceID, WorkspaceRevision, SettingsRevision, ReadinessToken string
		Mappings                                                         []DesiredMapping
		Restore                                                          bool
	}{req.WorkspaceID, req.ExpectedWorkspaceRevision, req.ExpectedSettingsRevision, req.ReadinessToken, mappings, req.Restore})
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func isMappingsEquivalent(current []TagMirrorMapping, desired map[string]TagMirrorMapping) bool {
	if len(current) != len(desired) {
		return false
	}
	for _, c := range current {
		d, exists := desired[c.TagID]
		if !exists || c.WorkspaceID != d.WorkspaceID || c.SourceRuleID != d.SourceRuleID ||
			c.SourceRuleRevision != d.SourceRuleRevision || c.TagID != d.TagID || c.MappingID != d.MappingID ||
			c.Register != d.Register || c.ShareStartRegister != d.ShareStartRegister ||
			c.ZeroBasedRegister != d.ZeroBasedRegister || c.DataType != d.DataType ||
			c.SpanRegisters != d.SpanRegisters || c.StrideRegisters != d.StrideRegisters ||
			c.CapacityRegisters != d.CapacityRegisters || c.TagKey != d.TagKey ||
			c.DisplayName != d.DisplayName {
			return false
		}
	}
	return true
}

func mappingProjectionChanged(current, desired TagMirrorMapping) bool {
	return current.WorkspaceID != desired.WorkspaceID ||
		current.SourceRuleID != desired.SourceRuleID ||
		current.SourceRuleRevision != desired.SourceRuleRevision ||
		current.MappingID != desired.MappingID ||
		current.TagKey != desired.TagKey || current.DisplayName != desired.DisplayName ||
		current.DataType != desired.DataType || current.ZeroBasedRegister != desired.ZeroBasedRegister ||
		current.ShareStartRegister != desired.ShareStartRegister || current.Register != desired.Register ||
		current.SpanRegisters != desired.SpanRegisters || current.StrideRegisters != desired.StrideRegisters ||
		current.CapacityRegisters != desired.CapacityRegisters
}
