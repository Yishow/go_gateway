package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

type firstStartTransition struct {
	name, kind             string
	sourceLong, targetLong bool
	mode                   workspace.WriteGroupStorageMode
}

func TestSQLiteSuccessorFirstStartActualReadback(t *testing.T) {
	for _, tc := range []firstStartTransition{{"wide-to-long", "sqlite", false, true, "plain"}, {"long-to-wide", "sqlite", true, false, "plain"}, {"wide-recreate", "sqlite", false, false, "plain"}, {"long-recreate", "sqlite", true, true, "plain"}} {
		t.Run(tc.name, func(t *testing.T) { firstStartActualReadback(t, tc) })
	}
}
func TestPostgresSuccessorFirstStartActualReadback(t *testing.T) {
	for _, tc := range postgresSuccessorRuntimeTransitions {
		mode := workspace.WriteGroupStorageModePlain
		if tc.timescale {
			mode = "timescale"
		}
		t.Run(tc.name, func(t *testing.T) {
			firstStartActualReadback(t, firstStartTransition{tc.name, "postgres", tc.sourceLong, tc.targetLong, mode})
		})
	}
}
func TestMySQLSuccessorFirstStartActualReadback(t *testing.T) {
	for _, tc := range mysqlSuccessorRuntimeTransitions {
		t.Run(tc.name, func(t *testing.T) {
			firstStartActualReadback(t, firstStartTransition{tc.name, tc.family, true, true, tc.mode})
		})
	}
}
func firstStartActualReadback(t *testing.T, tc firstStartTransition) {
	t.Helper()
	for _, counter := range []bool{false, true} {
		name := "gauge"
		if counter {
			name = "counter"
		}
		t.Run(name, func(t *testing.T) {
			r := newFirstStartAPIRig(t, tc.kind, tc.sourceLong, tc.targetLong, counter, tc.mode)
			request := r.confirmation(t)
			result := r.activationAPI().postActivation(t, request)
			require.Equal(t, "succeeded", result["status"])
			snapshot, err := r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary)
			require.NoError(t, err)
			r.snapshot = snapshot
			r.clock.Store(r.boundary.Add(6 * time.Second).UnixNano())
			require.NoError(t, r.pipe.Reconcile(t.Context()), "delayed reconciliation preserves explicit B")
			var value any = float32(7)
			if counter {
				value = uint64(107)
			}
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("first-at-boundary", r.boundary, value)))
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("first-after-boundary", r.boundary.Add(time.Millisecond), value)))
			window := 16 * time.Second
			if counter {
				window = 61 * time.Second
			}
			r.clock.Store(r.boundary.Add(window).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			r.deliver(t)
			var oldCount int
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=?`, r.groupID).Scan(&oldCount))
			require.Zero(t, oldCount, "G1 never acquired intake")
			if counter {
				rows := firstStartCounterRows(t, r)
				require.Len(t, rows, 1)
				require.Equal(t, "baseline", rows[0].status)
				require.Equal(t, "107", rows[0].raw)
				require.False(t, rows[0].hasDelta, "first G2 delta must be SQL NULL, never zero or a G1 carry")
				require.Empty(t, rows[0].delta)
				require.True(t, rows[0].deltaFrom.IsZero())
				require.True(t, rows[0].deltaTo.IsZero())
				later := r.boundary.Add(62 * time.Second)
				r.clock.Store(later.UnixNano())
				require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("next-110", later, uint64(110))))
				r.clock.Store(r.boundary.Add(121 * time.Second).UnixNano())
				r.pipe.TickAll(t.Context())
				require.Empty(t, r.reports)
				r.deliver(t)
				rows = firstStartCounterRows(t, r)
				require.Len(t, rows, 2)
				require.Equal(t, "ok", rows[1].status)
				require.Equal(t, "110", rows[1].raw)
				require.True(t, rows[1].hasDelta)
				require.Equal(t, "3", rows[1].delta)
			} else {
				require.Equal(t, []string{"7"}, firstStartValues(t, r))
			}
			before := successorAPIAllState(t, r.f.db)
			require.Equal(t, result, r.activationAPI().postActivation(t, request))
			require.Equal(t, result, successorStatus(t, r.successorAPIFixture)["first_start"].(map[string]any)["result"])
			require.Equal(t, before, successorAPIAllState(t, r.f.db), "post-delivery replay may not reset counter, runtime, or receipt state")
			for _, table := range []string{"write_group_cutovers", "write_group_revision_ends", "write_group_ownership_fences", "wg_delivery_sample_handoffs"} {
				var count int
				require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM `+table).Scan(&count))
				require.Zero(t, count, table)
			}
		})
	}
}
func firstStartCounterRows(t *testing.T, r *firstStartAPIRig) []counterWindow {
	t.Helper()
	g := r.prepared.Group
	switch r.kind {
	case "mysql", "mariadb":
		return mysqlSuccessorRuntimeCounterRows(t, &mysqlSuccessorRig{successorRuntimeRig: r.successorRuntimeRig, targetNamespace: r.targetNamespace, targetDB: r.target}, g)
	case "postgres":
		return postgresSuccessorRuntimeCounterRows(t, &postgresSuccessorRig{successorRuntimeRig: r.successorRuntimeRig, targetNamespace: r.targetNamespace}, g)
	default:
		return successorRuntimeCounterRows(t, r.successorRuntimeRig, g)
	}
}
func firstStartValues(t *testing.T, r *firstStartAPIRig) []string {
	t.Helper()
	switch r.kind {
	case "mysql", "mariadb":
		return mysqlSuccessorRuntimeValues(t, r.target, r.prepared.Group)
	case "postgres":
		return postgresSuccessorRuntimeValues(t, r.target, r.prepared.Group)
	default:
		return successorRuntimeValues(t, r.target, r.prepared.Group)
	}
}

// MySQL-family delivery must carry the production-frozen runtime catalog proof,
// exactly as the existing actual cutover fixture does; generic SQL is not enough.
func (r *firstStartAPIRig) deliver(t *testing.T) {
	t.Helper()
	if r.kind == "mysql" || r.kind == "mariadb" {
		(&mysqlSuccessorRig{successorRuntimeRig: r.successorRuntimeRig, targetNamespace: r.targetNamespace, targetDB: r.target}).deliver(t)
		return
	}
	r.successorRuntimeRig.deliver(t)
}
