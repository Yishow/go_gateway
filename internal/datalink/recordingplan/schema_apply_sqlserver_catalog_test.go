package recordingplan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func sqlServerCompleteGroupForTest(scope SchemaPreviewScope, objectID int64) TargetTableInspection {
	state := sqlServerSchemaStateForTest(*scope.GroupLayout, scope.GroupLayout.Columns, scope.GroupLayout.OwnerColumn, objectID)
	for _, index := range scope.GroupLayout.Indexes {
		state.Indexes = append(state.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	return state
}

func TestSQLServerSchemaConfirmationRejectsReplacedObject(t *testing.T) {
	svc, _ := newPreviewService()
	scope := sqlServerSchemaScopeForTest(t)
	objectID := int64(11)
	inspect := func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			return sqlServerReceiptStateForTest(), nil
		}
		return sqlServerCompleteGroupForTest(scope, objectID), nil
	}
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, inspect)
	require.NoError(t, err)
	objectID = 12 // same shape and owner marker, different physical object
	_, err = confirmPreviewMatchesTarget(t.Context(), token, inspect)
	require.ErrorIs(t, err, ErrPreviewTokenStale)
}

func TestSQLServerSchemaRecoveryPreservesExistingReceiptIdentity(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	scope := sqlServerSchemaScopeForTest(t)
	created, replacedReceipt := false, false
	inspect := func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == managedEffectReceiptTable {
			state := sqlServerReceiptStateForTest()
			if replacedReceipt {
				state.SQLServer.ObjectID++
			}
			return state, nil
		}
		if !created {
			return sqlServerMissingStateForTest(*scope.GroupLayout.SQLServer), nil
		}
		return sqlServerCompleteGroupForTest(scope, 11), nil
	}
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, inspect)
	require.NoError(t, err)
	_, claim, err := svc.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, claim)
	created, replacedReceipt = true, true
	ageLastSchemaOperation(t, repo, token.OperationID)
	op, _, err := svc.ApplySchemaPreview(t.Context(), token, SchemaApplyTarget{Inspect: inspect})
	require.NoError(t, err)
	require.Equal(t, SchemaOperationUnknown, op.Status, "a replacement receipt table cannot prove preservation of preexisting receipts")
	require.Equal(t, SchemaReasonUnverifiable, op.Reason)
}
