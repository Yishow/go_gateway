package dbtarget

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
	"testing"
)

func mixedScopedRows(t *testing.T, mode snapshot.Mode, scope string) []snapshot.Outcome {
	t.Helper()
	cfg := mixedEncodeConfig(mode)
	switch scope {
	case "workspace":
		cfg.WorkspaceID = "foreign"
	case "group":
		cfg.GroupID = "foreign"
	case "revision":
		cfg.GroupRevision = "foreign"
	}
	binding := cfg.Counters["c"]
	binding.WorkspaceID = cfg.WorkspaceID
	binding.GroupID = cfg.GroupID
	binding.ProfileRevision = cfg.GroupRevision
	cfg.Counters["c"] = binding
	a, err := snapshot.NewAssembler(cfg)
	require.NoError(t, err)
	mixedEncodeOffer(t, a, "c", 59, 1000)
	a.Commit(a.Plan(mixedEncodeAt(60)))
	mixedEncodeOffer(t, a, "a", 114, 22)
	mixedEncodeOffer(t, a, "b", 114, 12)
	mixedEncodeOffer(t, a, "c", 119, 1008)
	p := a.Plan(mixedEncodeAt(120))
	require.NoError(t, p.Err)
	var rows []snapshot.Outcome
	for _, o := range p.Outcomes {
		if o.DueAt.Equal(mixedEncodeAt(120)) {
			rows = append(rows, o)
		}
	}
	return rows
}
func TestMixedEncoderRejectsForeignStableIdentity(t *testing.T) {
	for _, mode := range []snapshot.Mode{snapshot.ModeMixedWide, snapshot.ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			wide, issues := NewGroupRowLayout(mixedWideSpec(SQLDialectSQLite))
			require.Empty(t, issues)
			long, issues := NewLongRowLayout(mixedLongSpec(t, SQLDialectSQLite))
			require.Empty(t, issues)
			encode := func(o snapshot.Outcome) error {
				if mode == snapshot.ModeMixedWide {
					_, err := wide.EncodeRow(o)
					return err
				}
				_, err := long.EncodeRow(o)
				return err
			}
			local := mixedScopedRows(t, mode, "")
			for _, o := range local {
				require.NoError(t, encode(o))
			}
			for _, scope := range []string{"workspace", "group", "revision"} {
				t.Run(scope, func(t *testing.T) {
					for i, o := range mixedScopedRows(t, mode, scope) {
						require.NotEqual(t, local[i].RecordID, o.RecordID)
						require.Error(t, encode(o))
					}
				})
			}
			o := local[0]
			o.RecordID = mixedSparseOutcomes(t, mode)[0].RecordID
			require.Error(t, encode(o), "a valid identity from a different dispatch cannot authorize this result")
			o = local[0]
			o.RecordID = "nonempty-but-not-this-dispatch"
			require.Error(t, encode(o))
			// Correct values cannot authorize borrowing another member's stable identity.
			if mode == snapshot.ModeMixedLong {
				o = local[0]
				o.RecordID = local[1].RecordID
				require.Error(t, encode(o))
			}
		})
	}
}
