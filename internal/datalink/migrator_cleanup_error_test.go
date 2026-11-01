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

type migrationCleanupSecretError struct{}

func (*migrationCleanupSecretError) Error() string {
	return "secret-token /private/db.sqlite SELECT private_column"
}

var errMigrationPrimaryFixture = errors.New("migration primary fixture failure")

type migrationRestoreFault struct {
	mode    string
	touched bool
	closes  int
}

type migrationCleanupConnector struct {
	restore       *migrationRestoreFault
	path          string
	failure       bool
	cleanup       *migrationCleanupSecretError
	rollbackCalls *int
}

func (c migrationCleanupConnector) Connect(context.Context) (driver.Conn, error) {
	conn, err := (&sqlite.Driver{}).Open(c.path)
	if err != nil {
		return nil, err
	}
	return &migrationCleanupConn{Conn: conn, fixture: c}, nil
}
func (migrationCleanupConnector) Driver() driver.Driver { return &sqlite.Driver{} }

type migrationCleanupConn struct {
	driver.Conn
	fixture migrationCleanupConnector
}

func (c *migrationCleanupConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return migrationCleanupTx{Tx: tx, fixture: c.fixture}, nil
}
func (c *migrationCleanupConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if fault := c.fixture.restore; fault != nil {
		if strings.HasPrefix(query, "CREATE TEMP TABLE") {
			fault.touched = true
		}
		if fault.touched && fault.mode == "confirmed" && strings.HasPrefix(query, "PRAGMA legacy_alter_table=") {
			return nil, errMigrationPrimaryFixture
		}
	}
	if c.fixture.failure {
		return nil, errMigrationPrimaryFixture
	}
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}
func (c *migrationCleanupConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if fault := c.fixture.restore; fault != nil && fault.touched && fault.mode == "unconfirmed" && query == "SELECT COUNT(*) FROM temp.sqlite_master" {
		return nil, driver.ErrBadConn
	}
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}

type migrationCleanupTx struct {
	driver.Tx
	fixture migrationCleanupConnector
}

func (tx migrationCleanupTx) Rollback() error {
	*tx.fixture.rollbackCalls++
	if err := tx.Tx.Rollback(); err != nil {
		return err
	}
	if tx.fixture.failure {
		return tx.fixture.cleanup
	}
	return nil
}

func TestMigrationRollbackCleanupPreservesPrimaryAndRedactsCause(t *testing.T) {
	for _, tc := range []struct {
		name    string
		migrate func(*sql.DB) error
	}{
		{"point-format", ensureSQLitePointFormatIntent},
		{"sample-handoffs", ensureSQLiteSampleHandoffs},
		{"lineages", ensureSQLiteWriteGroupLineages},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "migration.db")
			secret := &migrationCleanupSecretError{}
			rollbacks := 0
			db := sql.OpenDB(migrationCleanupConnector{path: path, failure: true, cleanup: secret, rollbackCalls: &rollbacks})
			defer func() { require.NoError(t, db.Close()) }()
			err := tc.migrate(db)
			require.ErrorIs(t, err, errMigrationPrimaryFixture)
			require.Equal(t, 1, rollbacks)
			require.ErrorIs(t, err, secret)
			var typed *migrationCleanupSecretError
			require.ErrorAs(t, err, &typed)
			require.NotContains(t, err.Error(), "secret-token")
			require.NotContains(t, err.Error(), "/private/")
			require.NotContains(t, err.Error(), "SELECT private_column")
		})
	}
}

func TestMigrationCommittedCleanupDoesNotReportTxDone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "migration.db")
	seed, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	require.NoError(t, NewMigrator().Migrate(seed))
	require.NoError(t, seed.Close())
	rollbacks := 0
	db := sql.OpenDB(migrationCleanupConnector{path: path, rollbackCalls: &rollbacks})
	defer func() { require.NoError(t, db.Close()) }()
	require.NoError(t, ensureSQLitePointFormatIntent(db))
	require.NoError(t, ensureSQLiteSampleHandoffs(db))
	require.NoError(t, ensureSQLiteWriteGroupLineages(db))
	require.Zero(t, rollbacks, "database/sql already completed each committed transaction")
}

func TestMigrationPragmaDiscardPreservesFailureAndSafeCause(t *testing.T) {
	for _, mode := range []string{"nil", "bad-connection", "unexpected", "restore-success"} {
		t.Run(mode, func(t *testing.T) {
			primary := errors.New("restore primary fixture failure")
			secret := &migrationCleanupSecretError{}
			discarded := 0
			err := restoreSQLiteMigrationPragmas([]sqliteMigrationPragma{{"foreign_keys", 1}}, func(string) error {
				if mode == "restore-success" {
					return nil
				}
				return primary
			}, func(string) (int, error) { return 1, nil }, func() error {
				discarded++
				switch mode {
				case "bad-connection":
					return driver.ErrBadConn
				case "unexpected":
					return secret
				default:
					return nil
				}
			})
			if mode == "restore-success" {
				require.NoError(t, err)
				require.Zero(t, discarded)
				return
			}
			require.ErrorIs(t, err, primary)
			require.Equal(t, 1, discarded)
			require.NotContains(t, err.Error(), "secret-token")
			if mode == "unexpected" {
				require.ErrorIs(t, err, secret)
				var typed *migrationCleanupSecretError
				require.ErrorAs(t, err, &typed)
			} else {
				require.NotErrorIs(t, err, driver.ErrBadConn)
			}
		})
	}
}

func (c *migrationCleanupConn) Close() error {
	if c.fixture.restore != nil {
		c.fixture.restore.closes++
	}
	return c.Conn.Close()
}

func TestMigrationRepeatedDiscardRequiresPriorConfirmation(t *testing.T) {
	for _, mode := range []string{"confirmed", "unconfirmed"} {
		t.Run(mode, func(t *testing.T) {
			fault := &migrationRestoreFault{mode: mode}
			rollbacks := 0
			db := sql.OpenDB(migrationCleanupConnector{path: filepath.Join(t.TempDir(), "restore.db"), rollbackCalls: &rollbacks, restore: fault})
			db.SetMaxOpenConns(1)
			defer func() { require.NoError(t, db.Close()) }()
			_, err := db.ExecContext(t.Context(), "CREATE TABLE write_groups(id TEXT)")
			require.NoError(t, err)
			err = runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
				_, execErr := tx.ExecContext(ctx, "CREATE TEMP TABLE cleanup_fixture(value TEXT)")
				return execErr
			})
			require.Error(t, err)
			require.True(t, fault.touched)
			require.Equal(t, 1, fault.closes, "only the original connection may be discarded")
			if mode == "confirmed" {
				require.ErrorIs(t, err, errMigrationPrimaryFixture)
				require.NotErrorIs(t, err, sql.ErrConnDone, "confirmed first discard makes the second closed-connection response idempotent")
				require.NotContains(t, err.Error(), "migration cleanup failed")
			} else {
				require.ErrorIs(t, err, sql.ErrConnDone)
				require.Contains(t, err.Error(), "migration cleanup failed", "unconfirmed unavailable connection is not proven discarded")
			}
		})
	}
}
