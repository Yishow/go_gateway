package groupdelivery

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func reviewSetCells(row *dbtarget.EncodedRow, values map[string]any) {
	for i := range row.Cells {
		if value, ok := values[row.Cells[i].Column]; ok {
			row.Cells[i].Value = value
		}
	}
}

func TestSourceRecordsReviewPreservesARealCutoverStartWindow(t *testing.T) {
	store, _, _ := newTestStore(t)
	row, outcome := (longSpec{source: sourceA, status: "partial", slots: map[string]any{"value_float": 21.5}}).row(t)
	actualStart := outcome.BucketStart.Add(7500 * time.Millisecond)
	outcome.Truncated, outcome.WindowStart, outcome.WindowEnd = "start", actualStart, outcome.BucketEnd
	row.BucketStart = actualStart
	provenance, err := json.Marshal(map[string]any{
		"type": "float64", "member": row.EntityKey, "profile": map[string]any{"kind": "instantaneous", "interval": 15, "unit": "second"},
		"window": map[string]any{"truncated": "start", "start": timeText(actualStart), "end": timeText(outcome.BucketEnd),
			"nominal_start": timeText(outcome.BucketStart), "nominal_end": timeText(outcome.BucketEnd)},
	})
	require.NoError(t, err)
	reviewSetCells(&row, map[string]any{"bucket_start": timeText(actualStart), "provenance": string(provenance)})
	require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
		Destination: testDestination, NextClose: outcome.BucketEnd, Buckets: []ClosedBucket{{Outcome: outcome, Row: &row}},
	}))
	records, err := store.SourceRecords(t.Context(), scope, sourceA)
	require.NoError(t, err)
	require.NotNil(t, records.Latest)
	require.True(t, records.Latest.Available, "outbox keeps nominal bucket; SQL payload keeps actual cutover bounds: %+v", records.Latest)
	require.Equal(t, actualStart, records.Latest.WindowStart)
}

func TestSourceRecordsReviewRefusesContradictoryFrozenMetadata(t *testing.T) {
	for _, name := range []string{"missing_with_value", "baseline_with_delta", "foreign_provenance_member"} {
		t.Run(name, func(t *testing.T) {
			store, _, _ := newTestStore(t)
			spec := longSpec{source: sourceA, status: "ok", slots: map[string]any{"value_float": 21.5}}
			wantScope := scope
			if name == "missing_with_value" {
				spec.status = "missing"
				spec.reason = "no-samples"
			}
			if name == "baseline_with_delta" {
				spec.status, spec.profile, spec.delta = "baseline", "cumulative", true
				spec.slots = map[string]any{"value_decimal": "100"}
				wantScope.ProfileKind = "cumulative"
			}
			row, outcome := spec.row(t)
			if name == "baseline_with_delta" {
				provenance, err := json.Marshal(map[string]any{"type": "decimal", "member": row.EntityKey, "reason": "first-reading",
					"profile": map[string]any{"kind": "cumulative", "interval": 10, "unit": "minute"}})
				require.NoError(t, err)
				outcome.BucketEnd = at(600)
				reviewSetCells(&row, map[string]any{"value_type": "decimal", "bucket_end": timeText(at(600)), "delta_to": timeText(at(600)), "provenance": string(provenance)})
			}
			if name == "foreign_provenance_member" {
				provenance, err := json.Marshal(map[string]any{"type": "float64", "member": "a-different-source",
					"profile": map[string]any{"kind": "instantaneous", "interval": 15, "unit": "second"}})
				require.NoError(t, err)
				reviewSetCells(&row, map[string]any{"provenance": string(provenance)})
			}
			require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
				Destination: testDestination, NextClose: outcome.BucketEnd, Buckets: []ClosedBucket{{Outcome: outcome, Row: &row}},
			}))
			records, err := store.SourceRecords(t.Context(), wantScope, sourceA)
			require.NoError(t, err)
			require.NotNil(t, records.Latest)
			require.False(t, records.Latest.Available, "contradictory source metadata must be isolated as unavailable: %+v", records.Latest)
			require.Nil(t, records.Latest.VerifiedSourceRecord)
		})
	}
}
