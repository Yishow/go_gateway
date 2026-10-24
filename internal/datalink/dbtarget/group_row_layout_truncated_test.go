package dbtarget

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

// A wide window a cutover cut short is written, not dropped, and says so. A
// table with a provenance column carries the evidence in it; a custom table
// without one gets no new column and no altered shape, and the cut is still
// told by the encoded row and by the real window start.
func TestWideRowEncodesAWindowACutoverCutShort(t *testing.T) {
	cut := func() (start, end time.Time) {
		return layoutBucket.Add(4 * time.Second), layoutBucket.Add(10 * time.Second)
	}

	t.Run("with a provenance column the evidence is in it and other rows keep their JSON", func(t *testing.T) {
		spec := baseSpec()
		spec.ProvenanceColumn, spec.RecordKeyColumn, spec.BucketStartColumn = "provenance", "record_id", "bucket_start"
		spec.EntityKeyed, spec.EntityKeyColumn = true, "entity"
		layout, issues := NewGroupRowLayout(spec)
		require.Empty(t, issues)
		outcome := rowOutcome(okMember("temp", measurement.NewText("a")), okMember("run", measurement.NewBool(true)), okMember("cnt", measurement.NewUint64(1)))
		outcome.EntityKey = "plant-1"

		whole, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		require.False(t, whole.Partial)
		require.NotContains(t, cellOf(whole, "provenance"), "window", "a complete row has exactly the JSON it always had")

		outcome.WindowStart, outcome.WindowEnd = cut()
		outcome.Truncated = "start"
		encoded, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		require.True(t, encoded.Partial, "the encoded row says it is cut")
		require.Equal(t, outcome.WindowStart.UTC().Format(time.RFC3339Nano), cellOf(encoded, "bucket_start"), "the row begins where the cutover put it")
		var entries []struct {
			Window map[string]any `json:"window"`
		}
		require.NoError(t, json.Unmarshal([]byte(cellOf(encoded, "provenance").(string)), &entries))
		require.Len(t, entries, 3)
		for _, entry := range entries {
			require.Equal(t, "start", entry.Window["truncated"])
			require.Equal(t, outcome.BucketStart.UTC().Format(longTimeText), entry.Window["nominal_start"])
		}
	})

	t.Run("without a provenance column nothing is added and the cut is still told", func(t *testing.T) {
		spec := baseSpec()
		spec.RecordKeyColumn, spec.BucketStartColumn = "record_id", "bucket_start"
		layout, issues := NewGroupRowLayout(spec)
		require.Empty(t, issues)
		outcome := rowOutcome(okMember("temp", measurement.NewText("a")), okMember("run", measurement.NewBool(true)), okMember("cnt", measurement.NewUint64(1)))
		outcome.WindowStart, outcome.WindowEnd = cut()
		outcome.Truncated = "start"
		encoded, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		require.True(t, encoded.Partial)
		for _, cell := range encoded.Cells {
			require.NotEqual(t, "provenance", cell.Column, "no column is invented")
		}
		require.Equal(t, outcome.WindowStart.UTC().Format(time.RFC3339Nano), cellOf(encoded, "bucket_start"))
	})
}

func cellOf(row EncodedRow, column string) any {
	for _, cell := range row.Cells {
		if cell.Column == column {
			return cell.Value
		}
	}
	return nil
}
