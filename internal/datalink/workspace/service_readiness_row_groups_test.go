package workspace

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type rowGroupReadinessDeviceService struct{}

func (rowGroupReadinessDeviceService) GetByID(context.Context, string) (*schema.Device, error) {
	return &schema.Device{ID: "dev-A", Status: schema.DeviceStatusDraft}, nil
}

func (rowGroupReadinessDeviceService) CheckReadiness(context.Context, string) (*schema.DeviceReadiness, error) {
	return &schema.DeviceReadiness{
		ActivationAllowed: true,
		ConnectStatus:     schema.ReadinessStageStatusSuccess,
		ProbeStatus:       schema.ReadinessStageStatusSuccess,
	}, nil
}

type rowGroupReadinessRuleService struct {
	links []*schema.SourceRuleLink
}

func (s rowGroupReadinessRuleService) ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error) {
	return []*schema.SourceRule{{ID: "rule-A", DeviceID: "dev-A"}}, nil
}

func (s rowGroupReadinessRuleService) ListLinks(context.Context, string) ([]*schema.SourceRuleLink, error) {
	return s.links, nil
}

type rowGroupReadinessConnectorService struct{}

func (rowGroupReadinessConnectorService) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return &schema.DatabaseConnector{
		ID:               "db-main",
		ConnectionConfig: `{"write_mode":"upsert"}`,
		Status:           schema.DatabaseConnectorStatusReady,
		Enabled:          true,
	}, nil
}

type rowGroupReadinessTargetService struct{}

func (rowGroupReadinessTargetService) List(_ context.Context, filter dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	tagID := ""
	if filter.TagID != nil {
		tagID = *filter.TagID
	}
	return []*schema.DatabaseTargetMapping{{
		ID:          "target-" + tagID,
		TagID:       tagID,
		ConnectorID: "db-main",
		TableSchema: "public",
		TableName:   "sensor_readings",
		ColumnName:  "temperature_c",
		WriteMode:   schema.DatabaseWriteModeUpsert,
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}}, nil
}

func TestService_ReadinessBlocksUnsafeRowGroupUpsertReuse(t *testing.T) {
	ctx := context.Background()
	pointA := "point-A"
	pointB := "point-B"
	tagA := "tag-A"
	tagB := "tag-B"
	mappingA := "mapping-A"
	mappingB := "mapping-B"
	workspaceSvc := NewService(NewMemoryRepository()).WithReadinessServices(
		rowGroupReadinessDeviceService{},
		rowGroupReadinessRuleService{links: []*schema.SourceRuleLink{
			{RuleID: "rule-A", PointID: pointA, TagID: &tagA, MappingID: &mappingA},
			{RuleID: "rule-A", PointID: pointB, TagID: &tagB, MappingID: &mappingB},
		}},
		rowGroupReadinessConnectorService{},
		rowGroupReadinessTargetService{},
	)
	_, err := workspaceSvc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-main")
	require.NoError(t, err)
	_, err = workspaceSvc.SaveDatabaseRowGroups(ctx, "db-main", "public", "sensor_readings", []DatabaseRowGroup{{
		ID:              "group-temp",
		ConnectorID:     "db-main",
		TableSchema:     "public",
		TableName:       "sensor_readings",
		MemberPointIDs:  []string{pointA, pointB},
		GroupKeyColumns: []string{"ts"},
	}})
	require.NoError(t, err)
	_, err = workspaceSvc.SaveDatabaseTargetReference(ctx, pointA, "group-temp")
	require.NoError(t, err)
	_, err = workspaceSvc.SaveDatabaseTargetReference(ctx, pointB, "group-temp")
	require.NoError(t, err)

	summary, err := workspaceSvc.Readiness(ctx)
	require.NoError(t, err)
	require.False(t, summary.Ready)
	requireReadinessIssue(t, summary, ReadinessIssue{
		Code:     "database-row-group-upsert-unsafe",
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep4,
		Scope:    "group-temp",
	})
}

func TestService_SaveDatabaseTargetReferenceRejectsNonMemberPoint(t *testing.T) {
	ctx := context.Background()
	workspaceSvc := NewService(NewMemoryRepository())
	_, err := workspaceSvc.SaveDatabaseRowGroups(ctx, "db-main", "public", "sensor_readings", []DatabaseRowGroup{{
		ID:              "group-temp",
		ConnectorID:     "db-main",
		TableSchema:     "public",
		TableName:       "sensor_readings",
		MemberPointIDs:  []string{"point-A"},
		GroupKeyColumns: []string{"ts"},
	}})
	require.NoError(t, err)

	_, err = workspaceSvc.SaveDatabaseTargetReference(ctx, "point-B", "group-temp")
	require.Error(t, err)
	require.ErrorContains(t, err, "database target row group does not contain point")
}

func TestService_SaveDatabaseRowGroupsPrunesRefsForRemovedMembers(t *testing.T) {
	ctx := context.Background()
	workspaceSvc := NewService(NewMemoryRepository())
	_, err := workspaceSvc.SaveDatabaseRowGroups(ctx, "db-main", "public", "sensor_readings", []DatabaseRowGroup{{
		ID:              "group-temp",
		ConnectorID:     "db-main",
		TableSchema:     "public",
		TableName:       "sensor_readings",
		MemberPointIDs:  []string{"point-A", "point-B"},
		GroupKeyColumns: []string{"ts"},
	}})
	require.NoError(t, err)
	_, err = workspaceSvc.SaveDatabaseTargetReference(ctx, "point-A", "group-temp")
	require.NoError(t, err)
	_, err = workspaceSvc.SaveDatabaseTargetReference(ctx, "point-B", "group-temp")
	require.NoError(t, err)

	record, err := workspaceSvc.SaveDatabaseRowGroups(ctx, "db-main", "public", "sensor_readings", []DatabaseRowGroup{{
		ID:              "group-temp",
		ConnectorID:     "db-main",
		TableSchema:     "public",
		TableName:       "sensor_readings",
		MemberPointIDs:  []string{"point-A"},
		GroupKeyColumns: []string{"ts"},
	}})
	require.NoError(t, err)
	require.Equal(t, []DatabaseTargetRef{{PointID: "point-A", RowGroupID: "group-temp"}}, record.DatabaseTargetRefs)
}
