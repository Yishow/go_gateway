package main

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestProductionManagedCrossDeviceRowsKeepMemberOriginsWithoutInventedDeviceID(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedPipelineTarget(t, kind)
			e := f.env
			_, err := e.db.ExecContext(t.Context(), `INSERT INTO devices(id,name,protocol,status,connection_config,readiness_status) VALUES('device-2','Second PLC','modbus_tcp','draft','{}','{}')`)
			require.NoError(t, err)
			_, err = e.db.ExecContext(t.Context(), `UPDATE points SET device_id='device-2' WHERE id='point-p-B'`)
			require.NoError(t, err)
			_, err = e.services.workspace.AttachDevice(t.Context(), "device-2")
			require.NoError(t, err)
			base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
			clock := &testClock{}
			clock.set(base.Add(-time.Second))
			e.services.writeGroups.WithClock(clock.now)
			group := f.prepare(t, []workspace.WriteGroupMember{
				{DeviceID: "device-1", PointID: "point-t-A", TagID: "tag-t-A", Required: true},
				{DeviceID: "device-2", PointID: "point-p-B", TagID: "tag-p-B", Required: true},
			}, workspace.WriteGroupRowPolicy{IntervalSeconds: 10})
			require.Empty(t, group.RowPolicy.DeviceIDColumn)
			inspection, err := e.services.dbTarget.InspectTable(t.Context(), "connector-A", f.namespace, group.Destination.TableName)
			require.NoError(t, err)
			for _, column := range inspection.Columns {
				require.NotEqual(t, "device_id", column.Name, "a cross-device group has no guessed scalar device identity")
			}
			pipe := e.pipeline(clock, "node-1/managed-multidevice")
			clock.set(base.Add(time.Second))
			require.NoError(t, pipe.Reconcile(t.Context()))
			observed := base.Add(2 * time.Second)
			clock.set(observed)
			require.NoError(t, pipe.AcceptSample(t.Context(), e.envelope(group, 0, "device-one", observed, 21.5)))
			require.NoError(t, pipe.AcceptSample(t.Context(), e.envelope(group, 1, "device-two", observed, int64(202))))
			clock.set(base.Add(11 * time.Second))
			pipe.TickAll(t.Context())
			var effect string
			require.NoError(t, e.db.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox`).Scan(&effect))
			sender := groupdelivery.NewSender(groupdelivery.NewStore(e.db), managedAckResolver{service: e.services.dbTarget, db: f.db, kind: f.kind}, groupdelivery.SenderConfig{})
			result, err := sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, result.State)
			var groupID, raw string
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT group_id,provenance FROM `+managedSQLTable(f, group.Destination.TableName)).Scan(&groupID, &raw))
			require.Equal(t, group.ID, groupID)
			var provenance []managedProvenance
			require.NoError(t, json.Unmarshal([]byte(raw), &provenance))
			require.Len(t, provenance, 2)
			require.Contains(t, provenance[0].Member, `"device-1"`)
			require.Contains(t, provenance[1].Member, `"device-2"`)
		})
	}
}
