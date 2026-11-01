package device

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionDeclarationRejectsInvalidBeforePersist(t *testing.T) {
	for name, declaration := range map[string]any{
		"null":             nil,
		"unknown field":    map[string]any{"version": 1, "profile": "safe-v1", "random_support": "supported"},
		"unknown version":  map[string]any{"version": 2, "profile": "safe-v1"},
		"unknown profile":  map[string]any{"version": 1, "profile": "all-mc-supported"},
		"fractional limit": map[string]any{"version": 1, "profile": "safe-v1", "limits": map[string]any{"native_units": 1.5}},
		"zero limit":       map[string]any{"version": 1, "profile": "safe-v1", "limits": map[string]any{"request_bytes": 0}},
		"negative limit":   map[string]any{"version": 1, "profile": "safe-v1", "limits": map[string]any{"response_bytes": -1}},
		"null limits":      map[string]any{"version": 1, "profile": "safe-v1", "limits": nil},
		"foreign family":   map[string]any{"version": 1, "profile": "fatek-fbs-ascii-v1"},
	} {
		t.Run(name, func(t *testing.T) {
			config := map[string]any{"host": "127.0.0.1", "port": 5000, "acquisition": declaration}
			require.Error(t, validateConnectionConfig(schema.ProtocolMC3E, config))
		})
	}
}

func TestAcquisitionDeclarationMissingPreservesLegacyValidation(t *testing.T) {
	require.NoError(t, validateConnectionConfig(schema.ProtocolMC3E, map[string]any{"host": "127.0.0.1", "port": 5000}))
}

func TestAcquisitionDeclarationServiceRejectsWithoutPersisting(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)
	config := map[string]any{"host": "127.0.0.1", "port": 5000, "acquisition": nil}
	_, err := svc.Create(ctx, CreateDeviceRequest{ID: "acquisition-device", Name: "PLC", Protocol: schema.ProtocolMC3E, ConnectionConfig: config})
	require.ErrorContains(t, err, connector.ErrAcquisitionConfig.Error())
	devices, err := repo.List(ctx, ListFilter{})
	require.NoError(t, err)
	require.Empty(t, devices)
	config["acquisition"] = map[string]any{"version": 1, "profile": "safe-v1"}
	created, err := svc.Create(ctx, CreateDeviceRequest{ID: "acquisition-device", Name: "PLC", Protocol: schema.ProtocolMC3E, ConnectionConfig: config})
	require.NoError(t, err)
	resolved, err := connector.ResolveAcquisitionConfig(created.Protocol, created.ConnectionConfig)
	require.NoError(t, err)
	require.Equal(t, "unsupported", resolved.Capability(16).RandomSupport)
	config["acquisition"] = map[string]any{"version": 1, "profile": "not-a-family"}
	_, err = svc.Update(ctx, created.ID, UpdateDeviceRequest{ConnectionConfig: config})
	require.ErrorContains(t, err, connector.ErrAcquisitionConfig.Error())
	persisted, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ConnectionConfig, persisted.ConnectionConfig)
}
