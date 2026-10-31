package dbtarget

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedWideLayoutRequiresVersionedFreshMappings(t *testing.T) {
	cases := map[string]func(*GroupRowSpec){
		"unknown shape":                func(s *GroupRowSpec) { s.ResultShape = "future" },
		"unversioned method":           func(s *GroupRowSpec) { s.ResultShape = "" },
		"legacy time reinterpretation": func(s *GroupRowSpec) { s.BucketStartColumn = "due_at"; s.DueAtColumn = "" },
		"required sparse not nullable": func(s *GroupRowSpec) { s.Columns[3].Nullable = false },
		"no provenance":                func(s *GroupRowSpec) { s.ProvenanceColumn = "" },
		"no partial proof":             func(s *GroupRowSpec) { s.Partial = false },
		"omitted method":               func(s *GroupRowSpec) { s.Members[0].Calculation = "" },
		"legacy reduction conflict":    func(s *GroupRowSpec) { s.Members[0].Reduction = "mean" },
		"missing cadence":              func(s *GroupRowSpec) { s.Members[0].Interval = 0 },
		"wrong result type":            func(s *GroupRowSpec) { s.Members[0].Type = measurement.ExactInt64; s.Columns[3].DataType = "INTEGER" },
		"unbound delta":                func(s *GroupRowSpec) { s.Members[2].Counter = nil },
		"foreign binding":              func(s *GroupRowSpec) { s.Members[2].Counter.GroupID = "foreign" },
		"duplicate member":             func(s *GroupRowSpec) { s.Members[1].MemberKey = "a" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			s := mixedWideSpec(SQLDialectSQLite)
			change(&s)
			l, issues := NewGroupRowLayout(s)
			require.Nil(t, l)
			require.NotEmpty(t, issues)
		})
	}
}
func TestMixedLongLayoutRejectsUnprovedShapesAndGroupWindow(t *testing.T) {
	for _, change := range []func(*LongRowSpec){func(s *LongRowSpec) { s.ResultShape = "" }, func(s *LongRowSpec) { s.ResultShape = "future" }, func(s *LongRowSpec) { s.ProfileInterval = 60; s.ProfileUnit = "second" }, func(s *LongRowSpec) { s.ProfileKind = "cumulative" }, func(s *LongRowSpec) { s.Dialect = SQLDialectMySQL }, func(s *LongRowSpec) { s.Dialect = SQLDialectSQLServer }, func(s *LongRowSpec) { s.Members[0].Calculation = "" }} {
		s := mixedLongSpec(t, SQLDialectSQLite)
		change(&s)
		l, issues := NewLongRowLayout(s)
		require.Nil(t, l)
		require.NotEmpty(t, issues)
	}
}
func TestMixedEncoderRejectsMismatchedOrForgedOutcome(t *testing.T) {
	wide, issues := NewGroupRowLayout(mixedWideSpec(SQLDialectSQLite))
	require.Empty(t, issues)
	cases := map[string]func(*snapshot.Outcome){
		"no due marker":             func(o *snapshot.Outcome) { o.DueAt = time.Time{} },
		"fake common window":        func(o *snapshot.Outcome) { o.BucketStart = mixedEncodeAt(60); o.BucketEnd = mixedEncodeAt(120) },
		"fake common actual window": func(o *snapshot.Outcome) { o.WindowStart = mixedEncodeAt(60); o.WindowEnd = mixedEncodeAt(120) },
		"duplicate result":          func(o *snapshot.Outcome) { o.Members = append(o.Members, o.Members[0]) },
		"changed method":            func(o *snapshot.Outcome) { o.Members[0].Window.Calculation = snapshot.CalculationLatest },
		"changed cadence":           func(o *snapshot.Outcome) { o.Members[0].Window.Interval = time.Minute },
		"wrong window":              func(o *snapshot.Outcome) { o.Members[0].Window.AlignedStart = mixedEncodeAt(60) },
		"required missing forged row": func(o *snapshot.Outcome) {
			o.Members[2].Status = snapshot.MemberMissing
			o.Members[2].Reason = "missing"
		},
		"mean raw fallback":     func(o *snapshot.Outcome) { o.Members[0].Reduction = nil },
		"wrong raw source type": func(o *snapshot.Outcome) { o.Members[0].Sample.Value = measurement.NewText("22") },
		"usable bad quality":    func(o *snapshot.Outcome) { o.Members[0].Sample.Quality = "bad" },
		"mean version drift":    func(o *snapshot.Outcome) { o.Members[0].Reduction.Version = "future" },
		"span as usable":        func(o *snapshot.Outcome) { o.Members[2].Counter.Delta.Span = true },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			o := mixedEncodeFixture(t, snapshot.ModeMixedWide)[0]
			mutate(&o)
			_, err := wide.EncodeRow(o)
			require.Error(t, err)
		})
	}
	long, issues := NewLongRowLayout(mixedLongSpec(t, SQLDialectSQLite))
	require.Empty(t, issues)
	o := mixedEncodeFixture(t, snapshot.ModeMixedLong)[0]
	o.Members[0].Status = snapshot.MemberNotDue
	_, err := long.EncodeRow(o)
	require.Error(t, err, "long not_due must never produce a row")
	o = mixedEncodeFixture(t, snapshot.ModeMixedLong)[0]
	o.Truncated = "start"
	o.WindowStart = mixedEncodeAt(115)
	o.WindowEnd = mixedEncodeAt(120)
	_, err = long.EncodeRow(o)
	require.Error(t, err, "SQL coverage must match the member's actual bounds")
	legacy := mixedWideSpec(SQLDialectSQLite)
	legacy.ResultShape = ""
	legacy.DueAtColumn = ""
	for i := range legacy.Members {
		legacy.Members[i].Calculation = ""
		legacy.Members[i].Interval = 0
		legacy.Members[i].Reduction = "mean"
	}
	l, issues := NewGroupRowLayout(legacy)
	require.Empty(t, issues)
	_, err = l.EncodeRow(mixedEncodeFixture(t, snapshot.ModeMixedWide)[0])
	require.Error(t, err, "legacy encoder cannot consume mixed evidence")
}
