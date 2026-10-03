package grouppipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestCanonicalManagedRecoveryRequiresFrozenRuntimeProof(t *testing.T) {
	cases := []struct {
		name       string
		remove     bool
		mutate     func(*workspace.WriteGroupAppliedSnapshot)
		inspected  bool
		wantStored int
	}{
		{
			name:       "descriptor missing",
			remove:     true,
			inspected:  true,
			wantStored: 0,
		},
		{
			name: "schema digest missing",
			mutate: func(snap *workspace.WriteGroupAppliedSnapshot) {
				snap.RuntimeLayout.SchemaDigest = ""
			},
			wantStored: 1,
		},
		{
			name: "owner proof missing",
			mutate: func(snap *workspace.WriteGroupAppliedSnapshot) {
				columns := make([]dbtarget.ColumnInfo, 0, len(snap.RuntimeLayout.Columns))
				owner := managedOwnerColumnForTest(snap.Group)
				for _, column := range snap.RuntimeLayout.Columns {
					if strings.EqualFold(column.Name, owner) {
						continue
					}
					columns = append(columns, column)
				}
				snap.RuntimeLayout.Columns = columns
			},
			wantStored: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, group, snap := canonicalManagedRecoveryFixture(t)
			p := f.pipeline()
			require.NoError(t, p.Reconcile(t.Context()))
			acceptComplete(t, f, p, group, "accepted", 2)
			var beforeOpen int
			require.NoError(t, rowCountWhere(f.store, &beforeOpen, "wg_delivery_samples", "consumed = 0"))
			require.Equal(t, 2, beforeOpen)

			if tc.remove {
				_, err := f.store.DB().ExecContext(t.Context(), `DELETE FROM wg_runtime_versions`)
				require.NoError(t, err)
			} else {
				require.NotNil(t, tc.mutate)
				payload, err := f.store.RuntimeVersion(t.Context(), groupKey(snap))
				require.NoError(t, err)
				var persisted workspace.WriteGroupAppliedSnapshot
				require.NoError(t, json.Unmarshal(payload, &persisted))
				tc.mutate(&persisted)
				payload, err = json.Marshal(&persisted)
				require.NoError(t, err)
				_, err = f.store.DB().ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload = ?
					WHERE workspace_id = ? AND group_id = ? AND group_revision = ?`,
					string(payload), snap.WorkspaceID, snap.GroupID, snap.AppliedRevision)
				require.NoError(t, err)
			}

			f.groups.groups = []*workspace.WriteGroup{liveGroup(group, workspace.WriteGroupStatusDeleted)}
			if tc.inspected {
				f.inspect.inspection = &dbtarget.TableInspection{
					Status:  dbtarget.TableInspectionExists,
					Columns: snap.RuntimeLayout.Columns,
				}
			} else {
				f.inspect.inspection = nil
			}
			reopenStore(t, f)

			restarted := f.pipeline()
			require.NoError(t, restarted.Reconcile(t.Context()))
			restarted.TickAll(t.Context())
			require.Equal(t, StateBlocked, statusOf(restarted).State)
			require.Equal(t, reasonSnapshotUnavailable, statusOf(restarted).Reason)

			var open int
			require.NoError(t, rowCountWhere(f.store, &open, "wg_delivery_samples", "consumed = 0"))
			require.Equal(t, 2, open, "unproven managed recovery must retain the accepted journal")
			var stored int
			require.NoError(t, f.store.DB().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_runtime_versions`).Scan(&stored))
			require.Equal(t, tc.wantStored, stored, "recovery must not replace or create an unknown managed descriptor")
		})
	}
}

func canonicalManagedRecoveryFixture(t *testing.T) (*fixture, *workspace.WriteGroup, *workspace.WriteGroupAppliedSnapshot) {
	t.Helper()
	f := newFixture(t)
	group := groupAt("rev-1")
	group.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyManaged
	group.Destination.TableName = "managed_readings"
	group.RowPolicy.RecordKeyColumn = "record_id"
	group.RowPolicy.GroupIDColumn = "group_id"
	group.RowPolicy.DeviceIDColumn = "device_id"
	group.RowPolicy.BucketStartColumn = "bucket_start"
	group.RowPolicy.ProvenanceColumn = "provenance"
	layout := &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLite,
		Columns: []dbtarget.ColumnInfo{
			{Name: "record_id", DataType: "TEXT", Nullable: false, PrimaryKey: true},
			{Name: "group_id", DataType: "TEXT", Nullable: false},
			{Name: "device_id", DataType: "TEXT", Nullable: false},
			{Name: "bucket_start", DataType: "TEXT", Nullable: false},
			{Name: "provenance", DataType: "TEXT", Nullable: false},
			{Name: "temperature", DataType: "REAL", Nullable: true},
			{Name: "pressure", DataType: "INTEGER", Nullable: true},
			{Name: managedOwnerColumnForTest(group), DataType: "TEXT", Nullable: false},
		},
		TagTypes: map[string]schema.DataType{
			"tag-t": schema.DataTypeFloat64,
			"tag-p": schema.DataTypeInt64,
		},
		ReceiptTableReady: true,
		SchemaDigest:      strings.Repeat("a", sha256.Size*2),
	}
	snap := snapshotOf(group, at(0))
	snap.RuntimeLayout = layout
	f.groups.groups = []*workspace.WriteGroup{liveGroup(group, workspace.WriteGroupStatusReady)}
	f.groups.snapshots[group.ID] = []*workspace.WriteGroupAppliedSnapshot{snap}
	payload, err := json.Marshal(snap)
	require.NoError(t, err)
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), groupKey(snap), payload))
	f.inspect.inspection.Columns = append([]dbtarget.ColumnInfo(nil), layout.Columns...)
	return f, group, snap
}

func managedOwnerColumnForTest(group *workspace.WriteGroup) string {
	digest := sha256.Sum256([]byte(group.WorkspaceID + "\x00" + group.ID))
	return "_gw_owner_" + hex.EncodeToString(digest[:12])
}
