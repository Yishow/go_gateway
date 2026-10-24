package grouppipeline

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func testCircuitIntegrationGap(t *testing.T) {
	for _, mode := range []struct {
		name  string
		build func(*testing.T) *longPipeline
	}{
		{"sqlite", newLongPipeline}, {"postgres", newLongPipelinePostgres}, {"timescale", newLongPipelineTimescale},
	} {
		t.Run(mode.name, func(t *testing.T) {
			f := mode.build(t)
			f.closeFirstWindow()
			f.deliver()
			_, baseline := f.baseline(0)
			f.closeUntil(at(1200))
			f.deliver()
			_, afterMissing := f.baseline(0)
			require.Equal(t, baseline, afterMissing, "the missing window cannot advance the baseline")
			missing := f.targetRows()[windowKey(longKwhFirst, 1)]
			require.Equal(t, "missing", missing.status.String)
			require.False(t, missing.raw.Valid)
			require.False(t, missing.delta.Valid)

			require.NoError(t, f.reading(0, "gap-first", at(29*60), firstStart+2))
			require.NoError(t, f.reading(1, "gap-second", at(29*60), secondStart+30))
			f.closeUntil(at(1800))
			f.deliver()
			rows := f.targetRows()
			require.Len(t, rows, 6)
			for i, source := range []string{longKwhFirst, longKwhSecond} {
				row := rows[windowKey(source, 2)]
				require.Equal(t, "gap", row.status.String)
				require.Equal(t, []string{"2", "30"}[i], row.delta.String)
				require.Equal(t, at(9*60).Format(longTime), row.deltaFrom.String)
				require.Equal(t, at(29*60).Format(longTime), row.deltaTo.String)
				var evidence struct {
					Counter struct {
						Span bool `json:"span_total"`
					} `json:"counter"`
				}
				require.NoError(t, json.Unmarshal([]byte(row.provenance.String), &evidence))
				require.True(t, evidence.Counter.Span)
				t.Logf("SQL source=%s status=%s delta=%s from=%s to=%s span_total=true", source, row.status.String, row.delta.String, row.deltaFrom.String, row.deltaTo.String)
			}
			f.requireBaselinesAt(2)
		})
	}
}
