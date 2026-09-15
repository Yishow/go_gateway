package workspace

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestServiceReadinessTreatsMissingConnectorBindingAsBlockingIssue(t *testing.T) {
	ctx := context.Background()
	workspaceSvc := NewService(NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	workspaceSvc.WithReadinessServices(deviceSvc, nil, connectorMissingStub{}, nil)

	lastTestSuccess := true
	require.NoError(t, deviceRepo.Create(ctx, &schema.Device{
		ID:               "dev-A",
		Name:             "Device A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		LastTestSuccess:  &lastTestSuccess,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
	_, err := workspaceSvc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-missing")
	require.NoError(t, err)

	summary, err := workspaceSvc.Readiness(ctx)
	require.NoError(t, err)
	require.False(t, summary.Ready)
	require.Equal(t, 1, summary.BlockingCount)
	requireReadinessIssue(t, summary, ReadinessIssue{
		Code:     "database-connector-missing",
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep4,
		Scope:    "db-missing",
	})
}

type connectorMissingStub struct{}

func (connectorMissingStub) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return nil, errors.New("資料庫連接器不存在")
}
