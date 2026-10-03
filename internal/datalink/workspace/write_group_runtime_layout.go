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
