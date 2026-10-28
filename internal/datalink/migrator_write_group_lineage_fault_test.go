package datalink

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

func TestWriteGroupLineagesMigrationAtomicFailureAndLostCommitAcknowledgement(t *testing.T) {
	for _, mode := range []string{"after-ddl", "after-backfill", "lost-commit-ack"} {
		t.Run(mode, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "fault.db")
			db := openLineageMigrationDB(t, path)
			createPre047LineageSchema(t, db)
			seedLegacyLineageRows(t, db)
			columns := lineageSnapshotColumns(t, db, lineagePreservedTables...)
			before := lineageMigrationSnapshot(t, db, columns)
			schema := lineageSchemaSnapshot(t, db)
			require.NoError(t, db.Close())
			fault := &lineageMigrationFault{mode: mode}
			wrapped := sql.OpenDB(lineageMigrationConnector{path: path, fault: fault})
			wrapped.SetMaxOpenConns(1)
			t.Cleanup(func() { require.NoError(t, wrapped.Close()) })
			require.ErrorIs(t, ensureSQLiteWriteGroupLineages(wrapped), errLineageMigrationInjected)
			require.True(t, fault.fired, "fault must intercept the real migration")
			require.True(t, fault.ddl, "the real 047 DDL must have executed")
			require.Equal(t, before, lineageMigrationSnapshot(t, wrapped, columns))
			var metadata map[string][]string
			var metaColumns map[string][]string
			if mode == "lost-commit-ack" {
				assertInitialLineageRows(t, wrapped)
				metaColumns = lineageSnapshotColumns(t, wrapped, "write_group_lineages", "write_group_generations", "write_group_basic_keys", "write_group_profile_keys")
				metadata = lineageMigrationSnapshot(t, wrapped, metaColumns)
			} else {
				require.Equal(t, schema, lineageSchemaSnapshot(t, wrapped), "rollback must include DDL, references and metadata")
			}
			if mode != "after-ddl" {
				require.Equal(t, 3, fault.lineages)
				require.Equal(t, 3, fault.generations)
				require.Equal(t, 6, fault.bindings, "both key tables were bound for each distinct group")
			}
			require.NoError(t, wrapped.Close())
			db = openLineageMigrationDB(t, path)
			require.NoError(t, NewMigrator().Migrate(db))
			assertInitialLineageRows(t, db)
			require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
			if metadata != nil {
				require.Equal(t, metadata, lineageMigrationSnapshot(t, db, metaColumns), "lost acknowledgement must not duplicate or reset identities")
			}
			require.NoError(t, NewMigrator().Migrate(db))
			assertInitialLineageRows(t, db)
		})
	}
}

func TestWriteGroupLineagesMigrationRechecksArtifactsUnderWriteLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "raced.db")
	db := openLineageMigrationDB(t, path)
	createPre047LineageSchema(t, db)
	seedLegacyLineageRows(t, db)
	require.NoError(t, db.Close())
	fault := &lineageMigrationFault{mode: "schema-race"}
	wrapped := sql.OpenDB(lineageMigrationConnector{path: path, fault: fault})
	wrapped.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, wrapped.Close()) })
	// The driver inserts an artifact after preflight but before the locked
	// schema inspection. Only the transaction's second classification sees it.
	require.ErrorContains(t, ensureSQLiteWriteGroupsTables(wrapped), "lineage migration is incomplete")
	require.True(t, fault.fired)
	initialized, err := sqliteWriteGroupLineageState(t.Context(), wrapped)
	require.NoError(t, err)
	require.False(t, initialized, "the injected transactional marker is rolled back too")
	require.NoError(t, NewMigrator().Migrate(wrapped))
	assertInitialLineageRows(t, wrapped)
}

var errLineageMigrationInjected = errors.New("injected lineage migration failure")

type lineageMigrationFault struct {
	mode                            string
	fired, ddl                      bool
	lineages, generations, bindings int
}

type lineageMigrationConnector struct {
	path  string
	fault *lineageMigrationFault
}

func (c lineageMigrationConnector) Connect(context.Context) (driver.Conn, error) {
	conn, err := (&sqlite.Driver{}).Open(c.path)
	if err != nil {
		return nil, err
	}
	return &lineageMigrationConn{Conn: conn, fault: c.fault}, nil
}

func (lineageMigrationConnector) Driver() driver.Driver { return &sqlite.Driver{} }

type lineageMigrationConn struct {
	driver.Conn
	fault *lineageMigrationFault
}

func (c *lineageMigrationConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer := c.Conn.(driver.ExecerContext)
	if c.fault.mode == "schema-race" && !c.fault.fired && strings.Contains(query, "UPDATE write_groups SET revision=revision WHERE 0") {
		c.fault.fired = true
		if _, err := execer.ExecContext(ctx, `ALTER TABLE write_group_basic_keys ADD COLUMN lineage_id TEXT`, nil); err != nil {
			return nil, err
		}
	}
	if c.fault.mode == "after-ddl" && !c.fault.fired && strings.HasPrefix(query, "INSERT INTO write_group_lineages") {
		c.fault.fired = true
		return nil, errLineageMigrationInjected
	}
	result, err := execer.ExecContext(ctx, query, args)
	if err == nil {
		switch {
		case strings.Contains(query, "CREATE TABLE write_group_lineages"):
			c.fault.ddl = true
		case strings.HasPrefix(query, "INSERT INTO write_group_lineages"):
			c.fault.lineages++
		case strings.HasPrefix(query, "INSERT INTO write_group_generations"):
			c.fault.generations++
		case strings.HasPrefix(query, "UPDATE write_group_") && strings.Contains(query, " SET lineage_id="):
			c.fault.bindings++
		}
	}
	return result, err
}

func (c *lineageMigrationConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.fault.mode == "after-backfill" && !c.fault.fired && c.fault.bindings > 0 && strings.HasPrefix(query, "SELECT EXISTS(") {
		c.fault.fired = true
		return nil, errLineageMigrationInjected
	}
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}

func (c *lineageMigrationConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return lineageMigrationTx{Tx: tx, fault: c.fault}, nil
}

type lineageMigrationTx struct {
	driver.Tx
	fault *lineageMigrationFault
}

func (tx lineageMigrationTx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		return err
	}
	if tx.fault.mode == "lost-commit-ack" && tx.fault.ddl && !tx.fault.fired {
		tx.fault.fired = true
		return errLineageMigrationInjected
	}
	return nil
}
