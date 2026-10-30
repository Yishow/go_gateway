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

type directProfileProtocol struct{ requests []connector.ReadRequest }

func (*directProfileProtocol) Connect(context.Context, string) error               { return nil }
func (*directProfileProtocol) Close() error                                        { return nil }
func (*directProfileProtocol) IsConnected() bool                                   { return true }
func (*directProfileProtocol) ProtocolType() schema.ProtocolType                   { return schema.ProtocolModbusTCP }
func (*directProfileProtocol) TestConnection(context.Context) error                { return nil }
func (*directProfileProtocol) Write(context.Context, connector.WriteRequest) error { return nil }
func (p *directProfileProtocol) Read(_ context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	p.requests = append(p.requests, req)
	return connector.ReadResult{Value: uint32(9), Quality: schema.QualityGood}, nil
}
func TestPointDirectReadUsesScopedProfileAndFailsClosedWithoutAuthority(t *testing.T) {
	old, exists := connector.GetFactory(schema.ProtocolModbusTCP)
	probe := &directProfileProtocol{}
	connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol { return probe })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolModbusTCP, old)
		} else {
			connector.Unregister(schema.ProtocolModbusTCP)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	dr := device.NewMemoryRepository()
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","data_format":"BADC"}`}
	require.NoError(t, dr.Create(t.Context(), d))
	ps := point.NewService(point.NewMemoryRepository(), nil)
	p, err := ps.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "P", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	ds := device.NewService(dr, manager)
	reader := NewPointDirectReader(ps, ds)
	values := reader.PollDirect(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityBad, values[0].Quality)
	require.Empty(t, probe.requests)
	db := openStudioV2CascadeTestDB(t)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
	require.NoError(t, err)
	live.SetWorkspaceProjectionReader(pointProfileProjectionFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
		current, err := ps.GetByID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		return &workspace.RuntimeProjection{WorkspaceID: "W", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: []*schema.Point{current}, Rules: []*schema.SourceRule{{ID: "R", DeviceID: d.ID, RevisionID: "rule-1", DataFormat: "CDAB"}}, RuleLinks: []*schema.SourceRuleLink{{RuleID: "R", PointID: p.ID, Address: p.Address}}}, nil
	}))
	reader.WithReadProfileResolver(live)
	values = reader.PollDirect(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityGood, values[0].Quality)
	require.Len(t, probe.requests, 1)
	require.Equal(t, "CDAB", probe.requests[0].Profile.Resolution().EffectiveFormat)
	// Existing direct preview is allowed for attached draft devices even when
	// the runtime is stopped and no scheduler slot exists.
	d.Status = schema.DeviceStatusDraft
	require.NoError(t, dr.Update(t.Context(), d))
	values = reader.PollDirect(t.Context(), []string{p.ID})
	require.Equal(t, schema.QualityGood, values[0].Quality)
	require.Len(t, probe.requests, 2)
	// The authoritative projection excludes the point: do not treat an explicit
	// per-point format as permission to read another scope/device.
	live.SetWorkspaceProjectionReader(pointProfileProjectionFunc(func(context.Context) (*workspace.RuntimeProjection, error) {
		return &workspace.RuntimeProjection{}, nil
	}))
	values = reader.PollDirect(t.Context(), []string{p.ID})
	require.Equal(t, schema.QualityBad, values[0].Quality)
	require.Len(t, probe.requests, 2)
	// Legacy callers intentionally retain the old direct API without a resolver.
	p.FormatMode = schema.PointFormatLegacy
	p.FormatVersion = 0
	require.NoError(t, ps.Restore(t.Context(), p))
	values = NewPointDirectReader(ps, ds).WithReadProfileResolver((*runtimepkg.Service)(nil)).PollDirect(t.Context(), []string{p.ID})
	require.Equal(t, schema.QualityGood, values[0].Quality)
	require.Len(t, probe.requests, 3)
	require.False(t, probe.requests[2].Profile.Present())
}
