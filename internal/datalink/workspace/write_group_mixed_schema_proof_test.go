package workspace

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func mixedSchemaFixture() (*WriteGroup, map[string]schema.DataType) {
	g := &WriteGroup{ID: "g", WorkspaceID: "w", Revision: "r", RecordID: "2f629763-f426-460f-a935-552d3a2e1ba8", LayoutVersion: WriteGroupLayoutManagedWideProfileV1,
		RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1},
		Destination:      WriteGroupDestination{ConnectorID: "db", ConnectorRevision: "db-r", Database: "target", TableSchema: "main", TableName: "mixed", StorageStrategy: WriteGroupStorageStrategyManaged},
		RowPolicy:        WriteGroupRowPolicy{IncompletePolicy: "partial", DueAtColumn: "due_at", RecordKeyColumn: "record_id", ProvenanceColumn: "provenance"},
		Members:          []WriteGroupMember{{DeviceID: "d", PointID: "p", TagID: "t", TargetColumn: "v", Required: true, Calculation: WriteGroupMemberCalculationMean, Cadence: &WriteGroupMemberCadence{Interval: 15, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}, SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactDecimal}}}
	return g, map[string]schema.DataType{"t": schema.DataTypeInt64}
}
func TestMixedConfigurationProofBindsAllCalculationSemantics(t *testing.T) {
	g, types := mixedSchemaFixture()
	inspection := &dbtarget.TableInspection{Columns: []dbtarget.ColumnInfo{{Name: "v", DataType: "TEXT", Nullable: true}}}
	source, err := managedSourceDigest(g, types)
	require.NoError(t, err)
	catalog, err := digestWriteGroupInspection(g, inspection)
	require.NoError(t, err)
	tests := map[string]func(*WriteGroup){
		"method":               func(g *WriteGroup) { g.Members[0].Calculation = WriteGroupMemberCalculationLatest },
		"cadence":              func(g *WriteGroup) { g.Members[0].Cadence.Interval = 60 },
		"required":             func(g *WriteGroup) { g.Members[0].Required = false },
		"result type":          func(g *WriteGroup) { g.Members[0].ResultValueType = measurement.ExactInt64 },
		"record identity":      func(g *WriteGroup) { g.RecordID = "3f629763-f426-460f-a935-552d3a2e1ba8" },
		"result shape":         func(g *WriteGroup) { g.RecordingProfile.ResultShape = "future" },
		"measurement identity": func(g *WriteGroup) { id := "new-binding"; g.Members[0].MeasurementID = &id },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			changed := cloneWriteGroup(g)
			mutate(changed)
			next, err := managedSourceDigest(changed, types)
			if name == "result shape" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEqual(t, source, next, "schema consent source digest")
			next, err = digestWriteGroupInspection(changed, inspection)
			require.NoError(t, err)
			require.NotEqual(t, catalog, next, "readiness digest")
			require.False(t, sameWriteGroupSchema(changed, g), "old saved confirmation must clear")
		})
	}
}
func TestMixedSchemaUsesFrozenResultNotRawType(t *testing.T) {
	for dialect, want := range map[string]string{"sqlite": "TEXT", "postgres": "NUMERIC"} {
		t.Run(dialect, func(t *testing.T) {
			g, types := mixedSchemaFixture()
			layout, err := managedSchemaLayout(g, types, dialect)
			require.NoError(t, err)
			for _, c := range layout.Columns {
				if c.Name == "v" {
					require.Equal(t, want, c.SQLType)
					require.True(t, c.Nullable)
					return
				}
			}
			t.Fatal("mixed result column missing")
		})
	}
}
