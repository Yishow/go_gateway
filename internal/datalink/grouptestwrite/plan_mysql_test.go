package grouptestwrite

import (
	"context"
	"strconv"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type mysqlTestGroups struct {
	*fakeGroups
	layout *workspace.WriteGroupRuntimeLayout
	gotID  string
	gotRev string
	err    error
}

func (f *mysqlTestGroups) MySQLTestWriteLayout(_ context.Context, groupID, revision string) (*workspace.WriteGroupRuntimeLayout, error) {
	f.gotID, f.gotRev = groupID, revision
	return f.layout, f.err
}

func mysqlTestGroup() *workspace.WriteGroup {
	return &workspace.WriteGroup{
		ID: "group-1", WorkspaceID: "ws-1", Revision: "rev-1", Status: workspace.WriteGroupStatusReady,
		LayoutVersion: workspace.WriteGroupLayoutManagedLongV1, StorageMode: workspace.WriteGroupStorageModePlain,
		RecordingProfile: &workspace.WriteGroupRecordingProfile{
			Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond,
			SampleReduction: workspace.WriteGroupSampleReductionLatest,
		},
		Members: []workspace.WriteGroupMember{
			{DeviceID: "dev", PointID: "p-int", TagID: "tag-int", TargetColumn: "value_int"},
			{DeviceID: "dev", PointID: "p-text", TagID: "tag-text", TargetColumn: "value_text"},
		},
		Destination: workspace.WriteGroupDestination{
			ConnectorID: "conn", ConnectorRevision: "crev", Database: "db", TableSchema: "db",
			TableName: "gw_long_abc", StorageStrategy: workspace.WriteGroupStorageStrategyManaged,
		},
		RowPolicy: workspace.WriteGroupRowPolicy{
			IntervalSeconds: 15, RecordKeyColumn: "record_id", BucketStartColumn: "bucket_start",
			ProvenanceColumn: "provenance", GroupIDColumn: "group_id", DeviceIDColumn: "device_id",
		},
		WritePolicy: workspace.WriteGroupWritePolicy{DedupeCapability: "receipt"},
	}
}

func mysqlTestInspection(table string) *dbtarget.MySQLManagedInspection {
	return &dbtarget.MySQLManagedInspection{
		Inspection: &dbtarget.TableInspection{
			Status: dbtarget.TableInspectionExists, Schema: "db", Table: table, Indexes: []dbtarget.IndexInfo{}, IndexesInspected: true,
		},
		Evidence: &dbtarget.MySQLManagedTableEvidence{
			ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "cap-1", Database: "db",
			Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
			PrimaryKey: []string{"bucket_start", "record_id"}, ColumnCollations: map[string]string{},
			ColumnDefaults: map[string]*string{}, JSONColumns: map[string]bool{},
		},
		Indexes: []dbtarget.IndexInfo{},
	}
}

func mysqlTestRuntimeLayout(t *testing.T) *workspace.WriteGroupRuntimeLayout {
	t.Helper()
	layout, err := recordingplan.ManagedMySQLLongGroupLayout("gw_long_abc", "_gw_owner_0123456789abcdef01234567", recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "cap-1", Engine: "InnoDB", Collation: "utf8mb4_0900_bin",
	})
	require.NoError(t, err)
	columns := make([]dbtarget.ColumnInfo, 0, len(layout.Columns))
	for _, column := range layout.Columns {
		columns = append(columns, dbtarget.ColumnInfo{Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey})
	}
	group := mysqlTestInspection("gw_long_abc")
	receipt := mysqlTestInspection(dbtarget.EffectReceiptTable)
	return &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectMySQL, Columns: columns,
		TagTypes:          map[string]schema.DataType{"tag-int": schema.DataTypeInt64, "tag-text": schema.DataTypeString},
		ReceiptTableReady: true, SchemaDigest: "schema-digest-1",
		MySQL: &dbtarget.MySQLManagedRowProof{ConnectorID: "conn", ConnectorRevision: "crev", Group: group, Receipt: receipt},
	}
}

func mysqlTestService(t *testing.T, group *workspace.WriteGroup, layout *workspace.WriteGroupRuntimeLayout) (*Service, *mysqlTestGroups) {
	t.Helper()
	groups := &mysqlTestGroups{fakeGroups: &fakeGroups{group: group}, layout: layout}
	svc := New(Dependencies{Groups: groups, Tags: fakeTags{"tag-int": schema.DataTypeInt64, "tag-text": schema.DataTypeString}}, Config{})
	return svc, groups
}

func TestBuildMySQLTestPlanUsesFrozenProof(t *testing.T) {
	group := mysqlTestGroup()
	layout := mysqlTestRuntimeLayout(t)
	svc, groups := mysqlTestService(t, group, layout)
	p, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.NoError(t, err)
	require.NotNil(t, p.mysql)
	require.Equal(t, layout.MySQL, p.mysql.proof)
	require.Equal(t, "group-1", groups.gotID)
	require.Equal(t, "rev-1", groups.gotRev)
	require.Len(t, p.mysql.members, 2)
	require.NotEmpty(t, p.digest)

	rows, err := p.mysqlRows("gw-test-op1", time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC))
	require.NoError(t, err)
	require.Len(t, rows, 2)
	require.NotEqual(t, rows[0].EffectKey, rows[1].EffectKey)
	require.Equal(t, measurement.ExactInt64, p.mysql.members[0].kind)
	require.Equal(t, measurement.ExactText, p.mysql.members[1].kind)
}

func TestBuildMySQLTestPlanRejectsForeignProof(t *testing.T) {
	for name, mutate := range map[string]func(*workspace.WriteGroupRuntimeLayout){
		"connector revision": func(layout *workspace.WriteGroupRuntimeLayout) { layout.MySQL.ConnectorRevision = "other" },
		"group table":        func(layout *workspace.WriteGroupRuntimeLayout) { layout.MySQL.Group.Inspection.Table = "foreign" },
		"missing receipt":    func(layout *workspace.WriteGroupRuntimeLayout) { layout.MySQL.Receipt = nil },
		"receipt not ready":  func(layout *workspace.WriteGroupRuntimeLayout) { layout.ReceiptTableReady = false },
	} {
		t.Run(name, func(t *testing.T) {
			group := mysqlTestGroup()
			layout := mysqlTestRuntimeLayout(t)
			mutate(layout)
			svc, _ := mysqlTestService(t, group, layout)
			_, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
			require.Error(t, err)
		})
	}
}

func TestBuildMySQLTestPlanNativeRequiresPartitions(t *testing.T) {
	group := mysqlTestGroup()
	group.StorageMode = workspace.WriteGroupStorageModeNativeMonthly
	layout := mysqlTestRuntimeLayout(t)
	svc, _ := mysqlTestService(t, group, layout)
	_, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.Error(t, err)
}

func TestBuildMySQLTestPlanNativeOverlayDoesNotChangeDigest(t *testing.T) {
	plainGroup := mysqlTestGroup()
	plainLayout := mysqlTestRuntimeLayout(t)
	plainService, _ := mysqlTestService(t, plainGroup, plainLayout)
	plainPlan, err := plainService.buildMySQLTestPlan(t.Context(), plainGroup, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.NoError(t, err)

	nativeGroup := mysqlTestGroup()
	nativeGroup.StorageMode = workspace.WriteGroupStorageModeNativeMonthly
	nativeLayout := mysqlTestRuntimeLayout(t)
	nativeLayout.MySQL.Group.Evidence.Partitions = &dbtarget.MySQLPartitionEvidence{
		Method: "RANGE COLUMNS", Expression: "bucket_start",
		Partitions: []dbtarget.MySQLPartitionBoundary{
			{Name: "p_past", LessThan: "2026-01-01T00:00:00Z", Ordinal: 1},
			{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 2},
		},
	}
	nativeService, _ := mysqlTestService(t, nativeGroup, nativeLayout)
	nativePlan, err := nativeService.buildMySQLTestPlan(t.Context(), nativeGroup, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.NoError(t, err)
	require.NotNil(t, nativePlan.mysql.proof.Group.Evidence.Partitions)
	require.Equal(t, plainPlan.digest, nativePlan.digest, "physical partition extension is not part of the immutable test-write digest")
}

func TestMySQLTestValueReadbackNormalizesDecimalScale(t *testing.T) {
	expected, err := measurement.NewDecimal("12000.4")
	require.NoError(t, err)
	require.True(t, mysqlTestValueMatches("DECIMAL(65,30)", measurement.ExactDecimal, expected, []byte("12000.400000000000000000000000000000")))
	require.False(t, mysqlTestValueMatches("DECIMAL(65,30)", measurement.ExactDecimal, expected, []byte("12000.400000000000000000000000000001")))
}

func TestMySQLPreviewShowsRealRowsAndExactCleanup(t *testing.T) {
	group := mysqlTestGroup()
	layout := mysqlTestRuntimeLayout(t)
	svc, _ := mysqlTestService(t, group, layout)
	p, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.NoError(t, err)
	values, recordIDs, cleanup, err := mysqlPreview(p, "gw-test-op1", time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, recordIDs, 2)
	require.NotEqual(t, recordIDs[0], recordIDs[1])
	require.Contains(t, cleanup, "previewed record IDs")
	require.LessOrEqual(t, len(cleanup), 256, "the shared client bounds explanatory text separately from exact IDs")
	require.NotContains(t, cleanup, p.ownerColumn+" =")
	ownerValue := ""
	recordValue := ""
	for _, value := range values {
		if value.Column == p.ownerColumn {
			ownerValue = value.Value
		}
		if value.Column == "record_id" {
			recordValue = value.Value
		}
	}
	require.Equal(t, "gateway", ownerValue, "the table owner marker is not the operation owner")
	require.Equal(t, recordIDs[0], recordValue)
}

type mysqlPreviewDestinations struct{}

func (mysqlPreviewDestinations) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"}, nil
}
func (mysqlPreviewDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	return nil, dbtarget.ErrDestinationUnreachable
}

func TestMySQLPreviewDoesNotClaimOperationOwnerInTableMarker(t *testing.T) {
	group := mysqlTestGroup()
	layout := mysqlTestRuntimeLayout(t)
	svc, _ := mysqlTestService(t, group, layout)
	svc.deps.Destinations = mysqlPreviewDestinations{}
	svc.deps.Ledger = recordingplan.NewService(recordingplan.NewMemoryRepository())
	preview, err := svc.Preview(t.Context(), Workspace{ID: "ws-1", Revision: "setup-1"}, group.ID)
	require.NoError(t, err)
	require.Equal(t, "gateway", preview.OwnerValue)
	require.Equal(t, "gateway", preview.Values[len(preview.Values)-1].Value)
	require.Len(t, preview.RecordIDs, len(group.Members))
	require.Contains(t, preview.Cleanup, "previewed record IDs")
	require.NotContains(t, preview.Cleanup, preview.OwnerColumn+" =")
}

func TestMySQLPreviewCleanupTextStaysBoundedForMultipleRows(t *testing.T) {
	group := mysqlTestGroup()
	layout := mysqlTestRuntimeLayout(t)
	group.Members = append(group.Members, workspace.WriteGroupMember{DeviceID: "dev", PointID: "p-third", TagID: "tag-int", TargetColumn: "value_int"})
	svc, _ := mysqlTestService(t, group, layout)
	p, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.NoError(t, err)
	_, ids, cleanup, err := mysqlPreview(p, "gw-test-op1", time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, ids, 3)
	require.LessOrEqual(t, len(cleanup), 256)
}

func TestMySQLTestWriteRejectsOversizedGroupBeforePreview(t *testing.T) {
	group := mysqlTestGroup()
	member := group.Members[0]
	group.Members = make([]workspace.WriteGroupMember, 257)
	for i := range group.Members {
		group.Members[i] = member
		group.Members[i].PointID = "point-" + strconv.Itoa(i)
	}
	svc, groups := mysqlTestService(t, group, mysqlTestRuntimeLayout(t))
	_, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.Error(t, err, "the shared client and cleanup identity arrays are bounded to 256 records")
	require.Empty(t, groups.gotID, "reject the size before external schema proof or any test write")
}
