package modbusshare

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type reconcileState struct {
	durable        DurableDesiredMappingStore
	key            string
	settings       Settings
	hydration      HydrationState
	currentRev     string
	lifecycleEpoch uint64
}

type reconcileTransaction struct {
	staged           map[string]TagMirrorMapping
	currentMappings  []TagMirrorMapping
	oldMemory        []byte
	oldDesired       []DesiredMapping
	newRev           string
	removedSpans     []SpanRange
	invalidatedSpans []SpanRange
	unknownOwnership bool
}

func (r *Reconciler) prepareReconcile(ctx context.Context, req ReconcileRequest, durable DurableDesiredMappingStore) (reconcileState, ReconcileOutcome, bool, error) {
	var err error
	key, err := reconcileKey(req)
	if err != nil {
		return reconcileState{}, ReconcileOutcome{Outcome: HydrationStateFailed}, true, fmt.Errorf("build reconcile idempotency key: %w", err)
	}
	settings, authoritativeSettings, lifecycleEpoch := r.svc.lifecycleSnapshot()
	hydration, hydrationErr := r.svc.CheckHydration(ctx)
	state := reconcileState{durable: durable, key: key, settings: settings, hydration: hydration, lifecycleEpoch: lifecycleEpoch}
	if !req.Restore && (hydrationErr != nil || hydration.WorkspaceID == "" || hydration.State == HydrationStatePending || hydration.State == HydrationStateFailed || !hydration.Readiness) {
		outcome := ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: hydration.WorkspaceRevision, SettingsRevision: settings.SettingsRevision}
		err := &Error{Code: ErrCodeHydrationRequired, Message: "share workspace hydration is not ready", Retryable: true, WorkspaceRevision: hydration.WorkspaceRevision, SettingsRevision: settings.SettingsRevision, Action: "hydrate the workspace again before reconcile"}
		return state, outcome, true, err
	}
	if authoritativeSettings {
		if !settings.Enabled {
			return state, ReconcileOutcome{Outcome: shareStateDisabled, SettingsRevision: settings.SettingsRevision}, true, NewError(ErrCodeDisabled, "modbus share is disabled in global settings", false)
		}
		if req.ExpectedSettingsRevision != "" && req.ExpectedSettingsRevision != settings.SettingsRevision {
			outcome := ReconcileOutcome{Outcome: HydrationStateFailed, SettingsRevision: settings.SettingsRevision}
			err := &Error{Code: ErrCodeRevisionConflict, Message: settingsRevisionConflictMessage, Retryable: true, SettingsRevision: settings.SettingsRevision}
			return state, outcome, true, err
		}
	}
	if r.revStore == nil {
		return state, ReconcileOutcome{}, false, nil
	}
	rev, isDirty, err := r.revStore.GetRevision(ctx, req.WorkspaceID)
	if err != nil {
		return state, ReconcileOutcome{Outcome: HydrationStateFailed}, true, fmt.Errorf("failed to get workspace revision: %w", err)
	}
	state.currentRev = rev
	if isDirty {
		outcome := ReconcileOutcome{Outcome: dirtyUnknownState, NewWorkspaceRevision: rev}
		err := &Error{Code: ErrCodeDirtyUnknown, Message: "workspace is in dirty_unknown state; recovery reconcile required", Retryable: true, WorkspaceRevision: rev}
		return state, outcome, true, err
	}
	return state, ReconcileOutcome{}, false, nil
}

func (r *Reconciler) stageTransaction(ctx context.Context, req ReconcileRequest, state reconcileState) (reconcileTransaction, ReconcileOutcome, error) {
	capRegs := r.svc.effectiveCapacityRegisters()
	allMappings := r.svc.ListMappings()
	foreignMappings := make([]TagMirrorMapping, 0, len(allMappings))
	for _, mapping := range allMappings {
		if mapping.WorkspaceID != req.WorkspaceID {
			foreignMappings = append(foreignMappings, mapping)
		}
	}
	for _, dm := range req.DesiredMappings {
		if dm.WorkspaceID != "" && dm.WorkspaceID != req.WorkspaceID {
			return reconcileTransaction{}, ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.currentRev}, NewError(ErrCodeWorkspaceScope, "desired mapping belongs to another workspace", false)
		}
		if r.ownership != nil {
			if err := r.ownership(ctx, dm); err != nil {
				return reconcileTransaction{}, ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.currentRev}, err
			}
		}
	}
	if err := ValidateDesiredMappingsAgainstExisting(req.DesiredMappings, foreignMappings, capRegs); err != nil {
		outcome := ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.currentRev}
		return reconcileTransaction{}, outcome, err
	}
	transaction := reconcileTransaction{staged: stagedMappings(req, capRegs, time.Now().UTC())}
	transaction.currentMappings = make([]TagMirrorMapping, 0, len(allMappings))
	for _, mapping := range allMappings {
		if mapping.WorkspaceID == req.WorkspaceID {
			transaction.currentMappings = append(transaction.currentMappings, mapping)
		}
	}
	transaction.oldMemory = r.svc.MemorySnapshot()
	var err error
	transaction.oldDesired, err = state.durable.GetDesiredMappings(ctx, req.WorkspaceID)
	if err != nil {
		outcome := ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.currentRev}
		failure := &Error{Code: ErrCodeReconcileFailed, Message: "failed to load durable desired mappings", Retryable: true, Action: retryReconcileAction}
		return reconcileTransaction{}, outcome, failure
	}
	transaction.newRev = uuid.NewString()
	if req.Restore {
		transaction.newRev = state.currentRev
	}
	for _, mapping := range transaction.currentMappings {
		newMapping, stillPresent := transaction.staged[mapping.TagID]
		if !stillPresent {
			if !r.mappingCanBeRemoved(ctx, mapping, req.WorkspaceID) {
				transaction.invalidatedSpans = append(transaction.invalidatedSpans, spanForMapping(mapping))
				transaction.unknownOwnership = true
				transaction.staged[mapping.TagID] = mapping
				continue
			}
			transaction.removedSpans = append(transaction.removedSpans, spanForMapping(mapping))
			continue
		}
		if mappingProjectionChanged(mapping, newMapping) {
			if !mappingOwnershipValid(ctx, mapping, req.WorkspaceID, r.ownership) {
				transaction.invalidatedSpans = append(transaction.invalidatedSpans, spanForMapping(mapping))
				transaction.unknownOwnership = true
				transaction.staged[mapping.TagID] = mapping
				continue
			}
			transaction.invalidatedSpans = append(transaction.invalidatedSpans, spanForMapping(mapping))
		}
	}
	if transaction.unknownOwnership {
		markHydrationFailed(r.svc, state.hydration)
		outcome := ReconcileOutcome{Outcome: "invalidated_unknown", NewWorkspaceRevision: state.currentRev, SettingsRevision: state.settings.SettingsRevision, Mappings: transaction.currentMappings, InvalidatedCount: len(transaction.invalidatedSpans), InvalidatedSpans: transaction.invalidatedSpans, Diagnostics: []Diagnostic{{Code: ErrCodeWorkspaceScope, Severity: "warning", Message: "obsolete mapping ownership could not be proven; span retained", Retryable: true, Action: "repair persisted workspace ownership before cleanup"}}}
		err := &Error{Code: ErrCodeWorkspaceScope, Message: "obsolete mapping ownership cannot be proven", Retryable: true, WorkspaceRevision: state.currentRev, Action: "repair persisted workspace ownership before retrying"}
		return transaction, outcome, err
	}
	return transaction, ReconcileOutcome{}, nil
}

func (r *Reconciler) mappingCanBeRemoved(ctx context.Context, mapping TagMirrorMapping, workspaceID string) bool {
	if r.ownership == nil {
		return true
	}
	ownershipErr := r.ownership(ctx, desiredMappingFromMirror(mapping))
	return ownershipErr == nil && mapping.WorkspaceID == workspaceID && mapping.SourceRuleID != "" && mapping.SourceRuleRevision != ""
}

func desiredMappingFromMirror(mapping TagMirrorMapping) DesiredMapping {
	return DesiredMapping{WorkspaceID: mapping.WorkspaceID, SourceRuleID: mapping.SourceRuleID, SourceRuleRevision: mapping.SourceRuleRevision, TagID: mapping.TagID, MappingID: mapping.MappingID, DataType: mapping.DataType, TagKey: mapping.TagKey, DisplayName: mapping.DisplayName, ZeroBasedRegister: mapping.ZeroBasedRegister, CapacityRegisters: mapping.CapacityRegisters, ShareStartRegister: mapping.ShareStartRegister, SpanRegisters: mapping.SpanRegisters, StrideRegisters: mapping.StrideRegisters}
}

func spanForMapping(mapping TagMirrorMapping) SpanRange {
	return SpanRange{TagID: mapping.TagID, Start: mapping.ZeroBasedRegister, Count: mapping.SpanRegisters, DataType: mapping.DataType}
}

func (r *Reconciler) applyTransaction(ctx context.Context, req ReconcileRequest, state reconcileState, transaction reconcileTransaction) (ReconcileOutcome, error) {
	// Serialize the durable commit and runtime projection with settings
	// transitions. The epoch check rejects a staged transaction that crossed a
	// disable or settings update, before it can mutate either boundary.
	r.svc.lifecycleMu.Lock()
	defer r.svc.lifecycleMu.Unlock()
	r.svc.mu.RLock()
	currentSettings := r.svc.settings
	authoritativeSettings := r.svc.authoritativeSettings
	currentEpoch := r.svc.lifecycleEpoch
	r.svc.mu.RUnlock()
	if currentEpoch != state.lifecycleEpoch || (authoritativeSettings && !currentSettings.Enabled) {
		outcome := ReconcileOutcome{Outcome: shareStateDisabled, NewWorkspaceRevision: state.currentRev, SettingsRevision: currentSettings.SettingsRevision}
		return outcome, NewError(ErrCodeDisabled, "modbus share is disabled in global settings", false)
	}
	if err := state.durable.CommitDesiredMappings(ctx, req.WorkspaceID, req.ExpectedWorkspaceRevision, transaction.newRev, req.DesiredMappings); err != nil {
		return ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.currentRev}, err
	}
	oldMappings := r.svc.ReplaceMappingsForWorkspace(req.WorkspaceID, transaction.staged)
	if err := r.svc.projectionSwapError(transaction.staged); err != nil {
		return r.rollbackTransaction(ctx, req.WorkspaceID, state, transaction, oldMappings, "runtime projection swap failed")
	}
	for _, old := range transaction.currentMappings {
		newMapping, stillPresent := transaction.staged[old.TagID]
		if stillPresent && newMapping.ZeroBasedRegister == old.ZeroBasedRegister && newMapping.SpanRegisters == old.SpanRegisters && newMapping.DataType == old.DataType {
			continue
		}
		if err := r.svc.ClearSpan(old.ZeroBasedRegister, old.SpanRegisters); err != nil {
			return r.rollbackTransaction(ctx, req.WorkspaceID, state, transaction, oldMappings, "failed to clear obsolete register span")
		}
	}
	return ReconcileOutcome{}, nil
}

func (r *Reconciler) rollbackTransaction(ctx context.Context, workspaceID string, state reconcileState, transaction reconcileTransaction, oldMappings map[string]TagMirrorMapping, message string) (ReconcileOutcome, error) {
	r.svc.ReplaceMappings(oldMappings)
	restoreErr := r.svc.RestoreMemory(transaction.oldMemory)
	var durableErr error
	if state.durable != nil {
		durableErr = state.durable.CommitDesiredMappings(ctx, workspaceID, transaction.newRev, state.currentRev, transaction.oldDesired)
	}
	if restoreErr != nil || durableErr != nil {
		return r.rollbackUnknown(ctx, workspaceID, state.hydration)
	}
	return ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.currentRev}, &Error{Code: ErrCodeReconcileFailed, Message: message, Retryable: true, Action: retryReconcileAction}
}

func (r *Reconciler) buildAppliedOutcome(state reconcileState, transaction reconcileTransaction, workspaceID string) ReconcileOutcome {
	outcome := ReconcileOutcome{Outcome: "applied", NewWorkspaceRevision: transaction.newRev, SettingsRevision: state.settings.SettingsRevision, AppliedCount: len(transaction.staged), RemovedCount: len(transaction.removedSpans), InvalidatedCount: len(transaction.invalidatedSpans), RemovedSpans: transaction.removedSpans, InvalidatedSpans: transaction.invalidatedSpans, Mappings: r.svc.ListMappingsByWorkspace(workspaceID)}
	if !transaction.unknownOwnership {
		outcome.NewReadinessToken = uuid.NewString()
		if state.hydration.ReadinessToken != "" {
			state.hydration.WorkspaceRevision = outcome.NewWorkspaceRevision
			state.hydration.ReadinessToken = outcome.NewReadinessToken
			r.svc.SetHydrationState(state.hydration)
		}
	}
	return outcome
}
