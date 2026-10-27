package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	// MySQLPartitionMaintenanceAction is the durable action recorded in the existing schema-operation ledger.
	MySQLPartitionMaintenanceAction     = "mysql_partition_maintenance"
	mysqlPartitionMaintenanceReady      = "ready"
	mysqlPartitionMaintenanceBlocked    = "blocked"
	mysqlPartitionMaintenanceCode       = "partition_maintenance_blocked"
	mysqlPartitionMaintenanceNext       = "inspect the destination and restore the last confirmed partition state"
	mysqlPartitionMaintenanceKind       = "mysql-native-monthly-partition-extension"
	mysqlPartitionMaintenanceBatch      = 4
	mysqlPartitionPastName              = "p_past"
	mysqlPartitionFutureName            = "p_future"
	mysqlPartitionLayoutVersion         = "managed_long_v1"
	mysqlPartitionStorageMode           = "native_monthly"
	mysqlPartitionReasonCatalogMismatch = "catalog_mismatch"
	mysqlPartitionReasonUnverifiable    = "catalog_unverifiable"
)

var (
	// ErrMySQLPartitionMaintenanceRequest marks an invalid maintenance scope or request.
	ErrMySQLPartitionMaintenanceRequest = errors.New("mysql partition maintenance request is invalid")
	// ErrMySQLPartitionMaintenanceStateMismatch marks an unverified or stale partition projection.
	ErrMySQLPartitionMaintenanceStateMismatch = errors.New("mysql partition maintenance state does not match confirmed storage")
	// ErrMySQLPartitionMaintenanceUnsupported marks a repository without durable maintenance support.
	ErrMySQLPartitionMaintenanceUnsupported = errors.New("mysql partition maintenance is unsupported by this repository")
)

// MySQLPartitionMaintenanceRequest identifies the confirmed native-monthly scope to maintain.
type MySQLPartitionMaintenanceRequest struct {
	Scope                      StorageScope
	GroupID, ConnectorRevision string
	At                         time.Time
}

// MySQLPartitionMaintenanceState is the mutable authorized partition projection. Its base digest and
// generation bind it to immutable storage confirmation; Current is advanced only with a verified ledger result.
type MySQLPartitionMaintenanceState struct {
	Scope                               StorageScope
	BaseLayoutDigest, ConnectorRevision string
	Generation                          int
	Current                             *MySQLPartitionEvidence
	Revision                            int64
	LastOperationID, Status, Reason     string
	UpdatedAt                           time.Time
}

// MySQLPartitionMaintenanceTarget exposes catalog inspection and a typed server-side reorganization.
// Reorganize must derive its SQL from the exact before/after evidence and must not accept caller SQL.
type MySQLPartitionMaintenanceTarget struct {
	Inspect    func(context.Context, string) (TargetTableInspection, error)
	Reorganize func(context.Context, *MySQLPartitionEvidence, *MySQLPartitionEvidence) error
}

// MySQLPartitionMaintenanceResult reports a verified state transition or a stable maintenance refusal.
type MySQLPartitionMaintenanceResult struct {
	State        *MySQLPartitionMaintenanceState
	Operation    *SchemaOperation
	Code, Reason string
	Changed      bool
}

type mysqlPartitionMaintenanceDetail struct {
	Kind              string                  `json:"kind"`
	Scope             StorageScope            `json:"scope"`
	GroupID           string                  `json:"group_id"`
	ConnectorRevision string                  `json:"connector_revision"`
	BaseLayoutDigest  string                  `json:"base_layout_digest"`
	Generation        int                     `json:"generation"`
	StateRevision     int64                   `json:"state_revision"`
	At                time.Time               `json:"at"`
	Before            *MySQLPartitionEvidence `json:"before"`
	After             *MySQLPartitionEvidence `json:"after"`
}

func validateMySQLPartitionMaintenanceScope(scope StorageScope) bool {
	return strings.TrimSpace(scope.WorkspaceID) != "" && strings.TrimSpace(scope.ConnectorID) != "" &&
		strings.TrimSpace(scope.Database) != "" && scope.Database == scope.TableSchema &&
		scope.StorageMode == mysqlPartitionStorageMode && scope.LayoutVersion == mysqlPartitionLayoutVersion
}

func validateMySQLPartitionMaintenanceState(state *MySQLPartitionMaintenanceState) error {
	if state == nil || !validateMySQLPartitionMaintenanceScope(state.Scope) || len(state.BaseLayoutDigest) != sha256.Size*2 ||
		state.ConnectorRevision == "" || state.Generation < 1 || state.Revision < 1 || state.UpdatedAt.IsZero() ||
		(state.Status != mysqlPartitionMaintenanceReady && state.Status != mysqlPartitionMaintenanceBlocked) {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	if _, err := hex.DecodeString(state.BaseLayoutDigest); err != nil {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	if err := validateCurrentMySQLPartitionEvidence(state.Current); err != nil {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	if state.Status == mysqlPartitionMaintenanceReady && state.Reason != "" {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	return nil
}

func digestGroupSchemaLayout(layout GroupSchemaLayout) string {
	normalized, err := ValidateGroupSchemaLayout(layout)
	if err != nil {
		return ""
	}
	encoded, err := json.Marshal(normalized)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

// MaintainMySQLPartitions advances one confirmed native-monthly scope by at most one bounded DDL operation.
