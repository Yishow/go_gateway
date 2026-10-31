package dbtarget

import (
	"go-gateway/internal/datalink/snapshot"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedMeanProvenancePreservesRoundedEvidence(t *testing.T) {
	for _, mode := range []snapshot.Mode{snapshot.ModeMixedWide, snapshot.ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			cfg := mixedEncodeConfig(mode)
			cfg.Members = cfg.Members[:1]
			cfg.Counters = nil
			a, err := snapshot.NewAssembler(cfg)
			require.NoError(t, err)
			mixedEncodeOffer(t, a, "a", 1, 0)
			mixedEncodeOffer(t, a, "a", 2, 0)
			mixedEncodeOffer(t, a, "a", 3, 1)
			out := a.Plan(mixedEncodeAt(15)).Outcomes[0]
			require.True(t, out.Members[0].Reduction.Rounded)
			var row EncodedRow
			if mode == snapshot.ModeMixedWide {
				s := mixedWideSpec(SQLDialectSQLite)
				s.Members = s.Members[:1]
				l, issues := NewGroupRowLayout(s)
				require.Empty(t, issues)
				row, err = l.EncodeRow(out)
			} else {
				s := mixedLongSpec(t, SQLDialectSQLite)
				s.Members = s.Members[:1]
				l, issues := NewLongRowLayout(s)
				require.Empty(t, issues)
				row, err = l.EncodeRow(out)
			}
			require.NoError(t, err)
			require.Contains(t, cellMap(row)["provenance"], `"rounded":true`)
		})
	}
}
