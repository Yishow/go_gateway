package recordingplan

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerSchemaConfirmation(t *testing.T) {
	for name, mutate := range map[string]func(*SchemaPreviewScope){
		"different schema": func(s *SchemaPreviewScope) { s.Schema = "other_schema" },
		"missing database": func(s *SchemaPreviewScope) { s.Database = "" },
		"wrong dialect":    func(s *SchemaPreviewScope) { s.Dialect = "postgres" },
	} {
		t.Run(name, func(t *testing.T) {
			svc, _ := newPreviewService()
			layout, err := ManagedSQLServerLongGroupLayout("sqlserver_long", "_gw_owner_0123456789abcdef01234567", sqlServerSchemaOptionsForTest())
			require.NoError(t, err)
			scope := groupPreviewScope()
			scope.Dialect, scope.Database, scope.Schema, scope.GroupLayout = "sqlserver", "recording", "dbo", &layout
			mutate(&scope)
			calls := 0
			token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (TargetTableInspection, error) {
				calls++
				return TargetTableInspection{Status: inspectionMissing}, nil
			})
			require.Error(t, err, "a SQL Server preview must bind one explicit database/schema/dialect")
			require.Nil(t, token)
			require.Zero(t, calls, "invalid saved scope must fail before target access")
		})
	}
}

func sqlServerSchemaScopeForTest(t *testing.T) SchemaPreviewScope {
	t.Helper()
	layout, err := ManagedSQLServerLongGroupLayout("sqlserver_long", "_gw_owner_0123456789abcdef01234567", sqlServerSchemaOptionsForTest())
	require.NoError(t, err)
	scope := groupPreviewScope()
	scope.Dialect, scope.Database, scope.Schema, scope.GroupLayout = "sqlserver", "recording", "dbo", &layout
	return scope
}

func sqlServerSchemaOptionsForTest() SQLServerGroupSchemaOptions {
	return SQLServerGroupSchemaOptions{
		Database: "recording", Schema: "dbo", Collation: sqlserverCollation,
		ServerVersion: "16.0.1000", ServerEdition: "Standard Edition", EngineEdition: 2,
		CompatibilityLevel: 160, CapabilityRevision: "sqlservercap-unit",
	}
}

func sqlServerSchemaEvidenceForTest(options SQLServerGroupSchemaOptions) *TargetSQLServerTableEvidence {
	return &TargetSQLServerTableEvidence{
		Database: options.Database, Schema: options.Schema, ServerVersion: options.ServerVersion,
		ServerEdition: options.ServerEdition, EngineEdition: options.EngineEdition,
		CompatibilityLevel: options.CompatibilityLevel, CapabilityRevision: options.CapabilityRevision,
		MetadataVerified: true,
	}
}

func sqlServerMissingStateForTest(options SQLServerGroupSchemaOptions) TargetTableInspection {
	return TargetTableInspection{Status: inspectionMissing, SQLServer: sqlServerSchemaEvidenceForTest(options)}
}

func sqlServerSchemaStateForTest(layout GroupSchemaLayout, columns []GroupSchemaColumn, owner string, objectID int64) TargetTableInspection {
	state := groupCompatibleState(GroupSchemaLayout{Columns: columns})
	state.IndexesInspected = true
	evidence := sqlServerSchemaEvidenceForTest(*layout.SQLServer)
	evidence.ObjectID, evidence.TableType, evidence.Plain = objectID, "USER_TABLE", true
	evidence.ColumnCollations, evidence.ColumnDefaults = map[string]string{}, map[string]*string{}
	for _, column := range columns {
		if column.PrimaryKey {
			evidence.PrimaryKey = append(evidence.PrimaryKey, column.Name)
		}
		if sqlserverCollatedSQLType(column.SQLType) {
			evidence.ColumnCollations[column.Name] = layout.SQLServer.Collation
		}
		if column.Name == sqlserverProvenanceColumn {
			evidence.ProvenanceJSONChecked = true
		}
	}
	gateway := managedOwnerDefault
	evidence.ColumnDefaults[owner] = &gateway
	state.SQLServer = evidence
	return state
}

func sqlServerReceiptStateForTest() TargetTableInspection {
	options := sqlServerSchemaOptionsForTest()
	layout := GroupSchemaLayout{SQLServer: &options}
	return sqlServerSchemaStateForTest(layout, sqlserverReceiptColumns(), managedEffectReceiptMarker, 22)
}

func TestSQLServerSchemaAckLost(t *testing.T) {
	// The catalog/commit event is simulated. This tests the existing operation
	// ledger's reconciliation, not actual SQL Server transaction acceptance.
	repo := NewMemoryRepository()
	svc := NewService(repo)
	scope := sqlServerSchemaScopeForTest(t)
	created, executions := false, 0
	inspect := func(_ context.Context, table string) (TargetTableInspection, error) {
		if !created {
			return sqlServerMissingStateForTest(*scope.GroupLayout.SQLServer), nil
		}
		if table == managedEffectReceiptTable {
			return sqlServerReceiptStateForTest(), nil
		}
		state := sqlServerSchemaStateForTest(*scope.GroupLayout, scope.GroupLayout.Columns, scope.GroupLayout.OwnerColumn, 11)
		for _, index := range scope.GroupLayout.Indexes {
			state.Indexes = append(state.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
		}
		return state, nil
	}
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, inspect)
	require.NoError(t, err)
	claimed, outcome, err := svc.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, token.OperationID, claimed.OperationID)
	created = true // external commit happened; acknowledgement/result was lost
	ageLastSchemaOperation(t, repo, token.OperationID)
	resolved, outcome, err := svc.ApplySchemaPreview(t.Context(), token, SchemaApplyTarget{
		Inspect: inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			executions++
			return SchemaExecution{}, errors.New("recovery must not execute DDL")
		},
	})
	require.NoError(t, err)
	require.Equal(t, ClaimCompleted, outcome)
	require.Equal(t, token.OperationID, resolved.OperationID)
	require.Equal(t, SchemaOperationSucceeded, resolved.Status)
	require.Equal(t, SchemaReasonRecovered, resolved.Reason)
	require.Zero(t, executions)
}

func TestSQLServerUnverifiedSchemaExecutionStaysClosed(t *testing.T) {
	svc, _ := newPreviewService()
	scope := sqlServerSchemaScopeForTest(t)
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (TargetTableInspection, error) {
		return sqlServerMissingStateForTest(*scope.GroupLayout.SQLServer), nil
	})
	require.NoError(t, err)
	executions := 0
	_, _, err = svc.ApplySchemaPreview(t.Context(), token, SchemaApplyTarget{
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			executions++
			return SchemaExecution{}, nil
		},
	})
	require.ErrorIs(t, err, ErrSchemaAdapterUnverified)
	require.Zero(t, executions)
}
