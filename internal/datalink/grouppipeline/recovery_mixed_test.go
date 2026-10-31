package grouppipeline

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"strings"
	"testing"
)

func mixedStorageProofFixture(t *testing.T, long bool, dialect dbtarget.SQLDialect) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	t.Helper()
	g, _ := canonicalProof(t)
	g.RecordID = "11111111-1111-4111-8111-111111111111"
	g.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
	g.StorageMode = workspace.WriteGroupStorageModePlain
	g.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileMixed, ResultShape: workspace.WriteGroupResultShapeMemberCalculationV1}
	g.RowPolicy.IntervalSeconds = 0
	g.RowPolicy.BucketStartColumn = ""
	g.RowPolicy.DueAtColumn = "due_at"
	g.RowPolicy.IncompletePolicy = "partial"
	g.Destination.SchemaRevision = "confirmed-operation"
	g.Destination.SchemaDigest = strings.Repeat("b", 64)
	g.Members = g.Members[:1]
	m := &g.Members[0]
	m.Calculation = workspace.WriteGroupMemberCalculationMean
	m.Cadence = &workspace.WriteGroupMemberCadence{Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond, Alignment: workspace.WriteGroupAlignmentUTCEpoch}
	m.SourceValueType = schema.DataTypeFloat64
	m.ResultValueType = measurement.ExactFloat64
	m.ValueRole = workspace.WriteGroupValueRoleGauge
	m.Required = false
	if long {
		g.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
		table, _, ok := workspace.LongStorageIdentity(g)
		require.True(t, ok)
		g.Destination.TableName = table
	}
	types := map[string]schema.DataType{m.TagID: schema.DataTypeFloat64}
	expected, err := workspace.MixedManagedSchemaLayout(g, types, string(dialect))
	require.NoError(t, err)
	layout := &workspace.WriteGroupRuntimeLayout{Dialect: dialect, TagTypes: types, ReceiptTableReady: true, SchemaDigest: strings.Repeat("a", 64)}
	for _, c := range expected.Columns {
		layout.Columns = append(layout.Columns, dbtarget.ColumnInfo{Name: c.Name, DataType: c.SQLType, Nullable: c.Nullable, PrimaryKey: c.PrimaryKey})
	}
	return g, layout
}
func TestMixedFrozenStorageProofUsesCanonicalSchema(t *testing.T) {
	for _, long := range []bool{false, true} {
		for _, dialect := range []dbtarget.SQLDialect{dbtarget.SQLDialectSQLite, dbtarget.SQLDialectPostgres} {
			g, l := mixedStorageProofFixture(t, long, dialect)
			require.True(t, validMixedRuntimeStorageProof(g, l))
			require.NotEqual(t, g.Destination.SchemaDigest, l.SchemaDigest, "schema consent and readiness inspection digests are different authorities")
		}
	}
}
func TestMixedFrozenStorageProofRejectsIncompleteFacts(t *testing.T) {
	for _, long := range []bool{false, true} {
		for name, mutate := range map[string]func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout){
			"owner missing": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				for i, c := range l.Columns {
					if strings.HasPrefix(c.Name, "_gw_owner_") {
						l.Columns = append(l.Columns[:i], l.Columns[i+1:]...)
						break
					}
				}
			},
			"provenance nullable": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				for i, c := range l.Columns {
					if c.Name == "provenance" {
						l.Columns[i].Nullable = true
					}
				}
			},
			"result nonnullable": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				for i, c := range l.Columns {
					if c.Name == g.Members[0].TargetColumn || c.Name == "value_float" {
						l.Columns[i].Nullable = false
					}
				}
			},
			"record key missing": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				for i := range l.Columns {
					l.Columns[i].PrimaryKey = false
				}
			},
			"unknown shape": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				g.RecordingProfile.ResultShape = "unknown"
			},
			"record UUID missing": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) { g.RecordID = "" },
			"digest absent":       func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) { l.SchemaDigest = "" },
			"consent absent":      func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) { g.Destination.SchemaRevision = "" },
			"mysql": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				l.Dialect = dbtarget.SQLDialectMySQL
			},
			"wrong source type": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				l.TagTypes[g.Members[0].TagID] = schema.DataTypeBool
			},
			"wrong result SQL": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				for i, c := range l.Columns {
					if c.Name == g.Members[0].TargetColumn || c.Name == "value_float" {
						l.Columns[i].DataType = "BLOB"
					}
				}
			},
		} {
			t.Run(name+map[bool]string{true: "/long", false: "/wide"}[long], func(t *testing.T) {
				g, l := mixedStorageProofFixture(t, long, dbtarget.SQLDialectSQLite)
				mutate(g, l)
				require.False(t, validMixedRuntimeStorageProof(g, l))
			})
		}
	}
}
func TestMixedFrozenStorageFixtureRoundTrip(t *testing.T) {
	g, l := mixedStorageProofFixture(t, true, dbtarget.SQLDialectSQLite)
	raw, err := json.Marshal(workspace.WriteGroupAppliedSnapshot{Group: g, RuntimeLayout: l})
	require.NoError(t, err)
	var round workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &round))
	require.True(t, validMixedRuntimeStorageProof(round.Group, round.RuntimeLayout))
}

func TestMixedFrozenStorageProofRequiresOwnedShapeAndNoForeignCatalog(t *testing.T) {
	require.False(t, validMixedRuntimeStorageProof(nil, nil))
	for _, long := range []bool{false, true} {
		for name, mutate := range map[string]func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout){
			"foreign owner identity": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) { g.WorkspaceID = "foreign" },
			"bad digest": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				l.SchemaDigest = strings.Repeat("z", 64)
			},
			"duplicate column": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				l.Columns = append(l.Columns, l.Columns[0])
			},
			"shared interval": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) { g.RowPolicy.IntervalSeconds = 60 },
			"sqlserver": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				l.Dialect = dbtarget.SQLDialectSQLServer
			},
			"timescale": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				g.StorageMode = workspace.WriteGroupStorageModeTimescale
			},
			"receipt absent": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
				g.WritePolicy.DedupeCapability = "receipt"
				l.ReceiptTableReady = false
			},
		} {
			t.Run(name+map[bool]string{true: "/long", false: "/wide"}[long], func(t *testing.T) {
				g, l := mixedStorageProofFixture(t, long, dbtarget.SQLDialectSQLite)
				mutate(g, l)
				require.False(t, validMixedRuntimeStorageProof(g, l))
			})
		}
	}
	g, l := mixedStorageProofFixture(t, true, dbtarget.SQLDialectSQLite)
	g.Destination.TableName = "foreign_table"
	require.False(t, validMixedRuntimeStorageProof(g, l))
	g, l = mixedStorageProofFixture(t, false, dbtarget.SQLDialectSQLite)
	g.RowPolicy.DueAtColumn = ""
	require.False(t, validMixedRuntimeStorageProof(g, l))
}
