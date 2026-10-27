package recordingplan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerManagedApplyUsesTypedTarget(t *testing.T) {
	scope := sqlServerSchemaScopeForTest(t)
	layout := *scope.GroupLayout
	created := false
	receipt := sqlServerReceiptStateForTest()
	groupState := func() TargetTableInspection {
		state := sqlServerSchemaStateForTest(layout, layout.Columns, layout.OwnerColumn, 101)
		for _, index := range layout.Indexes {
			state.Indexes = append(state.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
		}
		return state
	}
	svc := NewService(NewMemoryRepository())
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return receipt, nil
		}
		return sqlServerMissingStateForTest(*layout.SQLServer), nil
	})
	require.NoError(t, err)
	executions := 0
	op, outcome, err := svc.ApplySQLServerSchemaPreview(t.Context(), token, SQLServerSchemaApplyTarget{
		Inspect: func(_ context.Context, table string) (TargetTableInspection, error) {
			if table == managedEffectReceiptTable {
				return receipt, nil
			}
			if created {
				return groupState(), nil
			}
			return sqlServerMissingStateForTest(*layout.SQLServer), nil
		},
		Execute: func(_ context.Context, got *SchemaPreviewToken, statements []string) (SchemaExecution, error) {
			require.Same(t, token, got)
			require.NotEmpty(t, statements)
			executions++
			created = true
			return SchemaExecution{Committed: len(statements)}, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, SchemaOperationSucceeded, op.Status)
	require.Equal(t, 1, executions)
}
