package workspace

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type receiptRuntimeInspector struct{ inspected bool }

func (i *receiptRuntimeInspector) InspectTable(ctx context.Context, connectorID, tableSchema, table string) (*dbtarget.TableInspection, error) {
	if table == dbtarget.EffectReceiptTable {
		i.inspected = true
		return &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists}, nil
	}
	return writeGroupLifecycleInspector{}.InspectTable(ctx, connectorID, tableSchema, table)
}

func TestApplyPersistsVerifiedRuntimeVersionAndMigrationKeepsIt(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { _ = db.Close() })
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	workspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	applied, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	read := func() string {
		var payload string
		require.NoError(t, db.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions
			WHERE group_id = ? AND group_revision = ?`, applied.Group.ID, applied.Group.AppliedRevision).Scan(&payload))
		return payload
	}
	payload := read()
	var frozen WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(payload), &frozen))
	require.Equal(t, applied.Group.AppliedRevision, frozen.Group.Revision)
	require.Equal(t, dbtarget.SQLDialectPostgres, frozen.RuntimeLayout.Dialect)
	require.Equal(t, schema.DataTypeFloat32, frozen.RuntimeLayout.TagTypes[fixture.tagID])
	require.NotEmpty(t, frozen.RuntimeLayout.SchemaDigest)
	require.NotContains(t, payload, "connection_config")
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	require.Equal(t, payload, read(), "migration replay preserves the frozen descriptor")
}

func TestApplyRuntimeVersionRollsBackWithWorkspaceFailure(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { _ = db.Close() })
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	workspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `CREATE TRIGGER reject_runtime_version_workspace BEFORE UPDATE ON system_settings
		WHEN NEW.key = 'studio_v2_workspace' BEGIN SELECT RAISE(ABORT, 'save-failed'); END`)
	require.NoError(t, err)
	_, err = service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.Error(t, err)
	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_runtime_versions`).Scan(&count))
	require.Zero(t, count)
}

func TestDisplayOnlyApplyBackfillsOriginalRuntimeVersion(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { _ = db.Close() })
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	current, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	mutation := WriteGroupMutation{WorkspaceID: fixture.workspaceID,
		ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision:     created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1"}
	applied, err := service.Apply(ctx, created.Group.ID, mutation)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	draft := cloneWriteGroup(applied.Group)
	draft.Name = "Only a display rename"
	mutation.ExpectedWorkspaceRevision, mutation.ExpectedGroupRevision, mutation.Group = applied.WorkspaceRevision, applied.Group.Revision, draft
	renamed, err := service.Update(ctx, draft.ID, mutation)
	require.NoError(t, err)
	mutation.ExpectedWorkspaceRevision, mutation.ExpectedGroupRevision, mutation.Group = renamed.WorkspaceRevision, renamed.Group.Revision, nil
	result, err := service.Apply(ctx, renamed.Group.ID, mutation)
	require.NoError(t, err)
	require.Equal(t, applied.Group.AppliedRevision, result.Group.AppliedRevision)
	var payload string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE group_id = ?`, result.Group.ID).Scan(&payload))
	var frozen WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(payload), &frozen))
	require.Equal(t, applied.Group.Revision, frozen.Group.Revision)
	require.Equal(t, applied.Group.Name, frozen.Group.Name, "the descriptor keeps the original immutable group")
	require.NotNil(t, frozen.RuntimeLayout)
}

func TestApplyPersistsVerifiedReceiptCapability(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { _ = db.Close() })
	service, fixture, created := newReadinessGroup(ctx, t, db)
	inspector := &receiptRuntimeInspector{}
	service.WithTableInspector(inspector)
	current, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	draft := cloneWriteGroup(created.Group)
	draft.WritePolicy.DedupeCapability = "receipt"
	updated, err := service.Update(ctx, draft.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision: draft.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: draft,
	})
	require.NoError(t, err)
	applied, err := service.Apply(ctx, draft.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: updated.WorkspaceRevision,
		ExpectedGroupRevision: updated.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.True(t, inspector.inspected)
	var payload string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE group_id = ?`, applied.Group.ID).Scan(&payload))
	var frozen WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(payload), &frozen))
	require.True(t, frozen.RuntimeLayout.ReceiptTableReady)
}
