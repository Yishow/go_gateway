package dbtarget

import (
	"database/sql"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

func TestPrepareSQLiteDeliveryDSNUsesAtomicModeAndPreservesOptions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "target.db")
	plain, err := prepareSQLiteDeliveryDSN(path + "?cache=shared")
	require.NoError(t, err)
	parsed, err := url.Parse(plain)
	require.NoError(t, err)
	require.Equal(t, path, parsed.Path)
	require.Equal(t, "rw", parsed.Query().Get("mode"))
	require.Equal(t, "shared", parsed.Query().Get("cache"))

	plainReadOnly, err := prepareSQLiteDeliveryDSN(path + "?mode=ro&cache=shared")
	require.NoError(t, err)
	parsed, err = url.Parse(plainReadOnly)
	require.NoError(t, err)
	require.Equal(t, path, parsed.Path)
	require.Equal(t, "ro", parsed.Query().Get("mode"))
	require.Empty(t, parsed.Query()["_pragma"])

	ro := (&url.URL{Scheme: "file", Path: path, RawQuery: "mode=ro&cache=shared"}).String()
	prepared, err := prepareSQLiteDeliveryDSN(ro)
	require.NoError(t, err)
	require.Equal(t, ro, prepared, "an explicit read-only mode must not be upgraded")

	rwc := (&url.URL{Scheme: "file", Path: path, RawQuery: "mode=rwc&cache=shared"}).String()
	prepared, err = prepareSQLiteDeliveryDSN(rwc)
	require.NoError(t, err)
	parsed, err = url.Parse(prepared)
	require.NoError(t, err)
	require.Equal(t, path, parsed.Path)
	require.Equal(t, "rw", parsed.Query().Get("mode"))
	require.Equal(t, "shared", parsed.Query().Get("cache"))
}

func TestPrepareSQLiteDeliveryDSNAddsBoundedBusyTimeoutWithoutOverridingExplicitPragma(t *testing.T) {
	path := filepath.Join(t.TempDir(), "target.db")

	prepared, err := prepareSQLiteDeliveryDSN(path)
	require.NoError(t, err)
	parsed, err := url.Parse(prepared)
	require.NoError(t, err)
	require.Contains(t, parsed.Query()["_pragma"], "busy_timeout(15000)")

	prepared, err = prepareSQLiteDeliveryDSN(path + "?_pragma=busy_timeout(0)&_pragma=journal_mode(delete)")
	require.NoError(t, err)
	parsed, err = url.Parse(prepared)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"busy_timeout(0)", "journal_mode(delete)"}, parsed.Query()["_pragma"])

	uri := (&url.URL{Scheme: "file", Path: path}).String()
	prepared, err = prepareSQLiteDeliveryDSN(uri)
	require.NoError(t, err)
	parsed, err = url.Parse(prepared)
	require.NoError(t, err)
	require.Contains(t, parsed.Query()["_pragma"], "busy_timeout(15000)")
}

func TestPrepareSQLiteDeliveryDSNDoesNotAliasMemoryFilenameOrUnknownMode(t *testing.T) {
	prepared, err := prepareSQLiteDeliveryDSN("memory")
	require.NoError(t, err)
	parsed, err := url.Parse(prepared)
	require.NoError(t, err)
	require.Equal(t, "memory", parsed.Opaque)
	require.Equal(t, "rw", parsed.Query().Get("mode"))
	prepared, err = prepareSQLiteDeliveryDSN("memory?cache=shared")
	require.NoError(t, err)
	parsed, err = url.Parse(prepared)
	require.NoError(t, err)
	require.Equal(t, "memory", parsed.Opaque)
	require.Equal(t, "rw", parsed.Query().Get("mode"))

	_, err = prepareSQLiteDeliveryDSN(filepath.Join(t.TempDir(), "unknown.db") + "?mode=unsupported")
	require.Error(t, err)
}

func TestProductionGroupOutageRecoverySQLiteDeliveryWritesEscapedPlainPath(t *testing.T) {
	ctx := t.Context()
	targetPath := filepath.Join(t.TempDir(), "target#name.db")
	createSQLiteDeliveryTarget(t, targetPath)
	service, connector := newSQLiteDeliveryService(t, targetPath)
	preparedDSN, err := prepareSQLiteDeliveryDSN(targetPath)
	require.NoError(t, err)
	parsed, err := url.Parse(preparedDSN)
	require.NoError(t, err)
	require.Equal(t, targetPath, parsed.Path)

	opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, opened.Close()) })
	_, err = opened.DB.ExecContext(ctx, `INSERT INTO sensor_values(value) VALUES (?)`, 7)
	require.NoError(t, err)

	var count int
	require.NoError(t, opened.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values`).Scan(&count))
	require.Equal(t, 1, count)
}

func TestProductionGroupOutageRecoverySQLiteDeliveryHonorsURIOptions(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		mode      string
		wantWrite bool
	}{
		{name: "read-only", mode: "ro"},
		{name: "read-write-create", mode: "rwc", wantWrite: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := t.Context()
			targetPath := filepath.Join(t.TempDir(), testCase.name+".db")
			createSQLiteDeliveryTarget(t, targetPath)
			dsn := (&url.URL{Scheme: "file", Path: targetPath, RawQuery: "mode=" + testCase.mode + "&cache=shared"}).String()
			service, connector := newSQLiteDeliveryService(t, dsn)

			opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, opened.Close()) })
			var count int
			require.NoError(t, opened.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values`).Scan(&count))
			_, err = opened.DB.ExecContext(ctx, `INSERT INTO sensor_values(value) VALUES (?)`, 7)
			if testCase.wantWrite {
				require.NoError(t, err)
			} else {
				require.Error(t, err, "an explicit read-only connector must remain read-only")
			}
		})
	}
}

func TestProductionGroupOutageRecoverySQLiteDeliveryHonorsPlainPathReadOnly(t *testing.T) {
	ctx := t.Context()
	targetPath := filepath.Join(t.TempDir(), "plain-readonly.db")
	createSQLiteDeliveryTarget(t, targetPath)
	service, connector := newSQLiteDeliveryService(t, targetPath+"?mode=ro&cache=shared")

	opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, opened.Close()) })
	var count int
	require.NoError(t, opened.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values`).Scan(&count))
	_, err = opened.DB.ExecContext(ctx, `INSERT INTO sensor_values(value) VALUES (?)`, 7)
	require.Error(t, err, "an explicit read-only connector must remain read-only")
}

func TestProductionGroupOutageRecoverySQLiteDeliveryRejectsUnknownMode(t *testing.T) {
	targetPath := filepath.Join(t.TempDir(), "unknown-mode.db")
	service, connector := newSQLiteDeliveryService(t, targetPath+"?mode=unsupported")
	_, statErr := os.Stat(targetPath)
	require.ErrorIs(t, statErr, os.ErrNotExist, "connector probing must not create an unknown-mode SQLite target")
	_, err := service.OpenDestination(t.Context(), connector.ID, connector.IdentityRevision)
	require.ErrorIs(t, err, ErrDestinationBlocked)
	_, statErr = os.Stat(targetPath)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestProductionGroupOutageRecoverySQLiteDeliveryPreservesMemoryDSN(t *testing.T) {
	for _, dsn := range []string{":memory:", "file::memory:?mode=memory&cache=shared"} {
		t.Run(dsn, func(t *testing.T) {
			ctx := t.Context()
			service, connector := newSQLiteDeliveryService(t, dsn)
			opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, opened.Close()) })
			require.NoError(t, opened.DB.PingContext(ctx))
			_, err = opened.DB.ExecContext(ctx, `CREATE TABLE sensor_values (value INTEGER)`)
			require.NoError(t, err)
		})
	}
}

func TestProductionGroupOutageRecoverySQLiteDeliveryWaitsForDeleteJournalWriter(t *testing.T) {
	ctx := t.Context()
	targetPath := filepath.Join(t.TempDir(), "delete-journal-lock.db")
	createSQLiteDeliveryTarget(t, targetPath)
	service, connector := newSQLiteDeliveryService(t, targetPath)

	holder, err := sql.Open("sqlite", targetPath+"?_pragma=busy_timeout(0)&_pragma=journal_mode(delete)")
	require.NoError(t, err)
	holder.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, holder.Close()) })
	transaction, err := holder.BeginTx(ctx, nil)
	require.NoError(t, err)
	_, err = transaction.ExecContext(ctx, `INSERT INTO sensor_values(value) VALUES (?)`, 11)
	require.NoError(t, err)

	opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, opened.Close()) })

	releaseErr := make(chan error, 1)
	go func() {
		time.Sleep(150 * time.Millisecond)
		releaseErr <- transaction.Rollback()
	}()
	startedAt := time.Now()
	_, err = InsertGroupRow(ctx, opened.DB, GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLite, TableName: "sensor_values", Strategy: GroupEffectNone,
		Row: EncodedRow{
			RecordID: "delivery-lock", EffectKey: "delivery-lock-effect", EntityKey: "delivery-lock", Cells: []EncodedCell{{Column: "value", Value: int64(22)}},
		},
	})
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			t.Logf("raw sqlite lock error: code=%d error=%v", sqliteErr.Code(), sqliteErr)
		}
	}
	waited := time.Since(startedAt)
	releaseResult := <-releaseErr
	require.NoError(t, err, "delivery should wait for the real writer lock to clear")
	require.GreaterOrEqual(t, waited, 100*time.Millisecond, "delivery must wait for the held writer lock")
	t.Logf("delivery waited %s for the held SQLite writer lock", waited)
	require.NoError(t, releaseResult)

	var count int
	require.NoError(t, opened.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values WHERE value = 22`).Scan(&count))
	require.Equal(t, 1, count)
}

func createSQLiteDeliveryTarget(t *testing.T, path string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.PingContext(t.Context()))
	_, err = db.ExecContext(t.Context(), `CREATE TABLE sensor_values (value INTEGER)`)
	require.NoError(t, err)
}

func newSQLiteDeliveryService(t *testing.T, dsn string) (*ConnectorService, *schema.DatabaseConnector) {
	t.Helper()
	mainDB := openMigratedTestDB(t)
	connectorRepo := NewSQLConnectorRepository(mainDB)
	service := NewConnectorService(connectorRepo, NewSQLTargetMappingRepository(mainDB))
	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name: "sqlite-target", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: ConnectionConfig{"dsn": dsn},
	})
	require.NoError(t, err)
	return service, connector
}
