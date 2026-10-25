package grouppipeline

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// closureState is everything one closure transaction writes or consumes, read
// from the real local database.
type closureState struct {
	outbox, buckets, journaled, unconsumed, baselines int
	checkpoint                                        string
	baselinePayloads                                  [2]string
}

func (f *longPipeline) closureState() closureState {
	f.t.Helper()
	local := f.plainLocal()
	state := closureState{
		outbox:     f.count(local, `SELECT COUNT(*) FROM wg_delivery_outbox`),
		buckets:    f.count(local, `SELECT COUNT(*) FROM wg_delivery_buckets`),
		journaled:  f.count(local, `SELECT COUNT(*) FROM wg_delivery_samples`),
		unconsumed: f.count(local, `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed = 0`),
		baselines:  f.count(local, `SELECT COUNT(*) FROM wg_delivery_counter_baselines`),
	}
	require.NoError(f.t, local.QueryRowContext(f.t.Context(), `SELECT COALESCE(MAX(next_close), '') FROM wg_delivery_checkpoints`).Scan(&state.checkpoint))
	for member := range 2 {
		_, state.baselinePayloads[member] = f.baseline(member)
	}
	return state
}

// failClosureAt makes the local database refuse the closure statement that
// contains the text, and counts how often it did.
func (f *longPipeline) failClosureAt(statement string) *atomic.Int64 {
	hits := &atomic.Int64{}
	reject := func(query string) error {
		if strings.Contains(query, statement) {
			hits.Add(1)
			return errors.New("injected closure failure")
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	return hits
}

// countCommits counts the transactions a faulted database really committed.
func countCommits(faults interface{ Store(*func()) }) *atomic.Int64 {
	commits := &atomic.Int64{}
	hook := func() { commits.Add(1) }
	faults.Store(&hook)
	return commits
}

// requireBaselines checks the durable counter baselines: a counter shape has one
// per member at the end of the window; an average never writes any.
func (f *longPipeline) requireBaselines(shape averageShape, window int) {
	f.t.Helper()
	if shape.cumulative() {
		f.requireBaselinesAt(window)
		return
	}
	require.Zero(f.t, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_counter_baselines`), "an average has no counter baseline")
}

// uninterruptedDigests runs the same scripted windows without any fault: the
// frozen payloads every recovery below must end up with.
func uninterruptedDigests(t *testing.T, shape averageShape) map[string]string {
	t.Helper()
	f := newAveragePipeline(t, shape)
	for window := range 2 {
		f.feed(shape, window)
		f.closeUntil(shape.windowEnd(window))
		f.deliver()
	}
	require.Equal(t, shape.want(0, 1), f.averageRows(shape))
	return f.outboxDigests()
}

// resendOf is the one reading of the first member a collector resends after a
// lost journal ACK.
func resendOf(shape averageShape) (id string, seconds int, value uint64) {
	if shape.cumulative() {
		return "m1-1", 19 * 60, firstStart + 1
	}
	return "m0-1a", 16, 100
}

// lateSeconds is an instant inside window 1 that a reading reaches after the
// window closed.
func lateSeconds(shape averageShape) int {
	if shape.cumulative() {
		return 19*60 + 30
	}
	return 25
}

// TestAverageClosureRecovery injects the faults of tasks 2.4 into the production
// pipeline over a real local database and a real SQLite destination, for the
// long and wide averages and for the cumulative wide counters: a failure before,
// inside and after the closure transaction, a crash after the journal ACK and an
// offline destination across restarts. A closure that did not commit consumes
// nothing; one that did is never closed again; sums, counts, baselines and the
// frozen payloads are the ones an uninterrupted run freezes.
func TestAverageClosureRecovery(t *testing.T) {
	for _, shape := range averageShapes {
		t.Run(string(shape), func(t *testing.T) {
			control := uninterruptedDigests(t, shape)
			// ready has window 0 closed and delivered and window 1 accepted: ACKed
			// into the durable journal, not closed.
			ready := func(t *testing.T) *longPipeline {
				f := newAveragePipeline(t, shape)
				f.feed(shape, 0)
				f.closeUntil(shape.windowEnd(0))
				f.deliver()
				require.Equal(t, shape.want(0), f.averageRows(shape), "window 0 is stored once")
				f.requireBaselines(shape, 0)
				f.feed(shape, 1)
				return f
			}
			finish := func(t *testing.T, f *longPipeline) {
				t.Helper()
				f.deliver()
				require.Equal(t, shape.want(0, 1), f.averageRows(shape), "the exact sum, count, result and delta, stored once")
				require.Equal(t, control, f.outboxDigests(), "the frozen payloads are the uninterrupted run's")
				f.requireBaselines(shape, 1)
				require.Zero(t, f.closureState().unconsumed, "every closed reading is consumed")
			}

			for _, stage := range []struct{ name, statement string }{
				{"on the first row it writes", "INSERT INTO wg_delivery_buckets"},
				{"after the rows and before the readings are consumed", "UPDATE wg_delivery_samples SET consumed"},
				{"on its last statement", "INSERT INTO wg_delivery_checkpoints"},
			} {
				t.Run("a closure that fails "+stage.name+" changes nothing and retries identically", func(t *testing.T) {
					f := ready(t)
					before := f.closureState()
					require.Positive(t, before.unconsumed, "window 1's readings are waiting")
					hits := f.failClosureAt(stage.statement)

					f.closeUntil(shape.windowEnd(1))
					require.Positive(t, hits.Load(), "the injected fault fired")
					require.NotEmpty(t, f.reportedErrors(), "the failed closure is reported")
					require.Equal(t, before, f.closureState(), "no outbox, bucket, consumption, checkpoint or baseline of the failed closure is visible")
					require.Equal(t, shape.want(0), f.averageRows(shape))

					f.localFault.Reject.Store(nil)
					f.closeUntil(shape.windowEnd(1))
					finish(t, f)
				})
			}

			t.Run("a closure that committed but lost its response is not consumed twice", func(t *testing.T) {
				f := ready(t)
				commits := countCommits(&f.localFault.AfterCommit)
				f.localFault.Lose.Store(true)
				f.closeUntil(shape.windowEnd(1))
				require.Positive(t, commits.Load(), "the closure transaction really committed")
				require.NotEmpty(t, f.reportedErrors())
				committed := f.closureState()
				require.Equal(t, 2*shape.recordsPerWindow(), committed.outbox, "the commit happened")
				require.Zero(t, committed.unconsumed)
				digests := f.outboxDigests()

				f.localFault.Lose.Store(false)
				f.closeUntil(shape.windowEnd(1)) // the same process retries the closure it could not confirm
				require.Equal(t, committed, f.closureState(), "the retry commits the same frozen state, baseline and checkpoint")
				require.Equal(t, digests, f.outboxDigests())
				finish(t, f)
			})

			t.Run("a process that dies right after the closure commit restarts from the committed state", func(t *testing.T) {
				f := ready(t)
				f.localFault.Lose.Store(true)
				f.closeUntil(shape.windowEnd(1))
				committed := f.closureState()
				digests := f.outboxDigests()

				f.localFault.Lose.Store(false)
				f.restart()
				f.closeUntil(shape.windowEnd(1))
				require.Equal(t, committed, f.closureState(), "nothing is closed or consumed a second time")
				require.Equal(t, digests, f.outboxDigests())
				finish(t, f)
			})

			t.Run("a crash after the journal ACK keeps every reading and closes it once", func(t *testing.T) {
				f := ready(t)
				journaled := f.closureState().journaled
				f.restart() // dies before any closure of window 1
				require.Equal(t, journaled, f.closureState().journaled, "an ACKed reading is journaled")
				f.closeUntil(shape.windowEnd(1))
				finish(t, f)
			})

			t.Run("a reading resent after a lost journal ACK is the same reading", func(t *testing.T) {
				f := ready(t)
				journaled := f.closureState().journaled
				f.restart()
				id, seconds, value := resendOf(shape)
				m := f.group.Members[0]
				resend := func(value uint64) error {
					f.clock = at(seconds + 2)
					return f.pipe.AcceptSample(t.Context(), sampleEnvelope(f, m, id, at(seconds), at(seconds+30), value))
				}
				require.NoError(t, resend(value), "the same sample is acknowledged, not counted again")
				require.Equal(t, journaled, f.closureState().journaled)
				var refused *runtime.GroupSampleError
				require.ErrorAs(t, resend(value+9), &refused, "the same ID with different content is a visible conflict")
				require.Equal(t, "identity-conflict", refused.Reason)

				f.closeUntil(shape.windowEnd(1))
				finish(t, f)
			})

			t.Run("a pause of many windows is closed within the existing per-tick bound", func(t *testing.T) {
				f := newAveragePipeline(t, shape)
				f.feed(shape, 0)
				const pause = 1500
				f.closeUntil(shape.windowEnd(0))
				buckets := func() int { return f.closureState().buckets }
				closed := buckets()
				total := (pause + 1) * shape.recordsPerWindow()
				ticks := 0
				for closed < total && ticks < 8 {
					f.closeUntil(shape.windowEnd(pause))
					ticks++
					now := buckets()
					require.Positive(t, now-closed, "every tick makes progress")
					require.LessOrEqual(t, now-closed, snapshot.DefaultMaxBucketsPerTick, "one tick closes at most the existing bound")
					closed = now
				}
				require.Equal(t, total, closed, "the backlog drains completely, in order, with no window dropped")
				require.Greater(t, ticks, 1, "the pause really needed more than one bounded tick")
				require.Zero(t, f.closureState().unconsumed)
			})

			t.Run("an offline destination keeps the frozen pending results across restarts", func(t *testing.T) {
				f := newAveragePipeline(t, shape)
				f.dest.offline.Store(true)
				f.feed(shape, 0)
				f.closeUntil(shape.windowEnd(0))
				f.deliver()
				require.Empty(t, f.averageRows(shape), "the destination is offline")
				pending := `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state != 'sql_committed'`
				require.Equal(t, shape.recordsPerWindow(), f.count(f.plainLocal(), pending), "the closed window is kept as pending effects")
				require.Positive(t, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE last_error_code != ''`), "delivery really tried and failed")
				frozen := f.outboxDigests()

				f.restart()
				require.Equal(t, frozen, f.outboxDigests())
				f.feed(shape, 1)
				f.closeUntil(shape.windowEnd(1))
				f.deliver()
				require.Equal(t, 2*shape.recordsPerWindow(), f.count(f.plainLocal(), pending))
				f.requireBaselines(shape, 1)
				for key, digest := range frozen {
					require.Equal(t, digest, f.outboxDigests()[key], "a pending effect is never recomputed")
				}

				f.restart() // still offline
				f.dest.offline.Store(false)
				finish(t, f)
				require.Equal(t, 2*shape.recordsPerWindow(), f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = 'sql_committed'`))
			})
		})
	}
}

// TestAverageCommitAckLost loses the destination's commit acknowledgement of a
// finalized average or counter row: SQL holds one row and one receipt per effect,
// the retry resends the frozen payload after a restart, and neither the sum and
// count nor the counter baseline moves, even when a late reading arrives.
func TestAverageCommitAckLost(t *testing.T) {
	for _, shape := range averageShapes {
		t.Run(string(shape), func(t *testing.T) {
			control := uninterruptedDigests(t, shape)
			f := newAveragePipeline(t, shape)
			f.feed(shape, 0)
			f.closeUntil(shape.windowEnd(0))
			f.deliver()
			f.feed(shape, 1)
			f.closeUntil(shape.windowEnd(1))
			frozen := f.outboxDigests()
			closed := f.closureState()

			commits := countCommits(&f.targetFlt.AfterCommit)
			f.targetFlt.Lose.Store(true)
			f.deliverCycles(1)
			require.Positive(t, commits.Load(), "the destination committed a row whose reply was lost")
			require.Greater(t, len(f.averageRows(shape)), 2, "window 1 is in SQL even though the sender never heard back")
			require.Equal(t, shape.recordsPerWindow(), f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_receipts`), "no local proof until the destination confirms")

			// A reading that arrives for the closed window afterwards is late: it cannot
			// change what was frozen.
			err := f.reading(0, "late", at(lateSeconds(shape)), 999)
			var late *runtime.GroupSampleError
			require.ErrorAs(t, err, &late)
			require.Equal(t, "late-after-close", late.Reason)

			f.targetFlt.Lose.Store(false)
			f.restart() // the retry also survives a restart
			f.deliver()
			require.Equal(t, shape.want(0, 1), f.averageRows(shape), "one row per record with the original sum, count and delta")
			require.Equal(t, frozen, f.outboxDigests(), "retries resend the frozen payloads")
			require.Equal(t, control, f.outboxDigests())
			require.Equal(t, 2*shape.recordsPerWindow(), f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = 'sql_committed'`))
			require.Equal(t, 2*shape.recordsPerWindow(), f.count(f.plainTarget(), `SELECT COUNT(*) FROM gw_effect_receipts`), "one receipt per effect")
			after := f.closureState()
			require.Equal(t, closed.baselinePayloads, after.baselinePayloads, "delivery retries never move the baseline")
			require.Equal(t, closed.checkpoint, after.checkpoint)
			f.requireBaselines(shape, 1)
		})
	}
}
