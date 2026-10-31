package dbtarget

import (
	"fmt"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mixedEncodeAt(sec int) time.Time {
	return time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC).Add(time.Duration(sec) * time.Second)
}
func mixedEncodeBinding() aggregation.CounterBinding {
	return aggregation.CounterBinding{WorkspaceID: "w", GroupID: "g", ProfileRevision: "r", MemberKey: "c", MeasurementID: "energy", DefinitionRevision: "def", SourceBindingRevision: "source", SeriesEpoch: "epoch", Kind: aggregation.CounterKindMonotonic, Unit: "kWh"}
}
func mixedEncodeConfig(mode snapshot.Mode) snapshot.Config {
	return snapshot.Config{WorkspaceID: "w", GroupID: "g", GroupRevision: "r", DestinationScope: "db", Mode: mode, StartsAt: mixedEncodeAt(0), Storage: snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}, Counters: map[string]aggregation.CounterBinding{"c": mixedEncodeBinding()}, Members: []snapshot.Member{
		{Key: "a", Required: true, Calculation: snapshot.CalculationMean, Interval: 15 * time.Second, SourceValueType: schema.DataTypeInt64},
		{Key: "b", Required: true, Calculation: snapshot.CalculationMean, Interval: 15 * time.Second, SourceValueType: schema.DataTypeInt64},
		{Key: "c", Required: true, Calculation: snapshot.CalculationDelta, Interval: 60 * time.Second, SourceValueType: schema.DataTypeUint64},
	}}
}
func mixedEncodeOffer(t *testing.T, a *snapshot.Assembler, key string, sec int, n int64) {
	t.Helper()
	v := measurement.NewInt64(n)
	if key == "c" {
		v = measurement.NewUint64(uint64(n))
	}
	s := snapshot.Sample{MemberKey: key, SampleID: fmt.Sprintf("%s-%d", key, sec), ObservedAt: mixedEncodeAt(sec), ReceivedAt: mixedEncodeAt(sec).Add(time.Millisecond), Quality: schema.QualityGood, Value: v, SourceRevision: "src", MappingRevision: "map"}
	r := a.Offer(s, s.ObservedAt)
	require.Contains(t, []snapshot.OfferOutcome{snapshot.OfferSelected, snapshot.OfferNotSelected}, r.Outcome)
}
func mixedEncodeFixture(t *testing.T, mode snapshot.Mode) []snapshot.Outcome {
	t.Helper()
	a, err := snapshot.NewAssembler(mixedEncodeConfig(mode))
	require.NoError(t, err)
	mixedEncodeOffer(t, a, "c", 59, 1000)
	a.Commit(a.Plan(mixedEncodeAt(60)))
	for _, s := range []struct {
		k   string
		sec int
		n   int64
	}{{"a", 106, 20}, {"a", 110, 22}, {"a", 114, 24}, {"b", 106, 10}, {"b", 114, 14}, {"c", 119, 1008}} {
		mixedEncodeOffer(t, a, s.k, s.sec, s.n)
	}
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
func mixedWideSpec(dialect SQLDialect) GroupRowSpec {
	numeric := "TEXT"
	stamp := "TEXT"
	if dialect == SQLDialectPostgres {
		numeric = "NUMERIC(30,10)"
		stamp = "TIMESTAMPTZ"
	}
	b := mixedEncodeBinding()
	return GroupRowSpec{Dialect: dialect, ResultShape: MixedResultShapeV1, WorkspaceID: "w", GroupID: "g", GroupRevision: "r", Partial: true, RecordKeyColumn: "record_id", DueAtColumn: "due_at", ProvenanceColumn: "provenance", Columns: []ColumnInfo{{Name: "record_id", DataType: "TEXT", PrimaryKey: true}, {Name: "due_at", DataType: stamp}, {Name: "provenance", DataType: "TEXT"}, {Name: "a", DataType: numeric, Nullable: true}, {Name: "b", DataType: numeric, Nullable: true}, {Name: "c", DataType: numeric, Nullable: true}}, Members: []GroupRowMember{
		{MemberKey: "a", Column: "a", Type: measurement.ExactDecimal, Required: true, Calculation: snapshot.CalculationMean, Interval: 15 * time.Second, SourceValueType: schema.DataTypeInt64},
		{MemberKey: "b", Column: "b", Type: measurement.ExactDecimal, Required: true, Calculation: snapshot.CalculationMean, Interval: 15 * time.Second, SourceValueType: schema.DataTypeInt64},
		{MemberKey: "c", Column: "c", Type: measurement.ExactDecimal, Required: true, Calculation: snapshot.CalculationDelta, Interval: 60 * time.Second, SourceValueType: schema.DataTypeUint64, Counter: &b},
	}}
}
func mixedLongSpec(t *testing.T, dialect SQLDialect) LongRowSpec {
	t.Helper()
	d := "sqlite"
	if dialect == SQLDialectPostgres {
		d = "postgres"
	}
	s := LongRowSpec{Dialect: dialect, Columns: longColumnsFor(t, d), ResultShape: MixedResultShapeV1, WorkspaceID: "w", GroupID: "g", GroupRevision: "r", ProfileKind: "mixed"}
	for _, m := range mixedWideSpec(dialect).Members {
		s.Members = append(s.Members, LongRowMember{MemberKey: m.MemberKey, DeviceID: "device", PointID: "point-" + m.MemberKey, TagID: "tag-" + m.MemberKey, Type: m.Type, Required: m.Required, Calculation: m.Calculation, Interval: m.Interval, SourceValueType: m.SourceValueType, Counter: m.Counter})
	}
	return s
}
