package workspace

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type connectorIssueReadinessDeviceStub struct{}

func (connectorIssueReadinessDeviceStub) GetByID(context.Context, string) (*schema.Device, error) {
	return &schema.Device{ID: "dev-A"}, nil
}

func (connectorIssueReadinessDeviceStub) CheckReadiness(context.Context, string) (*schema.DeviceReadiness, error) {
	return &schema.DeviceReadiness{
		ActivationAllowed: true,
		ConnectStatus:     schema.ReadinessStageStatusSuccess,
		ProbeStatus:       schema.ReadinessStageStatusSuccess,
	}, nil
}

type connectorIssueReadinessConnectorStub struct {
	connector *schema.DatabaseConnector
}

func (s connectorIssueReadinessConnectorStub) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	cloned := *s.connector
	return &cloned, nil
}

func TestServiceReadinessConnectorIssuesUseFixedMessages(t *testing.T) {
	const raw = "dial tcp db.internal:5432: user=writer password=secret dsn=postgres://writer:secret@db.internal/app host=10.0.0.8"

	tests := []struct {
		status  schema.DatabaseConnectorStatus
		code    string
		message string
	}{
		{schema.DatabaseConnectorStatusUnreachable, "database-connector-unreachable", "database connector is currently unreachable"},
		{schema.DatabaseConnectorStatusAuthFailed, "database-connector-auth-failed", "database connector authentication failed"},
		{schema.DatabaseConnectorStatusError, "database-connector-error", "database connector is currently in an error state"},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			ctx := context.Background()
			workspaceSvc := NewService(NewMemoryRepository()).WithReadinessServices(
				connectorIssueReadinessDeviceStub{}, nil,
				connectorIssueReadinessConnectorStub{connector: &schema.DatabaseConnector{
					ID:             "db-main",
					Status:         test.status,
					Enabled:        true,
					LastCheckError: raw,
				}}, nil,
			)
			_, err := workspaceSvc.AttachDevice(ctx, "dev-A")
			require.NoError(t, err)
			_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-main")
			require.NoError(t, err)

			summary, err := workspaceSvc.Readiness(ctx)
			require.NoError(t, err)
			var issue ReadinessIssue
			for _, candidate := range summary.Issues {
				if candidate.Code == test.code {
					issue = candidate
					break
				}
			}
			require.Equal(t, test.code, issue.Code)
			require.Equal(t, test.message, issue.Message)
			payload, err := json.Marshal(summary)
			require.NoError(t, err)
			require.NotContains(t, string(payload), raw)
			require.NotContains(t, string(payload), "db.internal")
			require.NotContains(t, string(payload), "secret")
		})
	}
}
