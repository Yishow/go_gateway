package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
)

const sqlServerManagedSchemaMissingReason = "schema_missing"

var (
	// ErrSQLServerManagedTableForbidden means catalog metadata access was denied.
	ErrSQLServerManagedTableForbidden = errors.New("sqlserver managed table inspection forbidden")
	// ErrSQLServerManagedTableUnconfirmed means catalog facts could not be verified.
	ErrSQLServerManagedTableUnconfirmed = errors.New("sqlserver managed table inspection unconfirmed")
	// ErrSQLServerManagedScopeChanged means the saved database, schema, or capability binding changed.
	ErrSQLServerManagedScopeChanged = errors.New("sqlserver managed inspection scope changed")
	// ErrSQLServerManagedIndexUnusable means an index cannot satisfy the managed layout contract.
	ErrSQLServerManagedIndexUnusable = errors.New("sqlserver managed index unusable")
)

// SQLServerManagedTableEvidence is the verified SQL Server server/scope and catalog projection.
type SQLServerManagedTableEvidence struct {
	ServerVersion         string
	ServerEdition         string
	CapabilityRevision    string
	Database              string
	Schema                string
	EngineEdition         int64
	CompatibilityLevel    int64
	MetadataVerified      bool
	TableType             string
	ObjectID              int64
	Plain                 bool
	PrimaryKey            []string
	ColumnCollations      map[string]string
	ColumnDefaults        map[string]*string
	ProvenanceJSONChecked bool
	RepairOperationID     string                             `json:"RepairOperationID,omitempty"`
	NativeMonthly         *SQLServerManagedPartitionEvidence `json:"NativeMonthly,omitempty"`
	PartitionInspected    bool                               `json:"PartitionInspected,omitzero"`
}

// SQLServerManagedInspection combines a table inspection with its verified SQL Server catalog facts.
type SQLServerManagedInspection struct {
	Inspection *TableInspection
	Evidence   *SQLServerManagedTableEvidence
	Indexes    []IndexInfo
}

type sqlServerManagedSession interface {
	DB() *sql.DB
	Close() error
}

type sqlServerManagedDBManagerSession struct {
	manager *datalinkbase.DBManager
}

func (s *sqlServerManagedDBManagerSession) DB() *sql.DB {
	return s.manager.DB()
}

func (s *sqlServerManagedDBManagerSession) Close() error {
	return s.manager.Close()
}

var openSQLServerManagedSessionFunc = func(
	kind schema.DatabaseConnectorKind,
	config ConnectionConfig,
) (sqlServerManagedSession, error) {
	manager, err := openExternalDBManagerFunc(kind, config)
	if err != nil {
		return nil, err
	}
	if manager == nil || manager.DB() == nil {
		if manager != nil {
			_ = manager.Close()
		}
		return nil, errors.New("empty SQL Server managed connection")
	}
	return &sqlServerManagedDBManagerSession{manager: manager}, nil
}

// InspectSQLServerManagedTableAtRevision reads only the catalog for the saved SQL Server scope and capability revision.
func (i *ManagedTableInspector) InspectSQLServerManagedTableAtRevision(
	ctx context.Context,
	connectorID, expectedRevision, database, schemaName, expectedCapabilityRevision, tableName string,
) (*SQLServerManagedInspection, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := i.service.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, err
	}
	if connector.Kind != schema.DatabaseConnectorKindSQLServer {
		return nil, managedDestinationInvalid("saved connector kind is not sqlserver")
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, managedDestinationInvalid("saved sqlserver connector configuration is invalid")
	}
	settings, err := sqlServerConnectionSettingsFor(config)
	if err != nil {
		return nil, managedDestinationInvalid("saved sqlserver connector configuration is invalid")
	}
	if database == "" || schemaName == "" || database != settings.Database || schemaName != settings.Schema ||
		strings.TrimSpace(expectedCapabilityRevision) == "" {
		return nil, ErrSQLServerManagedScopeChanged
	}
	if strings.TrimSpace(tableName) == "" {
		return nil, managedDestinationInvalid("managed table name is required")
	}

	session, err := openSQLServerManagedSessionFunc(connector.Kind, config)
	if err != nil || session == nil || session.DB() == nil {
		if session != nil {
			_ = session.Close()
		}
		return nil, ErrSQLServerManagedTableUnconfirmed
	}
	defer session.Close()

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	metadata := inspectSQLServerRecordingMetadata(queryCtx, session.DB(), settings.Database, settings.Schema)
	projected := projectSQLServerRecordingCapability(connector, metadata, true)
	return inspectSQLServerManagedTable(
		queryCtx, session.DB(), schemaName, tableName, expectedCapabilityRevision, metadata, projected, nil,
	)
}

func inspectSQLServerManagedTable(
	ctx context.Context,
	db *sql.DB,
	schemaName, tableName, expectedCapabilityRevision string,
	metadata sqlServerRecordingMetadata,
	projected, expectedCapability *SQLServerRecordingCapability,
) (*SQLServerManagedInspection, error) {
	result := newSQLServerManagedInspection(schemaName, tableName, metadata, projected)
	if slices.Contains(metadata.reasons, "database_scope_mismatch") || slices.Contains(metadata.reasons, "schema_scope_mismatch") {
		return result, ErrSQLServerManagedScopeChanged
	}

	switch metadata.metadataStatus {
	case sqlServerRecordingMetadataForbidden:
		result.Inspection.Status = TableInspectionForbidden
		result.Inspection.Reason = inspectionReasonPermissionDenied
		return result, ErrSQLServerManagedTableForbidden
	case sqlServerRecordingMetadataMissing:
		if !metadata.permissions.DatabaseViewDefinition {
			return sqlServerManagedCatalogFailure(result, errors.New("database metadata visibility is unconfirmed"))
		}
		result.Inspection.Status = TableInspectionFailed
		result.Inspection.Reason = sqlServerManagedSchemaMissingReason
		return result, ErrSQLServerManagedTableUnconfirmed
	case sqlServerRecordingMetadataVerified:
		if metadata.actualDatabase != projected.Database || metadata.actualSchema != schemaName {
			return result, ErrSQLServerManagedScopeChanged
		}
		result.Evidence.MetadataVerified = true
	default:
		return sqlServerManagedCatalogFailure(result, errors.New("server or schema metadata is unconfirmed"))
	}
	if projected.CapabilityRevision != expectedCapabilityRevision ||
		(expectedCapability != nil && !sqlServerRecordingCapabilityScopeMatches(expectedCapability, projected)) {
		return result, ErrSQLServerManagedScopeChanged
	}

	facts, err := inspectSQLServerManagedCatalog(ctx, db, schemaName, tableName)
	if errors.Is(err, sql.ErrNoRows) {
		nativeMonthly, partitionInspected, partitionErr := readSQLServerManagedPartitions(ctx, db, schemaName, tableName, 0)
		if partitionErr != nil {
			return sqlServerManagedCatalogFailure(result, partitionErr)
		}
		result.Evidence.NativeMonthly = cloneSQLServerManagedPartitionEvidence(nativeMonthly)
		result.Evidence.PartitionInspected = partitionInspected
		result.Inspection.Status = TableInspectionMissing
		result.Inspection.Reason = ""
		result.Inspection.IndexesInspected = true
		result.Indexes = []IndexInfo{}
		result.Inspection.Indexes = []IndexInfo{}
		return result, nil
	}
	if errors.Is(err, ErrSQLServerManagedScopeChanged) {
		return result, err
	}
	if err != nil {
		return sqlServerManagedCatalogFailure(result, err)
	}

	result.Evidence.TableType = facts.object.tableType
	result.Evidence.ObjectID = facts.object.objectID
	result.Evidence.Plain = facts.plain
	result.Evidence.PrimaryKey = slices.Clone(facts.primaryKey)
	result.Evidence.ColumnCollations = facts.columnCollations
	result.Evidence.ColumnDefaults = facts.columnDefaults
	result.Evidence.ProvenanceJSONChecked = facts.provenanceJSONChecked
	result.Evidence.RepairOperationID = facts.repairOperationID
	result.Evidence.NativeMonthly = cloneSQLServerManagedPartitionEvidence(facts.nativeMonthly)
	result.Evidence.PartitionInspected = facts.partitionInspected
	result.Inspection.Status = TableInspectionExists
	result.Inspection.Reason = ""
	result.Inspection.Columns = facts.columns
	result.Inspection.Indexes = facts.indexes
	result.Inspection.IndexesInspected = true
	result.Indexes = facts.indexes
	if facts.indexUnusable {
		return result, ErrSQLServerManagedIndexUnusable
	}
	return result, nil
}

func newSQLServerManagedInspection(
	schemaName, tableName string,
	metadata sqlServerRecordingMetadata,
	capability *SQLServerRecordingCapability,
) *SQLServerManagedInspection {
	return &SQLServerManagedInspection{
		Inspection: &TableInspection{
			Status: TableInspectionFailed, Schema: schemaName, Table: tableName,
			Columns: []ColumnInfo{}, Reason: inspectionReasonExistenceUnconfirmed,
		},
		Evidence: &SQLServerManagedTableEvidence{
			ServerVersion: metadata.serverVersion, ServerEdition: metadata.serverEdition,
			CapabilityRevision: capability.CapabilityRevision, Database: metadata.actualDatabase, Schema: metadata.actualSchema,
			EngineEdition: metadata.engineEdition, CompatibilityLevel: metadata.compatibilityLevel,
			ColumnCollations: map[string]string{}, ColumnDefaults: map[string]*string{}, PrimaryKey: []string{},
		},
		Indexes: []IndexInfo{},
	}
}

func sqlServerManagedCatalogFailure(result *SQLServerManagedInspection, cause error) (*SQLServerManagedInspection, error) {
	if isSQLServerPermissionDeniedError(cause) || isSQLServerPermissionDeniedNumber(cause) {
		result.Inspection.Status = TableInspectionForbidden
		result.Inspection.Reason = inspectionReasonPermissionDenied
		return result, ErrSQLServerManagedTableForbidden
	}
	result.Inspection.Status = TableInspectionFailed
	result.Inspection.Reason = inspectionReasonExistenceUnconfirmed
	return result, ErrSQLServerManagedTableUnconfirmed
}
