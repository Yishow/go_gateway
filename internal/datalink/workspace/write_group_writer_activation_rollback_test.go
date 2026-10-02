package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/delivery"

	"github.com/stretchr/testify/require"
)

type writerActivationLegacyTarget struct {
	id, tagID, connectorID, tableSchema, tableName, columnName, writeMode string
	timestampColumn, groupKey                                             sql.NullString
	writeInterval                                                         sql.NullInt64
	enabled                                                               bool
	createdAt, updatedAt                                                  string
}

func readWriterActivationLegacyTarget(ctx context.Context, t *testing.T, db *sql.DB, id string) writerActivationLegacyTarget {
	t.Helper()
	var target writerActivationLegacyTarget
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT id, tag_id, connector_id, table_schema, table_name, column_name,
			write_mode, timestamp_column, group_key, write_interval_seconds, enabled,
			CAST(created_at AS TEXT), CAST(updated_at AS TEXT)
		FROM database_target_mappings WHERE id = ?
	`, id).Scan(
		&target.id, &target.tagID, &target.connectorID, &target.tableSchema, &target.tableName,
		&target.columnName, &target.writeMode, &target.timestampColumn, &target.groupKey,
		&target.writeInterval, &target.enabled, &target.createdAt, &target.updatedAt,
	))
	return target
}

func TestWriteGroupWriterOwnershipActivationBarrierWorkspaceSaveRollsBackOwnerAndAcceptedState(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newFixedWriterActivationGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	ownerKey := "test_writer_ownership_owner"
	legacyTargetID := "legacy-target-A"
	_, err := db.ExecContext(ctx, `
		INSERT INTO system_settings (key, value, updated_at) VALUES (?, ?, ?)
	`, ownerKey, "legacy", now)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO database_target_mappings (
			id, tag_id, connector_id, table_schema, table_name, column_name,
			write_mode, timestamp_column, group_key, write_interval_seconds, enabled
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, legacyTargetID, fixture.tagID, fixture.connectorID, "main", "raw_values", "temperature",
		"insert", "observed_at", nil, 15, true)
	require.NoError(t, err)
	targetBefore := readWriterActivationLegacyTarget(ctx, t, db, legacyTargetID)
	payload := []byte(`{"group_id":"group-A","group_revision":"group-revision-1","value":21.5}`)
	outbox := delivery.NewSQLOutbox(db)
	receipts := delivery.NewSQLReceiptLedger(db)
	require.NoError(t, outbox.Enqueue(&delivery.OutboxItem{
		ID: "activation-rollback-outbox", DestinationID: "group-A",
		DestinationRevision: "connector-revision-1", PlanRevision: "group-revision-1",
		RecordID: "accepted-record-1", CalculationRevision: 1, Table: "raw_values",
		Payload: payload, ObservedAt: now.Add(-time.Second),
	}))
	require.NoError(t, receipts.SaveReceipt(&delivery.Receipt{
		DestinationID: "group-A", RecordID: "accepted-record-1", CalculationRevision: 1,
		Table: "raw_values", DeliveredAt: now,
	}))
	var outboxBeforePayload, outboxBeforePlanRevision, outboxBeforeDestinationRevision, outboxBeforeRecordID string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT payload, plan_revision, destination_revision, record_id
		FROM gw_delivery_outbox WHERE id = ?
	`, "activation-rollback-outbox").Scan(
		&outboxBeforePayload, &outboxBeforePlanRevision, &outboxBeforeDestinationRevision, &outboxBeforeRecordID,
	))
	hasReceiptBefore, err := receipts.HasReceipt("group-A", "accepted-record-1", 1)
	require.NoError(t, err)
	require.True(t, hasReceiptBefore)

	beforeWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	beforeGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	var request WriterOwnershipActivationRequest
	var callbackCount int
	var ownerBefore, ownerAfter string
	var versionPayload string
	var versionEffectiveAt time.Time
	service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(
		ctx context.Context,
		tx *sql.Tx,
		activation WriterOwnershipActivationRequest,
	) error {
		callbackCount++
		request = activation
		if err := tx.QueryRowContext(ctx, `
			SELECT payload, effective_at FROM write_group_versions
			WHERE group_id = ? AND group_revision = ?
		`, activation.GroupID, activation.AppliedRevision).Scan(&versionPayload, &versionEffectiveAt); err != nil {
			return err
		}
		if err := tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, ownerKey).Scan(&ownerBefore); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ?
		`, "canonical", activation.EffectiveAt, ownerKey); err != nil {
			return err
		}
		return tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, ownerKey).Scan(&ownerAfter)
	}))

	_, err = db.ExecContext(ctx, `
		CREATE TRIGGER fail_writer_activation_workspace_save_exact
		BEFORE UPDATE ON system_settings
		WHEN NEW.key = 'studio_v2_workspace'
		BEGIN SELECT RAISE(ABORT, 'writer-activation-workspace-save-failed'); END
	`)
	require.NoError(t, err)
	_, err = service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.Error(t, err)

	require.Equal(t, 1, callbackCount)
	require.Equal(t, "legacy", ownerBefore)
	require.Equal(t, "canonical", ownerAfter)
	require.Equal(t, "group-A", request.GroupID)
	require.Equal(t, "group-revision-1", request.AppliedRevision)
	require.Equal(t, time.Date(2026, 10, 2, 12, 0, 15, 0, time.UTC), request.EffectiveAt)
	require.Equal(t, request.EffectiveAt, versionEffectiveAt.UTC())
	var version WriteGroup
	require.NoError(t, json.Unmarshal([]byte(versionPayload), &version))
	require.Equal(t, "group-A", version.ID)
	require.Equal(t, "group-revision-1", version.Revision)
	require.Equal(t, "group-revision-1", version.AppliedRevision)
	require.Equal(t, "connector-revision-1", version.Destination.ConnectorRevision)

	afterWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	afterGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace, afterWorkspace)
	require.Equal(t, beforeGroup, afterGroup)
	require.Empty(t, afterGroup.Group.AppliedRevision)
	var versionCount int
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, created.Group.ID,
	).Scan(&versionCount))
	require.Zero(t, versionCount)
	var ownerReloaded string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, ownerKey).Scan(&ownerReloaded))
	require.Equal(t, "legacy", ownerReloaded)
	require.Equal(t, targetBefore, readWriterActivationLegacyTarget(ctx, t, db, legacyTargetID))

	var persistedPayload, persistedPlanRevision, persistedDestinationRevision, persistedRecordID string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT payload, plan_revision, destination_revision, record_id
		FROM gw_delivery_outbox WHERE id = ?
	`, "activation-rollback-outbox").Scan(
		&persistedPayload, &persistedPlanRevision, &persistedDestinationRevision, &persistedRecordID,
	))
	require.Equal(t, outboxBeforePayload, persistedPayload)
	require.Equal(t, outboxBeforePlanRevision, persistedPlanRevision)
	require.Equal(t, outboxBeforeDestinationRevision, persistedDestinationRevision)
	require.Equal(t, outboxBeforeRecordID, persistedRecordID)
	hasReceipt, err := receipts.HasReceipt("group-A", "accepted-record-1", 1)
	require.NoError(t, err)
	require.Equal(t, hasReceiptBefore, hasReceipt)
}
