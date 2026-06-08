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

func TestServiceReadinessRecoversLostLinkIDsFromPersistedPointMapping(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "workspace-readiness-recovery.db")
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
		Count:        1,
		DataType:     schema.DataTypeFloat64,
		NamingPrefix: "LINE_",
		Enabled:      true,
	})
	require.NoError(t, err)

	links, err := ruleSvc.ListLinks(ctx, ruleRecord.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.a.saved",
		DisplayName: "Line A Saved",
		DataType:    schema.DataTypeFloat64,
		Unit:        "C",
		Labels: map[string]string{
			"source":              "source-rule",
			"source_rule_id":      "rule-A",
			"source_rule_address": "40001",
		},
	})
	require.NoError(t, err)
	_, err = mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID: links[0].PointID,
		TagID:   tagRecord.ID,
		Enabled: boolPtr(true),
	})
	require.NoError(t, err)

	summary, err := workspaceSvc.Readiness(ctx)
	require.NoError(t, err)
	require.True(t, summary.Ready)
	require.Zero(t, summary.BlockingCount)
	require.Zero(t, summary.WarningCount)
	for _, issue := range summary.Issues {
		require.NotEqual(t, links[0].PointID, issue.Scope)
	}
}

func TestServiceReadinessRejectsMismatchedOwnedPersistedBinding(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "workspace-readiness-ownership.db")
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
		Count:        1,
		DataType:     schema.DataTypeFloat64,
		NamingPrefix: "LINE_",
		Enabled:      true,
	})
	require.NoError(t, err)

	links, err := ruleSvc.ListLinks(ctx, ruleRecord.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.b.saved",
		DisplayName: "Line B Saved",
		DataType:    schema.DataTypeFloat64,
		Unit:        "C",
		Labels: map[string]string{
			"source":              "source-rule",
			"source_rule_id":      "rule-B",
			"source_rule_address": "40011",
		},
	})
	require.NoError(t, err)
	mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID: links[0].PointID,
		TagID:   tagRecord.ID,
		Enabled: boolPtr(true),
	})
	require.NoError(t, err)
	links[0].TagID = &tagRecord.ID
	links[0].MappingID = &mappingRecord.ID
	require.NoError(t, ruleSvc.ReplaceLinks(ctx, ruleRecord.ID, links))

	summary, err := workspaceSvc.Readiness(ctx)
	require.NoError(t, err)
	require.False(t, summary.Ready)
	require.Equal(t, 1, summary.BlockingCount)
	requireReadinessIssue(t, summary, ReadinessIssue{
		Code:     "tag-missing",
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep3,
		Scope:    links[0].PointID,
	})
}