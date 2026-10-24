package recordingplan

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// StorageScope is the identity of one shared managed storage, the key of the
// local write_group_storage_scopes row: the workspace and the confirmed
// destination, with the storage mode and layout that derive its table.
//
// The JSON names are part of the persisted detail of a storage repair operation.
type StorageScope struct {
	WorkspaceID   string `json:"workspace_id"`
	ConnectorID   string `json:"connector_id"`
	Database      string `json:"database"`
	TableSchema   string `json:"table_schema"`
	StorageMode   string `json:"storage_mode"`
	LayoutVersion string `json:"layout_version"`
}

// DBTX is the part of *sql.DB and *sql.Tx the lineage records need, so the same
// code runs inside the caller's transaction.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// ErrStorageLineageNotFound means no storage scope was ever saved for the key.
var ErrStorageLineageNotFound = errors.New("storage lineage not found")

// Storage event kinds and repair states. They name facts, never values.
const (
	StorageEventConfirmed             = "confirmed"
	StorageEventRecreatedByConfirm    = "recreated_by_confirmation"
	StorageEventReconstructed         = "reconstructed"
	StorageEventRepairBlocked         = "repair_blocked"
	StorageRepairStateBlocked         = "blocked"
	storageLineageTimeLayout          = time.RFC3339Nano
	storageLineageMaxDiagnosticEvents = 20
)

// StorageLineage is what the local database knows about the confirmed table of
// a scope: which generation it is, the layout and connector revision that were
// verified, and whether a rebuilt table means earlier history may be gone.
type StorageLineage struct {
	StorageScope
	TableName, OwnerColumn, LayoutDigest string
	Generation                           int
	// Layout is the layout verified at the last confirmation; nil before any.
	Layout            *GroupSchemaLayout
	ConnectorRevision string
	// ConfirmedOperationID and ConfirmedDigest identify that verification.
	ConfirmedOperationID, ConfirmedDigest string
	HistoryMayBeLost                      bool
	Reconstructions                       int
	LastReconstructedAt                   time.Time
	LastReconstructionOperationID         string
	// RepairState is non-empty (StorageRepairStateBlocked) when an automatic
	// restore ended in a result nobody can vouch for; RepairReason says why.
	RepairState, RepairReason string
}

// Confirmed reports whether a verified confirmation exists for the scope.
func (l *StorageLineage) Confirmed() bool {
	return l != nil && l.ConfirmedOperationID != "" && l.Layout != nil
}

const storageLineageColumns = `table_name, owner_column, layout_digest, generation, confirmed_layout,
	confirmed_connector_revision, confirmed_operation_id, confirmed_digest, history_may_be_lost, reconstruction_count,
	last_reconstructed_at, last_reconstruction_operation_id, repair_state, repair_reason`

const storageScopeWhere = `workspace_id = $1 AND connector_id = $2 AND database_name = $3 AND table_schema = $4
	AND storage_mode = $5 AND layout_version = $6`

// numberedPlaceholders turns $n into SQLite's ?n, which binds by number. The
// lineage statements list their parameters out of order (a SET clause before
// the scope key), which the in-order ? form of adaptPlaceholders cannot express.
func numberedPlaceholders(query string) string {
	for i := 30; i >= 1; i-- {
		query = strings.ReplaceAll(query, fmt.Sprintf("$%d", i), fmt.Sprintf("?%d", i))
	}
	return query
}

func scopeArgs(scope StorageScope) []any {
	return []any{scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion}
}

// ReadStorageLineage reads the lineage of a scope.
func ReadStorageLineage(ctx context.Context, q DBTX, scope StorageScope) (*StorageLineage, error) {
	lineage := &StorageLineage{StorageScope: scope}
	var layout, reconstructedAt string
	var lost int
	err := q.QueryRowContext(ctx, numberedPlaceholders(`SELECT `+storageLineageColumns+` FROM write_group_storage_scopes WHERE `+storageScopeWhere), scopeArgs(scope)...).Scan(
		&lineage.TableName, &lineage.OwnerColumn, &lineage.LayoutDigest, &lineage.Generation, &layout,
		&lineage.ConnectorRevision, &lineage.ConfirmedOperationID, &lineage.ConfirmedDigest, &lost, &lineage.Reconstructions,
		&reconstructedAt, &lineage.LastReconstructionOperationID, &lineage.RepairState, &lineage.RepairReason)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStorageLineageNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read storage lineage: %w", err)
	}
	lineage.HistoryMayBeLost = lost != 0
	if layout != "" {
		var decoded GroupSchemaLayout
		if err := json.Unmarshal([]byte(layout), &decoded); err != nil {
			return nil, fmt.Errorf("decode confirmed storage layout: %w", err)
		}
		validated, err := ValidateGroupSchemaLayout(decoded)
		if err != nil {
			return nil, fmt.Errorf("confirmed storage layout is invalid: %w", err)
		}
		lineage.Layout = &validated
	}
	if reconstructedAt != "" {
		if lineage.LastReconstructedAt, err = time.Parse(storageLineageTimeLayout, reconstructedAt); err != nil {
			return nil, fmt.Errorf("decode storage reconstruction time: %w", err)
		}
	}
	return lineage, nil
}

// StorageConfirmation is the evidence a successful, verified schema
// confirmation gives about the scope's table.
type StorageConfirmation struct {
	Layout            GroupSchemaLayout
	ConnectorRevision string
	OperationID       string
	VerifiedDigest    string
	// CreatedTable is true when the confirmed batch created the shared table
	// itself, as opposed to finding it already there.
	CreatedTable bool
	GroupID      string
}

// RecordStorageConfirmation freezes what a verified confirmation proved. The
// first confirmation of a scope keeps generation 1. A repeat of the same
// operation changes nothing. A confirmation whose batch had to create a table
// of an already confirmed scope means the table was gone: it starts the next
// generation and marks earlier history as possibly lost. Any verified
// confirmation clears a blocked repair state, because the table was just
// checked against the layout again.
func RecordStorageConfirmation(ctx context.Context, q DBTX, scope StorageScope, c StorageConfirmation, now time.Time) error {
	layout, err := ValidateGroupSchemaLayout(c.Layout)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(layout)
	if err != nil {
		return fmt.Errorf("encode confirmed storage layout: %w", err)
	}
	current, err := ReadStorageLineage(ctx, q, scope)
	if err != nil {
		return err
	}
	if current.ConfirmedOperationID == c.OperationID {
		return nil
	}
	timestamp := now.UTC().Format(storageLineageTimeLayout)
	if current.Confirmed() && c.CreatedTable {
		if _, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE write_group_storage_scopes SET generation = generation + 1,
			history_may_be_lost = 1, confirmed_layout = $7, confirmed_connector_revision = $8, confirmed_operation_id = $9,
			confirmed_digest = $10, repair_state = '', repair_reason = '' WHERE `+storageScopeWhere+` AND generation = $11`),
			append(scopeArgs(scope), string(encoded), c.ConnectorRevision, c.OperationID, c.VerifiedDigest, current.Generation)...); err != nil {
			return fmt.Errorf("record recreated storage: %w", err)
		}
		return recordStorageEvent(ctx, q, scope, current.Generation+1, c.GroupID, StorageEventRecreatedByConfirm, "", "table-created-again", c.OperationID, timestamp)
	}
	if _, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE write_group_storage_scopes SET confirmed_layout = $7,
		confirmed_connector_revision = $8, confirmed_operation_id = $9, confirmed_digest = $10, repair_state = '', repair_reason = ''
		WHERE `+storageScopeWhere),
		append(scopeArgs(scope), string(encoded), c.ConnectorRevision, c.OperationID, c.VerifiedDigest)...); err != nil {
		return fmt.Errorf("record storage confirmation: %w", err)
	}
	return recordStorageEvent(ctx, q, scope, current.Generation, c.GroupID, StorageEventConfirmed, "", "", c.OperationID, timestamp)
}

// RecordStorageReconstruction moves a scope from fromGeneration to the next
// after the repair operation rebuilt its table. It applies at most once per
// generation, so a repeat of the same completion is harmless. The returned
// bool says whether this call applied it.
func RecordStorageReconstruction(ctx context.Context, q DBTX, scope StorageScope, fromGeneration int, operationID, groupID string, now time.Time) (bool, error) {
	timestamp := now.UTC().Format(storageLineageTimeLayout)
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE write_group_storage_scopes SET generation = generation + 1,
		history_may_be_lost = 1, reconstruction_count = reconstruction_count + 1, last_reconstructed_at = $7,
		last_reconstruction_operation_id = $8, repair_state = '', repair_reason = '' WHERE `+storageScopeWhere+` AND generation = $9`),
		append(scopeArgs(scope), timestamp, operationID, fromGeneration)...)
	if err != nil {
		return false, fmt.Errorf("record storage reconstruction: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("read storage reconstruction: %w", err)
	}
	if affected != 1 {
		return false, nil
	}
	return true, recordStorageEvent(ctx, q, scope, fromGeneration+1, groupID, StorageEventReconstructed, "", "table-rebuilt", operationID, timestamp)
}

// BlockStorageRepair records that an automatic restore ended in a result
// nobody can vouch for. Delivery to the scope waits until an explicit
// confirmation checks the table again.
func BlockStorageRepair(ctx context.Context, q DBTX, scope StorageScope, generation int, reason, operationID, groupID string, now time.Time) error {
	if _, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE write_group_storage_scopes SET repair_state = $7, repair_reason = $8
		WHERE `+storageScopeWhere), append(scopeArgs(scope), StorageRepairStateBlocked, reason)...); err != nil {
		return fmt.Errorf("block storage repair: %w", err)
	}
	return recordStorageEvent(ctx, q, scope, generation, groupID, StorageEventRepairBlocked, "schema_repair_blocked", reason, operationID, now.UTC().Format(storageLineageTimeLayout))
}

// RecordStorageRefusal records a recovery that was refused without any change.
func RecordStorageRefusal(ctx context.Context, q DBTX, scope StorageScope, generation int, groupID, code, reason string, now time.Time) error {
	return recordStorageEvent(ctx, q, scope, generation, groupID, StorageEventRepairBlocked, code, reason, "", now.UTC().Format(storageLineageTimeLayout))
}

func recordStorageEvent(ctx context.Context, q DBTX, scope StorageScope, generation int, groupID, kind, code, reason, operationID, timestamp string) error {
	args := append(scopeArgs(scope), generation, groupID, kind, code, reason, operationID, timestamp)
	_, err := q.ExecContext(ctx, numberedPlaceholders(`INSERT INTO write_group_storage_events (workspace_id, connector_id, database_name,
		table_schema, storage_mode, layout_version, generation, group_id, kind, code, reason, operation_id, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version, generation, group_id, kind, reason)
		DO UPDATE SET occurrences = occurrences + 1, occurred_at = excluded.occurred_at, code = excluded.code`), args...)
	if err != nil {
		return fmt.Errorf("record storage event: %w", err)
	}
	return nil
}

// StorageEvent is one value-free diagnostic about managed long storage.
type StorageEvent struct {
	Generation  int       `json:"generation"`
	GroupID     string    `json:"group_id,omitempty"`
	Kind        string    `json:"kind"`
	Code        string    `json:"code,omitempty"`
	Reason      string    `json:"reason,omitempty"`
	OperationID string    `json:"operation_id,omitempty"`
	Occurrences int       `json:"occurrences"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// ReadStorageEvents returns the newest diagnostics of a scope that concern the
// group or the whole storage (no group), newest first and bounded.
func ReadStorageEvents(ctx context.Context, q DBTX, scope StorageScope, groupID string) ([]StorageEvent, error) {
	rows, err := q.QueryContext(ctx, numberedPlaceholders(`SELECT generation, group_id, kind, code, reason, operation_id, occurrences, occurred_at
		FROM write_group_storage_events WHERE `+storageScopeWhere+` AND (group_id = $7 OR group_id = '') ORDER BY id DESC LIMIT $8`),
		append(scopeArgs(scope), groupID, storageLineageMaxDiagnosticEvents)...)
	if err != nil {
		return nil, fmt.Errorf("read storage events: %w", err)
	}
	defer rows.Close()
	events := []StorageEvent{}
	for rows.Next() {
		var event StorageEvent
		var occurred string
		if err := rows.Scan(&event.Generation, &event.GroupID, &event.Kind, &event.Code, &event.Reason, &event.OperationID, &event.Occurrences, &occurred); err != nil {
			return nil, fmt.Errorf("scan storage event: %w", err)
		}
		if event.OccurredAt, err = time.Parse(storageLineageTimeLayout, occurred); err != nil {
			return nil, fmt.Errorf("decode storage event time: %w", err)
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
