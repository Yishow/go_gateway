//go:build !f_write_group_fixture

package dbtarget

import (
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestFixtureDestinationHookIsNoOpInNormalBuild(t *testing.T) {
	dsn := createTargetSQLite(t)
	manager := datalinkbase.NewDBManager(datalinkbase.DBConfig{
		Type: datalinkbase.DBTypeSQLite, DSN: dsn, MaxOpenConns: 1, MaxIdleConns: 1,
	})
	require.NoError(t, manager.Connect())
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	db, closeFn, err := openFixtureDestination("connector-1", schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": dsn}, manager)
	require.NoError(t, err)
	require.Nil(t, db)
	require.Nil(t, closeFn)
}
