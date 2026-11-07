package workspace

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/delivery"

	"github.com/stretchr/testify/require"
)

func TestPermanentDeleteRefusesLegacyDeliveryBacklog(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	disabled, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	outbox := delivery.NewSQLOutbox(db)
	require.NoError(t, outbox.Enqueue(&delivery.OutboxItem{ID: "Q-legacy", DestinationID: group.ID,
		DestinationRevision: group.Destination.ConnectorRevision, PlanRevision: group.AppliedRevision,
		RecordID: "legacy-record", CalculationRevision: 1, Table: "raw_values", Payload: []byte(`{"sample":42}`), ObservedAt: time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)}))
	before := lifecycleReadState(t, db)
	_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, disabled.Group))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_BACKLOG")
	require.Equal(t, before, lifecycleReadState(t, db))
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM gw_delivery_outbox WHERE id='Q-legacy'`).Scan(&count))
	require.Equal(t, 1, count)
}

func TestPermanentDeleteRemovesDeliveredLegacyRecords(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	disabled, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	outbox := delivery.NewSQLOutbox(db)
	require.NoError(t, outbox.Enqueue(&delivery.OutboxItem{ID: "Q-legacy", DestinationID: group.ID,
		DestinationRevision: group.Destination.ConnectorRevision, PlanRevision: group.AppliedRevision,
		RecordID: "legacy-record", CalculationRevision: 1, Table: "raw_values", Payload: []byte(`{"sample":42}`), ObservedAt: time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)}))
	_, err = db.ExecContext(t.Context(), `UPDATE gw_delivery_outbox SET status='delivered' WHERE id='Q-legacy'`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO gw_delivery_receipts(destination_id,record_id,calculation_revision,table_name,delivered_at) VALUES(?,'legacy-record',1,'raw_values',CURRENT_TIMESTAMP)`, group.ID)
	require.NoError(t, err)
	_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, disabled.Group))
	require.NoError(t, err)
	for _, table := range []string{"gw_delivery_outbox", "gw_delivery_receipts"} {
		var count int
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table+` WHERE destination_id=?`, group.ID).Scan(&count))
		require.Zero(t, count, table)
	}
}

func TestPermanentDeleteRefusesUnfinishedSuccessorOperation(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	disabled, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	detail := `{"version":"successor-operation/v1","scope":{"workspace_id":"` + group.WorkspaceID + `","lineage_id":"L1","predecessor_group_id":"G1","successor_group_id":"G2"}}`
	_, err = db.ExecContext(t.Context(), `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES('O2','T2',?,'successor-scope','owner','successor_first_start','running',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`, group.WorkspaceID, detail)
	require.NoError(t, err)
	before := lifecycleReadState(t, db)
	_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, disabled.Group))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_OPERATION_PENDING")
	require.Equal(t, before, lifecycleReadState(t, db))
}
