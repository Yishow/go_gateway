package dbtarget

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

func TestProductionGroupOutageRecoverySQLiteDeliveryWaitsForDeleteJournalReaderAtCommit(t *testing.T) {
	ctx := t.Context()
	targetPath := filepath.Join(t.TempDir(), "delete-journal-reader.db")
	createSQLiteDeliveryTarget(t, targetPath)
	service, connector := newSQLiteDeliveryService(t, targetPath)

	reader, err := sql.Open("sqlite", targetPath+"?_pragma=busy_timeout(0)&_pragma=journal_mode(delete)")
	require.NoError(t, err)
	reader.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, reader.Close()) })
	readerTx, err := reader.BeginTx(ctx, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = readerTx.Rollback() })
	var existing int
	require.NoError(t, readerTx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values`).Scan(&existing))
	require.Equal(t, 0, existing)

	opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, opened.Close()) })

	releaseErr := make(chan error, 1)
	go func() {
		time.Sleep(150 * time.Millisecond)
		releaseErr <- readerTx.Rollback()
	}()
	startedAt := time.Now()
	_, err = InsertGroupRow(ctx, opened.DB, GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLite, TableName: "sensor_values", Strategy: GroupEffectNone,
		Row: EncodedRow{
			RecordID: "delivery-reader-lock", EffectKey: "delivery-reader-lock-effect", EntityKey: "delivery-reader-lock", Cells: []EncodedCell{{Column: "value", Value: int64(33)}},
		},
	})
	if err != nil {
		var insertErr *GroupInsertError
		if errors.As(err, &insertErr) {
			t.Logf("group insert phase: %s", insertErr.Phase)
		}
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			t.Logf("raw sqlite lock error: code=%d error=%v", sqliteErr.Code(), sqliteErr)
		}
	}
	waited := time.Since(startedAt)
	releaseResult := <-releaseErr
	require.NoError(t, err, "delivery should wait for the reader lock at commit")
	require.GreaterOrEqual(t, waited, 100*time.Millisecond, "delivery must wait for the held reader lock")
	t.Logf("delivery waited %s for the held SQLite reader lock at commit", waited)
	require.NoError(t, releaseResult)

	var count int
	require.NoError(t, opened.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values WHERE value = 33`).Scan(&count))
	require.Equal(t, 1, count)
}
