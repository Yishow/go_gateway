package api

import (
	"context"
	"database/sql"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/testutil/lostcommit"
)

// The real connector is resolved first, retaining frozen revision validation;
// only its transport is wrapped to lose the response after PostgreSQL commits.
type postgresSuccessorLostACKResolver struct {
	rig *postgresSuccessorRig
	db  *sql.DB
}

func (r postgresSuccessorLostACKResolver) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	target, err := r.rig.Resolve(ctx, item)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	if target.Close != nil {
		if err := target.Close(); err != nil {
			return groupdelivery.Target{}, err
		}
	}
	target.DB, target.Close = r.db, nil
	return target, nil
}

func TestPostgresSuccessorActivationOldBacklogAndCommitACKLost(t *testing.T) {
	for _, transition := range postgresSuccessorRuntimeTransitions {
		t.Run(transition.name, func(t *testing.T) {
			r := newPostgresSuccessorRuntimeRig(t, transition.sourceLong, transition.targetLong, false, transition.timescale)
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
			require.Empty(t, postgresSuccessorRuntimeRows(t, r.target, r.snapshot.Group))
			r.clock.Store(r.boundary.Add(-2500 * time.Millisecond).UnixNano())
			r.activate(t)
			r.clock.Store(r.boundary.Add(6 * time.Second).UnixNano())
			require.NoError(t, r.pipe.Reconcile(t.Context()))
			frozen, err := r.store.GetOutbox(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, original, frozen, "head replacement must not touch finalized backlog")

			lose := &atomic.Bool{}
			dsn := os.Getenv("POSTGRES_DSN")
			if transition.timescale {
				dsn = os.Getenv("TIMESCALE_POSTGRES_DSN")
			}
			wrapped := lostcommit.Wrap(stdlib.GetDefaultDriver(), dsn, lose)
			t.Cleanup(func() { require.NoError(t, wrapped.Close()) })
			now := func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }
			sender := groupdelivery.NewSender(r.store, postgresSuccessorLostACKResolver{r, wrapped}, groupdelivery.SenderConfig{
				Owner: "pg-successor-lost-ack", Now: now,
				Backoff: func(int) (string, time.Time) { return groupdelivery.StateRetrying, now().Add(-time.Second) },
			})
			lose.Store(true)
			attempt, err := sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateRetrying, attempt.State)
			require.Equal(t, []string{"4"}, postgresSuccessorRuntimeValues(t, r.target, r.snapshot.Group))
			require.Empty(t, postgresSuccessorRuntimeRows(t, r.target, r.prepared.Group), "old backlog must never be redirected to B")
			committedSQL := postgresSuccessorRuntimeRows(t, r.target, r.snapshot.Group)
			receipt := postgresSuccessorRuntimeEvidenceRows(t, r.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.Len(t, receipt, 1)
			var localReceipts int
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, effect).Scan(&localReceipts))
			require.Zero(t, localReceipts)
			lose.Store(false)
			attempt, err = sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, attempt.State)
			require.Equal(t, committedSQL, postgresSuccessorRuntimeRows(t, r.target, r.snapshot.Group))
			require.Equal(t, receipt, postgresSuccessorRuntimeEvidenceRows(t, r.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, effect).Scan(&localReceipts))
			require.Equal(t, 1, localReceipts)
			settled, err := r.store.GetOutbox(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, original.Destination, settled.Destination)
			require.Equal(t, original.Payload, settled.Payload)
			require.Equal(t, original.PayloadDigest, settled.PayloadDigest)
			require.Equal(t, original.Key, settled.Key)
			require.Empty(t, postgresSuccessorRuntimeRows(t, r.target, r.prepared.Group))
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
			lose.Store(true)
			attempt, err = sender.Deliver(t.Context(), nextEffect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateRetrying, attempt.State)
			require.Equal(t, []string{"7"}, postgresSuccessorRuntimeValues(t, r.target, r.prepared.Group))
			newSQL := postgresSuccessorRuntimeRows(t, r.target, r.prepared.Group)
			newReceipt := postgresSuccessorRuntimeEvidenceRows(t, r.target, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.Len(t, newReceipt, 1)
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, nextEffect).Scan(&localReceipts))
			require.Zero(t, localReceipts)
			lose.Store(false)
			attempt, err = sender.Deliver(t.Context(), nextEffect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, attempt.State)
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=?`, nextEffect).Scan(&localReceipts))
			require.Equal(t, 1, localReceipts)
			require.Equal(t, newSQL, postgresSuccessorRuntimeRows(t, r.target, r.prepared.Group))
			require.Equal(t, newReceipt, postgresSuccessorRuntimeEvidenceRows(t, r.target, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
			require.Equal(t, committedSQL, postgresSuccessorRuntimeRows(t, r.target, r.snapshot.Group))
			require.Equal(t, receipt, postgresSuccessorRuntimeEvidenceRows(t, r.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))

		})
	}
}
