package main

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDestinationRowsWaitsForInFlightSQLiteCommit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "readings.db")
	holder, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	holder.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, holder.Close()) })
	_, err = holder.ExecContext(t.Context(), `CREATE TABLE readings (temperature REAL, pressure INTEGER); INSERT INTO readings VALUES (21.5,100)`)
	require.NoError(t, err)
	connection, err := holder.Conn(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, connection.Close()) })
	_, err = connection.ExecContext(t.Context(), `BEGIN EXCLUSIVE`)
	require.NoError(t, err)
	released := make(chan error, 1)
	t.Cleanup(func() { require.NoError(t, <-released) })
	go func() {
		time.Sleep(150 * time.Millisecond)
		_, releaseErr := connection.ExecContext(context.WithoutCancel(t.Context()), `ROLLBACK`)
		released <- releaseErr
	}()
	// A real SQL observer must wait for the in-flight writer, then verify the
	// unchanged known values. A lock error must never become fabricated data.
	require.Equal(t, [][2]float64{{21.5, 100}}, destinationRows(t, path))
}
