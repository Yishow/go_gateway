//go:build f_write_group_fixture

package dbtarget

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

func TestFixtureDestinationFaultWrapsRealSQLiteCommit(t *testing.T) {
	ctx := t.Context()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)
	prepareTargetFixture(t, targetDSN)

	connectorRepo := NewSQLConnectorRepository(mainDB)
	service := NewConnectorService(connectorRepo, NewSQLTargetMappingRepository(mainDB))
	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name: "fixture-target", Kind: "sqlite", ConnectionConfig: ConnectionConfig{"dsn": targetDSN},
	})
	require.NoError(t, err)

	controller := NewFixtureDestinationFaultController()
	InstallFixtureDestinationController(controller)
	t.Cleanup(func() {
		controller.Close()
		InstallFixtureDestinationController(nil)
	})
	_, err = controller.Configure(connector.ID, FixtureDestinationFaultTargetCommitResponseLost, true)
	require.NoError(t, err)

	opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, opened.Close()) })

	tx, err := opened.DB.BeginTx(ctx, nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `INSERT INTO fixture_effects (effect_key) VALUES (?)`, "effect-1")
	require.NoError(t, err)
	err = tx.Commit()
	require.ErrorIs(t, err, lostcommit.ErrResponseLost)

	var count int
	require.NoError(t, opened.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM fixture_effects`).Scan(&count))
	require.Equal(t, 1, count, "the real destination commit happened before the response was lost")
	state := controller.Snapshot()
	require.Len(t, state, 1)
	require.Equal(t, uint64(1), state[0].Commits)
	require.Equal(t, uint64(1), state[0].Reached)
}

func TestFixtureDestinationUnknownKindIsRejected(t *testing.T) {
	controller := NewFixtureDestinationFaultController()
	_, err := controller.Configure("connector-1", FixtureDestinationFaultKind("arbitrary_sql"), true)
	require.Error(t, err)
}

func TestFixtureDestinationCommitHoldReleasesAfterDisable(t *testing.T) {
	ctx := t.Context()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)
	prepareTargetFixture(t, targetDSN)
	connectorRepo := NewSQLConnectorRepository(mainDB)
	service := NewConnectorService(connectorRepo, NewSQLTargetMappingRepository(mainDB))
	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name: "fixture-hold-target", Kind: "sqlite", ConnectionConfig: ConnectionConfig{"dsn": targetDSN},
	})
	require.NoError(t, err)
	controller := NewFixtureDestinationFaultController()
	InstallFixtureDestinationController(controller)
	t.Cleanup(func() { controller.Close(); InstallFixtureDestinationController(nil) })
	_, err = controller.Configure(connector.ID, FixtureDestinationFaultTargetCommitHold, true)
	require.NoError(t, err)
	opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, opened.Close()) })

	tx, err := opened.DB.BeginTx(ctx, nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `INSERT INTO fixture_effects (effect_key) VALUES (?)`, "effect-hold")
	require.NoError(t, err)
	commitDone := make(chan error, 1)
	go func() { commitDone <- tx.Commit() }()
	require.Eventually(t, func() bool { return controller.Snapshot()[0].Holding }, time.Second, time.Millisecond)
	_, err = controller.Configure(connector.ID, FixtureDestinationFaultTargetCommitHold, false)
	require.NoError(t, err)
	require.NoError(t, <-commitDone)
	require.False(t, controller.Snapshot()[0].Holding)

	var count int
	require.NoError(t, opened.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM fixture_effects`).Scan(&count))
	require.Equal(t, 1, count)
}

func prepareTargetFixture(t *testing.T, dsn string) {
	t.Helper()
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	_, err = db.ExecContext(context.Background(), `CREATE TABLE fixture_effects (effect_key TEXT PRIMARY KEY)`)
	require.NoError(t, err)
}
