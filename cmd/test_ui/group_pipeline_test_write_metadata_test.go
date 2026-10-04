package main

import (
	"fmt"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestProductionTestWriteKeepsManagedMetadataAndNeighborRows(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			target := newManagedPipelineTarget(t, kind)
			members := addManagedMetadataTestPoints(t, target.env)
			group := target.prepare(t, members, workspace.WriteGroupRowPolicy{IntervalSeconds: 10, EntityKeyColumn: "entity"})
			table := managedTestWriteTable(group, target.namespace, target.kind)
			insertManagedNeighbor(t, target, group, table)

			router := target.env.router()
			path := testWriteBase + "/write-groups/" + group.ID
			status, preview := httpJSON(t, router, "POST", path+"/test-write-preview", nil)
			require.Equal(t, 200, status, "%v", preview)
			data := preview["data"].(map[string]any)
			status, confirmed := httpJSON(t, router, "POST", path+"/test-write", map[string]any{
				"token": data["token"], "operation_id": data["operation_id"],
			})
			require.Equal(t, 200, status, "%v", confirmed)
			result := confirmed["data"].(map[string]any)
			require.Equal(t, "written_verified", result["write_outcome"], "%v", result)
			require.Equal(t, "cleaned", result["cleanup_status"], "%v", result)

			var total, owned int
			require.NoError(t, target.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table).Scan(&total))
			require.NoError(t, target.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table+` WHERE "entity" LIKE 'gw-test-%'`).Scan(&owned))
			require.Equal(t, 1, total, "test cleanup must leave the production neighbor")
			require.Zero(t, owned, "operation-owned rows must be removed")

			var recordID, groupID, deviceID, bucketStart, entity string
			var temperature float64
			var pressure int64
			require.NoError(t, target.db.QueryRowContext(t.Context(), `SELECT "record_id", "group_id", "device_id", "bucket_start", "entity", `+quoteTestIdentifier(group.Members[0].TargetColumn)+`, `+quoteTestIdentifier(group.Members[1].TargetColumn)+` FROM `+table).Scan(
				&recordID, &groupID, &deviceID, &bucketStart, &entity, &temperature, &pressure,
			))
			require.Equal(t, "production-record", recordID)
			require.Equal(t, "production-group", groupID)
			require.Equal(t, "device-1", deviceID)
			require.Equal(t, "production-entity", entity)
			require.Equal(t, float64(1.25), temperature)
			require.Equal(t, int64(7), pressure)
		})
	}
}

func addManagedMetadataTestPoints(t *testing.T, env *outageEnv) []workspace.WriteGroupMember {
	t.Helper()
	members := []workspace.WriteGroupMember{
		{DeviceID: "device-1", PointID: "metadata-point-temperature", TagID: "metadata-tag-temperature", Required: true},
		{DeviceID: "device-1", PointID: "metadata-point-pressure", TagID: "metadata-tag-pressure", Required: true},
	}
	for i, member := range members {
		dataType, address := "float64", "60101"
		if i == 1 {
			dataType, address = "int64", "60102"
		}
		_, err := env.db.ExecContext(t.Context(), `INSERT INTO points(id,device_id,name,address,function,data_type,mode,enabled) VALUES(?,?,?,?, 'FC03',?,'read',1)`, member.PointID, member.DeviceID, member.PointID, address, dataType)
		require.NoError(t, err)
		_, err = env.db.ExecContext(t.Context(), `INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels) VALUES(?,?,?,?,?,'active','{}')`, member.TagID, member.TagID, member.TagID, member.TagID, dataType)
		require.NoError(t, err)
		_, err = env.db.ExecContext(t.Context(), `INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled) VALUES(?,?,?,'[]','active',1)`, "metadata-mapping-"+fmt.Sprint(i), member.PointID, member.TagID)
		require.NoError(t, err)
	}
	return members
}

func managedTestWriteTable(group *workspace.WriteGroup, tableSchema string, kind schema.DatabaseConnectorKind) string {
	table := quoteTestIdentifier(group.Destination.TableName)
	if kind == schema.DatabaseConnectorKindPostgres {
		return quoteTestIdentifier(tableSchema) + "." + table
	}
	return table
}

func quoteTestIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func insertManagedNeighbor(t *testing.T, target *managedPipelineTarget, group *workspace.WriteGroup, table string) {
	t.Helper()
	columns := []string{
		group.RowPolicy.RecordKeyColumn, group.RowPolicy.GroupIDColumn, group.RowPolicy.DeviceIDColumn,
		group.RowPolicy.BucketStartColumn, group.RowPolicy.ProvenanceColumn, group.RowPolicy.EntityKeyColumn,
		group.Members[0].TargetColumn, group.Members[1].TargetColumn,
	}
	values := []any{"production-record", "production-group", "device-1", "2026-01-01T00:00:00Z", "[]", "production-entity", 1.25, int64(7)}
	marks := make([]string, len(values))
	for i := range marks {
		if target.kind == schema.DatabaseConnectorKindPostgres {
			marks[i] = fmt.Sprintf("$%d", i+1)
		} else {
			marks[i] = "?"
		}
	}
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = quoteTestIdentifier(column)
	}
	statement := `INSERT INTO ` + table + ` (` + strings.Join(quoted, ", ") + `) VALUES (` + strings.Join(marks, ", ") + `)`
	_, err := target.db.ExecContext(t.Context(), statement, values...)
	require.NoError(t, err)
}
