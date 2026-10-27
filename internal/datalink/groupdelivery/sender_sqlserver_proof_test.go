package groupdelivery

import (
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func senderTestSQLServerProof(table string) *dbtarget.SQLServerManagedRowProof {
	return &dbtarget.SQLServerManagedRowProof{
		ConnectorID:       "connector-1",
		ConnectorRevision: "crev-1",
		Group:             senderTestSQLServerInspection(table, false, 91),
		Receipt:           senderTestSQLServerInspection(dbtarget.EffectReceiptTable, true, 92),
	}
}

func senderTestSQLServerInspection(table string, receipt bool, objectID int64) *dbtarget.SQLServerManagedInspection {
	const owner = sqlServerTestOwnerColumn
	const receiptOwner = "_gw_effect_receipts_v1"
	const collation = sqlServerTestCollation

	if receipt {
		cols := []dbtarget.ColumnInfo{
			{Name: "effect_key", DataType: "CHAR(64)", PrimaryKey: true, Unique: true},
			{Name: "payload_digest", DataType: "CHAR(64)"},
			{Name: "committed_at", DataType: "DATETIME2(6)"},
			{Name: receiptOwner, DataType: "NVARCHAR(16)", HasDefault: true},
		}
		collations := map[string]string{
			"effect_key":     collation,
			"payload_digest": collation,
			receiptOwner:     collation,
		}
		gateway := "gateway"
		defaults := map[string]*string{
			"effect_key":     nil,
			"payload_digest": nil,
			"committed_at":   nil,
			receiptOwner:     &gateway,
		}
		return &dbtarget.SQLServerManagedInspection{
			Inspection: &dbtarget.TableInspection{
				Status: dbtarget.TableInspectionExists, Schema: "dbo", Table: table,
				Columns: cols, Indexes: []dbtarget.IndexInfo{}, IndexesInspected: true,
			},
			Evidence: &dbtarget.SQLServerManagedTableEvidence{
				ServerVersion: "16.0.4180.1", ServerEdition: "Developer Edition (64-bit)",
				CapabilityRevision: "cap-1", Database: "gwtest", Schema: "dbo",
				EngineEdition: 3, CompatibilityLevel: 160, MetadataVerified: true,
				TableType: "USER_TABLE", ObjectID: objectID, Plain: true,
				PrimaryKey: []string{"effect_key"}, ColumnCollations: collations, ColumnDefaults: defaults,
				ProvenanceJSONChecked: false,
			},
			Indexes: []dbtarget.IndexInfo{},
		}
	}

	type colSpec struct {
		name     string
		dataType string
		pk       bool
		nullable bool
	}
	specs := []colSpec{
		{"bucket_start", "DATETIME2(6)", true, false},
		{"bucket_end", "DATETIME2(6)", false, false},
		{"record_id", "CHAR(64)", true, false},
		{"workspace_id", "NVARCHAR(191)", false, false},
		{"group_id", "NVARCHAR(191)", false, false},
		{"group_revision", "NVARCHAR(191)", false, false},
		{"circuit_id", "NVARCHAR(191)", false, true},
		{"device_id", "NVARCHAR(191)", false, false},
		{"point_id", "NVARCHAR(191)", false, false},
		{"tag_id", "NVARCHAR(191)", false, false},
		{"measurement_id", "NVARCHAR(191)", false, true},
		{"series_epoch", "NVARCHAR(191)", false, true},
		{"profile_kind", "NVARCHAR(191)", false, false},
		{"value_type", "NVARCHAR(191)", false, false},
		{"unit", "NVARCHAR(191)", false, true},
		{"status", "NVARCHAR(191)", false, false},
		{"quality", "NVARCHAR(191)", false, true},
		{"observed_at", "DATETIME2(6)", false, true},
		{"received_at", "DATETIME2(6)", false, true},
		{"value_int", "BIGINT", false, true},
		{"value_decimal", "DECIMAL(38,18)", false, true},
		{"value_float", "FLOAT(53)", false, true},
		{"value_bool", "BIT", false, true},
		{"value_text", "NVARCHAR(MAX)", false, true},
		{"delta_value", "DECIMAL(38,18)", false, true},
		{"delta_from", "DATETIME2(6)", false, true},
		{"delta_to", "DATETIME2(6)", false, true},
		{"provenance", "NVARCHAR(MAX)", false, false},
		{owner, "NVARCHAR(16)", false, false},
	}
	cols := make([]dbtarget.ColumnInfo, 0, len(specs))
	collations := make(map[string]string)
	defaults := make(map[string]*string)
	gateway := "gateway"
	for _, s := range specs {
		c := dbtarget.ColumnInfo{
			Name:       s.name,
			DataType:   s.dataType,
			PrimaryKey: s.pk,
			Unique:     s.pk,
			Nullable:   s.nullable,
			HasDefault: s.name == owner,
		}
		cols = append(cols, c)
		if strings.HasPrefix(s.dataType, "NVARCHAR(") || strings.HasPrefix(s.dataType, "CHAR(") {
			collations[s.name] = collation
		}
		if s.name == owner {
			defaults[s.name] = &gateway
		} else {
			defaults[s.name] = nil
		}
	}
	indexes := []dbtarget.IndexInfo{{
		Name:    "idx_" + table + "_src_time",
		Columns: []string{"workspace_id", "tag_id", "bucket_start"},
	}}
	return &dbtarget.SQLServerManagedInspection{
		Inspection: &dbtarget.TableInspection{
			Status: dbtarget.TableInspectionExists, Schema: "dbo", Table: table,
			Columns: cols, Indexes: indexes, IndexesInspected: true,
		},
		Evidence: &dbtarget.SQLServerManagedTableEvidence{
			ServerVersion: "16.0.4180.1", ServerEdition: "Developer Edition (64-bit)",
			CapabilityRevision: "cap-1", Database: "gwtest", Schema: "dbo",
			EngineEdition: 3, CompatibilityLevel: 160, MetadataVerified: true,
			TableType: "USER_TABLE", ObjectID: objectID, Plain: true,
			PrimaryKey:       []string{"bucket_start", "record_id"},
			ColumnCollations: collations, ColumnDefaults: defaults,
			ProvenanceJSONChecked: true,
		},
		Indexes: indexes,
	}
}

func enqueueCanonicalSQLServerLong(t *testing.T, store *Store, table, effectKey string) ClosedBucket {
	t.Helper()
	bucketTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	recordID := strings.Repeat("a", 64)
	if effectKey == "" {
		effectKey = strings.Repeat("b", 64)
	}
	cells := []dbtarget.EncodedCell{
		{Column: "bucket_start", Value: bucketTime},
		{Column: "bucket_end", Value: bucketTime.Add(15 * time.Second)},
		{Column: "record_id", Value: recordID},
		{Column: "workspace_id", Value: "ws-live"},
		{Column: "group_id", Value: "group-live"},
		{Column: "group_revision", Value: "rev-1"},
		{Column: "circuit_id", Value: nil},
		{Column: "device_id", Value: "device-live"},
		{Column: "point_id", Value: "point-live"},
		{Column: "tag_id", Value: "tag-live"},
		{Column: "measurement_id", Value: nil},
		{Column: "series_epoch", Value: nil},
		{Column: "profile_kind", Value: "instantaneous"},
		{Column: "value_type", Value: "int64"},
		{Column: "unit", Value: nil},
		{Column: "status", Value: "ok"},
		{Column: "quality", Value: nil},
		{Column: "observed_at", Value: nil},
		{Column: "received_at", Value: nil},
		{Column: "value_int", Value: int64(100)},
		{Column: "value_decimal", Value: nil},
		{Column: "value_float", Value: nil},
		{Column: "value_bool", Value: nil},
		{Column: "value_text", Value: nil},
		{Column: "delta_value", Value: nil},
		{Column: "delta_from", Value: nil},
		{Column: "delta_to", Value: nil},
		{Column: "provenance", Value: `{"type":"int64","member":"power","profile":{"kind":"instantaneous","interval":10,"unit":"W"}}`},
		{Column: sqlServerTestOwnerColumn, Value: "gateway"},
	}
	row := dbtarget.EncodedRow{
		RecordID:    recordID,
		EffectKey:   effectKey,
		BucketStart: bucketTime,
		Cells:       cells,
	}
	destination := FrozenDestination{
		ConnectorID:       "connector-1",
		ConnectorRevision: "crev-1",
		Database:          "gwtest",
		TableSchema:       "dbo",
		TableName:         table,
		DedupeCapability:  DedupeReceipt,
	}
	var err error
	destination.Scope, err = snapshot.DestinationScope(destination.ConnectorID, destination.ConnectorRevision,
		destination.Database, destination.TableSchema, destination.TableName)
	require.NoError(t, err)

	bucket := ClosedBucket{
		Outcome: snapshot.Outcome{
			Kind:        snapshot.OutcomeRow,
			RecordID:    recordID,
			EffectKey:   effectKey,
			BucketStart: bucketTime,
			BucketEnd:   bucketTime.Add(15 * time.Second),
			Members: []snapshot.MemberResult{{
				MemberKey: "power",
				Status:    snapshot.MemberOK,
				Sample: &snapshot.Sample{
					SampleID:   "s-1",
					MemberKey:  "power",
					ObservedAt: bucketTime,
					Quality:    schema.QualityGood,
					Value:      measurement.NewInt64(100),
				},
			}},
		},
		Row: &row,
	}
	closure := Closure{
		Destination: destination,
		NextClose:   bucketTime.Add(15 * time.Second),
		Buckets:     []ClosedBucket{bucket},
	}
	groupKey := GroupKey{
		WorkspaceID:   "ws-live",
		GroupID:       "group-live",
		GroupRevision: "rev-1",
	}
	require.NoError(t, store.CommitClosure(t.Context(), groupKey, closure))
	return bucket
}
