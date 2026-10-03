package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

const managedSQLiteSchemaName = "main"

// ManagedSchemaScope resolves a preview entirely from the canonical saved
// group and live source revisions. It performs no destination inspection.
func (s *WriteGroupService) ManagedSchemaScope(ctx context.Context, id string, mutation WriteGroupMutation) (recordingplan.SchemaPreviewScope, error) {
	if err := s.validate(); err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	snapshot, err := s.readinessSnapshot(ctx, id)
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	if snapshot.validateErr != nil {
		return recordingplan.SchemaPreviewScope{}, snapshot.validateErr
	}
	group := snapshot.validated
	if snapshot.record.ID != mutation.WorkspaceID || group.Status == WriteGroupStatusDeleted {
		return recordingplan.SchemaPreviewScope{}, writeGroupNotFound("prepare group schema")
	}
	if snapshot.record.DatabaseSetupRevision != mutation.ExpectedWorkspaceRevision || group.Revision != mutation.ExpectedGroupRevision || group.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
		return recordingplan.SchemaPreviewScope{}, writeGroupRevisionConflict("managed schema scope is stale")
	}
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged || group.RowPolicy.RecordKeyColumn == "" {
		return recordingplan.SchemaPreviewScope{}, fmt.Errorf("managed schema requires a saved managed layout: %w", ErrWriteGroupValidation)
	}
	if snapshot.connectorKind == writeGroupConnectorKindSQLite && group.Destination.TableSchema != managedSQLiteSchemaName {
		return recordingplan.SchemaPreviewScope{}, fmt.Errorf("managed SQLite schema must be main: %w", ErrWriteGroupValidation)
	}
	layout, err := managedSchemaLayout(group, snapshot.tagTypes, snapshot.connectorKind)
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	encoded, err := json.Marshal([]any{group.Members, snapshot.tagTypes})
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, fmt.Errorf("encode managed source scope: %w", err)
	}
	digest := sha256.Sum256(encoded)
	return recordingplan.SchemaPreviewScope{
		WorkspaceID: snapshot.record.ID, WorkspaceRevision: snapshot.record.DatabaseSetupRevision,
		PlanID: group.ID, PlanRevision: group.Revision, ConnectorID: group.Destination.ConnectorID,
		ConnectorRevision: group.Destination.ConnectorRevision, Dialect: snapshot.connectorKind,
		Database: group.Destination.Database, Schema: group.Destination.TableSchema, TablePrefix: group.Destination.TableName,
		GroupLayout: layout, SourceDigest: hex.EncodeToString(digest[:]),
		SchemaRevision: group.Destination.SchemaRevision, SchemaDigest: group.Destination.SchemaDigest,
	}, nil
}

func managedSchemaLayout(group *WriteGroup, types map[string]schema.DataType, dialect string) (*recordingplan.GroupSchemaLayout, error) {
	owner := sha256.Sum256([]byte(group.WorkspaceID + "\x00" + group.ID))
	layout := &recordingplan.GroupSchemaLayout{TableName: group.Destination.TableName, OwnerColumn: "_gw_owner_" + hex.EncodeToString(owner[:12])}
	add := func(name, sqlType string, nullable, primary bool) {
		if name != "" {
			layout.Columns = append(layout.Columns, recordingplan.GroupSchemaColumn{Name: name, SQLType: sqlType, Nullable: nullable, PrimaryKey: primary})
		}
	}
	add(group.RowPolicy.RecordKeyColumn, "TEXT", false, true)
	add(group.RowPolicy.GroupIDColumn, "TEXT", false, false)
	add(group.RowPolicy.DeviceIDColumn, "TEXT", false, false)
	add(group.RowPolicy.BucketStartColumn, "TEXT", false, false)
	add(group.RowPolicy.ProvenanceColumn, "TEXT", false, false)
	add(group.RowPolicy.EntityKeyColumn, "TEXT", false, false)
	add(layout.OwnerColumn, "TEXT", false, false)
	for _, member := range group.Members {
		sqlType, ok := managedSQLType(types[member.TagID], dialect)
		if !ok {
			return nil, fmt.Errorf("managed exact type is unsupported: %w", ErrWriteGroupValidation)
		}
		add(member.TargetColumn, sqlType, true, false)
	}
	return layout, nil
}

func managedSQLType(typ schema.DataType, dialect string) (string, bool) {
	exact, ok := measurement.ExactTypeForTag(typ)
	if !ok {
		return "", false
	}
	postgres := dialect == writeGroupConnectorKindPostgres
	switch exact {
	case measurement.ExactText:
		return "TEXT", true
	case measurement.ExactBool:
		if postgres {
			return "BOOLEAN", true
		}
		return "INTEGER", true
	case measurement.ExactInt64:
		if postgres {
			return "BIGINT", true
		}
		return "INTEGER", true
	case measurement.ExactUint64:
		if postgres {
			return "NUMERIC(20,0)", true
		}
		return "TEXT", true
	case measurement.ExactFloat64:
		if postgres {
			return "DOUBLE PRECISION", true
		}
		return "REAL", true
	}
	return "", false
}
