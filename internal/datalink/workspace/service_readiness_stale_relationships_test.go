package workspace

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestServiceReadinessTreatsStaleRelationshipsAsMissing(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "workspace-readiness-stale.db")
	db := openWorkspaceTestDB(t, dbPath)
	defer db.Close()

	ctx := context.Background()
	workspaceRepo := NewSQLRepository(db)
	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	sourceRuleRepo := sourcerule.NewSQLRepository(db)

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	ruleSvc := sourcerule.NewService(sourceRuleRepo, deviceSvc, pointSvc, nil)
	workspaceSvc := NewService(workspaceRepo).
		WithReadinessServices(deviceSvc, ruleSvc, nil, nil).
		WithReadinessSetupReaders(tagSvc, mappingSvc)

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

	ruleRecord, err := ruleSvc.Create(ctx, sourcerule.CreateRuleRequest{
		ID:           "rule-A",
		DeviceID:     "dev-A",
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeFloat64,
		NamingPrefix: "LINE_",
		Enabled:      true,
	})
	require.NoError(t, err)

	links, err := ruleSvc.ListLinks(ctx, ruleRecord.ID)
	require.NoError(t, err)
	require.Len(t, links, 2)

	validTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.a.second",
		DisplayName: "Line A Second",
		DataType:    schema.DataTypeFloat64,
		Unit:        "C",
		Labels: map[string]string{
			"source":              "source-rule",
			"source_rule_id":      "rule-A",
			"source_rule_address": "40001",
		},
	})
	require.NoError(t, err)

	staleMappingID := "mapping-stale"
	links = links[:1]
	links[0].TagID = &validTag.ID
	links[0].MappingID = &staleMappingID
	require.NoError(t, ruleSvc.ReplaceLinks(ctx, ruleRecord.ID, links))

	summary, err := workspaceSvc.Readiness(ctx)
	require.NoError(t, err)
	require.False(t, summary.Ready)
	require.Equal(t, 1, summary.BlockingCount)
	require.Zero(t, summary.WarningCount)
	requireReadinessIssue(t, summary, ReadinessIssue{
		Code:     "mapping-missing",
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep3,
		Scope:    links[0].PointID,
	})
}
