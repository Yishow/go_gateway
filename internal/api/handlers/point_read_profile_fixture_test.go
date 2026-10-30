package handlers

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	runtimepkg "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

type pointProfileProjectionFunc func(context.Context) (*workspace.RuntimeProjection, error)

func (f pointProfileProjectionFunc) RuntimeProjection(ctx context.Context) (*workspace.RuntimeProjection, error) {
	return f(ctx)
}

// Existing lifecycle fixtures now provide persisted point authority just as the
// production workspace reader does, rather than installing modern intent blind.
type profileFixtureDeviceReader interface {
	GetByID(context.Context, string) (*schema.Device, error)
}

func wirePointProfileFixture(t *testing.T, live *runtimepkg.Service, ps *point.Service, deviceID string, devices ...profileFixtureDeviceReader) {
	t.Helper()
	d := &schema.Device{ID: deviceID, Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","data_format":"ABCD"}`}
	if len(devices) > 0 {
		var err error
		d, err = devices[0].GetByID(t.Context(), deviceID)
		require.NoError(t, err)
	}
	base := workspace.RuntimeProjection{WorkspaceID: "fixture", Version: "fixture", DeviceIDs: []string{deviceID}, Devices: []*schema.Device{d}}
	require.NoError(t, live.ApplyWorkspaceProjection(t.Context(), &base))
	live.SetWorkspaceProjectionReader(pointProfileProjectionFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
		projection := base
		var err error
		if len(devices) > 0 {
			current, err := devices[0].GetByID(ctx, deviceID)
			if err != nil {
				return nil, err
			}
			projection.Devices = []*schema.Device{current}
		}
		projection.Points, err = ps.ListByDevice(ctx, deviceID)
		return &projection, err
	}))
}
