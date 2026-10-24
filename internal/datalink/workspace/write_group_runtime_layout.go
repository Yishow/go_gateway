package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
)

// WriteGroupRuntimeLayout freezes verified storage facts and applied tag types.
// It contains no connection configuration or credentials.
type WriteGroupRuntimeLayout struct {
	Dialect           dbtarget.SQLDialect        `json:"dialect"`
	Columns           []dbtarget.ColumnInfo      `json:"columns"`
	TagTypes          map[string]schema.DataType `json:"tag_types"`
	ReceiptTableReady bool                       `json:"receipt_table_ready,omitzero"`
	SchemaDigest      string                     `json:"schema_digest,omitempty"`
	// Hypertable is the frozen catalog evidence of a timescale group's table. It
	// is set only for that mode, after the catalog showed the owned hypertable.
	Hypertable *WriteGroupHypertableEvidence `json:"hypertable,omitempty"`
}

// WriteGroupHypertableEvidence is what the catalog showed when the layout was
// verified: the time dimension and chunk interval the group was applied on. The
// interval is kept in microseconds, the catalog's own unit, so that no coarser
// unit can hide a difference.
type WriteGroupHypertableEvidence struct {
	TimeColumn          string `json:"time_column"`
	ChunkIntervalMicros int64  `json:"chunk_interval_micros"`
}

// freezeHypertableEvidence turns the catalog's hypertable into the evidence a
// runtime layout carries, and only for exactly the confirmed layout: the
// partition column at precisely the confirmed interval. Anything else yields
// nil, so a differing hypertable is never recorded as a verified one.
func freezeHypertableEvidence(info *dbtarget.HypertableInfo) *WriteGroupHypertableEvidence {
	if info == nil || info.ChunkIntervalUnknown || info.TimeColumn != dbtarget.TimescalePartitionColumn ||
		info.ChunkInterval != dbtarget.TimescaleChunkInterval {
		return nil
	}
	return &WriteGroupHypertableEvidence{TimeColumn: info.TimeColumn, ChunkIntervalMicros: info.ChunkInterval.Microseconds()}
}

func (s *WriteGroupService) saveRuntimeVersionInTx(ctx context.Context, tx *sql.Tx, snap *WriteGroupAppliedSnapshot) error {
	if snap.RuntimeLayout == nil || (strings.EqualFold(strings.TrimSpace(snap.Group.WritePolicy.DedupeCapability), "receipt") && !snap.RuntimeLayout.ReceiptTableReady) {
		return nil
	}
	payload, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("encode applied runtime version: %w", err)
	}
	_, err = tx.ExecContext(ctx, s.repo.query(`INSERT INTO wg_runtime_versions
		(workspace_id, group_id, group_revision, payload) VALUES ($1, $2, $3, $4)
		ON CONFLICT (workspace_id, group_id, group_revision) DO NOTHING`),
		snap.WorkspaceID, snap.GroupID, snap.AppliedRevision, string(payload))
	if err != nil {
		return fmt.Errorf("save applied runtime version: %w", err)
	}
	return nil
}

// AppliedRuntimeTagTypes proves legacy recovery against the pinned source and
// mapping revisions before reading types. Draft members are never substituted.
func (s *WriteGroupService) AppliedRuntimeTagTypes(ctx context.Context, snap *WriteGroupAppliedSnapshot) (types map[string]schema.DataType, err error) {
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok || snap == nil || snap.Group == nil {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			types = nil
			err = errors.Join(err, fmt.Errorf("rollback applied runtime type read: %w", rollbackErr))
		}
	}()
	group := cloneWriteGroup(snap.Group)
	if err := s.repo.validate(ctx, setupTx.SQLTx(), group); err != nil {
		return nil, err
	}
	return loadWriteGroupTagTypes(ctx, setupTx.SQLTx(), s.repo.query, group)
}
