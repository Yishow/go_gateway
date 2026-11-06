package datalink

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func openDeviceAcquisitionDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	return db
}

func seedDeviceAcquisition(t *testing.T, db *sql.DB, stmts ...string) {
	t.Helper()
	for _, stmt := range stmts {
		_, err := db.ExecContext(t.Context(), stmt)
		require.NoError(t, err, stmt)
	}
}

func pointGroupInterval(t *testing.T, db *sql.DB, pointID string) (groupID string, owner sql.NullString, interval int) {
	t.Helper()
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT g.id, g.owner_device_id, g.interval_ms FROM points p JOIN polling_groups g ON g.id = p.polling_group_id WHERE p.id = ?`, pointID).Scan(&groupID, &owner, &interval))
	return groupID, owner, interval
}

func ownedGroupCount(t *testing.T, db *sql.DB, deviceID string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM polling_groups WHERE owner_device_id = ?`, deviceID).Scan(&n))
	return n
}

const deviceAcquisitionSeedDevices = `INSERT INTO devices (id, name, protocol) VALUES ('dev-a','A','modbus_tcp'),('dev-b','B','mc_3e'),('dev-empty','E','fatek_fbs')`

func TestDeviceAcquisitionMigrationAddsColumnsWithSafeDefaults(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	require.NoError(t, NewMigrator().Migrate(db))
	seedDeviceAcquisition(t, db, deviceAcquisitionSeedDevices)
	var stale sql.NullInt64
	var frozen, auto int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT stale_threshold_ms, frozen_threshold_ms, acquisition_auto_interval FROM devices WHERE id='dev-a'`).Scan(&stale, &frozen, &auto))
	require.False(t, stale.Valid, "NULL stale threshold means interval x5, at least 10s")
	require.Zero(t, frozen, "frozen detection defaults to disabled")
	require.Zero(t, auto)
	for _, bad := range []string{
		`UPDATE devices SET stale_threshold_ms = 0 WHERE id='dev-a'`,
		`UPDATE devices SET stale_threshold_ms = 'abc' WHERE id='dev-a'`,
		`UPDATE devices SET frozen_threshold_ms = -1 WHERE id='dev-a'`,
		`UPDATE devices SET frozen_threshold_ms = NULL WHERE id='dev-a'`,
		`UPDATE devices SET acquisition_auto_interval = 2 WHERE id='dev-a'`,
		`UPDATE devices SET acquisition_auto_interval = 'true' WHERE id='dev-a'`,
		`INSERT INTO polling_groups (name, interval_ms, owner_device_id) VALUES ('x', 1000, '')`,
		`INSERT INTO polling_groups (name, interval_ms, owner_device_id) VALUES ('x', 1000, 'missing-device')`,
	} {
		_, err := db.ExecContext(t.Context(), `PRAGMA foreign_keys=ON; `+bad)
		require.Error(t, err, bad)
	}
	seedDeviceAcquisition(t, db, `INSERT INTO polling_groups (name, interval_ms, owner_device_id) VALUES ('dup', 1000, 'dev-a')`)
	_, err := db.ExecContext(t.Context(), `INSERT INTO polling_groups (name, interval_ms, owner_device_id) VALUES ('dup2', 1000, 'dev-a')`)
	require.Error(t, err, "a device owns at most one polling group")
}

func TestDeviceAcquisitionBackfillKeepsIntervalsAndIsIdempotent(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	seedDeviceAcquisition(t, db,
		deviceAcquisitionSeedDevices,
		`INSERT INTO polling_groups (id, name, interval_ms, priority, enabled) VALUES ('auto','Auto Polling',1000,50,1),('slow','Slow',2500,70,1),('off','Off',300,10,0)`,
		`INSERT INTO points (id, device_id, name, address, data_type, polling_group_id) VALUES
			('a1','dev-a','a1','40001','uint16','auto'),('a2','dev-a','a2','40002','uint16','auto'),
			('a3','dev-a','a3','40003','uint16',NULL),
			('b1','dev-b','b1','D0','int16','slow'),('b2','dev-b','b2','D1','int16','off')`,
	)
	require.NoError(t, ensureSQLiteDeviceAcquisition(db))

	aGroup, aOwner, aInterval := pointGroupInterval(t, db, "a1")
	require.Equal(t, "dev-a", aOwner.String)
	require.Equal(t, 1000, aInterval, "upgrade keeps device A at the interval its points used")
	a2Group, _, _ := pointGroupInterval(t, db, "a2")
	require.Equal(t, aGroup, a2Group)
	var a3Group sql.NullString
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT polling_group_id FROM points WHERE id='a3'`).Scan(&a3Group))
	require.False(t, a3Group.Valid, "unassigned points are not given a cadence by the migration")

	_, bOwner, bInterval := pointGroupInterval(t, db, "b1")
	require.Equal(t, "dev-b", bOwner.String)
	require.Equal(t, 2500, bInterval)
	b2Group, _, _ := pointGroupInterval(t, db, "b2")
	require.Equal(t, "off", b2Group, "points in disabled groups must not start polling")

	require.Equal(t, 1, ownedGroupCount(t, db, "dev-empty"))
	var emptyInterval int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT interval_ms FROM polling_groups WHERE owner_device_id='dev-empty'`).Scan(&emptyInterval))
	require.Equal(t, 1000, emptyInterval, "device without points inherits the default group's interval")

	var autoInterval, autoPoints int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT interval_ms FROM polling_groups WHERE id='auto'`).Scan(&autoInterval))
	require.Equal(t, 1000, autoInterval, "Auto Polling stays as fallback")
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM points WHERE polling_group_id='auto'`).Scan(&autoPoints))
	require.Zero(t, autoPoints)

	var before int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM polling_groups`).Scan(&before))
	require.NoError(t, ensureSQLiteDeviceAcquisition(db))
	require.NoError(t, NewMigrator().Migrate(db))
	var after int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM polling_groups`).Scan(&after))
	require.Equal(t, before, after, "re-running must not create duplicate groups")
	for _, id := range []string{"dev-a", "dev-b", "dev-empty"} {
		require.Equal(t, 1, ownedGroupCount(t, db, id))
	}
	again, _, _ := pointGroupInterval(t, db, "a1")
	require.Equal(t, aGroup, again)
}

func TestDeviceAcquisitionBackfillUsesShortestEnabledInterval(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	seedDeviceAcquisition(t, db,
		deviceAcquisitionSeedDevices,
		`INSERT INTO polling_groups (id, name, interval_ms, enabled) VALUES ('fast','Fast',500,1),('slow','Slow',2000,1),('off','Off',100,0)`,
		`INSERT INTO points (id, device_id, name, address, data_type, polling_group_id, revision_id) VALUES
			('a1','dev-a','a1','40001','uint16','fast','a1-rev'),('a2','dev-a','a2','40002','uint16','slow','a2-rev'),
			('a-off','dev-a','a-off','40003','uint16','off','a-off-rev'),('a-none','dev-a','a-none','40004','uint16',NULL,'a-none-rev'),
			('b1','dev-b','b1','D0','int16','slow','b1-rev'),
			('b-off','dev-b','b-off','D1','int16','off','b-off-rev'),('b-none','dev-b','b-none','D2','int16',NULL,'b-none-rev')`,
		`UPDATE points SET enabled=0 WHERE id='a2'`,
	)
	type unchangedPoint struct {
		group     sql.NullString
		name      string
		revision  string
		updatedAt string
	}
	readUnchanged := func(id string) unchangedPoint {
		t.Helper()
		var state unchangedPoint
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT polling_group_id, name, revision_id, updated_at FROM points WHERE id = ?`, id).Scan(&state.group, &state.name, &state.revision, &state.updatedAt))
		return state
	}
	before := make(map[string]unchangedPoint)
	for _, id := range []string{"a-off", "a-none", "b-off", "b-none"} {
		before[id] = readUnchanged(id)
	}

	for range 2 {
		require.NoError(t, NewMigrator().Migrate(db))
		a1, aOwner, a1Interval := pointGroupInterval(t, db, "a1")
		a2, _, a2Interval := pointGroupInterval(t, db, "a2")
		require.Equal(t, "dev-a", aOwner.String)
		require.Equal(t, 500, a1Interval)
		require.Equal(t, a1, a2, "all enabled-group points use the same owned group")
		require.Equal(t, 500, a2Interval, "the longer prior interval adopts the shortest interval")
		var a2Enabled bool
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT enabled FROM points WHERE id='a2'`).Scan(&a2Enabled))
		require.False(t, a2Enabled, "moving a disabled point must not enable it")
		_, bOwner, bInterval := pointGroupInterval(t, db, "b1")
		require.Equal(t, "dev-b", bOwner.String)
		require.Equal(t, 2000, bInterval, "another device retains its interval")
		for id, original := range before {
			require.Equal(t, original, readUnchanged(id), "disabled or unassigned point %s must remain unchanged", id)
		}
		for _, id := range []string{"dev-a", "dev-b", "dev-empty"} {
			require.Equal(t, 1, ownedGroupCount(t, db, id), "migration remains idempotent")
		}
	}
}

func TestDeviceAcquisitionBackfillResolvesNameCollisionWithSuffix(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	seedDeviceAcquisition(t, db,
		deviceAcquisitionSeedDevices,
		`INSERT INTO polling_groups (id, name, interval_ms, enabled) VALUES
			('user-a','device:dev-a',500,1),
			('user-b1','device:dev-b',1000,1),
			('user-b2','device:dev-b#2',1000,1)`,
		`INSERT INTO points (id, device_id, name, address, data_type, polling_group_id) VALUES
			('a1','dev-a','a1','40001','uint16','user-a'),
			('b1','dev-b','b1','D0','int16','user-b1')`,
	)
	require.NoError(t, ensureSQLiteDeviceAcquisition(db))
	var nameA, nameB string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT name FROM polling_groups WHERE owner_device_id='dev-a'`).Scan(&nameA))
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT name FROM polling_groups WHERE owner_device_id='dev-b'`).Scan(&nameB))
	require.Equal(t, "device:dev-a#2", nameA)
	require.Equal(t, "device:dev-b#3", nameB)
	a1Group, a1Owner, a1Interval := pointGroupInterval(t, db, "a1")
	require.NotEqual(t, "user-a", a1Group)
	require.Equal(t, []any{sql.NullString{String: "dev-a", Valid: true}, 500}, []any{a1Owner, a1Interval})
}

func TestDeviceAcquisitionMigrationReAddsDroppedDeviceColumns(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	seedDeviceAcquisition(t, db,
		`ALTER TABLE devices DROP COLUMN stale_threshold_ms`,
		`ALTER TABLE devices DROP COLUMN frozen_threshold_ms`,
		`ALTER TABLE devices DROP COLUMN acquisition_auto_interval`,
	)
	require.NoError(t, ensureSQLiteDeviceAcquisition(db))
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pragma_table_info('devices') WHERE name IN ('stale_threshold_ms','frozen_threshold_ms','acquisition_auto_interval')`).Scan(&n))
	require.Equal(t, 3, n)
}

func TestDeviceAcquisitionMigrationRejectsUnknownColumnShape(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	seedDeviceAcquisition(t, db,
		`ALTER TABLE devices DROP COLUMN frozen_threshold_ms`,
		`ALTER TABLE devices ADD COLUMN frozen_threshold_ms TEXT`,
	)
	require.ErrorContains(t, ensureSQLiteDeviceAcquisition(db), "unknown devices.frozen_threshold_ms column")
}
