package workspace

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type writeGroupReadinessRules struct {
	deviceID string
	links    []*schema.SourceRuleLink
}

func (r writeGroupReadinessRules) ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error) {
	return []*schema.SourceRule{{ID: "source-rule", DeviceID: r.deviceID}}, nil
}

func (r writeGroupReadinessRules) ListLinks(context.Context, string) ([]*schema.SourceRuleLink, error) {
	return r.links, nil
}

func TestService_BasicGroupReadinessUsesCanonicalScopeAndRetainsLegacyGates(t *testing.T) {
	for _, strategy := range []WriteGroupStorageStrategy{WriteGroupStorageStrategyManaged, WriteGroupStorageStrategyCustom} {
		t.Run(string(strategy), func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, filepath.Join(t.TempDir(), "configuration.db"))
			defer db.Close()
			f := seedWriteGroupFixture(ctx, t, db)
			_, err := db.ExecContext(ctx, `UPDATE devices SET description = '', connection_config = '{"host":"127.0.0.1","port":502,"slave_id":1}', last_test_success = 1 WHERE id = ?`, f.deviceID)
			require.NoError(t, err)
			target := filepath.Join(t.TempDir(), "target.db")
			config, err := json.Marshal(map[string]string{"database": target, "schema": "main"})
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `UPDATE database_connectors SET connection_config = ? WHERE id = ?`, string(config), f.connectorID)
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled)
				VALUES ('unmigrated-point',?,'Unmigrated','40002','FC03','float32','read',1)`, f.deviceID)
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
				VALUES ('unmigrated-tag','unmigrated','unmigrated','Unmigrated','float32','active','{}')`)
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
				VALUES ('unmigrated-mapping','unmigrated-point','unmigrated-tag','[]','active',1)`)
			require.NoError(t, err)
			var mappingID string
			require.NoError(t, db.QueryRowContext(ctx, `SELECT id FROM mappings WHERE point_id = ? AND tag_id = ?`, f.pointID, f.tagID).Scan(&mappingID))
			unmigratedTag, unmigratedMapping := "unmigrated-tag", "unmigrated-mapping"
			rules := writeGroupReadinessRules{deviceID: f.deviceID, links: []*schema.SourceRuleLink{
				{RuleID: "source-rule", PointID: f.pointID, TagID: &f.tagID, MappingID: &mappingID},
				{RuleID: "source-rule", PointID: "unmigrated-point", TagID: &unmigratedTag, MappingID: &unmigratedMapping},
			}}
			connector := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(db), dbtarget.NewSQLTargetMappingRepository(db))
			ws := NewService(NewSQLRepository(db)).WithReadinessServices(
				device.NewService(device.NewSQLRepository(db), nil), rules, connector,
				dbtarget.NewMappingService(dbtarget.NewSQLTargetMappingRepository(db), dbtarget.NewSQLConnectorRepository(db), nil),
			)
			current, err := ws.BindDatabaseConnector(ctx, f.connectorID)
			require.NoError(t, err)
			groups := NewWriteGroupService(ws, NewSQLWriteGroupRepository(db))
			ws.WithWriteGroupReadiness(groups)
			group := newBasicWriteGroup(f)
			group.Destination.StorageStrategy = strategy
			group.Destination.TableSchema = "main"
			group.RowPolicy.IntervalSeconds = 15
			created, err := groups.Create(ctx, WriteGroupMutation{
				WorkspaceID: f.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
				ExpectedConnectorRevision: "connector-revision-1", Group: group,
			})
			require.NoError(t, err)
			require.Nil(t, created.Group.Members[0].MeasurementID)
			summary, err := ws.Readiness(ctx)
			require.NoError(t, err)
			var missingScopes []string
			for _, issue := range summary.Issues {
				if issue.Code == "database-target-missing" {
					missingScopes = append(missingScopes, issue.Scope)
				}
			}
			require.Equal(t, []string{"unmigrated-point"}, missingScopes)
			require.Contains(t, readinessIssueCodes(summary.Issues), "schema-unverified")
			require.Contains(t, readinessIssueCodes(summary.Issues), "write-group-apply-required")
			require.False(t, summary.Ready)
			reloaded, err := groups.Get(ctx, created.Group.ID)
			require.NoError(t, err)
			require.Equal(t, created, reloaded)
			require.NoFileExists(t, target)
			_, err = db.ExecContext(ctx, `UPDATE write_groups SET status = 'disabled' WHERE id = ?`, created.Group.ID)
			require.NoError(t, err)
			summary, err = ws.Readiness(ctx)
			require.NoError(t, err)
			missingScopes = nil
			for _, issue := range summary.Issues {
				if issue.Code == "database-target-missing" {
					missingScopes = append(missingScopes, issue.Scope)
				}
			}
			require.Equal(t, []string{"unmigrated-point"}, missingScopes)
			require.NotContains(t, readinessIssueCodes(summary.Issues), "write-group-apply-required")
		})
	}
}
