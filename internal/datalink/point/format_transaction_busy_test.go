package point

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestPointMutationCommitBusyLeavesNoNativeTransaction(t *testing.T) {
	for _, operation := range []string{"format", "ordinary"} {
		t.Run(operation, func(t *testing.T) {
			memory, _, req := formatFixture(t)
			path := filepath.Join(t.TempDir(), "busy.db")
			_, err := memory.Exec(`VACUUM INTO ?`, path)
			require.NoError(t, err)
			db, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			defer db.Close()
			db.SetMaxOpenConns(1)
			_, err = db.Exec(`PRAGMA journal_mode=DELETE; PRAGMA busy_timeout=0`)
			require.NoError(t, err)
			reader, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			defer reader.Close()
			readTx, err := reader.Begin()
			require.NoError(t, err)
			defer readTx.Rollback()
			var count int
			require.NoError(t, readTx.QueryRow(`SELECT count(*) FROM points`).Scan(&count))
			repo := NewSQLRepository(db)
			before, err := repo.GetByID(t.Context(), req.PointID)
			require.NoError(t, err)
			if operation == "format" {
				_, err = repo.MutateFormat(t.Context(), req)
			} else {
				after := *before
				after.RevisionID = "attempt"
				after.Name = "not committed"
				err = repo.UpdateRevision(t.Context(), &after, before.RevisionID)
			}
			require.Error(t, err)
			stored, err := repo.GetByID(t.Context(), req.PointID)
			require.NoError(t, err)
			require.Equal(t, before.RevisionID, stored.RevisionID)
			require.Equal(t, before.Name, stored.Name)
			require.NoError(t, readTx.Rollback())
			clean, err := db.BeginTx(t.Context(), nil)
			require.NoError(t, err, "failed commit must not poison next pool borrower")
			require.NoError(t, clean.Rollback())
			var receipts int
			require.NoError(t, db.QueryRow(`SELECT count(*) FROM system_settings WHERE key LIKE 'point_format_receipt:%'`).Scan(&receipts))
			require.Zero(t, receipts)
		})
	}
}

func TestFormatCompletionBusyLeavesPendingReceiptAndCleanPool(t *testing.T) {
	memory, _, req := formatFixture(t)
	path := filepath.Join(t.TempDir(), "finalize-busy.db")
	_, err := memory.Exec(`VACUUM INTO ?`, path)
	require.NoError(t, err)
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`PRAGMA journal_mode=DELETE; PRAGMA busy_timeout=0`)
	require.NoError(t, err)
	reader, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer reader.Close()
	var held *sql.Tx
	spy := &recordingFormatInvalidator{run: func() error {
		var err error
		held, err = reader.Begin()
		if err != nil {
			return err
		}
		var n int
		return held.QueryRow(`SELECT count(*) FROM system_settings`).Scan(&n)
	}}
	svc := NewService(NewSQLRepository(db), nil)
	svc.SetFormatInvalidator(spy)
	result, err := svc.MutateFormat(t.Context(), req)
	require.Nil(t, result)
	require.Error(t, err)
	require.Equal(t, 1, spy.calls)
	require.NotNil(t, held)
	defer held.Rollback()
	receipt, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
	require.NoError(t, err)
	require.Equal(t, "pending", receipt.InvalidationStatus)
	require.NoError(t, held.Rollback())
	next, err := db.Begin()
	require.NoError(t, err)
	require.NoError(t, next.Rollback())
}
