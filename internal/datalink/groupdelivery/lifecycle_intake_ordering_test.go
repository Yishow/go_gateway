package groupdelivery

import (
	"context"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

// Driver barriers arrange both writer orders without timing sleeps or retries.
func TestCanonicalLifecycleAppendAndStopLinearize(t *testing.T) {
	for _, status := range []string{"disabled", "deleted"} {
		t.Run(status+"-first", func(t *testing.T) {
			_, db, path := newTestStore(t)
			seedLifecycleIntake(t, db, testKey)
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			tx, err := db.BeginTx(ctx, nil)
			require.NoError(t, err)
			defer tx.Rollback()
			_, err = tx.ExecContext(ctx, `UPDATE write_groups SET status=? WHERE id=?`, status, testKey.GroupID)
			require.NoError(t, err)
			attempted := make(chan struct{})
			faults := &lostcommit.Faults{}
			hook := func(query string) error {
				if strings.Contains(query, "UPDATE wg_delivery_checkpoints") {
					close(attempted)
				}
				return nil
			}
			faults.Reject.Store(&hook)
			appends := lostcommit.OpenFaults(path+"?_pragma=busy_timeout(5000)", faults)
			defer appends.Close()
			result := make(chan error, 1)
			go func() {
				result <- NewStore(appends).WithCanonicalIntake().AppendSample(ctx, testKey, t0, sampleAt("new", "temp", 3, 30))
			}()
			awaitLifecycleBarrier(t, ctx, attempted)
			require.NoError(t, tx.Commit())
			err = <-result
			var typed *OwnershipError
			require.ErrorAs(t, err, &typed)
			require.Equal(t, "recording-intake-stopped", err.Error())
			require.Zero(t, countRows(t, db, "wg_delivery_samples"))
		})
		t.Run("append-before-"+status, func(t *testing.T) {
			_, db, path := newTestStore(t)
			seedLifecycleIntake(t, db, testKey)
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			atInsert, release := make(chan struct{}), make(chan struct{})
			appendFaults := &lostcommit.Faults{}
			hook := func(query string) error {
				if strings.Contains(query, "INSERT INTO wg_delivery_samples") {
					close(atInsert)
					select {
					case <-release:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				return nil
			}
			appendFaults.Reject.Store(&hook)
			appends := lostcommit.OpenFaults(path+"?_pragma=busy_timeout(5000)", appendFaults)
			defer appends.Close()
			store := NewStore(appends).WithCanonicalIntake()
			sample := sampleAt("accepted", "temp", 3, 30)
			appendResult := make(chan error, 1)
			stopFinished := make(chan struct{})
			// Delay the ACK until after stop commits: its earlier durable acceptance
			// must remain a success even though lifecycle has since changed.
			afterCommit := func() {
				select {
				case <-stopFinished:
				case <-ctx.Done():
				}
			}
			appendFaults.AfterCommit.Store(&afterCommit)
			go func() { appendResult <- store.AppendSample(ctx, testKey, t0, sample) }()
			awaitLifecycleBarrier(t, ctx, atInsert)
			stopAttempted := make(chan struct{})
			stopFaults := &lostcommit.Faults{}
			stopHook := func(query string) error {
				if strings.Contains(query, "UPDATE write_groups") {
					close(stopAttempted)
				}
				return nil
			}
			stopFaults.Reject.Store(&stopHook)
			stops := lostcommit.OpenFaults(path+"?_pragma=busy_timeout(5000)", stopFaults)
			defer stops.Close()
			stopResult := make(chan error, 1)
			go func() {
				_, err := stops.ExecContext(ctx, `UPDATE write_groups SET status=? WHERE id=?`, status, testKey.GroupID)
				stopResult <- err
				close(stopFinished)
			}()
			awaitLifecycleBarrier(t, ctx, stopAttempted)
			select {
			case <-stopFinished:
				t.Fatal("lifecycle committed through the in-flight journal transaction")
			default:
			}
			close(release)
			require.NoError(t, <-appendResult)
			require.NoError(t, <-stopResult)
			appendFaults.Reject.Store(nil)
			require.NoError(t, store.AppendSample(ctx, testKey, t0, sample))
			require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
			err := store.AppendSample(ctx, testKey, t0, sampleAt("later", "temp", 4, 40))
			var typed *OwnershipError
			require.ErrorAs(t, err, &typed)
		})
	}
}

func awaitLifecycleBarrier(t *testing.T, ctx context.Context, barrier <-chan struct{}) {
	t.Helper()
	select {
	case <-barrier:
	case <-ctx.Done():
		t.Fatal("lifecycle test barrier did not complete:", ctx.Err())
	}
}
