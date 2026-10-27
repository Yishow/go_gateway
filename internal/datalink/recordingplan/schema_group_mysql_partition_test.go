package recordingplan

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mysqlNativeTestEvidence(t *testing.T) *MySQLPartitionEvidence {
	t.Helper()
	partitions, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.February, 10, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	return partitions
}

func cloneMySQLPartitionTestEvidence(evidence *MySQLPartitionEvidence) *MySQLPartitionEvidence {
	if evidence == nil {
		return nil
	}
	cloned := *evidence
	cloned.Partitions = slices.Clone(evidence.Partitions)
	return &cloned
}

func mysqlNativeTestLayout(t *testing.T) GroupSchemaLayout {
	t.Helper()
	options := mysqlTestOptions()
	options.NativeMonthly = mysqlNativeTestEvidence(t)
	layout, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, options)
	require.NoError(t, err)
	return layout
}

func TestNewMySQLNativeMonthlyPartition(t *testing.T) {
	// The input's local date is March, but its UTC date is still February.
	at := time.Date(2024, time.March, 1, 0, 30, 0, 0, time.FixedZone("UTC+2", 2*60*60))
	got, err := NewMySQLNativeMonthlyPartition(at)
	require.NoError(t, err)
	require.Equal(t, MySQLPartitionEvidence{
		Method:     "RANGE COLUMNS",
		Expression: "bucket_start",
		Partitions: []MySQLPartitionBoundary{
			{Name: "p_past", LessThan: "2024-02-01T00:00:00Z", Ordinal: 1},
			{Name: "p_202402", LessThan: "2024-03-01T00:00:00Z", Ordinal: 2},
			{Name: "p_202403", LessThan: "2024-04-01T00:00:00Z", Ordinal: 3},
			{Name: "p_202404", LessThan: "2024-05-01T00:00:00Z", Ordinal: 4},
			{Name: "p_202405", LessThan: "2024-06-01T00:00:00Z", Ordinal: 5},
			{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 6},
		},
	}, *got)

	nextYear, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.December, 15, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, []MySQLPartitionBoundary{
		{Name: "p_past", LessThan: "2024-12-01T00:00:00Z", Ordinal: 1},
		{Name: "p_202412", LessThan: "2025-01-01T00:00:00Z", Ordinal: 2},
		{Name: "p_202501", LessThan: "2025-02-01T00:00:00Z", Ordinal: 3},
		{Name: "p_202502", LessThan: "2025-03-01T00:00:00Z", Ordinal: 4},
		{Name: "p_202503", LessThan: "2025-04-01T00:00:00Z", Ordinal: 5},
		{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 6},
	}, nextYear.Partitions)

	for name, input := range map[string]time.Time{
		"zero time":        {},
		"before DATETIME":  time.Date(999, time.December, 1, 0, 0, 0, 0, time.UTC),
		"after DATETIME":   time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC),
		"future overflows": time.Date(9999, time.September, 1, 0, 0, 0, 0, time.UTC),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := NewMySQLNativeMonthlyPartition(input)
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
			if name == "zero time" {
				require.Contains(t, err.Error(), "time is zero")
			} else {
				require.Contains(t, err.Error(), "MySQL DATETIME range")
			}
		})
	}
}

func TestMySQLNativeMonthlyDDL(t *testing.T) {
	plain, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mysqlTestOptions())
	require.NoError(t, err)
	plainStatements, err := GenerateGroupSchemaDDL(dialectMySQL, plain)
	require.NoError(t, err)

	partitioned := mysqlNativeTestLayout(t)
	statements, err := GenerateGroupSchemaDDL(dialectMySQL, partitioned)
	require.NoError(t, err)
	require.Len(t, statements, 3)
	require.Contains(t, statements[0], "PARTITION BY RANGE COLUMNS (`bucket_start`)")
	require.Contains(t, statements[0], "PRIMARY KEY (`bucket_start`, `record_id`)")
	for _, boundary := range []string{
		"PARTITION `p_past` VALUES LESS THAN ('2024-02-01 00:00:00')",
		"PARTITION `p_202402` VALUES LESS THAN ('2024-03-01 00:00:00')",
		"PARTITION `p_202405` VALUES LESS THAN ('2024-06-01 00:00:00')",
		"PARTITION `p_future` VALUES LESS THAN (MAXVALUE)",
	} {
		require.Contains(t, statements[0], boundary)
	}
	require.Equal(t, plainStatements[1:], statements[1:], "indexes and receipt DDL stay plain and statement-separated")
	require.NotContains(t, statements[2], "PARTITION BY")
	require.Contains(t, statements[2], "PRIMARY KEY (`effect_key`)")
	require.NotContains(t, statements[2], "FOREIGN KEY")
}

func TestMySQLNativeMonthlyOptionsNormalizationOwnsEvidence(t *testing.T) {
	input := mysqlNativeTestEvidence(t)
	options := mysqlTestOptions()
	options.NativeMonthly = input
	normalized, err := validateMySQLGroupSchemaOptions(options)
	require.NoError(t, err)

	input.Partitions[0].LessThan = "mutated"
	require.Equal(t, "2024-02-01T00:00:00Z", normalized.NativeMonthly.Partitions[0].LessThan)
}

func TestMySQLNativeMonthlyPreviewTokenOwnsPartitionEvidence(t *testing.T) {
	repo := NewMemoryRepository()
	layout := mysqlNativeTestLayout(t)
	token := mysqlTestToken(layout)
	token.Token = "mysql-native-partition-copy"
	require.NoError(t, repo.SavePreviewToken(t.Context(), &token))

	token.GroupLayout.MySQL.NativeMonthly.Partitions[0].LessThan = "caller-mutated"
	stored, err := repo.GetPreviewToken(t.Context(), token.Token)
	require.NoError(t, err)
	require.Equal(t, "2024-02-01T00:00:00Z", stored.GroupLayout.MySQL.NativeMonthly.Partitions[0].LessThan)

	stored.GroupLayout.MySQL.NativeMonthly.Partitions[1].Name = "returned-mutated"
	again, err := repo.GetPreviewToken(t.Context(), token.Token)
	require.NoError(t, err)
	require.Equal(t, "p_202402", again.GroupLayout.MySQL.NativeMonthly.Partitions[1].Name)
}

func TestMySQLNativeMonthlyPartitionEvidenceValidationAndJSON(t *testing.T) {
	layout := mysqlNativeTestLayout(t)
	validated, err := ValidateGroupSchemaLayout(layout)
	require.NoError(t, err)
	require.True(t, groupSchemaLayoutsEqual(layout, validated))

	mutations := map[string]func(*MySQLPartitionEvidence){
		"method":          func(e *MySQLPartitionEvidence) { e.Method = "RANGE" },
		"expression":      func(e *MySQLPartitionEvidence) { e.Expression = "DATE(bucket_start)" },
		"duplicate":       func(e *MySQLPartitionEvidence) { e.Partitions[2] = e.Partitions[1] },
		"ordinal":         func(e *MySQLPartitionEvidence) { e.Partitions[2].Ordinal = 9 },
		"name":            func(e *MySQLPartitionEvidence) { e.Partitions[2].Name = "p_injected" },
		"non-contiguous":  func(e *MySQLPartitionEvidence) { e.Partitions[3].LessThan = "2024-06-01T00:00:00Z" },
		"non-UTC":         func(e *MySQLPartitionEvidence) { e.Partitions[0].LessThan = "2024-02-01T00:00:00+00:00" },
		"past boundary":   func(e *MySQLPartitionEvidence) { e.Partitions[0].LessThan = "2024-02-02T00:00:00Z" },
		"future boundary": func(e *MySQLPartitionEvidence) { e.Partitions[5].LessThan = "2025-01-01T00:00:00Z" },
		"extra partition": func(e *MySQLPartitionEvidence) {
			e.Partitions = append(e.Partitions, MySQLPartitionBoundary{Name: "p_extra", LessThan: "MAXVALUE", Ordinal: 7})
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			changed := layout
			options := *layout.MySQL
			options.NativeMonthly = cloneMySQLPartitionTestEvidence(layout.MySQL.NativeMonthly)
			changed.MySQL = &options
			mutate(changed.MySQL.NativeMonthly)
			_, err := ValidateGroupSchemaLayout(changed)
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
			_, err = GenerateGroupSchemaDDL(dialectMySQL, changed)
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
		})
	}

	badColumn := layout
	badColumn.Columns = slices.Clone(layout.Columns)
	badColumn.Columns[0].SQLType = "DATETIME(3)"
	_, err = ValidateGroupSchemaLayout(badColumn)
	require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

	copyOptions := *layout.MySQL
	copyOptions.NativeMonthly = cloneMySQLPartitionTestEvidence(layout.MySQL.NativeMonthly)
	equivalent := layout
	equivalent.MySQL = &copyOptions
	require.True(t, groupSchemaLayoutsEqual(layout, equivalent), "pointer options compare by content")
	copyOptions.NativeMonthly.Partitions[1].Name = "p_changed"
	require.False(t, groupSchemaLayoutsEqual(layout, equivalent))

	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"native_monthly":{"method":"RANGE COLUMNS","expression":"bucket_start","partitions":[`)
	var roundtrip GroupSchemaLayout
	require.NoError(t, json.Unmarshal(encoded, &roundtrip))
	require.True(t, groupSchemaLayoutsEqual(layout, roundtrip))

	catalog, err := json.Marshal(TargetMySQLTableEvidence{Partitions: mysqlNativeTestEvidence(t)})
	require.NoError(t, err)
	require.Contains(t, string(catalog), `"partitions":{"method":"RANGE COLUMNS","expression":"bucket_start","partitions":[`)
	var catalogRoundtrip TargetMySQLTableEvidence
	require.NoError(t, json.Unmarshal(catalog, &catalogRoundtrip))
	require.True(t, sameMySQLPartitionEvidence(catalogRoundtrip.Partitions, layout.MySQL.NativeMonthly))

	plainCatalog, err := json.Marshal(TargetMySQLTableEvidence{
		ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "cap-1", Database: "gw",
		Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
		PrimaryKey: []string{"bucket_start", "record_id"}, ColumnCollations: map[string]string{},
		ColumnDefaults: map[string]*string{}, JSONColumns: map[string]bool{},
	})
	require.NoError(t, err)
	require.Equal(t,
		`{"ServerFamily":"mysql","ServerVersion":"8.4.5","CapabilityRevision":"cap-1","Database":"gw","Engine":"InnoDB","Charset":"utf8mb4","Collation":"utf8mb4_0900_bin","TableType":"BASE TABLE","PrimaryKey":["bucket_start","record_id"],"ColumnCollations":{},"ColumnDefaults":{},"JSONColumns":{}}`,
		string(plainCatalog), "nil partition evidence preserves plain catalog digest bytes")
	require.NotContains(t, string(plainCatalog), "partitions")
	require.True(t, strings.Contains(string(encoded), `"native_monthly"`))
}

func TestMySQLNativeMonthlyCatalogRequiresFrozenPartitionEvidence(t *testing.T) {
	layout, token, state, _ := mysqlCatalogFixture(t)
	options := *layout.MySQL
	options.NativeMonthly = mysqlNativeTestEvidence(t)
	layout, err := ManagedMySQLLongGroupLayout(layout.TableName, layout.OwnerColumn, options)
	require.NoError(t, err)
	token.GroupLayout = &layout
	state.MySQL.Partitions = cloneMySQLPartitionTestEvidence(layout.MySQL.NativeMonthly)
	statements, err := GenerateGroupSchemaDDL(dialectMySQL, layout)
	require.NoError(t, err)
	inspect := func(context.Context, string) (TargetTableInspection, error) { return state, nil }

	for _, step := range []mysqlApplyStepRef{
		mysqlStepReference(mysqlStepKindGroupTable, layout.TableName, statements[0]),
		mysqlStepReference(mysqlStepKindIndex, layout.Indexes[0].Name, statements[1]),
	} {
		evidence, err := inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
		require.NoError(t, err)
		require.Equal(t, mysqlStepStateExact, evidence.State)
	}

	for name, mutate := range map[string]func(*TargetTableInspection){
		"missing":  func(s *TargetTableInspection) { s.MySQL.Partitions = nil },
		"mismatch": func(s *TargetTableInspection) { s.MySQL.Partitions.Partitions[2].LessThan = "2024-05-01T00:00:00Z" },
	} {
		t.Run(name, func(t *testing.T) {
			bad := cloneTargetTableInspection(state)
			mutate(&bad)
			for _, step := range []mysqlApplyStepRef{
				mysqlStepReference(mysqlStepKindGroupTable, layout.TableName, statements[0]),
				mysqlStepReference(mysqlStepKindIndex, layout.Indexes[0].Name, statements[1]),
			} {
				_, err := inspectMySQLApplyStep(t.Context(), token, layout, step, func(context.Context, string) (TargetTableInspection, error) { return bad, nil })
				require.ErrorIs(t, err, ErrIncompatibleExistingTable)
			}
		})
	}

	plainLayout, plainToken, plainState, plainStatements := mysqlCatalogFixture(t)
	plainState.MySQL.Partitions = cloneMySQLPartitionTestEvidence(layout.MySQL.NativeMonthly)
	_, err = inspectMySQLApplyStep(t.Context(), plainToken, plainLayout,
		mysqlStepReference(mysqlStepKindGroupTable, plainLayout.TableName, plainStatements[0]),
		func(context.Context, string) (TargetTableInspection, error) { return plainState, nil })
	require.ErrorIs(t, err, ErrIncompatibleExistingTable, "an unconfigured native partition catalog is rejected")
}

func TestMySQLNativeMonthlyPartitionEvidenceRejectsReceiptPartitions(t *testing.T) {
	layout, token, group, statements := mysqlCatalogFixture(t)
	receipt := mysqlReceiptStateForApply(group.MySQL)
	step := mysqlStepReference(mysqlStepKindReceiptTable, managedEffectReceiptTable, statements[len(statements)-1])
	receipt.MySQL.Partitions = mysqlNativeTestEvidence(t)
	_, err := inspectMySQLApplyStep(t.Context(), token, layout, step, func(context.Context, string) (TargetTableInspection, error) { return receipt, nil })
	require.ErrorIs(t, err, ErrIncompatibleExistingTable)

	receipt.MySQL.Partitions = nil
	_, err = inspectMySQLApplyStep(t.Context(), token, layout, step, func(context.Context, string) (TargetTableInspection, error) { return receipt, nil })
	require.NoError(t, err)

	nativeOptions := *layout.MySQL
	nativeOptions.NativeMonthly = mysqlNativeTestEvidence(t)
	nativeLayout, err := ManagedMySQLLongGroupLayout(layout.TableName, layout.OwnerColumn, nativeOptions)
	require.NoError(t, err)
	nativeToken := *token
	nativeToken.GroupLayout = &nativeLayout
	nativeStatements, err := GenerateGroupSchemaDDL(dialectMySQL, nativeLayout)
	require.NoError(t, err)
	nativeReceiptStep := mysqlStepReference(mysqlStepKindReceiptTable, managedEffectReceiptTable, nativeStatements[len(nativeStatements)-1])
	receipt.MySQL.Partitions = nil
	_, err = inspectMySQLApplyStep(t.Context(), &nativeToken, nativeLayout, nativeReceiptStep, func(context.Context, string) (TargetTableInspection, error) { return receipt, nil })
	require.NoError(t, err, "the shared receipt remains unpartitioned for a native group layout")
}
