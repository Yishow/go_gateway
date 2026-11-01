package datalink

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationEmptyReadbacksPreserveNilSlices(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "empty.db"))
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()
	require.Nil(t, lineageSchemaSnapshot(t, db))
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		if cleanupErr := tx.Rollback(); cleanupErr != nil && !errors.Is(cleanupErr, sql.ErrTxDone) {
			t.Errorf("rollback empty fixture: %v", cleanupErr)
		}
	}()
	objects, err := mixedProfileObjects(t.Context(), tx)
	require.NoError(t, err)
	require.Nil(t, objects)
	objects, err = recordProfileKeyObjects(t.Context(), tx)
	require.NoError(t, err)
	require.Nil(t, objects)
	columns, err := mixedProfileColumns(t.Context(), tx)
	require.Error(t, err, "required historical columns must still reject an absent table")
	require.Nil(t, columns)
	require.NoError(t, tx.Rollback())
	_, err = db.ExecContext(t.Context(), "CREATE TABLE wg_delivery_mixed_outcomes(record_id TEXT)")
	require.NoError(t, err)
	require.Nil(t, mixedSourceIndexRows(t, db))
}
