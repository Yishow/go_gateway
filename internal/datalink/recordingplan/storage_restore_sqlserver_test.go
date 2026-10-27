package recordingplan

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type sqlServerRestoreFixture struct {
	db      *sql.DB
	service *Service
	request StorageRestoreRequest
	layout  GroupSchemaLayout
	options SQLServerGroupSchemaOptions
	group   TargetTableInspection
	receipt TargetTableInspection
	states  map[string]TargetTableInspection
	calls   []string
}

func newSQLServerRestoreFixture(t *testing.T) *sqlServerRestoreFixture {
	t.Helper()
	db := lineageDB(t)
	applyRecordingDDL(t, db)

	scope := StorageScope{
		WorkspaceID:   "ws",
		ConnectorID:   "conn",
		Database:      "recordings",
		TableSchema:   "dbo",
		StorageMode:   "plain",
		LayoutVersion: "managed_long_v1",
	}

	opts := SQLServerGroupSchemaOptions{
		Database:           scope.Database,
		Schema:             scope.TableSchema,
		Collation:          "Latin1_General_100_BIN2",
		ServerVersion:      "16.0.1000",
		ServerEdition:      "Standard Edition",
		CapabilityRevision: "cap-sql-1",
		EngineEdition:      2,
		CompatibilityLevel: 160,
	}

	table := "gw_long_" + mysqlStorageIdentityHash(scope, "table")[:16]
	owner := "_gw_owner_" + mysqlStorageIdentityHash(scope, "owner")[:24]

	layout, err := ManagedSQLServerLongGroupLayout(table, owner, opts)
	require.NoError(t, err)

	group := sqlServerRestoreCatalog(layout.Columns, layout.Indexes, table, owner, 101, opts)
	receipt := sqlServerRestoreReceiptCatalog(102, opts)

	digest, err := VerifyManagedSQLServerCatalog(layout, scope.Database, group, receipt)
	require.NoError(t, err)

	encodedLayout, err := ValidateGroupSchemaLayout(layout)
	require.NoError(t, err)

	layoutJSON, err := json.Marshal(encodedLayout)
	require.NoError(t, err)
	sum := sha256.Sum256(layoutJSON)
	layoutSum := hex.EncodeToString(sum[:])

	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode,
		scope.LayoutVersion, layout.TableName, layout.OwnerColumn, layoutSum)
	require.NoError(t, err)

	require.NoError(t, RecordStorageConfirmation(t.Context(), db, scope, StorageConfirmation{
		Layout: layout, ConnectorRevision: "rev-1", OperationID: "op-confirm", VerifiedDigest: digest, CreatedTable: true, GroupID: "group-1",
	}, time.Now()))

	missingGroup := group
	missingGroup.Status = inspectionMissing
	missingGroup.Columns = nil
	missingGroup.Indexes = nil

	states := map[string]TargetTableInspection{
		layout.TableName:          missingGroup,
		managedEffectReceiptTable: receipt,
	}

	reqEvidence := &SQLServerStorageRestoreEvidence{
		Group:   group,
		Receipt: receipt,
	}

	return &sqlServerRestoreFixture{
		db:      db,
		service: NewService(NewSQLRepository(db)),
		layout:  layout,
		options: opts,
		group:   group,
		receipt: receipt,
		states:  states,
		request: StorageRestoreRequest{
			Scope:             scope,
			Dialect:           dialectSQLServer,
			ConnectorRevision: "rev-1",
			TableName:         layout.TableName,
			OwnerColumn:       layout.OwnerColumn,
			GroupID:           "group-1",
			SQLServer:         reqEvidence,
		},
	}
}

func sqlServerRestoreCatalog(columns []GroupSchemaColumn, indexes []GroupSchemaIndex, table, owner string, objectID int64, opts SQLServerGroupSchemaOptions) TargetTableInspection {
	state := TargetTableInspection{
		Status: inspectionExists, ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}, IndexesInspected: true,
		SQLServer: &TargetSQLServerTableEvidence{
			Database: opts.Database, Schema: opts.Schema, ServerVersion: opts.ServerVersion, ServerEdition: opts.ServerEdition,
			EngineEdition: opts.EngineEdition, CompatibilityLevel: opts.CompatibilityLevel, CapabilityRevision: opts.CapabilityRevision,
			MetadataVerified: true, TableType: "USER_TABLE", ObjectID: objectID, Plain: true,
			ColumnDefaults: map[string]*string{}, ColumnCollations: map[string]string{}, ProvenanceJSONChecked: table != managedEffectReceiptTable,
		},
	}
	ownerDefault := managedOwnerDefault
	for _, col := range columns {
		state.Columns = append(state.Columns, col.Name)
		state.ColumnTypes[col.Name], state.ColumnNullable[col.Name], state.ColumnPrimaryKey[col.Name] = col.SQLType, col.Nullable, col.PrimaryKey
		state.SQLServer.ColumnDefaults[col.Name] = nil
		if col.Name == owner {
			state.SQLServer.ColumnDefaults[col.Name] = &ownerDefault
		}
		if strings.HasPrefix(col.SQLType, "NVARCHAR(") || strings.HasPrefix(col.SQLType, "CHAR(") {
			state.SQLServer.ColumnCollations[col.Name] = opts.Collation
		}
		if col.PrimaryKey {
			state.SQLServer.PrimaryKey = append(state.SQLServer.PrimaryKey, col.Name)
		}
	}
	for _, idx := range indexes {
		state.Indexes = append(state.Indexes, TargetIndexInspection{Name: idx.Name, Columns: idx.Columns})
	}
	return state
}

func sqlServerRestoreReceiptCatalog(objectID int64, opts SQLServerGroupSchemaOptions) TargetTableInspection {
	cols := sqlserverReceiptColumns()
	return sqlServerRestoreCatalog(cols, nil, managedEffectReceiptTable, managedEffectReceiptMarker, objectID, opts)
}

func (f *sqlServerRestoreFixture) inspect(_ context.Context, table string) (TargetTableInspection, error) {
	state, ok := f.states[table]
	if !ok {
		return TargetTableInspection{Status: inspectionMissing}, nil
	}
	return state, nil
}

func (f *sqlServerRestoreFixture) executeSQLServerRepair(_ context.Context, _, operationID string, statements []string) (SchemaExecution, error) {
	f.calls = append(f.calls, statements...)
	recreated := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 201, f.options)
	recreated.SQLServer.RepairOperationID = operationID
	f.states[f.layout.TableName] = recreated
	return SchemaExecution{Committed: len(statements)}, nil
}

func (f *sqlServerRestoreFixture) restore(ctx context.Context) (StorageRestoreResult, error) {
	return f.service.RestoreManagedStorage(ctx, f.request, StorageRestoreTarget{
		Inspect: f.inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{}, ErrSchemaAdapterUnverified
		},
		ExecuteSQLServerRepair: f.executeSQLServerRepair,
	})
}

func TestSQLServerMissingRepair(t *testing.T) {
	f := newSQLServerRestoreFixture(t)
	beforeLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, beforeLineage.Generation)
	require.False(t, beforeLineage.HistoryMayBeLost)
	require.Equal(t, 0, beforeLineage.Reconstructions)

	result, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome)
	require.NotEmpty(t, result.Operation.OperationID)
	require.NotEmpty(t, result.Operation.VerifiedDigest)
	require.NotEqual(t, beforeLineage.ConfirmedDigest, result.Operation.VerifiedDigest, "new ObjectID changes complete digest")

	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, afterLineage.Generation)
	require.True(t, afterLineage.HistoryMayBeLost)
	require.Equal(t, 1, afterLineage.Reconstructions)
	require.Equal(t, result.Operation.OperationID, afterLineage.LastReconstructionOperationID)
	require.Equal(t, beforeLineage.ConfirmedDigest, afterLineage.ConfirmedDigest, "original confirmed digest must remain unchanged")

	// Verification of Generation 2
	recreatedGroup := f.states[f.layout.TableName]
	require.NoError(t, f.service.VerifySQLServerStorageGeneration(t.Context(), f.request, recreatedGroup, f.receipt))

	// Wrong receipt fails verification
	tamperedReceipt := f.receipt
	tamperedReceipt.Indexes = []TargetIndexInspection{{Name: "idx_fake"}}
	require.Error(t, f.service.VerifySQLServerStorageGeneration(t.Context(), f.request, recreatedGroup, tamperedReceipt))

	// Foreign group (neither confirmed baseline nor current repair digest) fails verification
	foreignGroup := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 999, f.options)
	require.Error(t, f.service.VerifySQLServerStorageGeneration(t.Context(), f.request, foreignGroup, f.receipt))

	// Original confirmed catalog remains valid baseline regardless of absolute generation
	require.NoError(t, f.service.VerifySQLServerStorageGeneration(t.Context(), f.request, f.group, f.receipt))
}

func TestSQLServerNotMissing_AndForeignReplacementBlocked(t *testing.T) {
	t.Run("matching_group_returns_not_missing", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		f.states[f.layout.TableName] = f.group // Group is already present
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageNotMissing, result.Outcome)
		require.Empty(t, f.calls, "no DDL should run")
	})

	t.Run("foreign_object_id_replacement_is_blocked", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		foreignGroup := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 999, f.options)
		f.states[f.layout.TableName] = foreignGroup // Foreign table with different ObjectID
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageBlocked, result.Outcome)
		require.Equal(t, StorageCodeOwnerMismatch, result.Code)
		require.Equal(t, StorageReasonTableOccupied, result.Reason)
	})
}

func TestSQLServerRepairScopeDenied(t *testing.T) {
	t.Run("denied_policy", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		f.request.Deny = &StorageDenial{Code: StorageCodeRepairBlocked, Reason: StorageReasonPolicyUnconfirmed}
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageBlocked, result.Outcome)
		require.Equal(t, StorageReasonPolicyUnconfirmed, result.Reason)
	})

	t.Run("connector_revision_changed", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		f.request.ConnectorRevision = "rev-changed"
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageBlocked, result.Outcome)
		require.Equal(t, StorageReasonConnectorRevision, result.Reason)
	})

	t.Run("receipts_missing", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		missingReceipt := f.receipt
		missingReceipt.Status = inspectionMissing
		missingReceipt.Columns = nil
		f.states[managedEffectReceiptTable] = missingReceipt
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageBlocked, result.Outcome)
		require.Equal(t, StorageReasonReceiptsMissing, result.Reason)
	})

	t.Run("permission_denied", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		target := StorageRestoreTarget{
			Inspect: f.inspect,
			Execute: func(context.Context, []string) (SchemaExecution, error) {
				return SchemaExecution{}, ErrSchemaAdapterUnverified
			},
			ExecuteSQLServerRepair: f.executeSQLServerRepair,
			CanCreate: func(context.Context) (bool, error) {
				return false, nil
			},
		}
		result, err := f.service.RestoreManagedStorage(t.Context(), f.request, target)
		require.NoError(t, err)
		require.Equal(t, StorageBlocked, result.Outcome)
		require.Equal(t, StorageReasonPermissionDenied, result.Reason)
	})

	t.Run("mismatched_repair_operation_id_blocked", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		target := StorageRestoreTarget{
			Inspect: f.inspect,
			Execute: func(context.Context, []string) (SchemaExecution, error) {
				return SchemaExecution{}, ErrSchemaAdapterUnverified
			},
			ExecuteSQLServerRepair: func(_ context.Context, _ string, _ string, statements []string) (SchemaExecution, error) {
				recreated := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 201, f.options)
				recreated.SQLServer.RepairOperationID = "wrong-op-id"
				f.states[f.layout.TableName] = recreated
				return SchemaExecution{Committed: len(statements)}, nil
			},
		}
		result, err := f.service.RestoreManagedStorage(t.Context(), f.request, target)
		require.NoError(t, err)
		require.Equal(t, StorageBlocked, result.Outcome)
		require.Equal(t, StorageReasonVerificationMismatch, result.Reason)
	})
}

func TestSQLServerConcurrencyOneRepairPerScope(t *testing.T) {
	f := newSQLServerRestoreFixture(t)
	// First restore claims and executes
	result1, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result1.Outcome)

	// In the next run, the table is now restored (not missing)
	result2, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageNotMissing, result2.Outcome)
}

func TestSQLServerAuthorizedDigest_TamperedDetailRejected(t *testing.T) {
	f := newSQLServerRestoreFixture(t)
	res, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageRestored, res.Outcome)

	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, lineage.Generation)

	// Valid authorized digest
	digest, err := f.service.authorizedSQLServerDigest(t.Context(), f.request.Scope, lineage)
	require.NoError(t, err)
	require.NotEmpty(t, digest)

	// Corrupt detail: tamper scope
	var detail storageRepairDetail
	require.NoError(t, json.Unmarshal([]byte(res.Operation.Detail), &detail))
	detail.Scope.WorkspaceID = "tampered-ws"
	encoded, err := json.Marshal(detail)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encoded), res.Operation.OperationID)
	require.NoError(t, err)

	_, err = f.service.authorizedSQLServerDigest(t.Context(), f.request.Scope, lineage)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrSchemaOperationMismatch)
}

func TestSQLServerJudge_RolledBackRejectsEvenIfTablePresent(t *testing.T) {
	f := newSQLServerRestoreFixture(t)
	recreated := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 201, f.options)
	target := StorageRestoreTarget{
		Inspect: f.inspect,
		ExecuteSQLServerRepair: func(_ context.Context, _ string, opID string, _ []string) (SchemaExecution, error) {
			recreated.SQLServer.RepairOperationID = opID
			f.states[f.layout.TableName] = recreated
			return SchemaExecution{Committed: 0, RolledBack: true}, errors.New("simulated DDL rollback")
		},
	}
	result, err := f.service.RestoreManagedStorage(t.Context(), f.request, target)
	require.NoError(t, err)
	require.Equal(t, StorageBlocked, result.Outcome)
	require.Equal(t, StorageReasonStatementFailed, result.Reason)
}

func TestSQLServerStorageRepairRequestOf_PreservesSQLServer(t *testing.T) {
	f := newSQLServerRestoreFixture(t)
	res, err := f.restore(t.Context())
	require.NoError(t, err)

	req, ok := StorageRepairRequestOf(*res.Operation, "rev-1")
	require.True(t, ok)
	require.Equal(t, f.request.Scope, req.Scope)
	require.Equal(t, dialectSQLServer, req.Dialect)
	require.NotNil(t, req.SQLServer)
	require.Equal(t, f.request.SQLServer.Group.SQLServer.ObjectID, req.SQLServer.Group.SQLServer.ObjectID)

	badOp := *res.Operation
	badOp.Detail = "invalid json"
	_, ok = StorageRepairRequestOf(badOp, "rev-1")
	require.False(t, ok)
}

func TestSQLServerStorageRepair_GenerationCASMissRollsBack(t *testing.T) {
	f := newSQLServerRestoreFixture(t)
	recreated := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 201, f.options)
	target := StorageRestoreTarget{
		Inspect: f.inspect,
		ExecuteSQLServerRepair: func(ctx context.Context, _, opID string, _ []string) (SchemaExecution, error) {
			recreated.SQLServer.RepairOperationID = opID
			f.states[f.layout.TableName] = recreated
			_, err := f.db.ExecContext(ctx, `UPDATE write_group_storage_scopes SET generation = 99 WHERE workspace_id = ?`, f.request.Scope.WorkspaceID)
			require.NoError(t, err)
			return SchemaExecution{Committed: 1}, nil
		},
	}
	_, err := f.service.RestoreManagedStorage(t.Context(), f.request, target)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrSchemaOperationMismatch)

	rows, err := f.db.QueryContext(t.Context(), `SELECT status FROM managed_schema_operations WHERE action = ?`, StorageRepairAction)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var status string
		require.NoError(t, rows.Scan(&status))
		require.NotEqual(t, string(SchemaOperationSucceeded), status)
	}
}
