package grouppipeline

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLBoundaryOutboxMatchesFrozenProof(t *testing.T) {
	fixture := newFixture(t)
	pipeline := fixture.pipeline()
	snap := mysqlRuntimeSnapshot(t)
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
	built, reason := pipeline.buildFrozen(t.Context(), snap, at(0))
	require.Empty(t, reason)
	require.NotNil(t, built)
	require.NoError(t, built.boundary.AcceptSample(t.Context(), envelope(snap.Group, "sample", at(1), 21.5)))
	require.NoError(t, built.boundary.Tick(t.Context(), at(15)))
	var effectKey string
	require.NoError(t, fixture.store.DB().QueryRowContext(t.Context(),
		`SELECT effect_key FROM wg_delivery_outbox`).Scan(&effectKey))
	item, err := fixture.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)
	require.Equal(t, "record_id", item.Destination.RecordKeyColumn)
	frozen, err := pipeline.loadRuntimeSnapshot(t.Context(), item.Key)
	require.NoError(t, err)
	require.True(t, mysqlFrozenScopeMatches(item, frozen), "the production boundary must generate a destination the resolver accepts")
}
