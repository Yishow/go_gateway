package main

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestProductionManagedPartialRowsKeepNullAndQualityReasons(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedPipelineTarget(t, kind)
			e := f.env
			base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
			clock := &testClock{}
			clock.set(base.Add(-time.Second))
			e.services.writeGroups.WithClock(clock.now)
			group := f.prepare(t, []workspace.WriteGroupMember{
				{DeviceID: "device-1", PointID: "point-t-A", TagID: "tag-t-A", Required: true},
				{DeviceID: "device-1", PointID: "point-p-A", TagID: "tag-p-A"},
				{DeviceID: "device-1", PointID: "point-p-B", TagID: "tag-p-B"},
			}, workspace.WriteGroupRowPolicy{IntervalSeconds: 10, IncompletePolicy: "partial"})
			pipe := e.pipeline(clock, "node-1/managed-quality")
			clock.set(base.Add(time.Second))
			require.NoError(t, pipe.Reconcile(t.Context()))
			observed := base.Add(2 * time.Second)
			clock.set(observed)
			require.NoError(t, pipe.AcceptSample(t.Context(), e.envelope(group, 0, "good-required", observed, 21.5)))
			bad := e.envelope(group, 1, "bad-optional", observed, int64(999))
			bad.Quality, bad.QualityReason = schema.QualityBad, "read-failed"
			require.NoError(t, pipe.AcceptSample(t.Context(), bad))
			observed = base.Add(12 * time.Second)
			clock.set(observed)
			bad = e.envelope(group, 0, "bad-required", observed, 999.0)
			bad.Quality, bad.QualityReason = schema.QualityBad, "read-failed"
			require.NoError(t, pipe.AcceptSample(t.Context(), bad))
			clock.set(base.Add(31 * time.Second))
			pipe.TickAll(t.Context())
			for _, outcome := range []string{"row", "skipped", "no_data"} {
				var count int
				require.NoError(t, e.db.QueryRowContext(t.Context(), `SELECT count(*) FROM wg_delivery_buckets WHERE kind=?`, outcome).Scan(&count))
				require.Equal(t, 1, count, outcome)
			}
			var effect string
			require.NoError(t, e.db.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox`).Scan(&effect))
			store := groupdelivery.NewStore(e.db)
			sender := groupdelivery.NewSender(store, managedAckResolver{service: e.services.dbTarget, db: f.db, kind: f.kind}, groupdelivery.SenderConfig{})
			result, err := sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, result.State)
			var good float64
			var badValue, missingValue sql.NullInt64
			var raw string
			query := `SELECT "` + group.Members[0].TargetColumn + `","` + group.Members[1].TargetColumn + `","` + group.Members[2].TargetColumn + `",provenance FROM ` + managedSQLTable(f, group.Destination.TableName)
			require.NoError(t, f.db.QueryRowContext(t.Context(), query).Scan(&good, &badValue, &missingValue, &raw))
			require.Equal(t, 21.5, good)
			require.False(t, badValue.Valid, "a bad value must never become zero or an earlier good reading")
			require.False(t, missingValue.Valid)
			var provenance []managedProvenance
			require.NoError(t, json.Unmarshal([]byte(raw), &provenance))
			require.Len(t, provenance, 3)
			require.Equal(t, "bad", provenance[1].Status)
			require.Equal(t, "bad", provenance[1].Quality)
			require.Equal(t, "read-failed", provenance[1].Reason)
			require.Equal(t, base.Add(2*time.Second).Format(time.RFC3339Nano), provenance[1].ObservedAt)
			require.Equal(t, "missing", provenance[2].Status)
			require.NotEmpty(t, provenance[2].Reason)
			require.Empty(t, provenance[2].Quality)
			require.Empty(t, provenance[2].ObservedAt)
		})
	}
}
