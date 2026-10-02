package api

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func newLegacyWriteRouterFixture(t *testing.T, rowGroup bool) (fixture writeGroupRouterFixture, targetPath string) {
	t.Helper()
	var f writeGroupRouterFixture
	if rowGroup {
		f, targetPath = newRowGroupMigrationRouterFixture(t)
		review := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, rowGroupMigrationRequest(t, f))
		require.Equal(t, http.StatusOK, review.Code, review.Body.String())
	} else {
		f, targetPath = newSingleMigrationRouterFixture(t)
		review := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, migrationReviewRequest(t, f))
		require.Equal(t, http.StatusOK, review.Code, review.Body.String())
	}
	connectors := dbtarget.NewSQLConnectorRepository(f.db)
	mappings := dbtarget.NewSQLTargetMappingRepository(f.db)
	rules := sourcerule.NewSQLRepository(f.db)
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	f.router = NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)),
		DBTarget: dbtarget.NewConnectorService(connectors, mappings), DBMapping: dbtarget.NewMappingService(mappings, connectors, tags),
		SourceRule: sourcerule.NewService(rules, nil, nil, nil),
	})
	// Studio's point lookup uses persisted rule links, rather than the mapping
	// table alone. Add one real link so the request reaches the target writer.
	now := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	require.NoError(t, rules.Create(t.Context(), &schema.SourceRule{
		ID: "rule-A", DeviceID: "device-A", StartAddress: "40001", Count: 1,
		DataType: schema.DataTypeFloat32, NamingPrefix: "Temperature", Enabled: true,
		Origin: "manual", SkippedAddresses: "[]", RevisionID: "rule-1", CreatedAt: now, UpdatedAt: now,
	}))
	tagID, mappingID := "tag-A", "mapping-A"
	require.NoError(t, rules.CreateLinks(t.Context(), []*schema.SourceRuleLink{{
		ID: "link-A", RuleID: "rule-A", Address: "40001", PointID: "point-A", TagID: &tagID,
		MappingID: &mappingID, CreatedAt: now, UpdatedAt: now,
	}}))
	var err error
	f.record, err = f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	return f, targetPath
}

// Capture all local write authorities, including their revisions and original
// migration payloads. A rejected HTTP request must leave this snapshot intact.
func legacyWriteLocalSnapshot(t *testing.T, db *sql.DB) map[string][][]any {
	t.Helper()
	queries := map[string]string{
		"workspace":      `SELECT * FROM system_settings WHERE key = 'studio_v2_workspace' ORDER BY key`,
		"connectors":     `SELECT * FROM database_connectors ORDER BY id`,
		"targets":        `SELECT * FROM database_target_mappings ORDER BY id`,
		"groups":         `SELECT * FROM write_groups ORDER BY id`,
		"members":        `SELECT * FROM write_group_members ORDER BY group_id, member_index`,
		"versions":       `SELECT * FROM write_group_versions ORDER BY group_id, group_revision`,
		"migration_maps": `SELECT * FROM write_group_migration_maps ORDER BY workspace_id, source_kind, source_id`,
	}
	result := make(map[string][][]any, len(queries))
	for name, query := range queries {
		rows, err := db.QueryContext(t.Context(), query)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		values := make([][]any, 0)
		for rows.Next() {
			row := make([]any, len(columns))
			pointers := make([]any, len(columns))
			for index := range row {
				pointers[index] = &row[index]
			}
			require.NoError(t, rows.Scan(pointers...))
			for index, value := range row {
				if bytes, ok := value.([]byte); ok {
					row[index] = string(bytes)
				}
			}
			values = append(values, row)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		result[name] = values
	}
	return result
}

func assertLegacyWriteConflict(t *testing.T, code int, body map[string]any) {
	t.Helper()
	require.Equal(t, http.StatusConflict, code, body)
	assertWriteGroupError(t, body, "WRITE_GROUP_LEGACY_WRITE_CONFLICT")
	require.Equal(t, "open_write_groups", body["error"].(map[string]any)["action"])
}

func legacyDatabaseConfigRequest(f writeGroupRouterFixture, targetPath string) map[string]any {
	return map[string]any{
		"connector_id": "connector-A", "expected_connector_revision": "connector-1",
		"kind": "sqlite", "name": "Target", "database": targetPath, "schema": "main",
		"table": "raw_values", "write_mode": "insert", "write_interval_seconds": 15,
		"timestamp_column": "observed_at", "expected_setup_revision": f.record.DatabaseSetupRevision,
	}
}
