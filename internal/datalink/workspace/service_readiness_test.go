package workspace

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestService_ReadinessAggregatesBlockingAndWarningIssues(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "workspace-readiness.db")
	db := openWorkspaceTestDB(t, dbPath)
	defer db.Close()

	ctx := context.Background()
	workspaceRepo := NewSQLRepository(db)
	deviceRepo := device.NewSQLRepository(db)
	connectorRepo := dbtarget.NewSQLConnectorRepository(db)
	dbMappingRepo := dbtarget.NewSQLTargetMappingRepository(db)

	deviceSvc := device.NewService(deviceRepo, nil)
	connectorSvc := dbtarget.NewConnectorService(connectorRepo, dbMappingRepo)
	workspaceSvc := NewService(workspaceRepo).WithReadinessServices(deviceSvc, nil, connectorSvc, nil)

	lastTestSuccess := false
	require.NoError(t, deviceRepo.Create(ctx, &schema.Device{
		ID:               "dev-A",
		Name:             "Device A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		LastTestSuccess:  &lastTestSuccess,
		LastTestError:    "讀取探測失敗: timeout",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
	require.NoError(t, connectorRepo.Create(ctx, &schema.DatabaseConnector{
		ID:               "db-main",
		Name:             "db-main",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"database":"/tmp/readiness-target.db"}`,
		Status:           schema.DatabaseConnectorStatusUnreachable,
		Enabled:          true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))

	_, err := workspaceSvc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-main")
	require.NoError(t, err)

	summary, err := workspaceSvc.Readiness(ctx)
	require.NoError(t, err)
	require.False(t, summary.Ready)
	require.Equal(t, 1, summary.BlockingCount)
	require.Equal(t, 1, summary.WarningCount)
	require.Len(t, summary.Issues, 2)

	require.Equal(t, ReadinessSeverityBlocking, summary.Issues[0].Severity)
	require.Equal(t, "device-probe-required", summary.Issues[0].Code)
	require.Equal(t, ReadinessStep1, summary.Issues[0].Step)
	require.Equal(t, "dev-A", summary.Issues[0].Scope)

	require.Equal(t, ReadinessSeverityWarning, summary.Issues[1].Severity)
	require.Equal(t, "database-connector-unreachable", summary.Issues[1].Code)
	require.Equal(t, ReadinessStep4, summary.Issues[1].Step)
	require.Equal(t, "db-main", summary.Issues[1].Scope)
}

func TestService_ReadinessSurfacesMissingDownstreamRelationships(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "workspace-readiness-downstream.db")
	db := openWorkspaceTestDB(t, dbPath)
	defer db.Close()

	ctx := context.Background()
	workspaceRepo := NewSQLRepository(db)
	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	sourceRuleRepo := sourcerule.NewSQLRepository(db)
	connectorRepo := dbtarget.NewSQLConnectorRepository(db)
	dbMappingRepo := dbtarget.NewSQLTargetMappingRepository(db)

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	ruleSvc := sourcerule.NewService(sourceRuleRepo, deviceSvc, pointSvc, nil)
	connectorSvc := dbtarget.NewConnectorService(connectorRepo, dbMappingRepo)
	dbTargetSvc := dbtarget.NewMappingService(dbMappingRepo, connectorRepo, tagSvc)
	workspaceSvc := NewService(workspaceRepo).WithReadinessServices(deviceSvc, ruleSvc, connectorSvc, dbTargetSvc)

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

	require.NoError(t, connectorRepo.Create(ctx, &schema.DatabaseConnector{
		ID:               "db-main",
		Name:             "db-main",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"database":"/tmp/readiness-target.db"}`,
		Status:           schema.DatabaseConnectorStatusReady,
		Enabled:          true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
	_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-main")
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

	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.a.tag.2",
		DisplayName: "Line A Tag 2",
		DataType:    schema.DataTypeFloat64,
		Unit:        "C",
	})
	require.NoError(t, err)
	mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID: links[1].PointID,
		TagID:   tagRecord.ID,
		Enabled: boolPtr(true),
	})
	require.NoError(t, err)
	links[1].TagID = &tagRecord.ID
	links[1].MappingID = &mappingRecord.ID
	require.NoError(t, ruleSvc.ReplaceLinks(ctx, ruleRecord.ID, links))

	summary, err := workspaceSvc.Readiness(ctx)
	require.NoError(t, err)
	require.False(t, summary.Ready)
	require.Equal(t, 2, summary.BlockingCount)
	require.Zero(t, summary.WarningCount)

	requireReadinessIssue(t, summary, ReadinessIssue{
		Code:     "tag-missing",
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep3,
		Scope:    links[0].PointID,
	})
	requireReadinessIssue(t, summary, ReadinessIssue{
		Code:     "database-target-missing",
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep4,
		Scope:    links[1].PointID,
	})
}

func TestService_ReadinessSurfacesMissingMappingRelationship(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "workspace-readiness-mapping.db")
	db := openWorkspaceTestDB(t, dbPath)
	defer db.Close()

	ctx := context.Background()
	workspaceRepo := NewSQLRepository(db)
	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	sourceRuleRepo := sourcerule.NewSQLRepository(db)

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tagRepo)
	ruleSvc := sourcerule.NewService(sourceRuleRepo, deviceSvc, pointSvc, nil)
	workspaceSvc := NewService(workspaceRepo).WithReadinessServices(deviceSvc, ruleSvc, nil, nil)

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
		Key:         "line.a.tag.1",
		DisplayName: "Line A Tag 1",
		DataType:    schema.DataTypeFloat64,
		Unit:        "C",
	})
	require.NoError(t, err)
	links[0].TagID = &tagRecord.ID
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

func requireReadinessIssue(t *testing.T, summary *ReadinessSummary, expected ReadinessIssue) {
	t.Helper()

	for _, issue := range summary.Issues {
		if issue.Code != expected.Code {
			continue
		}
		if issue.Severity != expected.Severity {
			continue
		}
		if issue.Step != expected.Step {
			continue
		}
		if issue.Scope != expected.Scope {
			continue
		}
		return
	}

	t.Fatalf("expected readiness issue %+v, got %+v", expected, summary.Issues)
}

func boolPtr(value bool) *bool {
	return &value
}
