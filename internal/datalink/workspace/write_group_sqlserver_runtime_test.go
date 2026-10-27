package workspace

import (
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func validSQLServerTestFixture(t *testing.T) (*WriteGroup, *WriteGroupRuntimeLayout) {
	t.Helper()
	group := &WriteGroup{
		WorkspaceID: "ws-1", ID: "grp-1", Revision: "rev-1", AppliedRevision: "rev-1",
		LayoutVersion: WriteGroupLayoutManagedLongV1, StorageMode: WriteGroupStorageModePlain,
		RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond},
		Members: []WriteGroupMember{
			{DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", TargetColumn: "val1"},
			{DeviceID: "device-1", PointID: "point-2", TagID: "tag-2", TargetColumn: "val2"},
		},
		Destination: WriteGroupDestination{
			StorageStrategy:   WriteGroupStorageStrategyManaged,
			ConnectorID:       "conn-1",
			ConnectorRevision: "crev-1",
			Database:          "analytics_db",
			TableSchema:       "dbo",
			SchemaRevision:    "srev-1",
		},
		RowPolicy:   WriteGroupRowPolicy{IntervalSeconds: 15, RecordKeyColumn: "record_id"},
		WritePolicy: WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeReceipt},
	}
	table, owner, ok := LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table

	opts := recordingplan.SQLServerGroupSchemaOptions{
		Database:           group.Destination.Database,
		Schema:             group.Destination.TableSchema,
		Collation:          "Latin1_General_100_BIN2",
		ServerVersion:      "16.0.1000",
		ServerEdition:      "Standard Edition",
		CapabilityRevision: "sqlservercap-unit",
		EngineEdition:      2,
		CompatibilityLevel: 160,
	}

	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, opts)
	require.NoError(t, err)

	groupProof, groupState := sqlServerTestRuntimeCatalog(layout.Columns, layout.Indexes, table, owner, 101, opts)
	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptState := sqlServerTestRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 102, opts)

	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, opts.Database, groupState, receiptState)
	require.NoError(t, err)
	require.NotEmpty(t, digest)

	group.Destination.SchemaDigest = digest

	tagTypes := map[string]schema.DataType{
		"tag-1": schema.DataTypeInt64,
		"tag-2": schema.DataTypeFloat64,
	}

	frozen := &WriteGroupRuntimeLayout{
		Dialect:           dbtarget.SQLDialectSQLServer,
		Columns:           groupProof.Inspection.Columns,
		TagTypes:          tagTypes,
		ReceiptTableReady: true,
		SchemaDigest:      digest,
		SQLServer: &dbtarget.SQLServerManagedRowProof{
			ConnectorID:       group.Destination.ConnectorID,
			ConnectorRevision: group.Destination.ConnectorRevision,
			Group:             groupProof,
			Receipt:           receiptProof,
		},
	}

	return group, frozen
}

func sqlServerTestRuntimeCatalog(columns []recordingplan.GroupSchemaColumn, indexes []recordingplan.GroupSchemaIndex, table, owner string, objectID int64,
	options recordingplan.SQLServerGroupSchemaOptions,
) (*dbtarget.SQLServerManagedInspection, recordingplan.TargetTableInspection) {
	e := &dbtarget.SQLServerManagedTableEvidence{
		Database: options.Database, Schema: options.Schema, ServerVersion: options.ServerVersion, ServerEdition: options.ServerEdition,
		EngineEdition: options.EngineEdition, CompatibilityLevel: options.CompatibilityLevel, CapabilityRevision: options.CapabilityRevision,
		MetadataVerified: true, TableType: "USER_TABLE", ObjectID: objectID, Plain: true,
		ColumnDefaults: map[string]*string{}, ColumnCollations: map[string]string{}, ProvenanceJSONChecked: table != dbtarget.EffectReceiptTable,
	}
	state := recordingplan.TargetTableInspection{
		Status: "exists", ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}, IndexesInspected: true,
	}
	physical := &dbtarget.SQLServerManagedInspection{
		Inspection: &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Schema: options.Schema, Table: table, IndexesInspected: true},
		Evidence:   e, Indexes: []dbtarget.IndexInfo{},
	}
	ownerDefault := "gateway"
	for _, col := range columns {
		e.ColumnDefaults[col.Name] = nil
		if col.Name == owner {
			e.ColumnDefaults[col.Name] = &ownerDefault
		}
		if strings.HasPrefix(col.SQLType, "NVARCHAR(") || strings.HasPrefix(col.SQLType, "CHAR(") {
			e.ColumnCollations[col.Name] = options.Collation
		}
		if col.PrimaryKey {
			e.PrimaryKey = append(e.PrimaryKey, col.Name)
		}
		physical.Inspection.Columns = append(physical.Inspection.Columns, dbtarget.ColumnInfo{
			Name: col.Name, DataType: col.SQLType, Nullable: col.Nullable, PrimaryKey: col.PrimaryKey, Unique: col.PrimaryKey, HasDefault: col.Name == owner,
		})
		state.Columns = append(state.Columns, col.Name)
		state.ColumnTypes[col.Name], state.ColumnNullable[col.Name], state.ColumnPrimaryKey[col.Name] = col.SQLType, col.Nullable, col.PrimaryKey
	}
	for _, idx := range indexes {
		physical.Indexes = append(physical.Indexes, dbtarget.IndexInfo{Name: idx.Name, Columns: idx.Columns})
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: idx.Name, Columns: idx.Columns})
	}
	physical.Inspection.Indexes = physical.Indexes
	state.SQLServer = &recordingplan.TargetSQLServerTableEvidence{
		Database: e.Database, Schema: e.Schema, ServerVersion: e.ServerVersion, ServerEdition: e.ServerEdition,
		EngineEdition: e.EngineEdition, CompatibilityLevel: e.CompatibilityLevel, CapabilityRevision: e.CapabilityRevision,
		MetadataVerified: e.MetadataVerified, TableType: e.TableType, ObjectID: e.ObjectID, Plain: e.Plain, PrimaryKey: e.PrimaryKey,
		ColumnDefaults: e.ColumnDefaults, ColumnCollations: e.ColumnCollations, ProvenanceJSONChecked: e.ProvenanceJSONChecked,
	}
	return physical, state
}

func TestSQLServerRuntimeLayoutValidation(t *testing.T) {
	group, frozen := validSQLServerTestFixture(t)

	t.Run("valid_separate_db_and_schema_passes_and_json_roundtrips", func(t *testing.T) {
		require.NotEqual(t, group.Destination.Database, group.Destination.TableSchema, "Database and schema must be separate")
		require.NoError(t, ValidateSQLServerRuntimeLayout(group, frozen))

		encoded, err := json.Marshal(frozen)
		require.NoError(t, err)

		var persisted WriteGroupRuntimeLayout
		require.NoError(t, json.Unmarshal(encoded, &persisted))
		require.NoError(t, ValidateSQLServerRuntimeLayout(group, &persisted), "Proof must survive JSON roundtrip")
		require.Nil(t, persisted.MySQL, "MySQL proof must be omitted")
		require.Nil(t, persisted.Hypertable, "Hypertable evidence must be omitted")
	})

	t.Run("nil_sqlserver_field_omitted_and_old_json_bytes_regression", func(t *testing.T) {
		withoutSQLServer := &WriteGroupRuntimeLayout{
			Dialect:  dbtarget.SQLDialectSQLServer,
			Columns:  frozen.Columns,
			TagTypes: frozen.TagTypes,
		}
		raw, err := json.Marshal(withoutSQLServer)
		require.NoError(t, err)

		var rawMap map[string]any
		require.NoError(t, json.Unmarshal(raw, &rawMap))
		_, exists := rawMap["sqlserver"]
		require.False(t, exists, "nil SQLServer must be omitted in marshaled JSON")

		oldJSON := `{"dialect":"sqlserver","columns":[],"tag_types":{}}`
		var decoded WriteGroupRuntimeLayout
		require.NoError(t, json.Unmarshal([]byte(oldJSON), &decoded))
		require.Nil(t, decoded.SQLServer)
		require.Error(t, ValidateSQLServerRuntimeLayout(group, &decoded), "old JSON bytes without sqlserver proof must fail closed")
	})

	mutations := map[string]func(*WriteGroupRuntimeLayout){
		"missing proof":                      func(l *WriteGroupRuntimeLayout) { l.SQLServer = nil },
		"missing receipt":                    func(l *WriteGroupRuntimeLayout) { l.SQLServer.Receipt = nil },
		"unverified indexes":                 func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Indexes = nil },
		"receipt has secondary indexes":      func(l *WriteGroupRuntimeLayout) { l.SQLServer.Receipt.Indexes = []dbtarget.IndexInfo{{Name: "idx"}} },
		"changed connector id":               func(l *WriteGroupRuntimeLayout) { l.SQLServer.ConnectorID = "changed" },
		"changed connector revision":         func(l *WriteGroupRuntimeLayout) { l.SQLServer.ConnectorRevision = "changed" },
		"foreign database":                   func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.Database = "foreign" },
		"foreign schema":                     func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.Schema = "foreign" },
		"foreign table":                      func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Inspection.Table = "foreign" },
		"changed capability":                 func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.CapabilityRevision = "changed" },
		"schema digest mismatch":             func(l *WriteGroupRuntimeLayout) { l.SchemaDigest = "changed" },
		"missing tag type":                   func(l *WriteGroupRuntimeLayout) { delete(l.TagTypes, "tag-1") },
		"invalid tag type":                   func(l *WriteGroupRuntimeLayout) { l.TagTypes["tag-1"] = "invalid" },
		"wrong dialect":                      func(l *WriteGroupRuntimeLayout) { l.Dialect = dbtarget.SQLDialectMySQL },
		"hybrid timescale proof present":     func(l *WriteGroupRuntimeLayout) { l.Hypertable = &WriteGroupHypertableEvidence{} },
		"mysql proof present":                func(l *WriteGroupRuntimeLayout) { l.MySQL = &dbtarget.MySQLManagedRowProof{} },
		"unverified metadata":                func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.MetadataVerified = false },
		"evidence plain false":               func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.Plain = false },
		"evidence table type not USER_TABLE": func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.TableType = "VIEW" },
		"columns mismatch": func(l *WriteGroupRuntimeLayout) {
			l.Columns = append(l.Columns, dbtarget.ColumnInfo{Name: "extra"})
		},
		"physical column drift": func(l *WriteGroupRuntimeLayout) {
			for i := range l.SQLServer.Group.Inspection.Columns {
				if l.SQLServer.Group.Inspection.Columns[i].Name == "value_decimal" {
					l.SQLServer.Group.Inspection.Columns[i].DataType = "FLOAT"
				}
			}
		},
	}

	encoded, err := json.Marshal(frozen)
	require.NoError(t, err)

	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			var candidate WriteGroupRuntimeLayout
			require.NoError(t, json.Unmarshal(encoded, &candidate))
			mutate(&candidate)
			require.Error(t, ValidateSQLServerRuntimeLayout(group, &candidate))
		})
	}

	t.Run("storage mode not plain fails", func(t *testing.T) {
		nativeMonthlyGroup := *group
		nativeMonthlyGroup.StorageMode = WriteGroupStorageModeNativeMonthly
		require.Error(t, ValidateSQLServerRuntimeLayout(&nativeMonthlyGroup, frozen))
	})
}
