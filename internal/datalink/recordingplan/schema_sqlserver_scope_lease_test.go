package recordingplan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerSchemaLeaseKeepsPhysicalScope(t *testing.T) {
	svc, _ := newPreviewService()
	firstScope := sqlServerSchemaScopeForTest(t)
	prepare := func(scope SchemaPreviewScope) *SchemaPreviewToken {
		t.Helper()
		token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (TargetTableInspection, error) {
			return sqlServerMissingStateForTest(*scope.GroupLayout.SQLServer), nil
		})
		require.NoError(t, err)
		return token
	}
	first := prepare(firstScope)
	secondScope := firstScope
	secondScope.GroupLayout = clonePreviewToken(SchemaPreviewToken{GroupLayout: firstScope.GroupLayout}).GroupLayout
	secondScope.GroupLayout.SQLServer.CapabilityRevision = "sqlservercap-changed"
	second := prepare(secondScope)
	require.Equal(t, first.ScopeKey(), second.ScopeKey(), "new capability facts must not bypass the lease on the same physical table")
	firstOperation, result, err := svc.ClaimSchemaApply(t.Context(), first)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, result)
	secondOperation, result, err := svc.ClaimSchemaApply(t.Context(), second)
	require.NoError(t, err)
	require.Equal(t, ClaimScopeBusy, result)
	require.Equal(t, firstOperation.OperationID, secondOperation.OperationID)
}
