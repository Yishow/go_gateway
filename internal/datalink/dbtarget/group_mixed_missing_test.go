package dbtarget

import (
	"encoding/json"
	"fmt"
	"go-gateway/internal/datalink/snapshot"
	"testing"

	"github.com/stretchr/testify/require"
)

// Actual assembler output, not a forged CounterGap, exercises the encoder's
// distinction between a current-window value and retained span evidence.
func TestMixedUnavailableResultNeverFallsBackToRawOrSpan(t *testing.T) {
	for _, mode := range []snapshot.Mode{snapshot.ModeMixedWide, snapshot.ModeMixedLong} {
		for _, required := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/required=%t", mode, required), func(t *testing.T) {
				cfg := mixedEncodeConfig(mode)
				cfg.Members[2].Required = required
				assembler, err := snapshot.NewAssembler(cfg)
				require.NoError(t, err)
				ws := mixedWideSpec(SQLDialectSQLite)
				ws.Members[2].Required = required
				wide, issues := NewGroupRowLayout(ws)
				require.Empty(t, issues)
				ls := mixedLongSpec(t, SQLDialectSQLite)
				ls.Members[2].Required = required
				long, issues := NewLongRowLayout(ls)
				require.Empty(t, issues)
				check := func(p snapshot.Plan, due int, reason string, span bool) {
					t.Helper()
					require.NoError(t, p.Err)
					found := false
					for _, out := range p.Outcomes {
						if !out.DueAt.Equal(mixedEncodeAt(due)) {
							continue
						}
						if mode == snapshot.ModeMixedLong && out.EntityKey != "c" {
							continue
						}
						found = true
						var row EncodedRow
						var err error
						if mode == snapshot.ModeMixedWide {
							row, err = wide.EncodeRow(out)
						} else {
							row, err = long.EncodeRow(out)
						}
						if required {
							require.Equal(t, snapshot.OutcomeSkipped, out.Kind)
							require.Error(t, err)
							continue
						}
						require.NoError(t, err)
						cells := cellMap(row)
						var evidence map[string]any
						if mode == snapshot.ModeMixedWide {
							require.Nil(t, cells["c"])
							var entries []map[string]any
							require.NoError(t, json.Unmarshal([]byte(cells["provenance"].(string)), &entries))
							evidence = entries[2]
						} else {
							require.Nil(t, cells["value_decimal"])
							require.Nil(t, cells["delta_value"])
							require.NoError(t, json.Unmarshal([]byte(cells["provenance"].(string)), &evidence))
						}
						require.Equal(t, reason, evidence["reason"])
						if span {
							counter := evidence["counter"].(map[string]any)
							require.Equal(t, "8", counter["delta"])
							require.Equal(t, true, counter["span_total"])
						}
					}
					require.True(t, found)
				}
				mixedEncodeOffer(t, assembler, "a", 55, 20)
				mixedEncodeOffer(t, assembler, "b", 55, 12)
				mixedEncodeOffer(t, assembler, "c", 59, 1000)
				first := assembler.Plan(mixedEncodeAt(60))
				check(first, 60, "first-reading", false)
				assembler.Commit(first)
				mixedEncodeOffer(t, assembler, "a", 115, 20)
				mixedEncodeOffer(t, assembler, "b", 115, 12)
				missing := assembler.Plan(mixedEncodeAt(120))
				check(missing, 120, "no-samples", false)
				assembler.Commit(missing)
				mixedEncodeOffer(t, assembler, "a", 175, 20)
				mixedEncodeOffer(t, assembler, "b", 175, 12)
				mixedEncodeOffer(t, assembler, "c", 179, 1008)
				gap := assembler.Plan(mixedEncodeAt(180))
				check(gap, 180, "gap", true)
			})
		}
	}
}
