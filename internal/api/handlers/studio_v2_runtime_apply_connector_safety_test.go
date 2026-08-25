package handlers

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type runtimeApplyConnectorStub struct {
	connector *schema.DatabaseConnector
}

func (s runtimeApplyConnectorStub) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	cloned := *s.connector
	return &cloned, nil
}

func TestRuntimeApplyResponseDoesNotExposeConnectorDiagnostics(t *testing.T) {
	const raw = "dial tcp db.internal:5432: user=writer password=secret dsn=postgres://writer:secret@db.internal/app host=10.0.0.8"
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	lastTestSuccess := true
	require.NoError(t, deviceRepo.Create(ctx, &schema.Device{
		ID:               "dev-A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		LastTestSuccess:  &lastTestSuccess,
	}))

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(
		deviceSvc, nil,
		runtimeApplyConnectorStub{connector: &schema.DatabaseConnector{
			ID:             "db-main",
			Status:         schema.DatabaseConnectorStatusUnreachable,
			Enabled:        true,
			LastCheckError: raw,
		}}, nil,
	)
	_, err := workspaceSvc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-main")
	require.NoError(t, err)

	outcome := resolveStudioV2ScopedRuntimeApplyOutcome(ctx, workspaceSvc, deviceSvc, []string{"dev-A"}, []string{"db-main"})
	require.Equal(t, "applied", outcome.Status)
	require.Len(t, outcome.Issues, 1)
	require.Equal(t, "database connector is currently unreachable", outcome.Issues[0].Message)
	response := mapStudioV2RuntimeApplyResponse(outcome)
	payload, err := json.Marshal(response)
	require.NoError(t, err)
	require.NotContains(t, string(payload), raw)
	require.NotContains(t, string(payload), "db.internal")
	require.NotContains(t, string(payload), "secret")
	require.Contains(t, string(payload), "database connector is currently unreachable")
}
