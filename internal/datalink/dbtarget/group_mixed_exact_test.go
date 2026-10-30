package dbtarget

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"math"
	"testing"
)

func mixedExactOutcome(t *testing.T, mode snapshot.Mode) []snapshot.Outcome {
	t.Helper()
	a, err := snapshot.NewAssembler(mixedEncodeConfig(mode))
	require.NoError(t, err)
	mixedEncodeOffer(t, a, "c", 59, 9007199254740992)
	a.Commit(a.Plan(mixedEncodeAt(60)))
	mixedEncodeOffer(t, a, "a", 106, 9007199254740993)
	mixedEncodeOffer(t, a, "a", 114, 9007199254740995)
	mixedEncodeOffer(t, a, "b", 106, 1)
	mixedEncodeOffer(t, a, "b", 114, 2)
	mixedEncodeOffer(t, a, "c", 119, 9007199254740993)
	p := a.Plan(mixedEncodeAt(120))
	require.NoError(t, p.Err)
	var rows []snapshot.Outcome
	for _, row := range p.Outcomes {
		if row.DueAt.Equal(mixedEncodeAt(120)) {
			rows = append(rows, row)
		}
	}
	return rows
}
func TestMixedExactDecimalAndLargeIntegerDialectParams(t *testing.T) {
	for _, dialect := range []SQLDialect{SQLDialectSQLite, SQLDialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			l, issues := NewGroupRowLayout(mixedWideSpec(dialect))
			require.Empty(t, issues)
			row, err := l.EncodeRow(mixedExactOutcome(t, snapshot.ModeMixedWide)[0])
			require.NoError(t, err)
			cells := cellMap(row)
			require.Equal(t, "9007199254740994", cells["a"])
			require.Equal(t, "1.5", cells["b"])
			require.Equal(t, "1", cells["c"])
			long, issues := NewLongRowLayout(mixedLongSpec(t, dialect))
			require.Empty(t, issues)
			for i, out := range mixedExactOutcome(t, snapshot.ModeMixedLong) {
				row, err := long.EncodeRow(out)
				require.NoError(t, err)
				require.Equal(t, []string{"9007199254740994", "1.5", "1"}[i], cellMap(row)["value_decimal"])
			}
		})
	}
	s := mixedWideSpec(SQLDialectPostgres)
	s.Columns[3].DataType = "NUMERIC(18,4)"
	l, issues := NewGroupRowLayout(s)
	require.Empty(t, issues)
	_, err := l.EncodeRow(mixedExactOutcome(t, snapshot.ModeMixedWide)[0])
	require.Equal(t, rowErrSQLValueBlock, rowCode(t, err), "precision loss must fail, not round a large value")
	s = mixedWideSpec(SQLDialectSQLite)
	s.Columns[3].DataType = "REAL"
	_, issues = NewGroupRowLayout(s)
	require.Contains(t, issueCodes(issues), issueUnsupportedSQLType)
}
func TestMixedLatestUint64NeverPassesThroughFloat(t *testing.T) {
	c := mixedEncodeConfig(snapshot.ModeMixedWide)
	c.Counters = nil
	c.Members = c.Members[:1]
	c.Members[0].Calculation = snapshot.CalculationLatest
	c.Members[0].SourceValueType = schema.DataTypeUint64
	a, err := snapshot.NewAssembler(c)
	require.NoError(t, err)
	sample := snapshot.Sample{MemberKey: "a", SampleID: "max", ObservedAt: mixedEncodeAt(14), Quality: schema.QualityGood, Value: measurement.NewUint64(math.MaxUint64)}
	require.Equal(t, snapshot.OfferSelected, a.Offer(sample, sample.ObservedAt).Outcome)
	out := a.Plan(mixedEncodeAt(15)).Outcomes[0]
	for _, dialect := range []SQLDialect{SQLDialectSQLite, SQLDialectPostgres} {
		s := mixedWideSpec(dialect)
		s.Members = s.Members[:1]
		s.Members[0].Calculation = snapshot.CalculationLatest
		s.Members[0].SourceValueType = schema.DataTypeUint64
		s.Members[0].Type = measurement.ExactUint64
		l, issues := NewGroupRowLayout(s)
		require.Empty(t, issues)
		row, err := l.EncodeRow(out)
		require.NoError(t, err)
		require.Equal(t, "18446744073709551615", cellMap(row)["a"])
		if dialect == SQLDialectSQLite {
			target := newInspectionSQLiteTarget(t, `CREATE TABLE mixed_results(record_id TEXT PRIMARY KEY NOT NULL,due_at TEXT NOT NULL,provenance TEXT NOT NULL,a TEXT,b TEXT,c TEXT)`)
			db, err := sql.Open("sqlite", target)
			require.NoError(t, err)
			t.Cleanup(func() { _ = db.Close() })
			mixedInsert(t, db, "mixed_results", row)
			var value, kind string
			require.NoError(t, db.QueryRow(`SELECT a,typeof(a) FROM mixed_results`).Scan(&value, &kind))
			require.Equal(t, "18446744073709551615", value)
			require.Equal(t, "text", kind)
		}

	}
}
