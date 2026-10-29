package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRecordingLifecycleActionsClassifySavedState(t *testing.T) {
	for _, mode := range []string{"draft", "applied-owned", "applied-unknown", "applied-guard-error", "stale-connector", "removed-connector", "disabled", "disabled-stale", "disabled-unknown", "deleted"} {
		t.Run(mode, func(t *testing.T) {
			db, service, group := lifecycleReadFixture(t, mode != "draft" && mode != "deleted")
			guardCalls := 0
			service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(ctx context.Context, tx *sql.Tx, actual *WriteGroup) error {
				guardCalls++
				require.Equal(t, group.ID, actual.ID)
				var revision string
				require.NoError(t, tx.QueryRowContext(ctx, `SELECT revision FROM write_groups WHERE id=?`, actual.ID).Scan(&revision))
				require.Equal(t, group.Revision, revision)
				if mode == "applied-guard-error" || mode == "disabled-unknown" {
					return errors.New("private SQL endpoint password sample value")
				}
				return nil
			}))
			if mode == "applied-unknown" {
				service.WithBacklogOwnershipGuard(nil)
			}
			if mode == "stale-connector" || mode == "disabled-stale" {
				_, err := db.Exec(`UPDATE database_connectors SET identity_revision='C2' WHERE id=?`, group.Destination.ConnectorID)
				require.NoError(t, err)
			}
			if mode == "removed-connector" {
				_, err := db.Exec(`DELETE FROM database_connectors WHERE id=?`, group.Destination.ConnectorID)
				require.NoError(t, err)
			}
			if mode == "disabled" || mode == "disabled-stale" || mode == "disabled-unknown" {
				_, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
				require.NoError(t, err)
				group = lifecycleReadGroup(t, service, group.ID)
			}
			if mode == "deleted" {
				_, err := service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
				require.NoError(t, err)
				group = lifecycleReadGroup(t, service, group.ID)
			}
			inspector := &writeGroupReadinessInspector{}
			service.WithTableInspector(inspector)
			before := lifecycleReadState(t, db)
			ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
			defer cancel()
			view, err := service.Lifecycle(ctx, "  "+group.ID+"  ")
			require.NoError(t, err)
			require.Equal(t, group.WorkspaceID, view.WorkspaceID)
			require.Equal(t, group.Revision, view.GroupRevision)
			require.Equal(t, group.AppliedRevision, view.AppliedRevision)
			require.Equal(t, group.Destination.ConnectorID, view.ConnectorID)
			require.Equal(t, "connector-revision-1", view.ConnectorRevision)
			require.Equal(t, group.Status, view.Status)
			require.Equal(t, []string{"edit", "stop", "remove", "recreate", "transition"}, lifecycleActionNames(view))
			remove := view.Actions[2]
			switch mode {
			case "applied-unknown", "applied-guard-error":
				require.Equal(t, WriteGroupLifecycleAction{Action: "remove", ReasonCode: "ownership_unavailable", NextAction: "stop"}, remove)
			case "stale-connector":
				require.Equal(t, WriteGroupLifecycleAction{Action: "remove", ReasonCode: "connector_identity_stale", NextAction: "stop"}, remove)
			case "removed-connector":
				require.Equal(t, WriteGroupLifecycleAction{Action: "remove", ReasonCode: "connector_unavailable", NextAction: "stop"}, remove)
			case "disabled-stale":
				require.Equal(t, WriteGroupLifecycleAction{Action: "remove", ReasonCode: "connector_identity_stale", NextAction: "review_connector"}, remove)
			case "disabled-unknown":
				require.Equal(t, WriteGroupLifecycleAction{Action: "remove", ReasonCode: "ownership_unavailable", NextAction: "inspect_ownership"}, remove)
			case "deleted":
				for _, action := range view.Actions[:3] {
					require.False(t, action.Allowed)
					require.Equal(t, "recording_removed", action.ReasonCode)
				}
			default:
				require.Equal(t, WriteGroupLifecycleAction{Action: "remove", Allowed: true}, remove)
			}
			if group.Status != WriteGroupStatusDeleted {
				require.True(t, view.Actions[0].Allowed)
				require.Equal(t, group.Status != WriteGroupStatusDisabled, view.Actions[1].Allowed)
			}
			if group.Status == WriteGroupStatusDisabled {
				require.Equal(t, "already_stopped", view.Actions[1].ReasonCode)
			}
			for _, action := range view.Actions[3:] {
				require.False(t, action.Allowed)
				require.True(t, action.RequiresReview)
				require.Equal(t, "canonical_required", action.ReasonCode)
			}
			if group.AppliedRevision == "" || mode == "stale-connector" || mode == "disabled-stale" || mode == "removed-connector" {
				require.Zero(t, guardCalls)
			}
			require.Zero(t, inspector.calls)
			require.Equal(t, before, lifecycleReadState(t, db))
			raw, err := json.Marshal(view)
			require.NoError(t, err)
			for _, secret := range []string{"private", "password", "persisted-db", "persisted-schema", "raw_values", "connection_config"} {
				require.NotContains(t, string(raw), secret)
			}
			var decoded map[string]any
			require.NoError(t, json.Unmarshal(raw, &decoded))
			for _, action := range decoded["actions"].([]any) {
				require.Len(t, action.(map[string]any), 5)
			}
		})
	}
}

func TestRecordingLifecycleActionsPreserveAppliedDraftAndCAS(t *testing.T) {
	db, service, original := lifecycleReadFixture(t, true)
	var appliedBefore string
	require.NoError(t, db.QueryRow(`SELECT payload FROM write_group_versions WHERE group_id=? AND group_revision=?`, original.ID, original.AppliedRevision).Scan(&appliedBefore))
	draft := cloneWriteGroup(original)
	draft.RowPolicy.IntervalSeconds = 30
	mutation := lifecycleReadMutation(t, service, original)
	mutation.Group = draft
	saved, err := service.Update(t.Context(), original.ID, mutation)
	require.NoError(t, err)
	view, err := service.Lifecycle(t.Context(), original.ID)
	require.NoError(t, err)
	require.Equal(t, saved.Group.Revision, view.GroupRevision)
	require.Equal(t, original.AppliedRevision, view.AppliedRevision)
	require.NotEqual(t, view.GroupRevision, view.AppliedRevision)
	require.EqualValues(t, 30, saved.Group.RowPolicy.IntervalSeconds)
	applied, err := service.ResolveAppliedAt(t.Context(), original.ID, time.Date(2026, 10, 2, 12, 0, 20, 0, time.UTC))
	require.NoError(t, err)
	require.EqualValues(t, 15, applied.Group.RowPolicy.IntervalSeconds)
	var appliedAfter string
	require.NoError(t, db.QueryRow(`SELECT payload FROM write_group_versions WHERE group_id=? AND group_revision=?`, original.ID, original.AppliedRevision).Scan(&appliedAfter))
	require.Equal(t, appliedBefore, appliedAfter)
	before := lifecycleReadState(t, db)
	_, err = service.Disable(t.Context(), original.ID, mutation)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.Equal(t, before, lifecycleReadState(t, db), "availability never grants stale mutation authority")
}
