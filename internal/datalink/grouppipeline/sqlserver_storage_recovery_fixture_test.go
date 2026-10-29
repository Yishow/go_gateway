package grouppipeline

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// The catalog/DDL destination is scripted; only the gateway's local SQLite
// journal and operation ledger are real. This is not live SQL Server evidence.
type sqlServerRepairFixture struct {
	mu             sync.Mutex
	t              *testing.T
	store          *groupdelivery.Store
	pipeline       *Pipeline
	recovery       *StorageRecovery
	snapshot       *workspace.WriteGroupAppliedSnapshot
	item           groupdelivery.OutboxItem
	target         groupdelivery.Target
	current        *dbtarget.SQLServerManagedRowProof
	missing        bool
	receiptMissing bool
	inspectErr     error
	executeErr     error
	rolledBack     bool
	executions     int
	requests       []dbtarget.SQLServerManagedRepairRequest
}

func newSQLServerRepairFixture(t *testing.T) *sqlServerRepairFixture {
	return newSQLServerRepairFixtureWith(t, nil)
}

func newSQLServerRepairFixtureWith(t *testing.T, edit func(*workspace.WriteGroupAppliedSnapshot)) *sqlServerRepairFixture {
	t.Helper()
	base := newFixture(t)
	pipeline := base.pipeline()
	snap := sqlServerRuntimeSnapshot(t)
	enabled := true
	snap.Group.WritePolicy.RecreateMissing = &enabled
	if edit != nil {
		edit(snap)
	}
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
	table, owner, ok := workspace.LongStorageIdentity(snap.Group)
	require.True(t, ok)
	e := snap.RuntimeLayout.SQLServer.Group.Evidence
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, recordingplan.SQLServerGroupSchemaOptions{
		Database: e.Database, Schema: e.Schema, Collation: "Latin1_General_100_BIN2", ServerVersion: e.ServerVersion,
		ServerEdition: e.ServerEdition, EngineEdition: e.EngineEdition, CompatibilityLevel: e.CompatibilityLevel, CapabilityRevision: e.CapabilityRevision,
	})
	require.NoError(t, err)
	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	hash := sha256.Sum256(encoded)
	scope := workspace.StorageScopeOf(snap.Group)
	_, err = base.store.DB().ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema,
		scope.StorageMode, scope.LayoutVersion, table, owner, hex.EncodeToString(hash[:]))
	require.NoError(t, err)
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), base.store.DB(), scope, recordingplan.StorageConfirmation{
		Layout: layout, ConnectorRevision: snap.Group.Destination.ConnectorRevision, OperationID: snap.Group.Destination.SchemaRevision,
		VerifiedDigest: snap.Group.Destination.SchemaDigest, GroupID: snap.Group.ID,
	}, time.Now()))
	repo := recordingplan.NewSQLRepository(base.store.DB())
	now := time.Now().UTC()
	_, _, err = repo.ClaimSchemaOperation(t.Context(), &recordingplan.SchemaOperation{
		OperationID: snap.Group.Destination.SchemaRevision, Token: "fixture-confirmation", WorkspaceID: snap.WorkspaceID,
		ScopeKey: (&recordingplan.SchemaPreviewToken{WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID,
			Database: scope.Database, Schema: scope.TableSchema, GroupLayout: &layout}).ScopeKey(),
		Owner: "claim-fixture", Action: "schema_apply", Status: recordingplan.SchemaOperationRunning, CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = repo.FinishSchemaOperation(t.Context(), snap.Group.Destination.SchemaRevision, "claim-fixture", recordingplan.SchemaOperationResult{
		Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: snap.Group.Destination.SchemaDigest,
	})
	require.NoError(t, err)
	built, reason := pipeline.buildFrozen(t.Context(), snap, at(0), time.Time{})
	require.Empty(t, reason)
	require.NoError(t, built.boundary.AcceptSample(t.Context(), envelope(snap.Group, "sample", at(1), 21.5)))
	require.NoError(t, built.boundary.Tick(t.Context(), at(15)))
	var key string
	require.NoError(t, base.store.DB().QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox`).Scan(&key))
	item, err := base.store.GetOutbox(t.Context(), key)
	require.NoError(t, err)
	resolver, _, opened := sqlServerResolverFixture(t, snap)
	target, err := resolver.Resolve(t.Context(), item)
	require.NoError(t, err)
	opened.opened.SQLServer.SupportsMissingRepair = true
	recovery := NewStorageRecovery(base.store, recordingplan.NewService(repo), nil)
	f := &sqlServerRepairFixture{t: t, store: base.store, pipeline: pipeline, recovery: recovery, snapshot: snap, item: item,
		target: target, current: dbtarget.CloneSQLServerManagedRowProof(snap.RuntimeLayout.SQLServer), missing: true}
	recovery.sqlServerOps = &sqlServerStorageRestoreTargetOps{Inspect: f.inspect, Execute: f.execute}
	return f
}

func (f *sqlServerRepairFixture) inspect(_ context.Context, _ *sql.DB, _ *dbtarget.SQLServerRecordingCapability, table string) (*dbtarget.SQLServerManagedInspection, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	proof := dbtarget.CloneSQLServerManagedRowProof(f.current)
	physical := proof.Group
	missing := f.missing
	if table == dbtarget.EffectReceiptTable {
		physical, missing = proof.Receipt, f.receiptMissing
	}
	if f.inspectErr != nil {
		return physical, f.inspectErr
	}
	if missing {
		physical.Inspection.Status = dbtarget.TableInspectionMissing
		physical.Inspection.Columns, physical.Indexes, physical.Inspection.Indexes = nil, []dbtarget.IndexInfo{}, []dbtarget.IndexInfo{}
		physical.Evidence.ObjectID, physical.Evidence.TableType, physical.Evidence.Plain = 0, "", false
		physical.Evidence.PrimaryKey, physical.Evidence.RepairOperationID = nil, ""
		physical.Evidence.ColumnDefaults, physical.Evidence.ColumnCollations = map[string]*string{}, map[string]string{}
	}
	return physical, nil
}

func (f *sqlServerRepairFixture) execute(_ context.Context, _ *sql.DB, _ *dbtarget.SQLServerRecordingCapability, req dbtarget.SQLServerManagedRepairRequest) (*dbtarget.SchemaStatementExecution, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.executions++
	f.requests = append(f.requests, req)
	require.Equal(f.t, f.snapshot.Group.Destination.TableName, req.TableName)
	require.NotEmpty(f.t, req.OperationID)
	require.Len(f.t, req.Statements, 2, "repair creates only the group and its source index")
	for _, statement := range req.Statements {
		require.False(f.t, strings.Contains(statement, dbtarget.EffectReceiptTable), "the receipt table is never recreated")
	}
	if f.rolledBack {
		return &dbtarget.SchemaStatementExecution{RolledBack: true}, f.executeErr
	}
	f.missing = false
	f.current.Group.Evidence.ObjectID = int64(1000 + f.executions)
	f.current.Group.Evidence.RepairOperationID = req.OperationID
	if f.executeErr != nil {
		return &dbtarget.SchemaStatementExecution{}, f.executeErr
	}
	return &dbtarget.SchemaStatementExecution{Committed: len(req.Statements)}, nil
}

func (f *sqlServerRepairFixture) lineage() *recordingplan.StorageLineage {
	f.t.Helper()
	lineage, err := recordingplan.ReadStorageLineage(f.t.Context(), f.store.DB(), workspace.StorageScopeOf(f.snapshot.Group))
	require.NoError(f.t, err)
	return lineage
}

var errSQLServerFixtureAckLost = errors.New("scripted SQL DDL acknowledgement lost")
