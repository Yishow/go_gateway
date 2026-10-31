package workspace

import (
	"go-gateway/internal/datalink/dbtarget"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedPreparationReadinessHasMemberCadenceButNoRuntimeAuthority(t *testing.T) {
	g, types := mixedSchemaFixture()
	require.Empty(t, validateBasicWriteGroupConfig(g), "member cadence replaces inapplicable group interval")
	layout, err := managedSchemaLayout(g, types, "sqlite")
	require.NoError(t, err)
	inspection := &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Schema: "main", Table: g.Destination.TableName}
	for _, c := range layout.Columns {
		inspection.Columns = append(inspection.Columns, dbtarget.ColumnInfo{Name: c.Name, DataType: c.SQLType, Nullable: c.Nullable, PrimaryKey: c.PrimaryKey})
	}
	ready, issues := evaluateWriteGroupInspection(g, types, inspection, "sqlite")
	require.True(t, ready)
	require.Empty(t, issues)
	require.Empty(t, evaluateWriteGroupPartialInspection(g, types, inspection, "sqlite"))
	result := newWriteGroupReadiness(&Record{ID: "w"}, g)
	result.ConfigReady = true
	result.SchemaReady = true
	require.False(t, finalizeWriteGroupReadiness(result).Ready, "preparation is not mixed runtime support")
	require.False(t, recordingReductionUnsupported(g, types), "validated shape is now supported; absent coherent runtime proofs still block readiness")
	for i, c := range inspection.Columns {
		if c.Name == "v" {
			inspection.Columns[i].Nullable = false
		}
	}
	ready, issues = evaluateWriteGroupInspection(g, types, inspection, "sqlite")
	require.False(t, ready)
	require.NotEmpty(t, issues, "required but not-due still needs SQL NULL")
}
