package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// This fixture seeds an already accepted runtime version/outbox after real
// schema confirmation. It does not call or enable the gated product Apply.
type mysqlRepairRig struct {
	f        managedGroupFixture
	targets  *dbtarget.ConnectorService
	store    *groupdelivery.Store
	recovery *grouppipeline.StorageRecovery
	group    *workspace.WriteGroup
	frozen   *workspace.WriteGroupRuntimeLayout
	key      groupdelivery.GroupKey
	base     time.Time
}

func newMySQLRepairRig(t *testing.T, family string, recreate bool) *mysqlRepairRig {
	t.Helper()
	f := newMySQLSchemaRouterFixture(t, family, "float64")
	request := profileRequest(t, f.writeGroupRouterFixture, "instantaneous", 15, "second", "main", "")
	draft := request["group"].(map[string]any)
	draft["layout_version"] = "managed_long_v1"
	draft["destination"].(map[string]any)["table_schema"] = f.namespace
	draft["write_policy"] = map[string]any{"recreate_missing": recreate}
	response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	id := saved["id"].(string)
	table := destinationOf(saved)["table_name"].(string)
	var count int
	require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME IN (?,?)", table, dbtarget.EffectReceiptTable).Scan(&count))
	require.Zero(t, count, "this fixture exclusively creates its group and receipt tables")
	var owned []string
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		for _, name := range owned {
			require.True(t, longTablePattern.MatchString(name) || name == dbtarget.EffectReceiptTable)
			_, err := f.target.ExecContext(ctx, "DROP TABLE IF EXISTS `"+name+"`")
			require.NoError(t, err)
		}
	})
	token := previewLong(t, f, id)
	confirmRequest := scopeBody(t, f, id)
	confirmRequest["token"], confirmRequest["operation_id"] = token.Token, token.OperationID
	response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-apply", confirmRequest)
	for _, name := range []string{table, dbtarget.EffectReceiptTable} {
		require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?", name).Scan(&count))
		if count == 1 {
			owned = append(owned, name)
		}
	}
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	current, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	group := current.Group
	require.NotNil(t, group.WritePolicy.RecreateMissing)
	require.Equal(t, recreate, *group.WritePolicy.RecreateMissing)
	ready, err := f.groups.Readiness(t.Context(), id)
	require.NoError(t, err)
	require.True(t, ready.SchemaReady)
	require.True(t, ready.Ready, "verified destination is ready: %+v", ready.Issues)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	inspector := f.groups.ManagedTableInspector(targets)
	capability, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	require.NoError(t, err)
	physical, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID, capability.ConnectorRevision, capability.Database, capability.CapabilityRevision, table)
	require.NoError(t, err)
	receipt, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID, capability.ConnectorRevision, capability.Database, capability.CapabilityRevision, dbtarget.EffectReceiptTable)
	require.NoError(t, err)
	types := map[string]schema.DataType{}
	for _, member := range group.Members {
		var dataType schema.DataType
		require.NoError(t, f.db.QueryRowContext(t.Context(), "SELECT data_type FROM tags WHERE id=?", member.TagID).Scan(&dataType))
		types[member.TagID] = dataType
	}
	frozen := &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectMySQL, Columns: physical.Inspection.Columns, TagTypes: types, ReceiptTableReady: true,
		SchemaDigest: group.Destination.SchemaDigest,
		MySQL:        &dbtarget.MySQLManagedRowProof{ConnectorID: capability.ConnectorID, ConnectorRevision: capability.ConnectorRevision, Group: physical, Receipt: receipt},
	}
	require.NoError(t, workspace.ValidateMySQLRuntimeLayout(group, frozen))
	group.AppliedRevision = group.Revision
	base := time.Now().UTC().Truncate(15 * time.Second)
	version := &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.Revision,
		EffectiveAt: base, Group: group, RuntimeLayout: frozen,
	}
	payload, err := json.Marshal(version)
	require.NoError(t, err)
	store := groupdelivery.NewStore(f.db)
	key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision}
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
	return &mysqlRepairRig{
		f: f, targets: targets, store: store, group: group, frozen: frozen, key: key, base: base,
		recovery: grouppipeline.NewStorageRecovery(store, recordingplan.NewService(recordingplan.NewSQLRepository(f.db)), targets),
	}
}

func (r *mysqlRepairRig) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	opened, err := r.targets.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		if errors.Is(err, dbtarget.ErrDestinationBlocked) {
			err = errors.Join(groupdelivery.ErrTargetBlocked, err)
		}
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close, MySQL: r.frozen.MySQL, MySQLCapability: opened.MySQL}, nil
}

func (r *mysqlRepairRig) sender() *groupdelivery.Sender {
	return groupdelivery.NewSender(r.store, r, groupdelivery.SenderConfig{Recovery: r.recovery, Backoff: func(int) (string, time.Time) {
		return groupdelivery.StateRetrying, time.Now().UTC().Add(-time.Second)
	}})
}

func (r *mysqlRepairRig) enqueue(t *testing.T, index int) string {
	t.Helper()
	member := r.group.Members[0]
	start := r.base.Add(time.Duration(index) * 15 * time.Second)
	record := sha256.Sum256(fmt.Appendf(nil, "%s:%d", r.group.ID, index))
	effect := sha256.Sum256(record[:])
	recordID, effectKey := hex.EncodeToString(record[:]), hex.EncodeToString(effect[:])
	row := dbtarget.EncodedRow{RecordID: recordID, EffectKey: effectKey, BucketStart: start, Cells: []dbtarget.EncodedCell{
		{Column: "bucket_start", Value: start}, {Column: "bucket_end", Value: start.Add(15 * time.Second)},
		{Column: "record_id", Value: recordID}, {Column: "workspace_id", Value: r.group.WorkspaceID},
		{Column: "group_id", Value: r.group.ID}, {Column: "group_revision", Value: r.group.Revision},
		{Column: "device_id", Value: member.DeviceID}, {Column: "point_id", Value: member.PointID}, {Column: "tag_id", Value: member.TagID},
		{Column: "profile_kind", Value: "instantaneous"}, {Column: "value_type", Value: "float64"}, {Column: "status", Value: "ok"},
		{Column: "value_float", Value: float64(index) + 10.5}, {Column: "provenance", Value: `{"source":"mysql-repair-fixture"}`},
	}}
	d := r.group.Destination
	scope, err := snapshot.DestinationScope(d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName)
	require.NoError(t, err)
	require.NoError(t, r.store.CommitClosure(t.Context(), r.key, groupdelivery.Closure{
		Destination: groupdelivery.FrozenDestination{Scope: scope, ConnectorID: d.ConnectorID, ConnectorRevision: d.ConnectorRevision,
			Database: d.Database, TableSchema: d.TableSchema, TableName: d.TableName, DedupeCapability: r.group.WritePolicy.DedupeCapability, RecordKeyColumn: r.group.RowPolicy.RecordKeyColumn},
		NextClose: start.Add(15 * time.Second), Buckets: []groupdelivery.ClosedBucket{{
			Outcome: snapshot.Outcome{Kind: snapshot.OutcomeRow, RecordID: recordID, EffectKey: effectKey, BucketStart: start, BucketEnd: start.Add(15 * time.Second)}, Row: &row,
		}},
	}))
	return effectKey
}

func (r *mysqlRepairRig) count(t *testing.T, table string) int {
	t.Helper()
	require.True(t, longTablePattern.MatchString(table) || table == dbtarget.EffectReceiptTable)
	var count int
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+table+"`").Scan(&count))
	return count
}

func (r *mysqlRepairRig) exists(t *testing.T) bool {
	t.Helper()
	var count int
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?", r.group.Destination.TableName).Scan(&count))
	return count == 1
}

func (r *mysqlRepairRig) drop(t *testing.T) {
	t.Helper()
	table := r.group.Destination.TableName
	require.True(t, longTablePattern.MatchString(table))
	_, err := r.f.target.ExecContext(t.Context(), "DROP TABLE `"+table+"`")
	require.NoError(t, err)
}

func (r *mysqlRepairRig) receipts(t *testing.T) map[string][3]string {
	t.Helper()
	rows, err := r.f.target.QueryContext(t.Context(), "SELECT effect_key,payload_digest,committed_at,_gw_effect_receipts_v1 FROM gw_effect_receipts ORDER BY effect_key")
	require.NoError(t, err)
	defer rows.Close()
	result := map[string][3]string{}
	for rows.Next() {
		var key, digest, owner string
		var committed time.Time
		require.NoError(t, rows.Scan(&key, &digest, &committed, &owner))
		result[key] = [3]string{digest, committed.UTC().Format(time.RFC3339Nano), owner}
	}
	require.NoError(t, rows.Err())
	return result
}

func (r *mysqlRepairRig) lineage(t *testing.T) *recordingplan.StorageLineage {
	t.Helper()
	lineage, err := recordingplan.ReadStorageLineage(t.Context(), r.f.db, workspace.StorageScopeOf(r.group))
	require.NoError(t, err)
	return lineage
}
