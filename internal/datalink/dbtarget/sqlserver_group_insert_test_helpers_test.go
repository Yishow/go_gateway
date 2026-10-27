package dbtarget

import (
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	sqlServerInsertTestOwner        = "_gw_owner_0123456789abcdef01234567"
	sqlServerInsertTestReceiptOwner = "_gw_effect_receipts_v1"
	sqlServerInsertTestCollation    = "Latin1_General_100_BIN2"
)

func sqlServerGroupInsertTestRequest() GroupInsertRequest {
	const groupTable = "gw_long_fixture"
	return GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLServer, SchemaName: "dbo", TableName: groupTable,
		Strategy: GroupEffectReceipt, PayloadDigest: strings.Repeat("c", 64),
		CommittedAt: "2026-01-02T03:04:05.123456Z",
		Row:         sqlServerGroupInsertTestRow(),
		SQLServer: &SQLServerManagedRowProof{
			ConnectorID: "connector-fixture", ConnectorRevision: "revision-fixture",
			Group:   sqlServerInsertTestInspection(groupTable, false, 91),
			Receipt: sqlServerInsertTestInspection(EffectReceiptTable, true, 92),
		},
	}
}

func sqlServerGroupInsertTestRow() EncodedRow {
	bucketStart := time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC)
	values := map[string]any{
		longColBucketStart:       bucketStart,
		longColBucketEnd:         bucketStart.Add(time.Minute),
		longColRecordID:          strings.Repeat("a", 64),
		longColWorkspaceID:       "workspace-1",
		longColGroupID:           "group-1",
		longColGroupRev:          "revision-1",
		longColDeviceID:          "device-1",
		longColPointID:           "point-1",
		longColTagID:             "Tag-電力🟢",
		longColProfileKind:       "interval",
		longColValueType:         "int64",
		longColStatus:            "missing",
		longColProvenance:        `{"type":"int64","member":"member-1","profile":{"kind":"interval","interval":60,"unit":"minute"}}`,
		sqlServerInsertTestOwner: "gateway",
	}
	cells := make([]EncodedCell, 0, len(longColumnRoles)+1)
	for _, column := range longColumnRoles {
		cells = append(cells, EncodedCell{Column: column.name, Value: values[column.name]})
	}
	cells = append(cells, EncodedCell{Column: sqlServerInsertTestOwner, Value: values[sqlServerInsertTestOwner]})
	return EncodedRow{
		RecordID: strings.Repeat("a", 64), EffectKey: strings.Repeat("b", 64),
		BucketStart: bucketStart, Cells: cells,
	}
}

func sqlServerInsertTestInspection(table string, receipt bool, objectID int64) *SQLServerManagedInspection {
	columns := make([]ColumnInfo, 0, len(longColumnRoles)+1)
	primaryKey := []string{longColBucketStart, longColRecordID}
	indexes := []IndexInfo{{
		Name:    "idx_" + table + "_src_time",
		Columns: []string{longColWorkspaceID, longColTagID, longColBucketStart},
	}}
	owner := sqlServerInsertTestOwner
	if receipt {
		columns = []ColumnInfo{
			{Name: "effect_key", DataType: "CHAR(64)", PrimaryKey: true, Unique: true},
			{Name: "payload_digest", DataType: "CHAR(64)"},
			{Name: "committed_at", DataType: "DATETIME2(6)"},
			{Name: sqlServerInsertTestReceiptOwner, DataType: "NVARCHAR(16)", HasDefault: true},
		}
		primaryKey = []string{"effect_key"}
		indexes = []IndexInfo{}
		owner = sqlServerInsertTestReceiptOwner
	} else {
		for _, role := range longColumnRoles {
			column := ColumnInfo{
				Name: role.name, DataType: sqlServerInsertTestType(role.name, role.role), Nullable: !role.required,
			}
			if role.name == longColBucketStart || role.name == longColRecordID {
				column.PrimaryKey, column.Unique = true, true
			}
			columns = append(columns, column)
		}
		columns = append(columns, ColumnInfo{Name: owner, DataType: "NVARCHAR(16)", HasDefault: true})
	}
	collations := make(map[string]string)
	defaults := make(map[string]*string, len(columns))
	for _, column := range columns {
		if strings.HasPrefix(column.DataType, "NVARCHAR(") || strings.HasPrefix(column.DataType, "CHAR(") {
			collations[column.Name] = sqlServerInsertTestCollation
		}
		defaults[column.Name] = nil
	}
	ownerDefault := "gateway"
	defaults[owner] = &ownerDefault
	inspectionIndexes := append([]IndexInfo{}, indexes...)
	return &SQLServerManagedInspection{
		Inspection: &TableInspection{
			Status: TableInspectionExists, Schema: "dbo", Table: table,
			Columns: columns, Indexes: inspectionIndexes, IndexesInspected: true,
		},
		Evidence: &SQLServerManagedTableEvidence{
			ServerVersion: "16.0.1000.6", ServerEdition: "Developer Edition (64-bit)",
			CapabilityRevision: "capability-fixture", Database: "fixture_db", Schema: "dbo",
			EngineEdition: 3, CompatibilityLevel: 160, MetadataVerified: true,
			TableType: "USER_TABLE", ObjectID: objectID, Plain: true,
			PrimaryKey: primaryKey, ColumnCollations: collations, ColumnDefaults: defaults,
			ProvenanceJSONChecked: !receipt,
		},
		Indexes: indexes,
	}
}

func sqlServerInsertTestType(name string, role longColumnRole) string {
	switch role {
	case longTime:
		return "DATETIME2(6)"
	case longText:
		if name == longColRecordID {
			return "CHAR(64)"
		}
		return "NVARCHAR(191)"
	case longProvenanceRole, longSlotText:
		return "NVARCHAR(MAX)"
	case longSlotInt:
		return "BIGINT"
	case longSlotDecimal, longSlotDelta:
		return "DECIMAL(38,18)"
	case longSlotFloat:
		return "FLOAT(53)"
	case longSlotBool:
		return "BIT"
	default:
		return ""
	}
}
