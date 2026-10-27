package dbtarget

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
)

const (
	mysqlRecordingCapabilityProjectionVersion = "mysql-recording-capability-v1"
	mysqlRecordingMetadataVerified            = "verified"
	mysqlRecordingMetadataForbidden           = "forbidden"
	mysqlRecordingMetadataUnconfirmed         = "unconfirmed"
	mysqlRecordingFamilyUnknown               = "unknown"
	mysqlRecordingReasonRollbackFailed        = "metadata_rollback_failed"
)

// MySQLRecordingCapability is the backend-derived capability projection for a
// saved MySQL-kind connector. Unverified operations remain unavailable.
type MySQLRecordingCapability struct {
	ConnectorID             string   `json:"connector_id"`
	ConnectorRevision       string   `json:"connector_revision"`
	Database                string   `json:"database"`
	ServerFamily            string   `json:"server_family"`
	ServerVersion           string   `json:"server_version"`
	CapabilityRevision      string   `json:"capability_revision"`
	Connected               bool     `json:"connected"`
	MetadataStatus          string   `json:"metadata_status"`
	TransactionalEngine     bool     `json:"transactional_engine"`
	SupportsManagedLong     bool     `json:"supports_managed_long"`
	SupportsNativePartition bool     `json:"supports_native_partition"`
	SupportsGroupTestWrite  bool     `json:"supports_group_test_write"`
	SupportsMissingRepair   bool     `json:"supports_missing_repair"`
	Reasons                 []string `json:"reasons"`
}

type mysqlRecordingMetadata struct {
	serverFamily        string
	serverVersion       string
	database            string
	actualDatabase      string
	metadataStatus      string
	transactionalEngine bool
	reasons             []string
}

// ProbeMySQLRecordingAtRevision reads family and destination evidence from the
// current saved connector revision without using any database-provisioning path.
func (s *ConnectorService) ProbeMySQLRecordingAtRevision(
	ctx context.Context,
	connectorID, expectedRevision string,
) (*MySQLRecordingCapability, error) {
	if s == nil || s.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := s.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, err
	}
	if connector.Kind != schema.DatabaseConnectorKindMySQL {
		return nil, managedDestinationInvalid("saved connector kind is not mysql")
	}

	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, managedDestinationInvalid("saved mysql connector configuration is invalid")
	}
	dbConfig, err := buildExternalDBConfig(connector.Kind, config)
	if err != nil {
		return nil, managedDestinationInvalid("saved mysql connector configuration is invalid")
	}
	driverConfig, err := mysqldriver.ParseDSN(dbConfig.DSN)
	if err != nil {
		return nil, managedDestinationInvalid("saved mysql connector configuration is invalid")
	}

	metadata := mysqlRecordingMetadata{
		serverFamily:   mysqlRecordingFamilyUnknown,
		database:       driverConfig.DBName,
		metadataStatus: mysqlRecordingMetadataUnconfirmed,
		reasons:        []string{},
	}
	if !mysqlRecordingTransportVerified(driverConfig) {
		metadata.reasons = append(metadata.reasons, "transport_unverified")
		return projectMySQLRecordingCapability(connector, metadata, false), nil
	}

	// Open the saved connector directly. Do not call probeConnector: its MySQL
	// path may create a missing database, which is forbidden for this probe.
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		reason := "connection_failed"
		if classifyConnectorError(err) == schema.DatabaseConnectorStatusAuthFailed {
			reason = "authentication_failed"
		}
		metadata.reasons = append(metadata.reasons, reason)
		return projectMySQLRecordingCapability(connector, metadata, false), nil
	}
	if manager == nil || manager.DB() == nil {
		if manager != nil {
			_ = manager.Close()
		}
		metadata.reasons = append(metadata.reasons, "connection_failed")
		return projectMySQLRecordingCapability(connector, metadata, false), nil
	}
	defer manager.Close()

	metadata = inspectMySQLRecordingMetadata(ctx, manager.DB(), driverConfig.DBName)
	return projectMySQLRecordingCapability(connector, metadata, true), nil
}

func inspectMySQLRecordingMetadata(ctx context.Context, db *sql.DB, expectedDatabase string) (metadata mysqlRecordingMetadata) {
	metadata = mysqlRecordingMetadata{
		serverFamily:   mysqlRecordingFamilyUnknown,
		database:       expectedDatabase,
		metadataStatus: mysqlRecordingMetadataUnconfirmed,
		reasons:        []string{},
	}
	if db == nil {
		metadata.reasons = append(metadata.reasons, "metadata_unconfirmed")
		return metadata
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		setMySQLRecordingMetadataFailure(&metadata, err)
		return metadata
	}
	defer func() {
		if rollbackErr := rollbackReadOnly(tx); rollbackErr != nil {
			metadata.metadataStatus = mysqlRecordingMetadataUnconfirmed
			metadata.reasons = appendMySQLRecordingReason(metadata.reasons, mysqlRecordingReasonRollbackFailed)
		}
	}()

	var version, versionComment string
	var currentDatabase sql.NullString
	if err := tx.QueryRowContext(ctx, "SELECT VERSION(), @@version_comment, DATABASE()").Scan(
		&version, &versionComment, &currentDatabase,
	); err != nil {
		setMySQLRecordingMetadataFailure(&metadata, err)
		return metadata
	}
	metadata.serverVersion = version
	metadata.serverFamily = mysqlRecordingServerFamily(version, versionComment)
	metadata.actualDatabase = currentDatabase.String
	if metadata.database == "" && currentDatabase.Valid {
		metadata.database = currentDatabase.String
	}
	if !currentDatabase.Valid || currentDatabase.String == "" {
		metadata.reasons = append(metadata.reasons, "database_scope_unconfirmed")
		return metadata
	}
	if expectedDatabase == "" {
		metadata.reasons = append(metadata.reasons, "database_scope_unconfirmed")
		return metadata
	}
	if currentDatabase.String != expectedDatabase {
		metadata.reasons = append(metadata.reasons, "database_scope_mismatch")
		return metadata
	}

	var visibleDatabase string
	err = tx.QueryRowContext(ctx,
		"SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?",
		expectedDatabase,
	).Scan(&visibleDatabase)
	if errors.Is(err, sql.ErrNoRows) {
		metadata.reasons = append(metadata.reasons, "database_visibility_unconfirmed")
		return metadata
	}
	if err != nil {
		setMySQLRecordingMetadataFailure(&metadata, err)
		return metadata
	}
	if visibleDatabase != expectedDatabase {
		metadata.reasons = append(metadata.reasons, "database_scope_mismatch")
		return metadata
	}

	var engineSupport, engineTransactions string
	err = tx.QueryRowContext(ctx,
		"SELECT SUPPORT, TRANSACTIONS FROM information_schema.ENGINES WHERE ENGINE = 'InnoDB'",
	).Scan(&engineSupport, &engineTransactions)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		setMySQLRecordingMetadataFailure(&metadata, err)
		return metadata
	}
	metadata.transactionalEngine = err == nil &&
		(engineSupport == "YES" || engineSupport == "DEFAULT") && engineTransactions == "YES"
	metadata.metadataStatus = mysqlRecordingMetadataVerified

	if err := tx.Commit(); err != nil {
		setMySQLRecordingMetadataFailure(&metadata, err)
	}
	return metadata
}

func setMySQLRecordingMetadataFailure(metadata *mysqlRecordingMetadata, err error) {
	if isPermissionDeniedError(err) {
		metadata.metadataStatus = mysqlRecordingMetadataForbidden
		metadata.reasons = appendMySQLRecordingReason(metadata.reasons, "metadata_forbidden")
		return
	}
	metadata.metadataStatus = mysqlRecordingMetadataUnconfirmed
	metadata.reasons = appendMySQLRecordingReason(metadata.reasons, "metadata_unconfirmed")
}

func mysqlRecordingServerFamily(version, versionComment string) string {
	serverIdentity := strings.ToLower(version + " " + versionComment)
	switch {
	case strings.Contains(serverIdentity, "mariadb"):
		return "mariadb"
	case strings.Contains(strings.ToLower(versionComment), "mysql"):
		return "mysql"
	default:
		return mysqlRecordingFamilyUnknown
	}
}

func projectMySQLRecordingCapability(
	connector *schema.DatabaseConnector,
	metadata mysqlRecordingMetadata,
	connected bool,
) *MySQLRecordingCapability {
	capability := &MySQLRecordingCapability{
		ServerFamily:        metadata.serverFamily,
		ServerVersion:       metadata.serverVersion,
		Database:            metadata.database,
		Connected:           connected,
		MetadataStatus:      metadata.metadataStatus,
		TransactionalEngine: metadata.transactionalEngine,
		Reasons:             []string{},
	}
	if connector != nil {
		capability.ConnectorID = connector.ID
		capability.ConnectorRevision = connector.IdentityRevision
	}
	if capability.ServerFamily == "" {
		capability.ServerFamily = mysqlRecordingFamilyUnknown
	}
	if capability.MetadataStatus == "" {
		capability.MetadataStatus = mysqlRecordingMetadataUnconfirmed
	}
	for _, reason := range metadata.reasons {
		capability.Reasons = appendMySQLRecordingReason(capability.Reasons, reason)
	}
	if capability.MetadataStatus == mysqlRecordingMetadataVerified && !capability.TransactionalEngine {
		capability.Reasons = appendMySQLRecordingReason(capability.Reasons, "innodb_unavailable")
	}
	if !mysqlRecordingInitialBaseline(metadata.serverFamily, metadata.serverVersion) {
		capability.Reasons = appendMySQLRecordingReason(capability.Reasons, "server_version_unverified")
	}
	if mysqlRecordingOperationsVerified(metadata, capability) {
		capability.SupportsManagedLong = true
		capability.SupportsNativePartition = true
		capability.SupportsGroupTestWrite = true
		capability.SupportsMissingRepair = true
	} else {
		// Client values, driver presence, a major/minor match, or a saved
		// capability revision cannot open an operation gate without the exact
		// verified server pair and complete backend facts.
		capability.Reasons = appendMySQLRecordingReason(capability.Reasons,
			"adapter_unverified",
			"managed_long_unverified",
			"native_partition_unverified",
			"group_test_write_unverified",
			"missing_repair_unverified",
		)
	}
	capability.CapabilityRevision = mysqlRecordingCapabilityRevision(connector, metadata, capability)
	return capability
}

func mysqlRecordingOperationsVerified(metadata mysqlRecordingMetadata, capability *MySQLRecordingCapability) bool {
	return capability != nil && capability.Connected && capability.MetadataStatus == mysqlRecordingMetadataVerified &&
		capability.TransactionalEngine && metadata.database != "" && metadata.database == metadata.actualDatabase &&
		MySQLRecordingVersionVerified(metadata.serverFamily, metadata.serverVersion)
}

func mysqlRecordingInitialBaseline(family, version string) bool {
	version = strings.TrimPrefix(version, "5.5.5-")
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return false
	}
	switch family {
	case "mysql":
		return parts[0] == "8" && parts[1] == "4"
	case "mariadb":
		return parts[0] == "11" && parts[1] == "4"
	default:
		return false
	}
}

func mysqlRecordingCapabilityRevision(
	connector *schema.DatabaseConnector,
	metadata mysqlRecordingMetadata,
	capability *MySQLRecordingCapability,
) string {
	connectorID := ""
	connectorRevision := ""
	if connector != nil {
		connectorID = connector.ID
		connectorRevision = connector.IdentityRevision
	}
	// Operation rollout availability is not a physical destination fact. Keep
	// the four legacy false slots for byte compatibility with already prepared
	// schemas, while connector/server/metadata changes still invalidate proof.
	projection := struct {
		ProjectionVersion     string
		ConnectorID           string
		ConnectorRevision     string
		ServerFamily          string
		ServerVersion         string
		Database              string
		ActualDatabase        string
		MetadataStatus        string
		TransactionalEngine   bool
		Connected             bool
		SupportsManagedLong   bool
		SupportsPartition     bool
		SupportsGroupWrite    bool
		SupportsMissingRepair bool
	}{
		ProjectionVersion:     mysqlRecordingCapabilityProjectionVersion,
		ConnectorID:           connectorID,
		ConnectorRevision:     connectorRevision,
		ServerFamily:          metadata.serverFamily,
		ServerVersion:         metadata.serverVersion,
		Database:              metadata.database,
		ActualDatabase:        metadata.actualDatabase,
		MetadataStatus:        metadata.metadataStatus,
		TransactionalEngine:   capability.TransactionalEngine,
		Connected:             capability.Connected,
		SupportsManagedLong:   false,
		SupportsPartition:     false,
		SupportsGroupWrite:    false,
		SupportsMissingRepair: false,
	}
	// Quoted strings, fixed separators, and boolean literals form a stable,
	// unambiguous encoding without a serialization error path.
	canonical := fmt.Appendf(nil,
		"%q|%q|%q|%q|%q|%q|%q|%q|%t|%t|%t|%t|%t|%t",
		projection.ProjectionVersion,
		projection.ConnectorID,
		projection.ConnectorRevision,
		projection.ServerFamily,
		projection.ServerVersion,
		projection.Database,
		projection.ActualDatabase,
		projection.MetadataStatus,
		projection.TransactionalEngine,
		projection.Connected,
		projection.SupportsManagedLong,
		projection.SupportsPartition,
		projection.SupportsGroupWrite,
		projection.SupportsMissingRepair,
	)
	digest := sha256.Sum256(canonical)
	return "mysqlcap-v1-" + hex.EncodeToString(digest[:])
}

func appendMySQLRecordingReason(reasons []string, additions ...string) []string {
	for _, addition := range additions {
		if addition == "" {
			continue
		}
		if !slices.Contains(reasons, addition) {
			reasons = append(reasons, addition)
		}
	}
	return reasons
}

func mysqlRecordingTransportVerified(config *mysqldriver.Config) bool {
	return config != nil && config.TLS != nil && !config.TLS.InsecureSkipVerify && !config.AllowFallbackToPlaintext
}
