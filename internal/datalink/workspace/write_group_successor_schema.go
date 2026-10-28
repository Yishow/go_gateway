package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/dbtarget"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

type successorSchemaSnapshot struct {
	reserved     *WriteGroupSuccessorPreparation
	types        map[string]schema.DataType
	dialect      string
	sourceFamily string
}

// SuccessorSchemaScope resolves only the exact operation-bound reserved draft.
// It does not grant ordinary schema or intake authority to the successor.
func (s *WriteGroupService) SuccessorSchemaScope(ctx context.Context, id string, req WriteGroupSuccessorStageRequest) (recordingplan.SchemaPreviewScope, error) {
	snapshot, err := s.readSuccessorSchemaSnapshot(ctx, id, req)
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	group := snapshot.reserved.Group
	var layout *recordingplan.GroupSchemaLayout
	switch snapshot.dialect {
	case writeGroupConnectorKindMySQL:
		var capability *dbtarget.MySQLRecordingCapability
		layout, capability, err = s.managedMySQLSchemaLayoutWithCapability(ctx, group)
		if err == nil && (layout.MySQL.ServerFamily != snapshot.sourceFamily || !mysqlRecordingOperationGate(capability, group)) {
			err = transitionRejection("engine-evidence-unavailable")
		}
	case writeGroupConnectorKindSQLServer:
		layout, err = s.managedSQLServerSchemaLayout(ctx, group)
	default:
		layout, err = managedSchemaLayout(group, snapshot.types, snapshot.dialect)
	}
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	// Capability probes occur outside the local transaction. Their result grants
	// nothing when the reservation/source/current CAS changed while inspecting.
	current, err := s.readSuccessorSchemaSnapshot(ctx, id, req)
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	if !reflect.DeepEqual(snapshot, current) {
		return recordingplan.SchemaPreviewScope{}, recordingplan.ErrPreviewTokenStale
	}
	sourceDigest, err := successorSchemaSourceDigest(snapshot)
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	scope := recordingplan.SchemaPreviewScope{
		WorkspaceID: group.WorkspaceID, WorkspaceRevision: snapshot.reserved.WorkspaceRevision,
		PlanID: group.ID, PlanRevision: group.Revision, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
		Dialect: snapshot.dialect, Database: group.Destination.Database, Schema: group.Destination.TableSchema, TablePrefix: group.Destination.TableName,
		GroupLayout: layout, SourceDigest: sourceDigest, SchemaRevision: group.Destination.SchemaRevision, SchemaDigest: group.Destination.SchemaDigest,
	}
	if managedLongGroup(group) && group.WritePolicy.RecreateMissing != nil {
		v := *group.WritePolicy.RecreateMissing
		scope.RecreateMissing = &v
	}
	return scope, nil
}

func (s *WriteGroupService) readSuccessorSchemaSnapshot(ctx context.Context, id string, req WriteGroupSuccessorStageRequest) (*successorSchemaSnapshot, error) {
	return s.readSuccessorSchemaSnapshotChecked(ctx, id, req, nil)
}

func (s *WriteGroupService) readSuccessorSchemaSnapshotChecked(ctx context.Context, id string, req WriteGroupSuccessorStageRequest, check func(*sql.Tx, *successorSchemaSnapshot) error) (result *successorSchemaSnapshot, err error) {
	if err = s.validate(); err != nil {
		return nil, err
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setup, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := setup.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, e)
		}
	}()
	record, err := setup.Get(ctx)
	if err != nil {
		return nil, err
	}
	result, err = s.successorSchemaSnapshotInTx(ctx, setup.SQLTx(), record, id, req)
	if err == nil && check != nil {
		err = check(setup.SQLTx(), result)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WriteGroupService) successorSchemaSnapshotInTx(ctx context.Context, tx *sql.Tx, record *Record, id string, req WriteGroupSuccessorStageRequest) (*successorSchemaSnapshot, error) {
	reserved, err := s.successorReservationInTx(ctx, tx, record, id, req)
	if err != nil {
		return nil, err
	}
	group := reserved.Group
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged || group.RowPolicy.RecordKeyColumn == "" {
		return nil, ErrWriteGroupValidation
	}
	kind, err := s.successorDriverInTx(ctx, tx, group.Destination.ConnectorID)
	if err != nil {
		return nil, err
	}
	// Connector kind, unlike frozen source engine, uses mysql for MariaDB drivers.
	if kind == "mariadb" {
		kind = writeGroupConnectorKindMySQL
	}
	if kind == writeGroupConnectorKindSQLite && group.Destination.TableSchema != managedSQLiteSchemaName {
		return nil, ErrWriteGroupValidation
	}
	types, err := loadWriteGroupTagTypes(ctx, tx, s.repo.query, group)
	if err != nil {
		return nil, err
	}
	var family string
	if kind == writeGroupConnectorKindMySQL {
		pred, err := s.repo.GetInTx(ctx, tx, record.ID, id)
		if err != nil {
			return nil, err
		}
		family, err = s.successorSourceEngineInTx(ctx, tx, pred, reserved)
		if err != nil {
			return nil, err
		}
	}
	return &successorSchemaSnapshot{reserved: reserved, types: types, dialect: kind, sourceFamily: family}, nil
}

func successorSchemaSourceDigest(snapshot *successorSchemaSnapshot) (string, error) {
	source, err := managedSourceDigest(snapshot.reserved.Group, snapshot.types)
	if err != nil {
		return "", err
	}
	encoded, err := json.Marshal([]string{"write-group-successor-schema/v1", snapshot.reserved.LineageID, snapshot.reserved.PredecessorGroupID, snapshot.reserved.Group.ID, snapshot.reserved.OperationID, snapshot.reserved.IntentDigest, source})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

// successorDriverInTx selects an adapter only; mysql is not server-family proof.
// Its caller must independently validate immutable source and target evidence.
func (s *WriteGroupService) successorDriverInTx(ctx context.Context, tx *sql.Tx, connectorID string) (string, error) {
	var kind string
	if err := tx.QueryRowContext(ctx, s.repo.query(`SELECT kind FROM database_connectors WHERE id=$1`), connectorID).Scan(&kind); err != nil {
		return "", err
	}
	if strings.ToLower(strings.TrimSpace(kind)) == writeGroupConnectorKindMySQL {
		return writeGroupConnectorKindMySQL, nil
	}
	return s.repo.recreateConnectorEngineInTx(ctx, tx, connectorID)
}
