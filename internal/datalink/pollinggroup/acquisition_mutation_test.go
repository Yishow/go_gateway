package pollinggroup

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/point"

	"github.com/stretchr/testify/require"
)

type acquisitionDeleteLifecycle struct {
	t               *testing.T
	deviceIDs       []string
	begun, finished int
	active          bool
	begin           func() error
	onFinish        func()
}

func (l *acquisitionDeleteLifecycle) TryBeginAcquisitionConfigurationMutations(_ context.Context, devices []string) (func(), error) {
	l.begun++
	l.active = true
	l.deviceIDs = append([]string(nil), devices...)
	var err error
	if l.begin != nil {
		err = l.begin()
	}
	return func() {
		require.True(l.t, l.active)
		if l.onFinish != nil {
			l.onFinish()
		}
		l.active = false
		l.finished++
	}, err
}

func acquisitionDeleteFixture(t *testing.T) (*sql.DB, *SQLRepository, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "groups.db")
	dsn := path + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(0)"
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	for _, query := range []string{
		`INSERT INTO devices(id,name,protocol,status) VALUES('a','a','modbus_tcp','active'),('b','b','modbus_tcp','active'),('c','c','modbus_tcp','active')`,
		`INSERT INTO polling_groups(id,name) VALUES('g','delete-me'),('other','keep-me')`,
		`INSERT INTO points(id,device_id,name,address,data_type,polling_group_id,enabled) VALUES('a1','a','a1','40001','uint16','g',1),('a2','a','a2','40002','uint16','g',0),('b1','b','b1','40001','uint16','g',0),('c1','c','c1','40001','uint16','other',1)`,
	} {
		_, err := db.ExecContext(t.Context(), query)
		require.NoError(t, err)
	}
	return db, NewSQLRepository(db), dsn
}

func requireAcquisitionGroupMembers(t *testing.T, db *sql.DB, count int) {
	t.Helper()
	var actual int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM points WHERE polling_group_id='g'`).Scan(&actual))
	require.Equal(t, count, actual)
}

func TestAcquisitionGroupDeleteCoversBothServicePaths(t *testing.T) {
	for _, path := range []string{"point-service", "group-service"} {
		t.Run(path, func(t *testing.T) {
			db, repo, _ := acquisitionDeleteFixture(t)
			lifecycle := &acquisitionDeleteLifecycle{t: t}
			repo.SetAcquisitionMutationLifecycle(lifecycle)
			var err error
			if path == "point-service" {
				err = point.NewService(nil, repo).DeletePollingGroup(t.Context(), "g")
			} else {
				err = NewService(repo).Delete(t.Context(), "g")
			}
			require.NoError(t, err)
			require.Equal(t, []string{"a", "b"}, lifecycle.deviceIDs, "exact DISTINCT FK devices include disabled points")
			require.Equal(t, 1, lifecycle.begun)
			require.Equal(t, 1, lifecycle.finished)
			require.False(t, lifecycle.active)
			requireAcquisitionGroupMembers(t, db, 0)
		})
	}
}

func TestAcquisitionGroupDeleteBusyHookHasNoDurableEffect(t *testing.T) {
	db, repo, _ := acquisitionDeleteFixture(t)
	busy := errors.New("runtime lock busy")
	lifecycle := &acquisitionDeleteLifecycle{t: t, begin: func() error { return busy }}
	repo.SetAcquisitionMutationLifecycle(lifecycle)
	require.ErrorIs(t, repo.Delete(t.Context(), "g"), busy)
	requireAcquisitionGroupMembers(t, db, 3)
	require.Equal(t, 1, lifecycle.finished, "even a partially acquired callback must release its own tokens")
	_, err := repo.GetByID(t.Context(), "g")
	require.NoError(t, err)
}

func TestAcquisitionGroupDeleteReservationRejectsConcurrentAssignment(t *testing.T) {
	db, repo, dsn := acquisitionDeleteFixture(t)
	foreign, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	defer foreign.Close()
	rejected := errors.New("test callback refuses after reservation probe")
	lifecycle := &acquisitionDeleteLifecycle{t: t}
	lifecycle.begin = func() error {
		// Test-only lock probe. Production callbacks never perform SQL I/O.
		_, err := foreign.ExecContext(t.Context(), `UPDATE points SET polling_group_id='g' WHERE id='c1'`)
		require.Error(t, err, "affected-device enumeration must hold the writer reservation")
		return rejected
	}
	repo.SetAcquisitionMutationLifecycle(lifecycle)
	require.ErrorIs(t, repo.Delete(t.Context(), "g"), rejected)
	require.Equal(t, []string{"a", "b"}, lifecycle.deviceIDs)
	requireAcquisitionGroupMembers(t, db, 3)
	_, err = foreign.ExecContext(t.Context(), `UPDATE points SET polling_group_id='g' WHERE id='c1'`)
	require.NoError(t, err)
	lifecycle.begin = nil
	require.NoError(t, repo.Delete(t.Context(), "g"))
	require.Equal(t, []string{"a", "b", "c"}, lifecycle.deviceIDs, "retry must enumerate the now-committed assignment, not reuse old devices")
	require.Equal(t, 2, lifecycle.finished)
	requireAcquisitionGroupMembers(t, db, 0)
}
