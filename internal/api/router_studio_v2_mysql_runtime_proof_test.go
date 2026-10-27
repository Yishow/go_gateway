package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMySQLRuntimeProofFreezesCatalogBeforeRecordingGateOpens(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			f := newMySQLSchemaRouterFixture(t, family)
			group := ensureLong(t, f, "instantaneous", "second", 15, "")
			id, table := group["id"].(string), destinationOf(group)["table_name"].(string)
			var count int
			require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME IN (?,?)", table, dbtarget.EffectReceiptTable).Scan(&count))
			require.Zero(t, count)
			var created []string
			cleanupMySQLSchemaTables(t, f.target, &created)
			token := previewLong(t, f, id)
			request := scopeBody(t, f, id)
			request["token"], request["operation_id"] = token.Token, token.OperationID
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-apply", request)
			for _, name := range []string{table, dbtarget.EffectReceiptTable} {
				require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?", name).Scan(&count))
				if count == 1 {
					created = append(created, name)
				}
			}
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			saved, err := f.groups.Get(t.Context(), id)
			require.NoError(t, err)
			ready, err := f.groups.Readiness(t.Context(), id)
			require.NoError(t, err)
			require.True(t, ready.SchemaReady, "a fully verified MySQL schema is prepared: %+v", ready.Issues)
			require.True(t, ready.Ready, "verified version and catalog proof allow recording: %+v", ready.Issues)

			service := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
			inspector := f.groups.ManagedTableInspector(service)
			capability, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), saved.Group.Destination.ConnectorID, saved.Group.Destination.ConnectorRevision)
			require.NoError(t, err)
			physical, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID, capability.ConnectorRevision, capability.Database, capability.CapabilityRevision, table)
			require.NoError(t, err)
			receipt, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID, capability.ConnectorRevision, capability.Database, capability.CapabilityRevision, dbtarget.EffectReceiptTable)
			require.NoError(t, err)
			types := map[string]schema.DataType{}
			for _, member := range saved.Group.Members {
				var dataType schema.DataType
				require.NoError(t, f.db.QueryRowContext(t.Context(), "SELECT data_type FROM tags WHERE id=?", member.TagID).Scan(&dataType))
				types[member.TagID] = dataType
			}
			frozen := &workspace.WriteGroupRuntimeLayout{
				Dialect: dbtarget.SQLDialectMySQL, Columns: physical.Inspection.Columns, TagTypes: types,
				ReceiptTableReady: true, SchemaDigest: saved.Group.Destination.SchemaDigest,
				MySQL: &dbtarget.MySQLManagedRowProof{ConnectorID: capability.ConnectorID, ConnectorRevision: capability.ConnectorRevision, Group: physical, Receipt: receipt},
			}
			require.NoError(t, workspace.ValidateMySQLRuntimeLayout(saved.Group, frozen))
			require.Equal(t, frozen.SchemaDigest, ready.SchemaDigest)
			mutations := map[string]func(*workspace.WriteGroupRuntimeLayout){
				"missing proof":      func(l *workspace.WriteGroupRuntimeLayout) { l.MySQL = nil },
				"missing receipt":    func(l *workspace.WriteGroupRuntimeLayout) { l.MySQL.Receipt = nil },
				"unverified indexes": func(l *workspace.WriteGroupRuntimeLayout) { l.MySQL.Receipt.Indexes = nil },
				"changed connector":  func(l *workspace.WriteGroupRuntimeLayout) { l.MySQL.ConnectorRevision = "changed" },
				"foreign database":   func(l *workspace.WriteGroupRuntimeLayout) { l.MySQL.Group.Evidence.Database = "foreign" },
				"foreign table":      func(l *workspace.WriteGroupRuntimeLayout) { l.MySQL.Group.Inspection.Table = "foreign" },
				"changed capability": func(l *workspace.WriteGroupRuntimeLayout) { l.MySQL.Group.Evidence.CapabilityRevision = "changed" },
				"schema digest":      func(l *workspace.WriteGroupRuntimeLayout) { l.SchemaDigest = "changed" },
				"missing type":       func(l *workspace.WriteGroupRuntimeLayout) { delete(l.TagTypes, saved.Group.Members[0].TagID) },
				"physical type drift": func(l *workspace.WriteGroupRuntimeLayout) {
					for i := range l.MySQL.Group.Inspection.Columns {
						if l.MySQL.Group.Inspection.Columns[i].Name == "value_decimal" {
							l.MySQL.Group.Inspection.Columns[i].DataType = "DOUBLE"
						}
					}
				},
			}
			encoded, err := json.Marshal(frozen)
			require.NoError(t, err)
			var persisted workspace.WriteGroupRuntimeLayout
			require.NoError(t, json.Unmarshal(encoded, &persisted))
			require.NoError(t, workspace.ValidateMySQLRuntimeLayout(saved.Group, &persisted), "real catalog proof must remain valid after persisted JSON reload")
			for name, mutate := range mutations {
				t.Run(name, func(t *testing.T) {
					var candidate workspace.WriteGroupRuntimeLayout
					require.NoError(t, json.Unmarshal(encoded, &candidate))
					mutate(&candidate)
					require.Error(t, workspace.ValidateMySQLRuntimeLayout(saved.Group, &candidate))
				})
			}
		})
	}
}
