package dbtarget

import (
	"encoding/json"
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloneSQLServerManagedRowProofClonesNativeMonthlyPartitionEvidence(t *testing.T) {
	original := &SQLServerManagedRowProof{Group: &SQLServerManagedInspection{
		Evidence: &SQLServerManagedTableEvidence{NativeMonthly: sqlServerInsertTestNativeMonthlyEvidence("gw_long_fixture")},
	}}

	cloned := CloneSQLServerManagedRowProof(original)
	cloned.Group.Evidence.NativeMonthly.Partitions[1].Name = "changed"
	require.Equal(t, "p_202601", original.Group.Evidence.NativeMonthly.Partitions[1].Name)
	require.NotSame(t, original.Group.Evidence.NativeMonthly, cloned.Group.Evidence.NativeMonthly)
	require.NotSame(t, &original.Group.Evidence.NativeMonthly.Partitions[0], &cloned.Group.Evidence.NativeMonthly.Partitions[0])

	empty := &SQLServerManagedRowProof{Group: &SQLServerManagedInspection{
		Evidence: &SQLServerManagedTableEvidence{NativeMonthly: &SQLServerManagedPartitionEvidence{Partitions: []SQLServerPartitionBoundary{}}},
	}}
	clonedEmpty := CloneSQLServerManagedRowProof(empty)
	require.NotNil(t, clonedEmpty.Group.Evidence.NativeMonthly.Partitions)
	require.Empty(t, clonedEmpty.Group.Evidence.NativeMonthly.Partitions)
}

func TestSQLServerGroupInsertAcceptsNativeMonthlyProof(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	req.SQLServer.Group.Evidence.Plain = false
	req.SQLServer.Group.Evidence.PartitionInspected = true
	req.SQLServer.Group.Evidence.NativeMonthly = sqlServerInsertTestNativeMonthlyEvidence(req.TableName)
	require.NoError(t, validateSQLServerManagedRowProof(req))
}

func TestSQLServerGroupInsertRejectsInvalidNativeMonthlyProof(t *testing.T) {
	mutations := map[string]func(*SQLServerManagedTableEvidence){
		"missing proof":             func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly = nil },
		"uninspected proof":         func(evidence *SQLServerManagedTableEvidence) { evidence.PartitionInspected = false },
		"contradictory plain shape": func(evidence *SQLServerManagedTableEvidence) { evidence.Plain = true },
		"missing function":          func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.FunctionExists = false },
		"missing scheme":            func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.SchemeExists = false },
		"foreign ownership":         func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.Exclusive = false },
		"unaligned indexes":         func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.IndexesAligned = false },
		"wrong range direction":     func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.Method = "RANGE LEFT" },
		"wrong partition column":    func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.Column = "record_id" },
		"wrong scoped function":     func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.FunctionName = "pf_foreign" },
		"wrong scoped scheme":       func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.SchemeName = "ps_foreign" },
		"wrong filegroup":           func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.Filegroup = "SECONDARY" },
		"noncanonical UTC boundary": func(evidence *SQLServerManagedTableEvidence) {
			evidence.NativeMonthly.Partitions[0].LessThan += "Z"
		},
		"nonmonthly boundary": func(evidence *SQLServerManagedTableEvidence) {
			evidence.NativeMonthly.Partitions[2].LessThan = "2026-03-02T00:00:00.000000"
		},
		"wrong ordinal": func(evidence *SQLServerManagedTableEvidence) { evidence.NativeMonthly.Partitions[1].Ordinal = 3 },
		"missing future catchall": func(evidence *SQLServerManagedTableEvidence) {
			evidence.NativeMonthly.Partitions = evidence.NativeMonthly.Partitions[:5]
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			req := sqlServerGroupInsertTestRequest()
			req.SQLServer.Group.Evidence.Plain = false
			req.SQLServer.Group.Evidence.PartitionInspected = true
			req.SQLServer.Group.Evidence.NativeMonthly = sqlServerInsertTestNativeMonthlyEvidence(req.TableName)
			mutate(req.SQLServer.Group.Evidence)
			require.ErrorIs(t, validateSQLServerManagedRowProof(req), ErrGroupInsertUnsupported)
		})
	}

	t.Run("receipt rejects native proof", func(t *testing.T) {
		req := sqlServerGroupInsertTestRequest()
		req.SQLServer.Receipt.Evidence.NativeMonthly = sqlServerInsertTestNativeMonthlyEvidence(EffectReceiptTable)
		require.ErrorIs(t, validateSQLServerManagedRowProof(req), ErrGroupInsertUnsupported)
	})
}

func sqlServerInsertTestNativeMonthlyEvidence(table string) *SQLServerManagedPartitionEvidence {
	return &SQLServerManagedPartitionEvidence{
		Method: "RANGE RIGHT",
		Column: "bucket_start", FunctionName: "pf_" + table, SchemeName: "ps_" + table, Filegroup: "PRIMARY",
		Partitions: []SQLServerPartitionBoundary{
			{Name: "p_past", LessThan: "2026-01-01T00:00:00.000000", Ordinal: 1},
			{Name: "p_202601", LessThan: "2026-02-01T00:00:00.000000", Ordinal: 2},
			{Name: "p_202602", LessThan: "2026-03-01T00:00:00.000000", Ordinal: 3},
			{Name: "p_202603", LessThan: "2026-04-01T00:00:00.000000", Ordinal: 4},
			{Name: "p_202604", LessThan: "2026-05-01T00:00:00.000000", Ordinal: 5},
			{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 6},
		},
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func TestSQLServerMissingTableInspectsNativePartitionObjectsAsUnowned(t *testing.T) {
	script := nativeSQLServerCatalogFixture()
	script.base.missingObject = true
	physical, err := inspectNativeSQLServerCatalog(t, script)
	require.NoError(t, err)
	require.Equal(t, TableInspectionMissing, physical.Inspection.Status)
	require.True(t, physical.Evidence.PartitionInspected)
	require.NotNil(t, physical.Evidence.NativeMonthly)
	require.False(t, physical.Evidence.NativeMonthly.Exclusive)
	require.False(t, physical.Evidence.NativeMonthly.IndexesAligned)
	require.Equal(t, 1, script.nativeReads)
}

func TestSQLServerManagedCatalogComparisonBindsNativePartitionProof(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	group := req.SQLServer.Group
	group.Evidence.Plain = false
	group.Evidence.PartitionInspected = true
	group.Evidence.NativeMonthly = sqlServerInsertTestNativeMonthlyEvidence(req.TableName)

	require.True(t, sqlServerManagedCatalogMatches(group, sqlServerInsertTestCatalogFacts(group)))

	mutations := map[string]func(*sqlServerManagedCatalogFacts){
		"partition inspection missing": func(facts *sqlServerManagedCatalogFacts) { facts.partitionInspected = false },
		"partition function missing":   func(facts *sqlServerManagedCatalogFacts) { facts.nativeMonthly.FunctionExists = false },
		"partition scheme missing":     func(facts *sqlServerManagedCatalogFacts) { facts.nativeMonthly.SchemeExists = false },
		"partition ownership lost":     func(facts *sqlServerManagedCatalogFacts) { facts.nativeMonthly.Exclusive = false },
		"partition index misaligned":   func(facts *sqlServerManagedCatalogFacts) { facts.nativeMonthly.IndexesAligned = false },
		"partition boundary changed": func(facts *sqlServerManagedCatalogFacts) {
			facts.nativeMonthly.Partitions[1].LessThan = "2026-02-02T00:00:00.000000"
		},
		"table reported plain": func(facts *sqlServerManagedCatalogFacts) { facts.plain = true },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			facts := sqlServerInsertTestCatalogFacts(group)
			mutate(facts)
			require.False(t, sqlServerManagedCatalogMatches(group, facts))
		})
	}

	plain := sqlServerInsertTestInspection("gw_long_plain_fixture", false, 93)
	plain.Evidence.PartitionInspected = false
	plainFacts := sqlServerInsertTestCatalogFacts(plain)
	plainFacts.partitionInspected = false
	require.True(t, sqlServerManagedCatalogMatches(plain, plainFacts), "legacy plain proof remains valid")
}

func sqlServerInsertTestCatalogFacts(physical *SQLServerManagedInspection) *sqlServerManagedCatalogFacts {
	evidence := physical.Evidence
	return &sqlServerManagedCatalogFacts{
		object: sqlServerManagedObjectFacts{
			objectID: evidence.ObjectID, schema: evidence.Schema, name: physical.Inspection.Table,
			tableType: evidence.TableType, userTable: evidence.TableType == sqlServerManagedUserTableType,
		},
		columns: slices.Clone(physical.Inspection.Columns), columnCollations: maps.Clone(evidence.ColumnCollations),
		columnDefaults: maps.Clone(evidence.ColumnDefaults), primaryKey: slices.Clone(evidence.PrimaryKey),
		indexes: cloneSQLServerIndexes(physical.Indexes), plain: evidence.Plain,
		provenanceJSONChecked: evidence.ProvenanceJSONChecked, repairOperationID: evidence.RepairOperationID,
		nativeMonthly:      cloneSQLServerManagedPartitionEvidence(evidence.NativeMonthly),
		partitionInspected: evidence.PartitionInspected,
	}
}

func TestSQLServerManagedTableEvidencePreservesNativePartitionProofJSON(t *testing.T) {
	const nativeMonthly = `{"Method":"RANGE RIGHT","Column":"bucket_start","FunctionName":"pf_gw_long_fixture","SchemeName":"ps_gw_long_fixture","Filegroup":"PRIMARY","Partitions":[{"Name":"p_past","LessThan":"2026-10-01T00:00:00.000000","Ordinal":1},{"Name":"p_202610","LessThan":"2026-11-01T00:00:00.000000","Ordinal":2},{"Name":"p_202611","LessThan":"2026-12-01T00:00:00.000000","Ordinal":3},{"Name":"p_202612","LessThan":"2027-01-01T00:00:00.000000","Ordinal":4},{"Name":"p_202701","LessThan":"2027-02-01T00:00:00.000000","Ordinal":5},{"Name":"p_future","LessThan":"MAXVALUE","Ordinal":6}],"FunctionExists":true,"SchemeExists":true,"Exclusive":true,"IndexesAligned":true}`
	input := `{"NativeMonthly":` + nativeMonthly + `,"PartitionInspected":true}`

	var evidence SQLServerManagedTableEvidence
	require.NoError(t, json.Unmarshal([]byte(input), &evidence))

	encoded, err := json.Marshal(evidence)
	require.NoError(t, err)
	var roundTrip map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &roundTrip))

	native, exists := roundTrip["NativeMonthly"]
	require.True(t, exists, "native partition proof must survive SQLServerManagedTableEvidence round-trip")
	require.JSONEq(t, nativeMonthly, string(native))
	partitionInspected, exists := roundTrip["PartitionInspected"]
	require.True(t, exists, "partition inspection status must survive SQLServerManagedTableEvidence round-trip")
	require.JSONEq(t, `true`, string(partitionInspected))
}
