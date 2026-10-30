package grouppipeline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// Both destination DDL/catalog and local operation/lineage/outbox are real
// disposable SQLite. Faults are injected only in negative tests.
type storageObservationFixture struct {
	*longPipeline
	recovery *StorageRecovery
	ledger   *recordingplan.Service
	item     groupdelivery.OutboxItem
	target   groupdelivery.Target
	scope    recordingplan.StorageScope
	events   []StorageOperationObservation
	reads    int
}

func newStorageObservationFixture(t *testing.T) *storageObservationFixture {
	t.Helper()
	f := &storageObservationFixture{longPipeline: newAveragePipeline(t, shapeMeanLong)}
	f.target = groupdelivery.Target{DB: f.dest.target, Kind: schema.DatabaseConnectorKindSQLite}
	f.ledger = recordingplan.NewService(recordingplan.NewSQLRepository(f.local))
	layout, err := workspace.LongStorageLayout(f.group, "sqlite")
	require.NoError(t, err)
	inspector := func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		return inspectOpenTarget(ctx, f.target, "main", table, false)
	}
	token, err := f.ledger.PrepareSchemaPreview(t.Context(), recordingplan.SchemaPreviewScope{
		WorkspaceID: f.group.WorkspaceID, WorkspaceRevision: "observed-wrev", PlanID: f.group.ID, PlanRevision: f.group.Revision,
		ConnectorID: f.group.Destination.ConnectorID, ConnectorRevision: f.group.Destination.ConnectorRevision, Dialect: "sqlite", Schema: "main", Database: f.group.Destination.Database, GroupLayout: layout, SourceDigest: "source-fixture",
	}, inspector)
	require.NoError(t, err)
	op, _, err := f.ledger.ApplySchemaPreview(t.Context(), token, recordingplan.SchemaApplyTarget{Inspect: inspector, Execute: func(context.Context, []string) (recordingplan.SchemaExecution, error) {
		return recordingplan.SchemaExecution{}, errors.New("already-compatible preview must not execute")
	}})
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status)
	enabled := true
	f.group.WritePolicy.RecreateMissing = &enabled
	f.group.Destination.SchemaRevision = op.OperationID
	f.group.Destination.SchemaDigest = op.VerifiedDigest
	f.snapshot.RuntimeLayout.SchemaDigest = op.VerifiedDigest
	raw, err := json.Marshal(f.snapshot)
	require.NoError(t, err)
	_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=? AND group_revision=?`, string(raw), f.group.ID, f.group.Revision)
	require.NoError(t, err)
	f.scope = workspace.StorageScopeOf(f.group)
	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	digest := sha256.Sum256(encoded)
	_, err = f.local.Exec(`INSERT INTO write_group_storage_scopes(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation) VALUES(?,?,?,?,?,?,?,?,?,1)`, f.scope.WorkspaceID, f.scope.ConnectorID, f.scope.Database, f.scope.TableSchema, f.scope.StorageMode, f.scope.LayoutVersion, layout.TableName, layout.OwnerColumn, hex.EncodeToString(digest[:]))
	require.NoError(t, err)
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), f.local, f.scope, recordingplan.StorageConfirmation{Layout: *layout, ConnectorRevision: f.group.Destination.ConnectorRevision, OperationID: op.OperationID, VerifiedDigest: op.VerifiedDigest, GroupID: f.group.ID}, time.Now()))
	require.NoError(t, f.reading(0, "input-before-table-loss", at(1), 100))
	f.closeUntil(at(15))
	var effect string
	require.NoError(t, f.local.QueryRow(`SELECT effect_key FROM wg_delivery_outbox ORDER BY effect_key LIMIT 1`).Scan(&effect))
	f.item, err = f.store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	f.recovery = NewStorageRecoveryWithObservation(f.store, f.ledger, f.dest, StorageObservationConfig{
		RecordLimit: 1024, InfoEnabled: true, WarnEnabled: true,
		ScopeRevision: func(_ context.Context, key groupdelivery.GroupKey) (string, bool) {
			f.reads++
			return "observed-wrev", key == f.item.Key
		},
		Observe: func(o StorageOperationObservation) { f.events = append(f.events, o) },
	})
	return f
}
func (f *storageObservationFixture) drop() {
	f.t.Helper()
	_, err := f.target.DB.Exec(`DROP TABLE "` + f.table + `"`)
	require.NoError(f.t, err)
}
func (f *storageObservationFixture) restore() (groupdelivery.RestoreDecision, error) {
	return f.recovery.Restore(f.t.Context(), f.item, f.target, errors.New("no such table: private-destination"))
}
func (f *storageObservationFixture) lineage() *recordingplan.StorageLineage {
	f.t.Helper()
	l, err := recordingplan.ReadStorageLineage(f.t.Context(), f.local, f.scope)
	require.NoError(f.t, err)
	return l
}
