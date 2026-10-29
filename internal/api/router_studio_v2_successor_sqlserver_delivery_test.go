package api

import (
	"context"
	"database/sql"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
)

// The real connector is resolved first, retaining frozen revision validation;
// only its transport is wrapped to lose the response after the actual SQL Server commits.
type sqlServerSuccessorLostACKResolver struct {
	rig      *sqlServerSuccessorRig
	dbA, dbB *sql.DB
}

func (r sqlServerSuccessorLostACKResolver) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	target, err := r.rig.Resolve(ctx, item)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	if target.Close != nil {
		if err := target.Close(); err != nil {
			return groupdelivery.Target{}, err
		}
	}
	target.DB, target.Close = r.dbA, nil
	if item.Destination.ConnectorID == "connector-B" {
		target.DB = r.dbB
	}
	return target, nil
}

func TestSQLServerSuccessorActivationOldBacklogAndCommitACKLost(t *testing.T) {
	sqlServerSuccessorFixtures(t)
	for _, transition := range sqlServerSuccessorRuntimeTransitions {
		t.Run(transition.name, func(t *testing.T) {
			r := newSQLServerSuccessorRig(t, transition.mode, false, true)
			sqlServerSuccessorRegisterEffectCleanup(t, r)
			priorAReceipts := sqlServerSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			priorBReceipts := sqlServerSuccessorRuntimeEvidenceRows(t, r.targetDB, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			first := r.snapshot.EffectiveAt.Add(time.Second)
			r.clock.Store(first.UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("old-backlog", first, float32(4))))
			r.clock.Store(r.snapshot.EffectiveAt.Add(16 * time.Second).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			var effect string
			require.NoError(t, r.f.db.QueryRow(`SELECT effect_key FROM wg_delivery_outbox WHERE group_id=?`, r.groupID).Scan(&effect))
			original, err := r.store.GetOutbox(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, "connector-A", original.Destination.ConnectorID)
			require.Empty(t, sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group))
			r.clock.Store(r.boundary.Add(-2500 * time.Millisecond).UnixNano())
			r.activate(t)
			r.clock.Store(r.boundary.Add(6 * time.Second).UnixNano())
			require.NoError(t, r.pipe.Reconcile(t.Context()))
			frozen, err := r.store.GetOutbox(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, original, frozen, "head replacement must not touch finalized backlog")

			wrappedA, faultsA := sqlServerCircuitFaultDB(t, r.fixtures[0])
			wrappedB, faultsB := sqlServerCircuitFaultDB(t, r.fixtures[1])
			var commitsA, commitsB atomic.Int32
			afterA, afterB := func() { commitsA.Add(1) }, func() { commitsB.Add(1) }
			faultsA.AfterCommit.Store(&afterA)
			faultsB.AfterCommit.Store(&afterB)
			t.Cleanup(func() { require.NoError(t, wrappedA.Close()); require.NoError(t, wrappedB.Close()) })
			now := func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }
			sender := groupdelivery.NewSender(r.store, sqlServerSuccessorLostACKResolver{r, wrappedA, wrappedB}, groupdelivery.SenderConfig{
				Owner: "sqlserver-successor-lost-ack", Now: now,
				Backoff: func(int) (string, time.Time) { return groupdelivery.StateRetrying, now().Add(-time.Second) },
			})
			faultsA.Lose.Store(true)
			attempt, err := sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateRetrying, attempt.State)
			require.Equal(t, "commit-ambiguous", attempt.Code)
			require.EqualValues(t, 1, commitsA.Load(), "fault fires only after actual SQL COMMIT")
			require.Zero(t, commitsB.Load())
			require.Equal(t, []string{"4"}, sqlServerSuccessorRuntimeValues(t, r.f.target, r.snapshot.Group))
			require.Empty(t, sqlServerSuccessorRuntimeRows(t, r.targetDB, r.prepared.Group), "old backlog must never be redirected to B")
			committedSQL := sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group)
			receipt := sqlServerSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.Len(t, receipt, len(priorAReceipts)+1)
			successorRuntimePreservedRows(t, priorAReceipts, receipt)
			var localReceipts int
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, effect).Scan(&localReceipts))
			require.Zero(t, localReceipts)
			faultsA.Lose.Store(false)
			attempt, err = sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, attempt.State)
			require.EqualValues(t, 1, commitsA.Load(), "receipt replay must not commit another effect")
			require.Equal(t, committedSQL, sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group))
			require.Equal(t, receipt, sqlServerSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, effect).Scan(&localReceipts))
			require.Equal(t, 1, localReceipts)
			settled, err := r.store.GetOutbox(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, original.Destination, settled.Destination)
			require.Equal(t, original.Payload, settled.Payload)
			require.Equal(t, original.PayloadDigest, settled.PayloadDigest)
			require.Equal(t, original.Key, settled.Key)
			require.Empty(t, sqlServerSuccessorRuntimeRows(t, r.targetDB, r.prepared.Group))
			// The same ambiguous-commit protection applies to the new generation at B.
			observed := r.boundary.Add(1500 * time.Millisecond)
			r.clock.Store(r.boundary.Add(6 * time.Second).UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("new-after-backlog", observed, float32(7))))
			r.clock.Store(r.boundary.Add(16 * time.Second).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			var nextCount int
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_outbox WHERE group_id=?`, r.prepared.Group.ID).Scan(&nextCount))
			require.Equal(t, 1, nextCount, "close only the successor window containing the sample")
			var nextEffect string
			require.NoError(t, r.f.db.QueryRow(`SELECT effect_key FROM wg_delivery_outbox WHERE group_id=?`, r.prepared.Group.ID).Scan(&nextEffect))
			require.NotEqual(t, effect, nextEffect)
			nextOriginal, err := r.store.GetOutbox(t.Context(), nextEffect)
			require.NoError(t, err)
			require.Equal(t, "connector-B", nextOriginal.Destination.ConnectorID)
			faultsB.Lose.Store(true)
			attempt, err = sender.Deliver(t.Context(), nextEffect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateRetrying, attempt.State)
			require.Equal(t, "commit-ambiguous", attempt.Code)
			require.EqualValues(t, 1, commitsB.Load(), "B has its own actual commit-lost-response attempt")
			require.Equal(t, []string{"7"}, sqlServerSuccessorRuntimeValues(t, r.targetDB, r.prepared.Group))
			newSQL := sqlServerSuccessorRuntimeRows(t, r.targetDB, r.prepared.Group)
			newReceipt := sqlServerSuccessorRuntimeEvidenceRows(t, r.targetDB, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.Len(t, newReceipt, len(priorBReceipts)+1)
			successorRuntimePreservedRows(t, priorBReceipts, newReceipt)
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, nextEffect).Scan(&localReceipts))
			require.Zero(t, localReceipts)
			faultsB.Lose.Store(false)
			attempt, err = sender.Deliver(t.Context(), nextEffect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, attempt.State)
			require.EqualValues(t, 1, commitsB.Load(), "B receipt replay must not commit another effect")
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, nextEffect).Scan(&localReceipts))
			require.Equal(t, 1, localReceipts)
			nextSettled, err := r.store.GetOutbox(t.Context(), nextEffect)
			require.NoError(t, err)
			require.Equal(t, nextOriginal.Destination, nextSettled.Destination)
			require.Equal(t, nextOriginal.PayloadDigest, nextSettled.PayloadDigest)
			require.Equal(t, nextOriginal.Payload, nextSettled.Payload)
			require.Equal(t, nextOriginal.Key, nextSettled.Key)
			require.Equal(t, newSQL, sqlServerSuccessorRuntimeRows(t, r.targetDB, r.prepared.Group))
			require.Equal(t, newReceipt, sqlServerSuccessorRuntimeEvidenceRows(t, r.targetDB, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
			require.Equal(t, committedSQL, sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group))
			require.Equal(t, receipt, sqlServerSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))

		})
	}
}
