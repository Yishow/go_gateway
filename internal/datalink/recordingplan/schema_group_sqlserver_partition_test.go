package recordingplan

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionDDL(t *testing.T) {
	layout := sqlServerNativeMonthlyLayoutForTest(t)
	statements, err := GenerateGroupSchemaDDL(dialectSQLServer, layout)
	require.NoError(t, err)
	joined := strings.Join(statements, "\n")

	require.Contains(t, joined, "CREATE PARTITION FUNCTION [pf_sqlserver_long] (datetime2(6)) AS RANGE RIGHT FOR VALUES")
	require.Contains(t, joined, "'2026-10-01T00:00:00.000000'")
	require.Contains(t, joined, "'2026-11-01T00:00:00.000000'")
	require.Contains(t, joined, "'2026-12-01T00:00:00.000000'")
	require.Contains(t, joined, "'2027-01-01T00:00:00.000000'")
	require.Contains(t, joined, "'2027-02-01T00:00:00.000000'")
	require.Contains(t, joined, "CREATE PARTITION SCHEME [ps_sqlserver_long] AS PARTITION [pf_sqlserver_long] ALL TO ([PRIMARY])")
	require.Contains(t, joined, "gateway.recording.partition_owner")
	require.Contains(t, joined, "N'dbo.sqlserver_long'")
	require.Contains(t, joined, "N'PARTITION FUNCTION'")
	require.Contains(t, joined, "N'PARTITION SCHEME'")
	require.Contains(t, joined, "ON [ps_sqlserver_long] ([bucket_start])")
	require.Contains(t, joined, "PRIMARY KEY CLUSTERED ([bucket_start], [record_id])")
	require.Contains(t, joined, "CREATE INDEX [idx_sqlserver_long_src_time] ON [dbo].[sqlserver_long]")
	require.Contains(t, joined, "ON [ps_sqlserver_long] ([bucket_start])")
	require.Contains(t, joined, "CREATE TABLE [dbo].[gw_effect_receipts]")
	if strings.Contains(joined, "gw_effect_receipts") {
		receiptPart := joined[strings.LastIndex(joined, "CREATE TABLE [dbo].[gw_effect_receipts]"):]
		require.NotContains(t, receiptPart, "PARTITION")
	}
}

func TestSQLServerPartitionUnavailable(t *testing.T) {
	scope := sqlServerSchemaScopeForTest(t)
	layout := sqlServerNativeJSONLayoutForTest(t, *scope.GroupLayout)
	scope.GroupLayout = &layout

	receipt := sqlServerReceiptStateForTest()
	group := sqlServerMissingStateForTest(*scope.GroupLayout.SQLServer)
	raw, err := json.Marshal(group)
	require.NoError(t, err)
	var object map[string]any
	require.NoError(t, json.Unmarshal(raw, &object))
	sqlServer := object["SQLServer"].(map[string]any)
	sqlServer["partition_inspected"] = true
	sqlServer["native_monthly"] = map[string]any{
		"function_name": "pf_sqlserver_long", "scheme_name": "ps_sqlserver_long",
		"method": "RANGE RIGHT", "column": "bucket_start", "filegroup": "PRIMARY",
		"function_exists": true, "scheme_exists": true, "exclusive": false, "indexes_aligned": false,
		"partitions": sqlServerNativePartitionsJSON(),
	}
	raw, err = json.Marshal(object)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &group))

	svc := NewService(NewMemoryRepository())
	_, err = svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return receipt, nil
		}
		return group, nil
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "partition_unavailable")
}

func sqlServerNativeMonthlyLayoutForTest(t *testing.T) GroupSchemaLayout {
	t.Helper()
	layout, err := ManagedSQLServerLongGroupLayout("sqlserver_long", "_gw_owner_0123456789abcdef01234567", sqlServerSchemaOptionsForTest())
	require.NoError(t, err)
	return sqlServerNativeJSONLayoutForTest(t, layout)
}

func sqlServerNativeJSONLayoutForTest(t *testing.T, layout GroupSchemaLayout) GroupSchemaLayout {
	t.Helper()
	raw, err := json.Marshal(layout)
	require.NoError(t, err)
	var object map[string]any
	require.NoError(t, json.Unmarshal(raw, &object))
	sqlServer := object["sqlserver"].(map[string]any)
	sqlServer["native_monthly"] = map[string]any{
		"method": "RANGE RIGHT", "column": "bucket_start",
		"function_name": "pf_sqlserver_long", "scheme_name": "ps_sqlserver_long", "filegroup": "PRIMARY",
		"partitions": sqlServerNativePartitionsJSON(),
	}
	raw, err = json.Marshal(object)
	require.NoError(t, err)
	var decoded GroupSchemaLayout
	require.NoError(t, json.Unmarshal(raw, &decoded))
	return decoded
}

func sqlServerNativePartitionsJSON() []any {
	partitions := []any{
		map[string]any{"name": "p_past", "less_than": "2026-10-01T00:00:00.000000", "ordinal": 1},
		map[string]any{"name": "p_202610", "less_than": "2026-11-01T00:00:00.000000", "ordinal": 2},
		map[string]any{"name": "p_202611", "less_than": "2026-12-01T00:00:00.000000", "ordinal": 3},
		map[string]any{"name": "p_202612", "less_than": "2027-01-01T00:00:00.000000", "ordinal": 4},
		map[string]any{"name": "p_202701", "less_than": "2027-02-01T00:00:00.000000", "ordinal": 5},
		map[string]any{"name": "p_future", "less_than": "MAXVALUE", "ordinal": 6},
	}
	return partitions
}
