package workspace

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"testing"
	"time"
)

func TestRecreateFrozenSourceEngineAllowsEditedOrRemovedOldConnector(t *testing.T) {
	for _, state := range []string{"edited", "removed"} {
		t.Run(state, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			old := cloneWriteGroup(req.Group)
			old.Status = WriteGroupStatusReady
			old.AppliedRevision = old.Revision
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			effective := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
			require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, old, effective))
			require.NoError(t, h.service.saveRuntimeVersionInTx(t.Context(), tx, &WriteGroupAppliedSnapshot{WorkspaceID: old.WorkspaceID, GroupID: old.ID, GroupRevision: old.Revision, AppliedRevision: old.AppliedRevision, EffectiveAt: effective, Group: old, RuntimeLayout: &WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialect("sqlite"), ReceiptTableReady: true}}))
			_, err = tx.Exec(`UPDATE write_groups SET applied_revision=? WHERE id=?`, old.AppliedRevision, old.ID)
			require.NoError(t, err)
			require.NoError(t, tx.Commit())
			_, err = h.db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('new-target','new-target','sqlite','{"database":"new-database","schema":"main"}','target-1','ok',1)`)
			require.NoError(t, err)
			if state == "edited" {
				_, err = h.db.Exec(`UPDATE database_connectors SET kind='postgres',identity_revision='edited-old',connection_config='{}' WHERE id=?`, old.Destination.ConnectorID)
			} else {
				_, err = h.db.Exec(`DELETE FROM database_connectors WHERE id=?`, old.Destination.ConnectorID)
			}
			require.NoError(t, err)
			req.Group.Destination.ConnectorID = "new-target"
			req.Group.Destination.ConnectorRevision = "target-1"
			req.Group.Destination.Database = "new-database"
			req.Group.Destination.TableSchema = "main"
			req.Group.Destination.TableName = ""
			req.ExpectedTargetConnectorRevision = "target-1"
			calls := 0
			h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { calls++; return nil }))
			probe := &reservedProbeCounter{}
			h.service.WithManagedTableInspector(probe)
			oldBefore, err := h.service.Get(t.Context(), old.ID)
			require.NoError(t, err)
			versions := lifecyclePreservedData(t, h.db)
			delete(versions, "write_group_storage_scopes")
			result, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, "new-target", result.Group.Destination.ConnectorID)
			require.Empty(t, result.Group.Destination.SchemaRevision)
			require.Equal(t, 1, calls)
			require.Zero(t, probe.calls)
			oldAfter, err := h.service.Get(t.Context(), old.ID)
			require.NoError(t, err)
			require.Equal(t, oldBefore.Group, oldAfter.Group)
			after := lifecyclePreservedData(t, h.db)
			delete(after, "write_group_storage_scopes")
			require.Equal(t, versions, after)
		})
	}
}
