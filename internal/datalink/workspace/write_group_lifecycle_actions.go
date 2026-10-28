package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// WriteGroupLifecycleAction is advisory availability, never mutation authority.
// Existing mutation endpoints still recheck their scope, revisions and guards.
type WriteGroupLifecycleAction struct {
	Action         string `json:"action"`
	Allowed        bool   `json:"allowed"`
	ReasonCode     string `json:"reason_code"`
	NextAction     string `json:"next_action"`
	RequiresReview bool   `json:"requires_review"`
}

// WriteGroupLifecycle distinguishes the saved draft from its applied revision.
// Status describes saved configuration, not runtime health or delivery success.
// Connector identity belongs to the saved group, not the live connector draft.
type WriteGroupLifecycle struct {
	WorkspaceID       string                      `json:"workspace_id"`
	WorkspaceRevision string                      `json:"workspace_revision"`
	GroupID           string                      `json:"group_id"`
	GroupRevision     string                      `json:"group_revision"`
	AppliedRevision   string                      `json:"applied_revision"`
	Status            WriteGroupStatus            `json:"status"`
	ConnectorID       string                      `json:"connector_id"`
	ConnectorRevision string                      `json:"connector_revision"`
	Actions           []WriteGroupLifecycleAction `json:"actions"`
	History           *WriteGroupLifecycleHistory `json:"history"`
	SuccessorSeed     *WriteGroupSuccessorSeed    `json:"successor_seed,omitempty"`
}

// Lifecycle reads local configuration and ownership evidence in one transaction.
// It never probes the destination, invokes DDL, or changes a saved/applied group.
func (s *WriteGroupService) Lifecycle(ctx context.Context, id string) (result *WriteGroupLifecycle, err error) {
	if err = s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, writeGroupNotFound("read recording lifecycle")
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read recording lifecycle", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback recording lifecycle read: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("read recording lifecycle")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read recording lifecycle", err)
	}
	group, err := s.repo.GetInTx(ctx, setupTx.SQLTx(), record.ID, id)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read recording lifecycle", err)
	}
	if err := s.repo.requireUnreservedWriteGroupInTx(ctx, setupTx.SQLTx(), record.ID, id); err != nil && !isSuccessorReservedError(err) {
		return nil, err
	}
	switch group.Status {
	case WriteGroupStatusDraft, WriteGroupStatusReady, WriteGroupStatusRunning, WriteGroupStatusDisabled, WriteGroupStatusDeleted:
	default:
		return nil, ErrWriteGroupLifecycleBlocked
	}
	history, err := s.repo.lifecycleHistoryInTx(ctx, setupTx.SQLTx(), record.ID, id)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read recording lifecycle history", err)
	}
	result = &WriteGroupLifecycle{
		WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision,
		GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.AppliedRevision,
		Status: group.Status, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
		Actions: s.lifecycleActions(ctx, setupTx.SQLTx(), group),
		History: history,
	}
	result.SuccessorSeed = s.lifecycleSuccessorSeed(ctx, setupTx.SQLTx(), group, history, result.Actions)
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WriteGroupService) lifecycleActions(ctx context.Context, tx *sql.Tx, group *WriteGroup) []WriteGroupLifecycleAction {
	actions := []WriteGroupLifecycleAction{
		{Action: "edit", Allowed: true},
		{Action: "stop", Allowed: true},
		{Action: "remove", Allowed: true},
		{Action: "recreate", ReasonCode: "source_definition_unavailable", NextAction: "review_lifecycle", RequiresReview: true},
		{Action: "transition", ReasonCode: "source_definition_unavailable", NextAction: "review_lifecycle", RequiresReview: true},
	}
	if isSuccessorReservedError(s.repo.requireUnreservedWriteGroupInTx(ctx, tx, group.WorkspaceID, group.ID)) {
		for i := range actions {
			actions[i].Allowed, actions[i].ReasonCode, actions[i].NextAction = false, "successor_reserved", "review_lifecycle"
		}
		return actions
	}
	if group.Status == WriteGroupStatusDeleted {
		for i := 0; i < 3; i++ {
			actions[i].Allowed, actions[i].ReasonCode, actions[i].NextAction = false, "recording_removed", "review_lifecycle"
		}
		return actions
	}
	if group.Status == WriteGroupStatusDisabled {
		actions[1].Allowed, actions[1].ReasonCode = false, "already_stopped"
	}
	remove := &actions[2]
	connectorErr := validateLiveWriteGroupConnector(ctx, tx, s.repo, group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	switch {
	case errors.Is(connectorErr, ErrWriteGroupConnectorRevisionConflict):
		remove.ReasonCode, remove.NextAction = "connector_identity_stale", "review_connector"
	case connectorErr != nil:
		remove.ReasonCode, remove.NextAction = "connector_unavailable", "review_connector"
	case group.AppliedRevision != "":
		if s.backlogGuard == nil || s.backlogGuard.CheckWriteGroupBacklog(ctx, tx, group) != nil {
			remove.ReasonCode, remove.NextAction = "ownership_unavailable", "inspect_ownership"
		}
	}
	if remove.ReasonCode != "" {
		remove.Allowed = false
		if actions[1].Allowed {
			remove.NextAction = "stop"
		}
	}
	return actions
}
