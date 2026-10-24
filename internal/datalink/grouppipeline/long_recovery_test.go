package grouppipeline

import (
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const (
	firstStart  = uint64(9007199254740992)
	firstDelta  = "1"
	secondStart = uint64(100)
	secondDelta = "30"
)

func sampleEnvelope(f *longPipeline, m workspace.WriteGroupMember, id string, observed, received time.Time, value uint64) measurement.SampleEnvelope {
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: f.group.WorkspaceID, DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID,
		SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision, ObservedAt: observed, ReceivedAt: received,
		TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}
}

// eachBackend runs a scenario against a real SQLite destination and a real,
// disposable PostgreSQL schema. The PostgreSQL case reports itself as not run,
// never as passed, when no database is configured.
func eachBackend(t *testing.T, run func(t *testing.T, build func(*testing.T) *longPipeline)) {
	t.Helper()
	t.Run("sqlite", func(t *testing.T) { run(t, newLongPipeline) })
	t.Run("postgres", func(t *testing.T) { run(t, newLongPipelinePostgres) })
}

// closeFirstWindow accepts one reading per counter in window 0 and closes it.
func (f *longPipeline) closeFirstWindow() {
	f.t.Helper()
	require.NoError(f.t, f.reading(0, "m1-0", at(9*60), firstStart))
	require.NoError(f.t, f.reading(1, "m2-0", at(9*60), secondStart))
	f.closeUntil(at(600))
}

// acceptSecondWindow accepts window 1's readings; the counters advanced by 1
// and by 30.
func (f *longPipeline) acceptSecondWindow() {
	f.t.Helper()
	require.NoError(f.t, f.reading(0, "m1-1", at(19*60), firstStart+1))
	require.NoError(f.t, f.reading(1, "m2-1", at(19*60), secondStart+30))
}

// requireDestinationExact reads the destination's long table and checks both
// windows: a baseline then an exact difference, exactly one row per record.
func (f *longPipeline) requireDestinationExact() {
	f.t.Helper()
	rows := f.targetRows()
	require.Len(f.t, rows, 4, "one stored record per counter and window")
	for _, tag := range []string{longKwhFirst, longKwhSecond} {
		require.Equal(f.t, "baseline", rows[windowKey(tag, 0)].status.String, tag)
		require.False(f.t, rows[windowKey(tag, 0)].delta.Valid, "the first window has no delta, not zero")
		require.Equal(f.t, "ok", rows[windowKey(tag, 1)].status.String, tag)
	}
	require.Equal(f.t, "9007199254740992", rows[windowKey(longKwhFirst, 0)].raw.String)
	require.Equal(f.t, "9007199254740993", rows[windowKey(longKwhFirst, 1)].raw.String, "exact above 2^53")
	require.Equal(f.t, firstDelta, rows[windowKey(longKwhFirst, 1)].delta.String)
	require.Equal(f.t, secondDelta, rows[windowKey(longKwhSecond, 1)].delta.String)
	require.Equal(f.t, at(9*60).Format(longTime), rows[windowKey(longKwhFirst, 1)].deltaFrom.String, "the delta spans the real observation times")
	require.Equal(f.t, at(19*60).Format(longTime), rows[windowKey(longKwhFirst, 1)].deltaTo.String)
}

// requireBaselineAt checks that both members' durable baselines sit at the end
// of the given window, each committed exactly once.
func (f *longPipeline) requireBaselinesAt(window int) {
	f.t.Helper()
	want := at((window + 1) * int(longWindow/time.Second)).Format("2006-01-02T15:04:05.000000000Z")
	for member := range 2 {
		end, payload := f.baseline(member)
		require.Equal(f.t, want, end, "member %d baseline window", member)
		require.NotEmpty(f.t, payload)
	}
	require.Equal(f.t, 2, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_counter_baselines`), "one baseline row per member")
}

func (f *longPipeline) failClosureBaselineWrite() {
	reject := func(query string) error {
		if strings.Contains(query, "wg_delivery_counter_baselines") {
			return errors.New("injected closure failure")
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
}

// TestProfileClosureCrashRecovery injects a failure before, in the middle of and
// after the closure transaction, and after the journal ACK, through the
// production pipeline into a real SQLite destination.
func TestProfileClosureCrashRecovery(t *testing.T) {
	eachBackend(t, func(t *testing.T, build func(*testing.T) *longPipeline) {
		t.Run("a closure that fails midway changes nothing and retries identically", func(t *testing.T) {
			f := build(t)
			f.closeFirstWindow()
			f.deliver()
			f.requireBaselinesAt(0)
			f.acceptSecondWindow()

			f.failClosureBaselineWrite()
			f.closeUntil(at(1200))
			require.NotEmpty(t, f.reportedErrors(), "the failed closure is reported")
			local := f.plainLocal()
			require.Equal(t, 2, f.count(local, `SELECT COUNT(*) FROM wg_delivery_outbox`), "no outbox intent without its baseline")
			require.Equal(t, 2, f.count(local, `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed = 0`), "the window's readings are still unconsumed")
			f.requireBaselinesAt(0)

			f.localFault.Reject.Store(nil)
			f.closeUntil(at(1200))
			f.deliver()
			f.requireBaselinesAt(1)
			f.requireDestinationExact()
			require.Equal(t, 4, f.count(local, `SELECT COUNT(DISTINCT effect_key) FROM wg_delivery_outbox`), "one effect per record")
		})

		t.Run("a closure that committed but lost its response is not consumed twice", func(t *testing.T) {
			f := build(t)
			f.closeFirstWindow()
			f.deliver()
			f.acceptSecondWindow()

			f.localFault.Lose.Store(true)
			f.closeUntil(at(1200))
			require.NotEmpty(t, f.reportedErrors())
			local := f.plainLocal()
			require.Equal(t, 4, f.count(local, `SELECT COUNT(*) FROM wg_delivery_outbox`), "the commit really happened")
			f.requireBaselinesAt(1)
			_, committed := f.baseline(0)
			digestsBefore := f.outboxDigests()

			f.localFault.Lose.Store(false)
			f.closeUntil(at(1200)) // the same process retries the closure it could not confirm
			_, retried := f.baseline(0)
			require.Equal(t, committed, retried, "the retry commits the same frozen baseline")
			require.Equal(t, digestsBefore, f.outboxDigests(), "and the same frozen payloads")
			f.deliver()
			f.requireDestinationExact()
			f.requireBaselinesAt(1)
		})

		t.Run("a process that dies right after the closure commit restarts from the committed state", func(t *testing.T) {
			f := build(t)
			f.closeFirstWindow()
			f.deliver()
			f.acceptSecondWindow()
			f.localFault.Lose.Store(true)
			f.closeUntil(at(1200))
			digestsBefore := f.outboxDigests()

			f.localFault.Lose.Store(false)
			f.restart()
			f.closeUntil(at(1200))
			require.Equal(t, digestsBefore, f.outboxDigests(), "nothing is closed or consumed a second time")
			f.deliver()
			f.requireDestinationExact()
			f.requireBaselinesAt(1)
		})

		t.Run("a crash after the journal ACK keeps the reading and the baseline", func(t *testing.T) {
			f := build(t)
			f.closeFirstWindow()
			f.deliver()
			f.acceptSecondWindow() // ACKed: journaled before the intake returned
			f.restart()            // dies before any closure of window 1
			f.closeUntil(at(1200))
			f.deliver()
			f.requireDestinationExact()
			f.requireBaselinesAt(1)
		})

		t.Run("a reading resent after a lost journal ACK is the same reading", func(t *testing.T) {
			f := build(t)
			f.closeFirstWindow()
			f.deliver()
			f.acceptSecondWindow()
			local := f.plainLocal()
			journaled := f.count(local, `SELECT COUNT(*) FROM wg_delivery_samples`)
			f.restart()

			// The collector never saw the ACK and resends with a later receive time.
			m := f.group.Members[0]
			resend := func(value uint64) error {
				f.clock = at(19*60 + 2)
				return f.pipe.AcceptSample(t.Context(), sampleEnvelope(f, m, "m1-1", at(19*60), at(19*60+30), value))
			}
			require.NoError(t, resend(firstStart+1), "the same sample is acknowledged, not duplicated")
			require.Equal(t, journaled, f.count(local, `SELECT COUNT(*) FROM wg_delivery_samples`))
			var refused *runtime.GroupSampleError
			require.ErrorAs(t, resend(firstStart+9), &refused, "the same ID with different content is a visible conflict")
			require.Equal(t, "identity-conflict", refused.Reason)

			f.closeUntil(at(1200))
			f.deliver()
			f.requireDestinationExact()
			var receivedAt string
			require.NoError(t, local.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_samples WHERE sample_id = 'm1-1'`).Scan(&receivedAt))
			require.Contains(t, receivedAt, at(19*60).Add(40*time.Millisecond).Format(time.RFC3339Nano), "the first receive time is kept")
		})
	})
}

// TestProfileOfflineRestart keeps accepted work while the destination is down
// across restarts: baselines survive, pending records are neither lost nor
// recomputed, and they arrive unchanged once the destination returns.
func TestProfileOfflineRestart(t *testing.T) {
	eachBackend(t, offlineRestart)
}

func offlineRestart(t *testing.T, build func(*testing.T) *longPipeline) {
	f := build(t)
	f.dest.offline.Store(true)
	f.closeFirstWindow()
	f.deliver()
	require.Empty(t, f.targetRows(), "the destination is offline")
	require.Equal(t, 2, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state != 'sql_committed'`), "pending records are kept")
	frozen := f.outboxDigests()

	f.restart()
	require.Equal(t, frozen, f.outboxDigests())
	f.acceptSecondWindow()
	f.closeUntil(at(1200))
	f.deliver()
	f.requireBaselinesAt(1)
	require.Equal(t, 4, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state != 'sql_committed'`))
	for key, digest := range frozen {
		require.Equal(t, digest, f.outboxDigests()[key], "a pending record is never rewritten")
	}

	f.restart() // still offline
	f.dest.offline.Store(false)
	f.deliver()
	f.requireDestinationExact()
	require.Equal(t, 4, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = 'sql_committed'`))
	f.requireBaselinesAt(1)
}

// TestProfileCommitAckLost loses the destination's commit acknowledgement: the
// row is stored once, the retry resends the frozen payload, and neither the
// baseline nor a late reading changes what was closed.
func TestProfileCommitAckLost(t *testing.T) {
	eachBackend(t, commitAckLost)
}

func commitAckLost(t *testing.T, build func(*testing.T) *longPipeline) {
	f := build(t)
	f.closeFirstWindow()
	f.deliver()
	f.acceptSecondWindow()
	f.closeUntil(at(1200))
	frozen := f.outboxDigests()
	_, baselineBefore := f.baseline(0)

	f.targetFlt.Lose.Store(true)
	f.deliverCycles(1)
	require.NotEmpty(t, f.targetRows(), "the destination committed even though the response was lost")
	require.Equal(t, 0, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key IN (SELECT effect_key FROM wg_delivery_outbox WHERE bucket_start > '2026-01-01T00:10:00.000000000Z')`),
		"no local proof until the destination confirms")

	// A reading that arrives for the closed window afterwards is late: it cannot change what was frozen.
	err := f.reading(0, "m1-late", at(19*60+30), firstStart+500)
	var late *runtime.GroupSampleError
	require.ErrorAs(t, err, &late)
	require.Equal(t, "late-after-close", late.Reason)

	f.targetFlt.Lose.Store(false)
	f.restart() // the retry also survives a restart
	f.deliver()
	f.requireDestinationExact()
	require.Equal(t, frozen, f.outboxDigests(), "retries resend the frozen payloads")
	require.Equal(t, 4, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = 'sql_committed'`))
	_, baselineAfter := f.baseline(0)
	require.Equal(t, baselineBefore, baselineAfter, "delivery retries never move the baseline")
	f.requireBaselinesAt(1)
	var states int
	require.NoError(t, f.plainLocal().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = ?`, groupdelivery.StateCommitted).Scan(&states))
	require.Equal(t, 4, states)
}
