package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// This exercises the production constructor, not an optionally configured test
// Store: forgetting its canonical intake option must permit a failing ACK here.
func TestProductionServicesFenceCachedIntakeAfterStop(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	services, err := wireGatewayServices(db, connector.NewConnectionManager(connector.ConnectionManagerConfig{}))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	record, err := services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	now := time.Now().UTC()
	effective := now.Truncate(10 * time.Second).Add(-10 * time.Second)
	group := &workspace.WriteGroup{
		ID: "G1", WorkspaceID: record.ID, Revision: "R1", AppliedRevision: "R1", Name: "Temperature", Status: workspace.WriteGroupStatusReady,
		Members:     []workspace.WriteGroupMember{{DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", SourceRevision: "source-1", MappingRevision: "mapping-1", TargetColumn: "temperature", Required: true}},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-1", ConnectorRevision: "C1", TableName: "readings"},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10},
	}
	_, err = db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('connector-1','Target','sqlite','{}','C1','ready',1)`)
	require.NoError(t, err)
	policy, err := json.Marshal(group.RowPolicy)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,status,destination_connector_id,destination_connector_revision,destination_table_name,row_policy) VALUES('G1',?,'R1','R1','Temperature','ready','connector-1','C1','readings',?)`, record.ID, string(policy))
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_members(group_id,member_index,device_id,point_id,tag_id,source_revision,mapping_revision,target_column,required) VALUES('G1',0,'device-1','point-1','tag-1','source-1','mapping-1','temperature',1)`)
	require.NoError(t, err)
	payload, err := json.Marshal(group)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES('G1','R1',?,?)`, effective, string(payload))
	require.NoError(t, err)
	snap := workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: record.ID, GroupID: "G1", GroupRevision: "R1", AppliedRevision: "R1", EffectiveAt: effective, Group: group,
		RuntimeLayout: &workspace.WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectSQLite, Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "REAL"}}, TagTypes: map[string]schema.DataType{"tag-1": schema.DataTypeFloat64}},
	}
	payload, err = json.Marshal(snap)
	require.NoError(t, err)
	key := groupdelivery.GroupKey{WorkspaceID: record.ID, GroupID: "G1", GroupRevision: "R1"}
	require.NoError(t, groupdelivery.NewStore(db).SaveRuntimeVersion(t.Context(), key, payload))
	require.NoError(t, services.groupPipe.Reconcile(t.Context()))
	require.Equal(t, []grouppipeline.GroupStatus{{GroupID: "G1", Revision: "R1", State: grouppipeline.StateActive}}, services.groupPipe.Status())
	envelope := measurement.SampleEnvelope{SampleID: "accepted", WorkspaceID: record.ID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", SourceRevision: "source-1", MappingRevision: "mapping-1", ObservedAt: now, ReceivedAt: now, TimeOrigin: "gateway", Value: 12.5, Quality: schema.QualityGood}
	require.NoError(t, services.groupPipe.AcceptSample(t.Context(), envelope))
	current, err := services.writeGroups.Get(t.Context(), "G1")
	require.NoError(t, err)
	_, err = services.writeGroups.Disable(t.Context(), "G1", workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: "C1"})
	require.NoError(t, err)
	require.NoError(t, services.groupPipe.AcceptSample(t.Context(), envelope), "accepted retry remains an ACK without reconcile")
	envelope.SampleID = "refused"
	err = services.groupPipe.AcceptSample(t.Context(), envelope)
	var refusal *groupdelivery.OwnershipError
	require.ErrorAs(t, err, &refusal, "the production Store must fence the stale active boundary")
	var count int
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&count))
	require.Equal(t, 1, count)
}
