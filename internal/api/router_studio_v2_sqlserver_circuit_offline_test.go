package api

import (
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

// The actual target transport rejects writes while offline; journal, closure,
// frozen outbox and counter baseline remain the production implementations.
func TestSQLServerCircuitOfflineRestartIntegration(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	f := newSQLServerCircuitPlainFixture(t, fixture)
	require.NoError(t, f.pipe.Stop(time.Second))
	faultDB, faults := sqlServerCircuitFaultDB(t, fixture)
	t.Cleanup(func() { require.NoError(t, faultDB.Close()) })
	capability, err := f.targetService.ProbeSQLServerRecordingAtRevision(t.Context(), "connector-A", "connector-1")
	require.NoError(t, err)
	offline := errors.New("fixture destination write transport unavailable")
	faults.ExecError.Store(&offline)
	start := func() *grouppipeline.Pipeline {
		pipe := grouppipeline.New(grouppipeline.Dependencies{
			Groups: sqlServerCircuitOneGroup{f.groups, f.cumulativeID}, Tags: f.tags,
			Destinations: sqlServerFaultDestinations{inner: f.targetService, db: faultDB, capability: capability},
			Inspector:    dbtarget.NewReadOnlyTableInspector(f.targetService), Store: f.store,
			Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db)),
		}, f.pipeConfig)
		require.NoError(t, pipe.Start(t.Context()))
		t.Cleanup(func() { _ = pipe.Stop(time.Second) })
		f.pipe = pipe
		return pipe
	}
	pipe := start()
	f.accept(t, f.cumulativeID, sqlServerCircuitMeasurementID, sqlServerCircuitInputs[0])
	f.closeAt(t, 10*time.Minute)
	require.Eventually(t, func() bool {
		var state string
		return f.db.QueryRowContext(t.Context(), `SELECT state FROM wg_delivery_outbox WHERE group_id=?`, f.cumulativeID).Scan(&state) == nil && state == groupdelivery.StateRetrying
	}, 20*time.Second, 25*time.Millisecond)
	var beforePayload, beforeDigest, beforeConnector, beforeRevision, beforeTable string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload,payload_digest,connector_id,connector_revision,table_name FROM wg_delivery_outbox WHERE group_id=?`,
		f.cumulativeID).Scan(&beforePayload, &beforeDigest, &beforeConnector, &beforeRevision, &beforeTable))
	var baseline string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines WHERE group_id=?`, f.cumulativeID).Scan(&baseline))
	require.NotEmpty(t, baseline)
	require.NoError(t, pipe.Stop(time.Second))
	pipe = start() // Restart from persisted runtime/outbox/checkpoints, still offline.
	var afterPayload, afterDigest, afterConnector, afterRevision, afterTable string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload,payload_digest,connector_id,connector_revision,table_name FROM wg_delivery_outbox WHERE group_id=?`,
		f.cumulativeID).Scan(&afterPayload, &afterDigest, &afterConnector, &afterRevision, &afterTable))
	require.Equal(t, []string{beforePayload, beforeDigest, beforeConnector, beforeRevision, beforeTable},
		[]string{afterPayload, afterDigest, afterConnector, afterRevision, afterTable}, "restart cannot retarget or recalculate a frozen effect")
	faults.ExecError.Store(nil)
	require.Eventually(t, func() bool {
		var state string
		return f.db.QueryRowContext(t.Context(), `SELECT state FROM wg_delivery_outbox WHERE group_id=?`, f.cumulativeID).Scan(&state) == nil && state == groupdelivery.StateCommitted
	}, 20*time.Second, 25*time.Millisecond)
	require.NoError(t, pipe.Stop(time.Second))
	var resumedBaseline string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines WHERE group_id=?`, f.cumulativeID).Scan(&resumedBaseline))
	require.Equal(t, baseline, resumedBaseline, "delivery recovery must not fold the counter again")
	row, err := groupdelivery.DecodeRowPayload([]byte(beforePayload))
	require.NoError(t, err)
	var records, receipts int
	require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+sqlServerQualified(f.namespace, beforeTable)+" WHERE record_id=@p1", row.RecordID).Scan(&records))
	require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+sqlServerQualified(f.namespace, dbtarget.EffectReceiptTable)+
		" WHERE effect_key=@p1 AND payload_digest=@p2", row.EffectKey, beforeDigest).Scan(&receipts))
	require.Equal(t, 1, records)
	require.Equal(t, 1, receipts)
	cleanupSQLServerCircuitEffects(t, f.target, f.namespace, beforeTable, []string{row.RecordID},
		map[string]string{row.EffectKey: beforeDigest}, map[string]dbtarget.EncodedRow{row.RecordID: row})
}
