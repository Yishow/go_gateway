//go:build f_write_group_fixture

package main

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

type fixtureRetryTargetResolver struct{ db *sql.DB }

func (r fixtureRetryTargetResolver) Resolve(context.Context, groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	return groupdelivery.Target{DB: r.db, Kind: schema.DatabaseConnectorKindSQLite}, nil
}

func TestFixtureRetryOverrideInjectsPipelineSender(t *testing.T) {
	t.Setenv(fixtureMaxRetriesEnv, "2")
	fixture := newGroupFixture().(*groupFixture)
	config := fixture.pipelineConfig("node", "owner")
	require.Equal(t, 2, config.Sender.MaxRetries)
}

func TestFixtureRetryPreflightRejectsInvalidOverrides(t *testing.T) {
	for name, value := range map[string]string{
		"zero":  "0",
		"minus": "-1",
		"text":  "not-a-number",
		"large": strconv.Itoa(fixtureMaxRetriesMax + 1),
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(fixtureMaxRetriesEnv, value)
			require.Error(t, preflightFixtureRetries())
		})
	}
}

func TestFixtureRetryInvalidOverrideFailsDatabasePreflight(t *testing.T) {
	fixtureOwnedDatabasePath(t, fixtureMarkerValue)
	t.Setenv(fixtureMaxRetriesEnv, "0")
	require.Error(t, preflightGroupFixtureDatabase())
}

func TestFixtureRetryPreflightLeavesUnsetAtProductionUnlimited(t *testing.T) {
	t.Setenv(fixtureMaxRetriesEnv, "")
	require.NoError(t, preflightFixtureRetries())
	require.Zero(t, groupDeliverySenderConfig().MaxRetries)
}

func TestFixtureRetryBoundExhaustsTheRealStoreAndSender(t *testing.T) {
	t.Setenv(fixtureMaxRetriesEnv, "2")
	env := newOutageEnv(t)
	store := groupdelivery.NewStore(env.db)
	targetFaults := &lostcommit.Faults{Lose: &atomic.Bool{}}
	target := lostcommit.OpenFaults(filepath.Join(t.TempDir(), "retry-target.db")+"?_pragma=busy_timeout(5000)", targetFaults)
	t.Cleanup(func() { require.NoError(t, target.Close()) })
	_, err := target.ExecContext(t.Context(), `CREATE TABLE readings (record_id TEXT, counter INTEGER NOT NULL)`)
	require.NoError(t, err)

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	row := dbtarget.EncodedRow{
		RecordID: "record-1", EffectKey: "effect-1", EntityKey: "entity-1", BucketStart: start,
		Cells: []dbtarget.EncodedCell{{Column: "record_id", Value: "record-1"}, {Column: "counter", Value: int64(1)}},
	}
	outcome := snapshot.Outcome{
		Kind: snapshot.OutcomeRow, EntityKey: row.EntityKey, RecordID: row.RecordID,
		EffectKey: row.EffectKey, BucketStart: start, BucketEnd: start.Add(10 * time.Second),
	}
	key := groupdelivery.GroupKey{WorkspaceID: "workspace-retry", GroupID: "group-retry", GroupRevision: "revision-1"}
	require.NoError(t, store.CommitClosure(t.Context(), key, groupdelivery.Closure{
		Destination: groupdelivery.FrozenDestination{
			Scope: "scope-retry", ConnectorID: "connector-retry", ConnectorRevision: "connector-revision-1",
			TableName: "readings", DedupeCapability: groupdelivery.DedupeNone,
		},
		NextClose: start.Add(10 * time.Second), Buckets: []groupdelivery.ClosedBucket{{Outcome: outcome, Row: &row}},
	}))

	failure := errors.New("connection reset by peer")
	targetFaults.ExecError.Store(&failure)
	sender := groupdelivery.NewSender(store, fixtureRetryTargetResolver{db: target}, groupDeliverySenderConfig())
	first, err := sender.Deliver(t.Context(), row.EffectKey)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateRetrying, first.State)
	second, err := sender.Deliver(t.Context(), row.EffectKey)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateBlocked, second.State)
	item, err := store.GetOutbox(t.Context(), row.EffectKey)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateBlocked, item.State)
	require.Equal(t, 2, item.RetryCount)
	require.Equal(t, "insert-failed", item.LastErrorCode)
	require.Equal(t, 1, countFixtureRetryRows(t, env.db, "wg_delivery_outbox"))
	_, err = sender.Deliver(t.Context(), row.EffectKey)
	require.ErrorIs(t, err, groupdelivery.ErrNotDeliverable)
}

func countFixtureRetryRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+table).Scan(&count))
	return count
}
