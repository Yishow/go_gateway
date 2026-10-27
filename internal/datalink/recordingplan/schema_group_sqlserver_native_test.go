package recordingplan

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func sqlServerNativeTestLayout(t *testing.T) GroupSchemaLayout {
	t.Helper()
	options := sqlServerSchemaOptionsForTest()
	native, err := NewSQLServerNativeMonthlyPartition(
		time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
		"pf_sqlserver_long", "ps_sqlserver_long",
	)
	require.NoError(t, err)
	options.NativeMonthly = native
	layout, err := ManagedSQLServerLongGroupLayout("sqlserver_long", "_gw_owner_0123456789abcdef01234567", options)
	require.NoError(t, err)
	return layout
}

func sqlServerNativeTargetEvidence(layout GroupSchemaLayout, exclusive, aligned bool) *TargetSQLServerPartitionEvidence {
	native := layout.SQLServer.NativeMonthly
	return &TargetSQLServerPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: append([]SQLServerPartitionBoundary(nil), native.Partitions...),
		FunctionExists: true, SchemeExists: true, Exclusive: exclusive, IndexesAligned: aligned,
	}
}

func TestSQLServerNativeMonthlyLayoutAndDigest(t *testing.T) {
	layout := sqlServerNativeTestLayout(t)
	require.Len(t, layout.SQLServer.NativeMonthly.Partitions, 6)
	require.Equal(t, "p_past", layout.SQLServer.NativeMonthly.Partitions[0].Name)
	require.Equal(t, "2026-10-01T00:00:00.000000", layout.SQLServer.NativeMonthly.Partitions[0].LessThan)
	require.Equal(t, "p_202610", layout.SQLServer.NativeMonthly.Partitions[1].Name)
	require.Equal(t, "2026-11-01T00:00:00.000000", layout.SQLServer.NativeMonthly.Partitions[1].LessThan)
	require.Equal(t, "p_202701", layout.SQLServer.NativeMonthly.Partitions[4].Name)
	require.Equal(t, "2027-02-01T00:00:00.000000", layout.SQLServer.NativeMonthly.Partitions[4].LessThan)
	require.Equal(t, "p_future", layout.SQLServer.NativeMonthly.Partitions[5].Name)
	require.Equal(t, "MAXVALUE", layout.SQLServer.NativeMonthly.Partitions[5].LessThan)

	group := sqlServerSchemaStateForTest(layout, layout.Columns, layout.OwnerColumn, 101)
	for _, index := range layout.Indexes {
		group.Indexes = append(group.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	group.SQLServer.NativeMonthly = sqlServerNativeTargetEvidence(layout, true, true)
	group.SQLServer.Plain = false
	group.SQLServer.PartitionInspected = true
	nativeDigest, err := sqlserverCatalogDigest(group)
	require.NoError(t, err)
	plainTarget := cloneTargetTableInspectionForTest(group.SQLServer)
	plainTarget.NativeMonthly = nil
	plainTarget.PartitionInspected = false
	plainDigest, err := sqlserverCatalogDigest(TargetTableInspection{Status: group.Status, Columns: group.Columns, SQLServer: plainTarget})
	require.NoError(t, err)
	require.NotEqual(t, nativeDigest, plainDigest)
	receipt := sqlServerReceiptStateForTest()
	digest, err := VerifyManagedSQLServerCatalog(layout, layout.SQLServer.Database, group, receipt)
	require.NoError(t, err)
	require.NotEmpty(t, digest)

	foreign := group
	foreign.SQLServer = cloneTargetTableInspectionForTest(group.SQLServer)
	foreign.SQLServer.NativeMonthly.Exclusive = false
	_, err = VerifyManagedSQLServerCatalog(layout, layout.SQLServer.Database, foreign, receipt)
	require.ErrorIs(t, err, ErrSQLServerPartitionUnavailable)

	plain := sqlServerNativeTestLayout(t)
	plain.SQLServer.NativeMonthly = nil
	encoded, err := json.Marshal(plain)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "native_monthly")
}

func TestSQLServerNativePartitionPermissionUnavailable(t *testing.T) {
	layout := sqlServerNativeTestLayout(t)
	scope := sqlServerSchemaScopeForTest(t)
	scope.GroupLayout = &layout
	svc := NewService(NewMemoryRepository())
	_, err := svc.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (TargetTableInspection, error) {
		return TargetTableInspection{}, ErrTargetPermissionDenied
	})
	require.ErrorIs(t, err, ErrSQLServerPartitionUnavailable)
}

func TestSQLServerNativeRejectsPlainExistingTable(t *testing.T) {
	layout := sqlServerNativeTestLayout(t)
	scope := sqlServerSchemaScopeForTest(t)
	scope.GroupLayout = &layout
	group := sqlServerSchemaStateForTest(layout, layout.Columns, layout.OwnerColumn, 101)
	for _, index := range layout.Indexes {
		group.Indexes = append(group.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	group.SQLServer.NativeMonthly = nil
	group.SQLServer.PartitionInspected = true
	receipt := sqlServerReceiptStateForTest()
	svc := NewService(NewMemoryRepository())
	_, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return receipt, nil
		}
		return group, nil
	})
	require.ErrorIs(t, err, ErrSQLServerPartitionUnavailable)
}

func TestSQLServerNativePreviewAdoptsOnlyExactOwnedPartitionObjects(t *testing.T) {
	layout := sqlServerNativeTestLayout(t)
	scope := sqlServerSchemaScopeForTest(t)
	scope.GroupLayout = &layout
	group := sqlServerMissingStateForTest(*layout.SQLServer)
	group.SQLServer.PartitionInspected = true
	group.SQLServer.NativeMonthly = sqlServerNativeTargetEvidence(layout, true, false)
	group.SQLServer.Plain = false
	receipt := sqlServerReceiptStateForTest()
	svc := NewService(NewMemoryRepository())
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return receipt, nil
		}
		return group, nil
	})
	require.NoError(t, err)
	joined := joinStatements(token.Statements)
	require.NotContains(t, joined, "CREATE PARTITION FUNCTION")
	require.NotContains(t, joined, "CREATE PARTITION SCHEME")
	require.Contains(t, joined, "CREATE TABLE [dbo].[sqlserver_long]")
	require.NotNil(t, token.Tables[0].SQLServerPartition)
	require.Equal(t, SchemaTableActionUnchanged, token.Tables[0].SQLServerPartition.Action)
}

func TestSQLServerNativePreviewCreatesOwnedPartitionObjects(t *testing.T) {
	layout := sqlServerNativeTestLayout(t)
	scope := sqlServerSchemaScopeForTest(t)
	scope.GroupLayout = &layout
	group := sqlServerMissingStateForTest(*layout.SQLServer)
	group.SQLServer.PartitionInspected = true
	receipt := sqlServerReceiptStateForTest()

	svc := NewService(NewMemoryRepository())
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return receipt, nil
		}
		return group, nil
	})
	require.NoError(t, err)
	require.Len(t, token.Tables, 2)
	require.NotNil(t, token.Tables[0].SQLServerPartition)
	require.Equal(t, SchemaTableActionCreate, token.Tables[0].SQLServerPartition.Action)
	require.Equal(t, "2026-10-01T00:00:00.000000", token.Tables[0].SQLServerPartition.Partitions[0].LessThan)
	require.Contains(t, joinStatements(token.Statements), "CREATE PARTITION FUNCTION")
	require.Contains(t, joinStatements(token.Statements), "CREATE PARTITION SCHEME")
	require.Contains(t, joinStatements(token.Statements), "ON [ps_sqlserver_long] ([bucket_start])")
	require.NotContains(t, joinStatements(token.Statements), "CREATE TABLE [dbo].[gw_effect_receipts]")

	existing := sqlServerSchemaStateForTest(layout, layout.Columns, layout.OwnerColumn, 101)
	for _, index := range layout.Indexes {
		existing.Indexes = append(existing.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	existing.SQLServer.NativeMonthly = sqlServerNativeTargetEvidence(layout, true, true)
	existing.SQLServer.Plain = false
	existing.SQLServer.PartitionInspected = true
	token, err = svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return receipt, nil
		}
		return existing, nil
	})
	require.NoError(t, err)
	require.Empty(t, token.Statements)
	require.NotNil(t, token.Tables[0].SQLServerPartition)
	require.Equal(t, SchemaTableActionUnchanged, token.Tables[0].SQLServerPartition.Action)
}

func TestSQLServerNativeForeignPartitionUnavailable(t *testing.T) {
	layout := sqlServerNativeTestLayout(t)
	scope := sqlServerSchemaScopeForTest(t)
	scope.GroupLayout = &layout
	group := sqlServerSchemaStateForTest(layout, layout.Columns, layout.OwnerColumn, 101)
	for _, index := range layout.Indexes {
		group.Indexes = append(group.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	group.SQLServer.NativeMonthly = sqlServerNativeTargetEvidence(layout, false, true)
	group.SQLServer.Plain = false
	group.SQLServer.PartitionInspected = true
	receipt := sqlServerReceiptStateForTest()
	svc := NewService(NewMemoryRepository())
	_, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return receipt, nil
		}
		return group, nil
	})
	require.ErrorIs(t, err, ErrSQLServerPartitionUnavailable)
}

func joinStatements(statements []string) string {
	result := ""
	for _, statement := range statements {
		result += statement + "\n"
	}
	return result
}

func cloneTargetTableInspectionForTest(evidence *TargetSQLServerTableEvidence) *TargetSQLServerTableEvidence {
	cloned := *evidence
	if evidence.NativeMonthly != nil {
		native := *evidence.NativeMonthly
		native.Partitions = append([]SQLServerPartitionBoundary(nil), evidence.NativeMonthly.Partitions...)
		cloned.NativeMonthly = &native
	}
	return &cloned
}
