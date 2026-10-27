package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
)

// MySQLManagedTableEvidence is the backend-derived catalog binding of one
// MySQL/MariaDB managed table. recordingplan mirrors these fields; the packages
// deliberately do not import each other.
type MySQLManagedTableEvidence struct {
	ServerFamily       string
	ServerVersion      string
	CapabilityRevision string
	Database           string

	Engine    string
	Charset   string
	Collation string
	TableType string

	PrimaryKey       []string
	ColumnCollations map[string]string
	ColumnDefaults   map[string]*string
	JSONColumns      map[string]bool
	Partitions       *MySQLPartitionEvidence `json:"partitions,omitempty"`
}

// MySQLManagedInspection is the exact catalog state plus the server binding it
// was read under. Evidence is always present once server facts were verified,
// including for a missing table.
type MySQLManagedInspection struct {
	Inspection *TableInspection
	Evidence   *MySQLManagedTableEvidence
	Indexes    []IndexInfo
}

var (
	// ErrMySQLManagedTableForbidden means metadata was refused, not that the
	// table is missing.
	ErrMySQLManagedTableForbidden = errors.New("mysql managed table inspection forbidden")
	// ErrMySQLManagedTableUnconfirmed means the catalog answer is not proof of
	// absence or presence.
	ErrMySQLManagedTableUnconfirmed = errors.New("mysql managed table inspection unconfirmed")
	// ErrMySQLManagedIndexUnusable means an index exists but cannot satisfy a
	// declared managed index (prefix, functional, invisible/ignored, non-BTREE).
	ErrMySQLManagedIndexUnusable = errors.New("mysql managed index unusable")
	// ErrMySQLManagedCapabilityChanged means the saved capability revision no
	// longer matches the server facts read for this access.
	ErrMySQLManagedCapabilityChanged = errors.New("mysql managed capability revision changed")
	// ErrMySQLManagedServerMismatch means the actual server, database, engine or
	// transport does not match the saved managed scope.
	ErrMySQLManagedServerMismatch = errors.New("mysql managed server facts do not match the saved scope")
	// ErrMySQLManagedStatementRefused means the statement is outside the
	// server-generated CREATE TABLE / CREATE INDEX contract.
	ErrMySQLManagedStatementRefused = errors.New("mysql managed statement refused")
)

type mysqlManagedServerFacts struct {
	family              string
	version             string
	capabilityRevision  string
	database            string
	metadataStatus      string
	transactionalEngine bool
}

type mysqlManagedColumnFacts struct {
	name         string
	dataType     string
	columnType   string
	collation    string
	nullable     bool
	primaryKey   bool
	hasDefault   bool
	defaultValue *string
	jsonValid    bool
}

type mysqlManagedIndexFacts struct {
	name       string
	columns    []string
	unique     bool
	indexType  string
	prefix     bool
	functional bool
	visible    bool
	ignored    bool
}

type mysqlManagedTableFacts struct {
	status      TableInspectionStatus
	reason      string
	tableType   string
	engine      string
	charset     string
	collation   string
	columns     []mysqlManagedColumnFacts
	primaryKey  []string
	indexes     []mysqlManagedIndexFacts
	jsonColumns map[string]bool
	partitions  *MySQLPartitionEvidence
}

type mysqlManagedSession interface {
	ServerFacts(ctx context.Context, database string) (mysqlManagedServerFacts, error)
	TableFacts(ctx context.Context, database, table string) (mysqlManagedTableFacts, error)
	Exec(ctx context.Context, statement string) error
	Close() error
}

var openMySQLManagedSessionFunc = func(connector *schema.DatabaseConnector, config ConnectionConfig, database string) (mysqlManagedSession, error) {
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		return nil, err
	}
	if manager == nil || manager.DB() == nil {
		if manager != nil {
			_ = manager.Close()
		}
		return nil, fmt.Errorf("open mysql destination: empty connection")
	}
	return &sqlMySQLManagedSession{manager: manager, connector: connector}, nil
}

// ProbeMySQLRecordingAtRevision exposes the existing readonly 1.1 probe
// through the managed inspector, so workspace/API code can use the current
// optional inspector interface without opening its own connection.
func (i *ManagedTableInspector) ProbeMySQLRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*MySQLRecordingCapability, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	return i.service.ProbeMySQLRecordingAtRevision(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
}

// MySQLRecordingDatabaseScope resolves only the MySQL database/schema from a
// saved connector configuration. It never opens a connection and never returns
// a DSN, user or password.
func MySQLRecordingDatabaseScope(rawConfig string) (string, error) {
	config, err := parseConnectionConfig(rawConfig)
	if err != nil {
		return "", fmt.Errorf("%w: mysql connector configuration is invalid", ErrManagedDestinationInvalid)
	}
	dbConfig, err := buildExternalDBConfig(schema.DatabaseConnectorKindMySQL, config)
	if err != nil {
		return "", fmt.Errorf("%w: mysql connector configuration is incomplete", ErrManagedDestinationInvalid)
	}
	driverConfig, err := mysqldriver.ParseDSN(dbConfig.DSN)
	if err != nil {
		return "", fmt.Errorf("%w: mysql connector configuration is invalid", ErrManagedDestinationInvalid)
	}
	database := strings.TrimSpace(driverConfig.DBName)
	if database == "" {
		return "", fmt.Errorf("%w: mysql database is required", ErrMySQLManagedServerMismatch)
	}
	return database, nil
}

// InspectMySQLManagedTableAtRevision reads the exact managed-table catalog at
// the saved connector revision. It never creates a database or table.
func (i *ManagedTableInspector) InspectMySQLManagedTableAtRevision(
	ctx context.Context,
	connectorID, expectedRevision, database, expectedCapabilityRevision, tableName string,
) (*MySQLManagedInspection, error) {
	session, server, err := i.openMySQLManagedForRevision(ctx, connectorID, expectedRevision, database, expectedCapabilityRevision)
	evidence := mysqlManagedBindingEvidence(server, database)
	if err != nil {
		if evidence.ServerFamily != "" && evidence.ServerVersion != "" && evidence.CapabilityRevision != "" {
			return &MySQLManagedInspection{
				Inspection: mysqlManagedFailedInspection(database, tableName, err),
				Evidence:   evidence,
			}, err
		}
		return nil, err
	}
	defer session.Close()

	facts, tableErr := session.TableFacts(ctx, database, tableName)
	if tableErr != nil {
		return &MySQLManagedInspection{
			Inspection: mysqlManagedFailedInspection(database, tableName, tableErr),
			Evidence:   evidence,
		}, fmt.Errorf("%w: %w", ErrMySQLManagedTableUnconfirmed, tableErr)
	}
	inspection, buildErr := buildMySQLManagedInspection(database, tableName, server, facts)
	return inspection, buildErr
}

// ExecuteMySQLManagedStatementAtRevision executes one server-generated MySQL
// DDL statement. It never opens a transaction and never claims a rollback.
func (i *ManagedTableInspector) ExecuteMySQLManagedStatementAtRevision(
	ctx context.Context,
	connectorID, expectedRevision, database, expectedCapabilityRevision, statement string,
) error {
	if err := validateMySQLManagedStatement(statement); err != nil {
		return err
	}
	session, _, err := i.openMySQLManagedForRevision(ctx, connectorID, expectedRevision, database, expectedCapabilityRevision)
	if err != nil {
		return err
	}
	defer session.Close()
	if err := session.Exec(ctx, statement); err != nil {
		// The statement may already be applied: never report a rollback here.
		return fmt.Errorf("execute mysql managed statement: %w", err)
	}
	return nil
}

func (i *ManagedTableInspector) openMySQLManagedForRevision(
	ctx context.Context,
	connectorID, expectedRevision, database, expectedCapabilityRevision string,
) (mysqlManagedSession, mysqlManagedServerFacts, error) {
	var server mysqlManagedServerFacts
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, server, managedDestinationInvalid("saved connector service is unavailable")
	}
	if strings.TrimSpace(database) == "" {
		return nil, server, fmt.Errorf("%w: mysql database is required", ErrMySQLManagedServerMismatch)
	}
	connector, err := i.service.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, server, err
	}
	if connector.Kind != schema.DatabaseConnectorKindMySQL {
		return nil, server, managedDestinationInvalid("saved connector kind is not mysql")
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, server, managedDestinationInvalid("saved mysql connector configuration is invalid")
	}
	dbConfig, err := buildExternalDBConfig(connector.Kind, config)
	if err != nil {
		return nil, server, managedDestinationInvalid("saved mysql connector configuration is incomplete")
	}
	driverConfig, err := mysqldriver.ParseDSN(dbConfig.DSN)
	if err != nil {
		return nil, server, managedDestinationInvalid("saved mysql connector configuration is invalid")
	}
	if strings.TrimSpace(driverConfig.DBName) == "" || driverConfig.DBName != database {
		return nil, server, fmt.Errorf("%w: mysql database scope does not match the saved connection", ErrMySQLManagedServerMismatch)
	}
	if !mysqlRecordingTransportVerified(driverConfig) {
		return nil, server, fmt.Errorf("%w: mysql managed access requires verified TLS", ErrMySQLManagedServerMismatch)
	}

	session, err := openMySQLManagedSessionFunc(connector, config, database)
	if err != nil {
		return nil, server, err
	}
	server, err = session.ServerFacts(ctx, database)
	if err != nil {
		_ = session.Close()
		return nil, server, err
	}
	if err := validateMySQLManagedServer(server, database, expectedCapabilityRevision); err != nil {
		_ = session.Close()
		return nil, server, err
	}
	return session, server, nil
}

func validateMySQLManagedServer(server mysqlManagedServerFacts, database, expectedCapabilityRevision string) error {
	switch server.metadataStatus {
	case mysqlRecordingMetadataVerified:
	case mysqlRecordingMetadataForbidden:
		return ErrMySQLManagedTableForbidden
	default:
		return ErrMySQLManagedTableUnconfirmed
	}
	if server.family != mysqlManagedFamilyMySQL && server.family != mysqlManagedFamilyMariaDB {
		return fmt.Errorf("%w: server family", ErrMySQLManagedServerMismatch)
	}
	if !MySQLRecordingVersionVerified(server.family, server.version) {
		return fmt.Errorf("%w: server version %s is not a verified recording version", ErrMySQLManagedServerMismatch, server.version)
	}
	if server.database != database {
		return fmt.Errorf("%w: actual database does not match the saved scope", ErrMySQLManagedServerMismatch)
	}
	if !server.transactionalEngine {
		return fmt.Errorf("%w: InnoDB transactional support is unverified", ErrMySQLManagedServerMismatch)
	}
	if strings.TrimSpace(server.capabilityRevision) == "" {
		return ErrMySQLManagedCapabilityChanged
	}
	if server.capabilityRevision != expectedCapabilityRevision {
		return ErrMySQLManagedCapabilityChanged
	}
	return nil
}

func mysqlManagedBindingEvidence(server mysqlManagedServerFacts, database string) *MySQLManagedTableEvidence {
	return &MySQLManagedTableEvidence{
		ServerFamily:       server.family,
		ServerVersion:      server.version,
		CapabilityRevision: server.capabilityRevision,
		Database:           defaultString(server.database, database),
	}
}

func mysqlManagedFailedInspection(database, table string, err error) *TableInspection {
	status, reason := TableInspectionFailed, "inspection_failed"
	switch {
	case errors.Is(err, ErrMySQLManagedTableForbidden):
		status, reason = TableInspectionForbidden, "permission_denied"
	case errors.Is(err, ErrMySQLManagedTableUnconfirmed):
		reason = "existence_unconfirmed"
	}
	return &TableInspection{
		Status: status, Schema: database, Table: table, Columns: []ColumnInfo{}, Reason: reason,
	}
}

func buildMySQLManagedInspection(database, table string, server mysqlManagedServerFacts, facts mysqlManagedTableFacts) (*MySQLManagedInspection, error) {
	evidence := mysqlManagedBindingEvidence(server, database)
	inspection := &TableInspection{Schema: database, Table: table, Columns: []ColumnInfo{}, Status: facts.status, Reason: facts.reason}
	switch facts.status {
	case TableInspectionMissing:
		return &MySQLManagedInspection{Inspection: inspection, Evidence: evidence}, nil
	case TableInspectionForbidden:
		return &MySQLManagedInspection{Inspection: inspection, Evidence: evidence}, ErrMySQLManagedTableForbidden
	case TableInspectionExists:
	default:
		return &MySQLManagedInspection{Inspection: inspection, Evidence: evidence}, ErrMySQLManagedTableUnconfirmed
	}

	evidence.Engine = facts.engine
	evidence.Charset = facts.charset
	evidence.Collation = facts.collation
	evidence.TableType = facts.tableType
	evidence.Partitions = cloneMySQLPartitionEvidence(facts.partitions)

	primaryKey := make(map[string]struct{}, len(facts.primaryKey))
	for _, name := range facts.primaryKey {
		if _, duplicate := primaryKey[name]; duplicate {
			return &MySQLManagedInspection{Inspection: inspection, Evidence: evidence}, ErrMySQLManagedTableUnconfirmed
		}
		primaryKey[name] = struct{}{}
	}
	columnNames := make(map[string]struct{}, len(facts.columns))
	evidence.ColumnCollations = map[string]string{}
	evidence.ColumnDefaults = map[string]*string{}
	evidence.JSONColumns = map[string]bool{}
	for _, column := range facts.columns {
		columnNames[column.name] = struct{}{}
		_, isPrimary := primaryKey[column.name]
		dataType := column.columnType
		if strings.TrimSpace(dataType) == "" {
			dataType = column.dataType
		}
		inspection.Columns = append(inspection.Columns, ColumnInfo{
			Name: column.name, DataType: dataType, Nullable: column.nullable,
			PrimaryKey: isPrimary, HasDefault: column.hasDefault,
		})
		if column.collation != "" {
			evidence.ColumnCollations[column.name] = column.collation
		}
		if column.defaultValue != nil {
			value := *column.defaultValue
			evidence.ColumnDefaults[column.name] = &value
		} else {
			evidence.ColumnDefaults[column.name] = nil
		}
		if column.jsonValid || facts.jsonColumns[column.name] {
			evidence.JSONColumns[column.name] = true
		}
	}
	for name := range primaryKey {
		if _, exists := columnNames[name]; !exists {
			return &MySQLManagedInspection{Inspection: inspection, Evidence: evidence}, ErrMySQLManagedTableUnconfirmed
		}
	}
	evidence.PrimaryKey = slices.Clone(facts.primaryKey)
	inspection.IndexesInspected = true
	inspection.Indexes = make([]IndexInfo, 0, len(facts.indexes))
	for _, index := range facts.indexes {
		if reason := unusableMySQLManagedIndex(index); reason != "" {
			return &MySQLManagedInspection{Inspection: inspection, Evidence: evidence},
				fmt.Errorf("%w: index %s (%s)", ErrMySQLManagedIndexUnusable, index.name, reason)
		}
		entry := IndexInfo{Name: index.name, Columns: slices.Clone(index.columns), Unique: index.unique}
		inspection.Indexes = append(inspection.Indexes, entry)
	}
	return &MySQLManagedInspection{Inspection: inspection, Evidence: evidence, Indexes: inspection.Indexes}, nil
}

func unusableMySQLManagedIndex(index mysqlManagedIndexFacts) string {
	switch {
	case strings.TrimSpace(index.name) == "":
		return "missing_name"
	case index.functional:
		return mysqlManagedIndexFunctional
	case len(index.columns) == 0:
		return mysqlManagedIndexFunctional
	case index.prefix:
		return "prefix"
	case index.ignored:
		return mysqlManagedIndexIgnored
	case !index.visible:
		return "invisible"
	case !strings.EqualFold(index.indexType, mysqlManagedIndexTypeBTREE):
		return "non_btree"
	}
	for _, column := range index.columns {
		if strings.TrimSpace(column) == "" || strings.ContainsAny(column, "()") {
			return mysqlManagedIndexFunctional
		}
	}
	return ""
}

func validateMySQLManagedStatement(statement string) error {
	normalized := strings.TrimSpace(statement)
	if normalized == "" {
		return ErrMySQLManagedStatementRefused
	}
	withoutTrailing := strings.TrimSuffix(normalized, ";")
	if strings.Contains(withoutTrailing, ";") {
		return ErrMySQLManagedStatementRefused
	}
	upper := strings.ToUpper(withoutTrailing)
	if !strings.HasPrefix(upper, "CREATE TABLE ") && !strings.HasPrefix(upper, "CREATE INDEX ") {
		return ErrMySQLManagedStatementRefused
	}
	for _, forbidden := range []string{
		"IF NOT EXISTS", "CREATE DATABASE", "DROP ", "TRUNCATE ", "ALTER ", "RENAME ",
		"REPLACE ", "INSERT ", "UPDATE ", "DELETE ",
	} {
		if strings.Contains(upper, forbidden) {
			return ErrMySQLManagedStatementRefused
		}
	}
	return nil
}
