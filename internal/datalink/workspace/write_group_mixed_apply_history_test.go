package workspace

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// mixedHistoryConfirm obtains consent through real preview, DDL and operation
// ledger verification. It does not manufacture a schema or runtime proof.
func mixedHistoryConfirm(t *testing.T, h *circuitProfileHarness, group *WriteGroup) *WriteGroup {
	t.Helper()
	ctx := t.Context()
	request := h.mutation(ctx, group)
	request.ExpectedGroupRevision = group.Revision
	scope, err := h.service.ManagedSchemaScope(ctx, group.ID, request)
	require.NoError(t, err)
	inspector := dbtarget.NewManagedTableInspector(dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(h.db)), h.db)
	h.service.WithManagedTableInspector(inspector)
	inspect := func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		found, err := inspector.InspectTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		return recordingPlanTableInspection(found), nil
	}
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(h.db))
	token, err := ledger.PrepareSchemaPreview(ctx, scope, inspect)
	require.NoError(t, err)
	verified, err := ledger.ValidatePreviewTokenForApply(ctx, token.Token, scope)
	require.NoError(t, err)
	operation, _, err := ledger.ApplySchemaPreview(ctx, verified, recordingplan.SchemaApplyTarget{Inspect: inspect, Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
		done, err := inspector.ExecuteSchemaStatements(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, statements)
		if done == nil {
			return recordingplan.SchemaExecution{}, err
		}
		return recordingplan.SchemaExecution{Committed: done.Committed, RolledBack: done.RolledBack}, err
	}})
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status)
	confirmed, err := h.service.ConfirmManagedSchemaWithStorage(ctx, group.ID, request, operation.OperationID, operation.VerifiedDigest, &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: len(token.Statements) > 0})
	require.NoError(t, err)
	ready, err := h.service.Readiness(ctx, group.ID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready)
	return confirmed.Group
}

func mixedHistoryRequest(t *testing.T, h *circuitProfileHarness, group *WriteGroup) WriteGroupMutation {
	t.Helper()
	request := h.mutation(t.Context(), nil)
	request.ExpectedGroupRevision = group.Revision
	return request
}

func mixedHistoryFixture(t *testing.T, now *time.Time) (*circuitProfileHarness, *WriteGroup) {
	t.Helper()
	h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "workspace.db"))
	target := filepath.Join(t.TempDir(), "target.db")
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `PRAGMA user_version=0`)
	require.NoError(t, err)
	require.NoError(t, db.Close())
	setSQLiteConnectorConfig(t, h.db, h.fixture.connectorID, target)
	h.service.WithClock(func() time.Time { return *now })
	draft := mixedDraft(h, WriteGroupMemberCalculationLatest, 15, "mixed_history")
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	draft.Destination.Database = target
	created, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.NoError(t, err)
	group := mixedHistoryConfirm(t, h, created.Group)
	applied, err := h.service.Apply(t.Context(), group.ID, mixedHistoryRequest(t, h, group))
	require.NoError(t, err)
	version, err := h.service.repo.versionByRevisionInTx(t.Context(), h.db, h.fixture.workspaceID, group.ID, applied.Group.AppliedRevision)
	require.NoError(t, err)
	require.Equal(t, *now, version.EffectiveAt, "only never-applied records use the transaction instant")
	return h, applied.Group
}

func TestMixedScopedApplyHistoryAndRevisionCutover(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		name := "active revision"
		if disabled {
			name = "disabled history"
		}
		t.Run(name, func(t *testing.T) {
			now := time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC)
			h, original := mixedHistoryFixture(t, &now)
			current := original
			if disabled {
				stopped, err := h.service.Disable(t.Context(), current.ID, mixedHistoryRequest(t, h, current))
				require.NoError(t, err)
				require.Equal(t, WriteGroupStatusDisabled, stopped.Group.Status)
				require.Equal(t, original.AppliedRevision, stopped.Group.AppliedRevision)
				current = stopped.Group
			}
			now = now.Add(10 * time.Second)
			changed := cloneWriteGroup(current)
			changed.Members[0].Cadence.Interval = 60
			update := h.mutation(t.Context(), changed)
			update.ExpectedGroupRevision = current.Revision
			saved, err := h.service.Update(t.Context(), current.ID, update)
			require.NoError(t, err)
			fresh := mixedHistoryConfirm(t, h, saved.Group)
			request := mixedHistoryRequest(t, h, fresh)
			applied, err := h.service.Apply(t.Context(), fresh.ID, request)
			require.NoError(t, err)
			version, err := h.service.repo.versionByRevisionInTx(t.Context(), h.db, h.fixture.workspaceID, fresh.ID, applied.Group.AppliedRevision)
			require.NoError(t, err)
			require.Equal(t, time.Date(2026, 10, 31, 12, 1, 0, 0, time.UTC), version.EffectiveAt)
			require.NotEqual(t, now, version.EffectiveAt, "existing history cannot take first-start shortcut")
			require.Zero(t, version.Group.RowPolicy.IntervalSeconds)
			require.Equal(t, int64(60), version.Group.Members[0].Cadence.Interval)
			old, err := h.service.repo.versionByRevisionInTx(t.Context(), h.db, h.fixture.workspaceID, original.ID, original.AppliedRevision)
			require.NoError(t, err)
			require.Equal(t, int64(15), old.Group.Members[0].Cadence.Interval)
			require.Equal(t, time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC), old.EffectiveAt)
			before := h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, original.ID)
			stale := mixedHistoryRequest(t, h, fresh)
			stale.ExpectedGroupRevision = original.Revision
			_, err = h.service.Apply(t.Context(), fresh.ID, stale)
			require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
			require.Equal(t, before, h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, original.ID))
			got, err := h.service.Get(t.Context(), fresh.ID)
			require.NoError(t, err)
			require.Equal(t, applied.Group, got.Group)
		})
	}
}

func TestMixedScopedApplyRejectsOrphanedAppliedHistory(t *testing.T) {
	for _, missing := range []string{"snapshot", "applied identity"} {
		t.Run(missing, func(t *testing.T) {
			now := time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC)
			h, original := mixedHistoryFixture(t, &now)
			stopped, err := h.service.Disable(t.Context(), original.ID, mixedHistoryRequest(t, h, original))
			require.NoError(t, err)
			stopped.Group = mixedHistoryConfirm(t, h, stopped.Group)
			// Deliberate corruption after a real successful Apply: neither missing
			// current identity nor a missing snapshot erases durable application history.
			if missing == "snapshot" {
				_, err = h.db.ExecContext(t.Context(), `DELETE FROM write_group_versions WHERE group_id=?`, original.ID)
			} else {
				_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET applied_revision='' WHERE id=?`, original.ID)
			}
			require.NoError(t, err)
			before := h.countRows(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, original.ID)
			now = now.Add(10 * time.Second)
			_, err = h.service.Apply(t.Context(), original.ID, mixedHistoryRequest(t, h, stopped.Group))
			require.ErrorIs(t, err, ErrWriteGroupAppliedRevisionUnavailable)
			require.Equal(t, before, h.countRows(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, original.ID))
			got, err := h.service.Get(t.Context(), original.ID)
			require.NoError(t, err)
			require.Equal(t, WriteGroupStatusDisabled, got.Group.Status)
		})
	}
}
