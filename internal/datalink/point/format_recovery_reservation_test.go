package point

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
	"path/filepath"
	"testing"
)

// This verifies the gateway's SQLite-only acceptance primitive without touching
// a row or firing UPDATE triggers. It is not a portable database lock claim.
func TestFormatRecoveryZeroRowReservationBlocksForeignWriter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reservation.db")
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE system_settings(key TEXT PRIMARY KEY,value TEXT,updated_at TEXT);INSERT INTO system_settings VALUES('scope','original','unchanged');CREATE TRIGGER any_update AFTER UPDATE ON system_settings BEGIN UPDATE system_settings SET updated_at='triggered' WHERE key=NEW.key;END`)
	require.NoError(t, err)
	foreign, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer foreign.Close()
	foreign.SetMaxOpenConns(1)
	_, err = foreign.Exec(`PRAGMA busy_timeout=0`)
	require.NoError(t, err)
	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	result, err := tx.Exec(`UPDATE system_settings SET value=value WHERE 0`)
	require.NoError(t, err)
	rows, err := result.RowsAffected()
	require.NoError(t, err)
	require.Zero(t, rows)
	_, err = foreign.Exec(`UPDATE system_settings SET value='foreign' WHERE key='scope'`)
	var busy *sqlite.Error
	require.ErrorAs(t, err, &busy)
	require.Equal(t, 5, busy.Code()&0xff)
	var stamp string
	require.NoError(t, tx.QueryRow(`SELECT updated_at FROM system_settings WHERE key='scope'`).Scan(&stamp))
	require.Equal(t, "unchanged", stamp)
	require.NoError(t, tx.Commit())
	_, err = foreign.Exec(`UPDATE system_settings SET value='foreign' WHERE key='scope'`)
	require.NoError(t, err)
}
