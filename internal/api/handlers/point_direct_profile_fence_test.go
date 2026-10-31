package handlers

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	runtimepkg "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

func TestDirectHTTPReadAcceptanceUsesCapturedProof(t *testing.T) {
	for _, batch := range []bool{false, true} {
		for _, scenario := range []string{"normal draft", "point changed", "absent invalidated", "foreign install", "same revision rebind", "legacy to modern"} {
			t.Run(scenario+map[bool]string{false: "/single", true: "/batch"}[batch], func(t *testing.T) {
				old, exists := connector.GetFactory(schema.ProtocolModbusTCP)
				probe := &lateManualProfileProtocol{hook: func() {}}
				connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol { return probe })
				t.Cleanup(func() {
					if exists {
						connector.Register(schema.ProtocolModbusTCP, old)
					} else {
						connector.Unregister(schema.ProtocolModbusTCP)
					}
				})
				manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
				t.Cleanup(func() { manager.CloseAll() })
				dr := device.NewMemoryRepository()
				d := &schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","data_format":"BADC"}`}
				require.NoError(t, dr.Create(t.Context(), d))
				ps := point.NewService(point.NewMemoryRepository(), nil)
				p, err := ps.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "P", Address: "HR0", DataType: schema.DataTypeUint32})
				require.NoError(t, err)
				if scenario == "legacy to modern" {
					p.FormatMode = schema.PointFormatLegacy
					p.FormatVersion = 0
					require.NoError(t, ps.Restore(t.Context(), p))
				}
				db := openStudioV2CascadeTestDB(t)
				t.Cleanup(func() { db.Close() })
				live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
				require.NoError(t, err)
				live.SetWorkspaceProjectionReader(pointProfileProjectionFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
					current, err := ps.GetByID(ctx, p.ID)
					if err != nil {
						return nil, err
					}
					return &workspace.RuntimeProjection{WorkspaceID: "W", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: []*schema.Point{current}}, nil
				}))
				require.NoError(t, live.UpsertDevice(t.Context(), d))
				if scenario == "same revision rebind" {
					require.NotZero(t, live.InstallPointForMutation(p))
				}
				if scenario == "normal draft" {
					d.Status = schema.DeviceStatusDraft
					require.NoError(t, dr.Update(t.Context(), d))
				}
				probe.hook = func() {
					switch scenario {
					case "point changed":
						p.FormatMode = schema.PointFormatOverride
						p.FormatOverride = "DCBA"
						p.RevisionID = "next"
						require.NoError(t, ps.Restore(t.Context(), p))
					case "absent invalidated":
						require.True(t, live.SuspendPointRevision(p.ID, p.RevisionID))
					case "foreign install", "same revision rebind":
						require.NotZero(t, live.InstallPointForMutation(p))
					case "legacy to modern":
						p.FormatMode = schema.PointFormatInherit
						p.FormatVersion = 1
						require.NoError(t, ps.Restore(t.Context(), p))
					}
				}
				h := NewPointHandler(ps, nil).WithDirectReader(NewPointDirectReader(ps, device.NewService(dr, manager)).WithReadProfileResolver(live))
				var got PollResult
				if batch {
					results, err := h.pollBatch(t.Context(), []string{p.ID})
					require.NoError(t, err)
					require.Len(t, results, 1)
					got = results[0]
				} else {
					got, err = h.pollPoint(t.Context(), p.ID)
					require.NoError(t, err)
				}
				persisted, err := ps.GetByID(t.Context(), p.ID)
				require.NoError(t, err)
				if scenario == "normal draft" {
					require.Equal(t, pollQualityGood, got.Quality)
					require.Equal(t, uint32(77), got.Value)
					require.NotNil(t, persisted.LastValue)
				} else {
					require.NotEqual(t, pollQualityGood, got.Quality)
					require.Nil(t, got.Value)
					require.Nil(t, persisted.LastValue)
				}
			})
		}
	}
}
