package workspace

import (
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupReadinessAllowsExplicitPartialOnlyWithVerifiedStorage(t *testing.T) {
	for _, test := range []struct {
		name       string
		nullable   bool
		provenance string
		ready      bool
	}{
		{"nullable-with-provenance", true, "TEXT", true},
		{"not-null", false, "TEXT", false},
		{"no-provenance", true, "", false},
		{"incompatible-provenance", true, "INTEGER", false},
		{"entity-without-column", true, "TEXT", false},
		{"entity-with-column", true, "TEXT", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture, created := newReadinessGroup(ctx, t, db)
			candidate := cloneWriteGroup(created.Group)
			candidate.RowPolicy.IncompletePolicy = "partial"
			candidate.RowPolicy.ProvenanceColumn = "quality_info"
			candidate.Members[0].Required = false
			if test.name == "entity-without-column" || test.name == "entity-with-column" {
				candidate.Members[0].EntityKey = "line-A"
			}
			if test.name == "entity-with-column" {
				candidate.RowPolicy.EntityKeyColumn = "line"
			}
			current, err := service.workspaceSvc.GetOrCreate(ctx)
			require.NoError(t, err)
			updated, err := service.Update(ctx, created.Group.ID, WriteGroupMutation{
				WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
				ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: candidate,
			})
			require.NoError(t, err)
			columns := []dbtarget.ColumnInfo{{Name: "temperature", DataType: "DOUBLE PRECISION", Nullable: test.nullable}}
			if test.provenance != "" {
				columns = append(columns, dbtarget.ColumnInfo{Name: "quality_info", DataType: test.provenance})
			}
			columns = append(columns, dbtarget.ColumnInfo{Name: "line", DataType: "TEXT"})
			service.WithTableInspector(&writeGroupReadinessInspector{inspection: &dbtarget.TableInspection{
				Status: dbtarget.TableInspectionExists, Schema: "persisted-schema", Table: "raw_values", Columns: columns,
			}})
			readiness, err := service.Readiness(ctx, updated.Group.ID)
			require.NoError(t, err)
			require.Equal(t, test.ready, readiness.Ready, "only verified storage may apply: %+v", readiness.Issues)
			if !test.ready {
				require.NotEmpty(t, readiness.Issues)
			}
			require.Empty(t, updated.Group.AppliedRevision)
		})
	}
}
