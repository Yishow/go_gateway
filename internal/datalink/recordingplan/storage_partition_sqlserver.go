package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	SQLServerPartitionMaintenanceAction      = "sqlserver_partition_maintenance"
	sqlServerPartitionMaintenanceReady       = "ready"
	sqlServerPartitionMaintenanceBlocked     = "blocked"
	sqlServerPartitionMaintenanceCode        = "partition_maintenance_blocked"
	sqlServerPartitionMaintenanceNext        = "inspect the destination and restore the last confirmed partition state"
	sqlServerPartitionMaintenanceKind        = "sqlserver-native-monthly-partition-extension"
	sqlServerPartitionMaintenanceBatch       = 4
	sqlServerPartitionPastName               = "p_past"
	sqlServerPartitionFutureName             = "p_future"
	sqlServerPartitionLayoutVersion          = "managed_long_v1"
	sqlServerPartitionStorageMode            = "native_monthly"
	sqlServerPartitionReasonUnverifiable     = "catalog_unverifiable"
	sqlServerPartitionReasonCatalogMismatch  = "catalog_mismatch"
	sqlServerPartitionReasonPermissionLost   = "permission_lost"
	sqlServerPartitionReasonStatementFailed  = "statement_failed"
	sqlServerPartitionReasonAttemptAbandoned = "attempt_abandoned"
)

var (
	ErrSQLServerPartitionMaintenanceRequest       = errors.New("sqlserver partition maintenance request is invalid")
	ErrSQLServerPartitionMaintenanceStateMismatch = errors.New("sqlserver partition maintenance state does not match confirmed storage")
	ErrSQLServerPartitionMaintenanceUnsupported   = errors.New("sqlserver partition maintenance is unsupported by this repository")
)

// SQLServerPartitionMaintenanceRequest identifies one confirmed native scope.
type SQLServerPartitionMaintenanceRequest struct {
	Scope                      StorageScope
	GroupID, ConnectorRevision string
	At                         time.Time
}

// SQLServerPartitionMaintenanceState is the mutable authorized current
// projection. The initial frozen layout remains untouched.
type SQLServerPartitionMaintenanceState struct {
	Scope                               StorageScope
	BaseLayoutDigest, ConnectorRevision string
	Generation                          int
	Current                             *SQLServerPartitionEvidence
	Revision                            int64
	LastOperationID, Status, Reason     string
	UpdatedAt                           time.Time
}

// SQLServerPartitionMaintenanceTarget inspects the real catalog and performs a
// typed, append-only split derived from before/after evidence. It never accepts
// caller SQL.
type SQLServerPartitionMaintenanceTarget struct {
	Inspect TargetInspector
	Split   func(ctx context.Context, before, after *SQLServerPartitionEvidence) error
}

// SQLServerPartitionMaintenanceResult reports one verified transition or a
// stable maintenance refusal.
type SQLServerPartitionMaintenanceResult struct {
	State        *SQLServerPartitionMaintenanceState
	Operation    *SchemaOperation
	Code, Reason string
	Changed      bool
}

type sqlServerPartitionMaintenanceDetail struct {
	Kind              string                      `json:"kind"`
	Scope             StorageScope                `json:"scope"`
	GroupID           string                      `json:"group_id"`
	ConnectorRevision string                      `json:"connector_revision"`
	BaseLayoutDigest  string                      `json:"base_layout_digest"`
	Generation        int                         `json:"generation"`
	StateRevision     int64                       `json:"state_revision"`
	At                time.Time                   `json:"at"`
	Before            *SQLServerPartitionEvidence `json:"before"`
	After             *SQLServerPartitionEvidence `json:"after"`
}

func validateSQLServerPartitionMaintenanceScope(scope StorageScope) bool {
	return strings.TrimSpace(scope.WorkspaceID) != "" && strings.TrimSpace(scope.ConnectorID) != "" &&
		strings.TrimSpace(scope.Database) != "" && strings.TrimSpace(scope.TableSchema) != "" &&
		scope.StorageMode == sqlServerPartitionStorageMode && scope.LayoutVersion == sqlServerPartitionLayoutVersion
}

func validateSQLServerPartitionMaintenanceState(state *SQLServerPartitionMaintenanceState) error {
	if state == nil || !validateSQLServerPartitionMaintenanceScope(state.Scope) || len(state.BaseLayoutDigest) != sha256.Size*2 ||
		state.ConnectorRevision == "" || state.Generation < 1 || state.Revision < 1 || state.UpdatedAt.IsZero() ||
		(state.Status != sqlServerPartitionMaintenanceReady && state.Status != sqlServerPartitionMaintenanceBlocked) {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	if _, err := hex.DecodeString(state.BaseLayoutDigest); err != nil {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	if err := validateCurrentSQLServerPartitionEvidence(state.Current); err != nil {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	if state.Status == sqlServerPartitionMaintenanceReady && state.Reason != "" {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	return nil
}

// PlanSQLServerNativeMonthlyExtension retains the whole finite historical
// prefix and adds one gap boundary plus at most four monthly boundaries. A
// long outage therefore produces at most five splits, never an unbounded loop.
func PlanSQLServerNativeMonthlyExtension(current *SQLServerPartitionEvidence, at time.Time) (*SQLServerPartitionEvidence, error) {
	if at.IsZero() || validateCurrentSQLServerPartitionEvidence(current) != nil {
		return nil, fmt.Errorf("%w: invalid current evidence or zero planning time", ErrInvalidGroupSchemaLayout)
	}
	planned := CloneSQLServerPartitionEvidence(current)
	finite := planned.Partitions[:len(planned.Partitions)-1]
	lastBoundary, err := parseSQLServerPartitionBoundary(finite[len(finite)-1].LessThan)
	if err != nil {
		return nil, fmt.Errorf("%w: current final finite boundary", ErrInvalidGroupSchemaLayout)
	}
	utc := at.UTC()
	monthStart := time.Date(utc.Year(), utc.Month(), 1, 0, 0, 0, 0, time.UTC)
	targetEnd := monthStart.AddDate(0, sqlserverNativeMonthlyCount, 0)
	if !sqlServerPartitionYearInRange(monthStart.Year()) || !sqlServerPartitionYearInRange(targetEnd.Year()) {
		return nil, fmt.Errorf("%w: sqlserver maintenance boundary outside datetime2 range", ErrInvalidGroupSchemaLayout)
	}
	appendBoundary := func(part SQLServerPartitionBoundary) {
		part.Ordinal = len(finite) + 1
		finite = append(finite, part)
	}
	if lastBoundary.Before(monthStart) {
		name := "p_gap_" + monthStart.Format("200601")
		if lastBoundary.AddDate(0, 1, 0).Equal(monthStart) {
			name = "p_" + lastBoundary.Format("200601")
		}
		appendBoundary(SQLServerPartitionBoundary{Name: name, LessThan: monthStart.Format(sqlServerPartitionTimeLayout)})
		lastBoundary = monthStart
	}
	addedMonthly := 0
	for lastBoundary.Before(targetEnd) && addedMonthly < sqlServerPartitionMaintenanceBatch {
		next := lastBoundary.AddDate(0, 1, 0)
		if !sqlServerPartitionYearInRange(next.Year()) {
			return nil, fmt.Errorf("%w: sqlserver maintenance middle boundary outside datetime2 range", ErrInvalidGroupSchemaLayout)
		}
		appendBoundary(SQLServerPartitionBoundary{
			Name: "p_" + lastBoundary.Format("200601"), LessThan: next.Format(sqlServerPartitionTimeLayout),
		})
		lastBoundary = next
		addedMonthly++
	}
	if lastBoundary.Before(targetEnd) {
		return nil, fmt.Errorf("%w: sqlserver extension exceeded bounded monthly batch", ErrInvalidGroupSchemaLayout)
	}
	planned.Partitions = append(finite, SQLServerPartitionBoundary{
		Name: sqlServerPartitionFutureName, LessThan: sqlserverPartitionFutureValue, Ordinal: len(finite) + 1,
	})
	for i := range planned.Partitions {
		planned.Partitions[i].Ordinal = i + 1
	}
	if err := validateCurrentSQLServerPartitionEvidence(planned); err != nil {
		return nil, err
	}
	return planned, nil
}

func validateCurrentSQLServerPartitionEvidence(evidence *SQLServerPartitionEvidence) error {
	invalid := func(reason string) error {
		return fmt.Errorf("%w: sqlserver native monthly partition %s", ErrInvalidGroupSchemaLayout, reason)
	}
	if evidence == nil || evidence.Method != sqlserverPartitionMethod || evidence.Column != sqlserverPartitionColumn ||
		evidence.FunctionName == "" || evidence.SchemeName == "" || evidence.Filegroup != sqlserverPartitionFilegroup ||
		len(evidence.Partitions) < sqlserverPartitionBoundaryCount {
		return invalid("current partition definition is incomplete")
	}
	parts := evidence.Partitions
	past, err := parseSQLServerPartitionBoundary(parts[0].LessThan)
	if err != nil || parts[0] != (SQLServerPartitionBoundary{Name: sqlServerPartitionPastName, LessThan: past.Format(sqlServerPartitionTimeLayout), Ordinal: 1}) {
		return invalid("past catch-all boundary is not canonical")
	}
	seen := map[string]struct{}{strings.ToLower(parts[0].Name): {}}
	previous := past
	for index, part := range parts[1 : len(parts)-1] {
		boundary, parseErr := parseSQLServerPartitionBoundary(part.LessThan)
		if parseErr != nil || !previous.Before(boundary) || part.Ordinal != index+2 || part.Name == "" {
			return invalid("finite boundaries are not strictly increasing UTC month starts")
		}
		name := strings.ToLower(part.Name)
		if _, exists := seen[name]; exists {
			return invalid("partition names are duplicated")
		}
		seen[name] = struct{}{}
		expected := "p_" + previous.Format("200601")
		if !boundary.Equal(previous.AddDate(0, 1, 0)) {
			expected = "p_gap_" + boundary.Format("200601")
		}
		if part.Name != expected {
			return invalid("finite range name does not match its previous month or one-month gap")
		}
		previous = boundary
	}
	future := parts[len(parts)-1]
	if future != (SQLServerPartitionBoundary{Name: sqlServerPartitionFutureName, LessThan: sqlserverPartitionFutureValue, Ordinal: len(parts)}) {
		return invalid("future catch-all must remain final")
	}
	return nil
}

func sameSQLServerPartitionPrefix(before, after *SQLServerPartitionEvidence) bool {
	if before == nil || after == nil || len(after.Partitions) <= len(before.Partitions) {
		return false
	}
	return slices.Equal(before.Partitions[:len(before.Partitions)-1], after.Partitions[:len(before.Partitions)-1])
}

func sqlServerPartitionEvidenceDigest(evidence *SQLServerPartitionEvidence) string {
	encoded, err := json.Marshal(evidence)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}
